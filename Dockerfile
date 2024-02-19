FROM golang:1.22

WORKDIR /usr/local/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/nature-remo-prometheus ./...

CMD ["nature-remo-prometheus"]