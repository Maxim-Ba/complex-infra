FROM golang:1.24-alpine as bulder

WORKDIR /usr/local/src/

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . ./

RUN go build -o ./bin/app ./cmd/main.go

FROM alpine

RUN apk --no-cache update && \
  apk --no-cache upgrade && \
  apk --no-cache add ca-certificates \
  rm -rf /var/cache/*

COPY --from=bulder /usr/local/src/bin/app /

CMD [ "/app" ]
