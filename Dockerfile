FROM golang:1.14.4-alpine3.12

WORKDIR /go/src/app

ENV PORT 8080

RUN apk add --no-cache git \
    && go get -u github.com/Fukkatsuso/sudoku \
    && go get github.com/oxequa/realize

EXPOSE 8080