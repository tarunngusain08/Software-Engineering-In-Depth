FROM golang:1.20-alpine


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o app .

EXPOSE 8081

CMD ["./app"]
