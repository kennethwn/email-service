FROM golang:1.25 as builder

WORKDIR /app

COPY go.mod go.sum config.yml ./

RUN go mod tidy

COPY . .

RUN go build -o main cmd/main.go

EXPOSE 8083

CMD ["./main"]