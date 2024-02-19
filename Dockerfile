FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /app/nature-remo-prometheus

CMD ["/app/nature-remo-prometheus"]