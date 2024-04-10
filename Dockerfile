# Base image for building the go project
FROM golang:alpine AS build

# Updates the repository and installs git
RUN apk update && apk upgrade && \
    apk add --no-cache git

# Switches to /tmp/app as the working directory, similar to 'cd'
WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Builds the current project to a binary file
RUN GOOS=linux go build -o ./out/nodebin .

# lightweight alpine image to run the server 
FROM alpine:latest

# root ONLY FOR TESTING MOUNT PERMISSIONS
USER root 
RUN apk add ca-certificates nfs-utils iputils
RUN set -ex && apk --no-cache add sudo
RUN echo '%wheel ALL=(ALL) ALL' > /etc/sudoers.d/wheel
RUN adduser $USER wheel

# Copies the binary file from the BUILD container to /app folder
COPY --from=build /tmp/app/out/nodebin /app/nodebin

# Switches working directory to /app
WORKDIR "/app"

# Exposes the 5000 port from the container
EXPOSE 5000

# Runs the binary once the container starts
CMD ["./nodebin"]

