# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

ENV GO111MODULE=on

WORKDIR /backend

COPY ./backend .
RUN go mod tidy

RUN go build ./cmd/backend

EXPOSE 8080

CMD ["./backend"]