FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o goverifier .

EXPOSE 8080

ENTRYPOINT ["/app/goverifier"]
