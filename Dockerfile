# syntax=docker/dockerfile:1

FROM golang:1.22

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# https://docs.docker.com/reference/dockerfile/#copy
# Copy the entire project directory into the Docker image
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /file-manager-node


# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 50051

# Run
CMD ["/file-manager-node"]