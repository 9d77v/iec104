FROM golang:1.15.6-alpine3.12

LABEL maintainer="9d77v@9d77v.me"

WORKDIR /app
ENV GOPROXY=https://goproxy.cn

RUN echo "http://mirrors.aliyun.com/alpine/v3.12/main/" > /etc/apk/repositories \
    && echo "http://mirrors.aliyun.com/alpine/v3.12/community/" >> /etc/apk/repositories \
    && apk add --no-cache upx tzdata\
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" >  /etc/timezone \
    && apk del tzdata 

COPY go.mod /app
COPY go.sum /app

RUN go mod download \
    && rm -rf /app
