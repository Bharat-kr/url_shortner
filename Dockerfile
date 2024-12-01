FROM golang:1.23.3

WORKDIR /app

# Install air
RUN go install github.com/air-verse/air@latest

COPY . .

RUN go mod tidy

RUN go mod download

# Ensure air is in the PATH
ENV PATH="/go/bin:${PATH}"