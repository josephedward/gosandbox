# FROM golang:onbuild
FROM ubuntu
FROM golang 

# Add a work directory
WORKDIR /app

# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
# Install Reflex for development
RUN go install github.com/cespare/reflex@latest

# Expose port
EXPOSE 4000
# Start app
CMD reflex -g '*.go' go run main.go --start-service