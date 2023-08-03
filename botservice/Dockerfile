FROM golang:1.18

WORKDIR /bot-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build

CMD ["./go-eats"]