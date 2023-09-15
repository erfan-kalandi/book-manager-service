FROM golang:1.20-alpine as builder

ENV GO111MODULE=on
ENV GOOS=linux
ENV CGO_ENABLED=0

RUN mkdir Bookmanager
WORKDIR /Bookmanager
COPY go.sum .
COPY go.mod .
RUN go mod download
COPY . .
Run go build -o Bookmanager


FROM alpine:3.18.2

RUN mkdir executable
WORKDIR /executable
COPY --from=builder /Bookmanager/Bookmanager .
CMD "./Bookmanager"

