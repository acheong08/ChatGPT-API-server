FROM golang:1.19-alpine AS build

RUN apk add --no-cache git
RUN apk add gcc build-base

RUN mkdir -p /go/src/ChatGPT-API-server
WORKDIR /go/src/ChatGPT-API-server

RUN git clone https://github.com/ChatGPT-Hackers/ChatGPT-API-server/ .
RUN go install .

FROM alpine:latest
COPY --from=build /go/bin/ChatGPT-API-server /usr/local/bin/

RUN apk add --no-cache curl
