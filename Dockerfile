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

# update apt and install libx11-dev and anything needed for chromium to run
RUN apt-get update && apt-get install -y libx11-dev
RUN apt-get install -y libxss1 libappindicator1 
RUN apt-get update && apt-get install -y build-essential && apt-get install -y libgl1-mesa-dev libnss3-dev xvfb libx11-dev
RUN apt-get install -y chromium
RUN Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &
RUN export DISPLAY=:0.0

# Expose port
EXPOSE 4000
EXPOSE 7317 

# Start app
CMD reflex -g '*.go' go run main.go --start-service

