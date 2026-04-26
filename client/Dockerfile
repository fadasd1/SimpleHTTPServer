FROM golang:1.26

WORKDIR /app

COPY . .

RUN go build -o client .

CMD ["./client"]