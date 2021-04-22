FROM alpine:3.6

RUN sed -i 's/dl-cdn\.alpinelinux\.org/mirrors\.aliyun\.com/g' /etc/apk/repositories

RUN apk update --no-cache

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

RUN mkdir /main
COPY main /main/
COPY config-debug.yaml /main/

WORKDIR /main
ENTRYPOINT ["./main"]
