FROM golang:1.14

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    BOT_TOKEN=__TOKEN__

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

EXPOSE 3000

CMD ["/dist/main"]