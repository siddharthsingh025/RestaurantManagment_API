# use golang official image
FROM golang:1.20.6-alpine

# set working directory 
WORKDIR /app

# Copy the source code 
COPY  . .

# Download and install dependencies
RUN go get -d -v  ./...

# Build the Go application
RUN go build -o api .

# EXPOSE the port
EXPOSE 8000

# Run the executable
CMD [ "./api" ]


