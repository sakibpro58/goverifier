FROM golang:1.21-alpine

WORKDIR /app

# Install go dependencies if go.mod does not exist
RUN if [ ! -f go.mod ]; then go mod init goverifier; fi

# Copy the entire project
COPY . .

# Install dependencies and build the Go app
RUN go mod tidy && go build -o goverifier .

EXPOSE 8080

ENTRYPOINT ["/app/goverifier"]
