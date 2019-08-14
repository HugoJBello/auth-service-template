
# build it:
# docker build -t go-docker .
# run it:
# docker run -d -p 8000:8000 go-docker
# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.12 base image
FROM golang:1.12

# Add Maintainer Info
LABEL maintainer="Hugo J. Bello <hjbello.wk@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR $home/test

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8000

# Run the executable
CMD ["auth-service-template"]