# Start from golang base image
FROM golang:alpine as builder

# ENV GO111MODULE=on

# Add Maintainer info
LABEL maintainer="Udhayan G <udhaysekar@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container 
WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .       

# Expose port 3000 to the outside world
EXPOSE 3000

#Command to run the executable
CMD ["./main"]























## Old setting ##

## We specify the base image we need for our
## go application
##FROM golang:1.16 
## We create an /app directory within our
## image that will hold our application source
## files
##RUN mkdir /app
## We copy everything in the root directory
## into our /app directory
##ADD . /app
## We specify that we now wish to execute 
## any further commands inside our /app
## directory
##WORKDIR /app
## we run go build to compile the binary
## executable of our Go program
##RUN go build -o main .
## Our start command which kicks off
## our newly created binary executable
##CMD ["/app/main"]