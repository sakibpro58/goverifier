FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o goverify .

EXPOSE 8080

ENTRYPOINT ["/app/goverify"]
