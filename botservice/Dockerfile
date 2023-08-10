FROM golang:alpine

WORKDIR /bot-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build

CMD ["./botservice"]