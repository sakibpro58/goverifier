package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	SMTP_PORT    = "25"
	SMTP_TIMEOUT = 30 * time.Second
)

type VerifyResult struct {
	Result       string `json:"result,omitempty"`
	MailboxExist bool   `json:"mailbox_exists"`
	IsCatchAll   bool   `json:"is_catch_all"`
	IsDisposable bool   `json:"is_disposable"`
	Email        string `json:"email"`
	Domain       string `json:"domain"`
	User         string `json:"user"`

	Client *smtp.Client `json:"-"`
}

// Add Proxy Dialer
func dialWithProxy(network, addr string) (net.Conn, error) {
	proxyURL := os.Getenv("SOCKS5_PROXY_URL")
	if proxyURL == "" {
		// Default Dialer if no proxy is set
		return net.DialTimeout(network, addr, SMTP_TIMEOUT)
	}

	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	dialer, err := net.Dial("tcp", parsedURL.Host)
	if err != nil {
		return nil, err
	}

	return dialer, nil
}

func (self *VerifyResult) ConnectSmtp() error {
	mx, err := net.LookupMX(self.Domain)
	if err != nil || len(mx) == 0 {
		self.Result = "NoMxServersFound"
		return err
	}

	addr := mx[0].Host + ":" + SMTP_PORT
	conn, err := dialWithProxy("tcp", addr)
	if err != nil {
		self.Result = "ConnectionRefused"
		return err
	}

	client, err := smtp.NewClient(conn, mx[0].Host)
	if err != nil {
		self.Result = "NoMxServersFound"
		return err
	}

	// Attempt STARTTLS if supported
	tlsConfig := &tls.Config{ServerName: mx[0].Host}
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(tlsConfig); err != nil {
			self.Result = "TLSFailed"
			return err
		}
	}

	self.Client = client
	err = self.Client.Hello("example.com")
	if err != nil {
		self.Result = "HelloFailed"
		return err
	}

	return nil
}

func (self *VerifyResult) ParseEmailAddress() error {
	pieces := strings.Split(self.Email, "@")
	if len(pieces) == 2 {
		self.User = pieces[0]
		self.Domain = pieces[1]
		return nil
	}

	self.Result = "InvalidEmailAddress"
	return errors.New("Invalid email address")
}

func (self *VerifyResult) CheckMailboxExist() {
	self.MailboxExist = addressExists(self.Client, self.Email)
}

func (self *VerifyResult) CheckIsCatchAll() {
	randomAddress := "n0n3x1st1ng4ddr355@" + self.Domain
	self.IsCatchAll = addressExists(self.Client, randomAddress)
}

func (self *VerifyResult) Verify() {
	if err := self.ParseEmailAddress(); err != nil {
		return
	}

	if err := self.ConnectSmtp(); err != nil {
		log.Printf("%s\n", err.Error())
		return
	}

	defer self.Client.Quit()
	self.CheckMailboxExist()
	if self.MailboxExist {
		self.CheckIsCatchAll()
	}
	self.CheckIsDisposable()
}

func (self *VerifyResult) CheckIsDisposable() {
	b, err := Asset("list.txt")
	if err != nil {
		panic(err)
	}
	self.IsDisposable = strings.Contains(string(b), self.Domain)
}

func addressExists(client *smtp.Client, address string) bool {
	if err := client.Mail(address); err != nil {
		return false
	}
	if err := client.Rcpt(address); err != nil {
		return false
	}
	return true
}
