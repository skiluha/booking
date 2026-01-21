FROM golang:1.24-bookworm

WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o booking ./cmd/api


CMD ["./booking"]
