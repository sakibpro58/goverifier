# Set the Go version and container image name
GO_VERSION=1.21
IMAGE_NAME=goverifier

# Build the Go app and the Docker image
build:
    go mod tidy
    go build -o goverifier .

docker-build:
    docker build -t $(IMAGE_NAME) .

# Run the app
run:
    ./goverifier

# Clean the project
clean:
    rm goverifier
