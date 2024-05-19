# Use the official Golang image to create a build environment
FROM golang:1.21.1 as builder


# Set the Current Working Directory inside the container
WORKDIR /workdir


# Copy go mod and sum files
COPY go.mod go.sum ./


# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download


# Copy the source code into the container
COPY . .


# Environment variables for cross-compilation
ENV GOOS=linux GOARCH=amd64
RUN go build -o /workdir/output/taskapp-linux-amd64


# Use a smaller base image for the final stage (not necessary to run the app, just to copy the binaries)
FROM scratch


# Set the Current Working Directory inside the container
WORKDIR /workdir


# Copy the binaries from the builder stage
COPY --from=builder /workdir .

CMD ["/workdir/output/taskapp-linux-amd64", "/workdir/test_file.txt"]
