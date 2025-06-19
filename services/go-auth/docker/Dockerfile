FROM golang:1.24-alpine as bulder

WORKDIR /usr/local/src/

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . ./

RUN go build -o ./bin/app main.go

FROM alpine

RUN apk update && \
    apk upgrate -U && \
    apk add ca-certificares && \
    rm -rf /var/caches/*

COPY --from=bulder /usr/local/src/bin/app /

CMD [ "/app" ]
