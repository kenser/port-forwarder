FROM golang:alpine as builder
RUN apk update && apk upgrade && apk add --no-cache git build-base make
COPY . /go/src
WORKDIR /go/src
RUN go mod download
RUN buildflags="-X 'main.BuildTime=`date`' -X 'main.GitHead=`git rev-parse --short HEAD`' -X 'main.GoVersion=$(go version)'" && go build -ldflags "$buildflags" -o go-portforward
FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates tzdata
COPY from builder /go/src/go-portforward /app/
ENV ENV prod
EXPOSE 80
VOLUME /app/data
CMD ["/app/go-portforward"]
