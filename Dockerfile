FROM 9d77v/iec104:base-v0.0.1

ENV APP_NAME=iec104

COPY . /app

RUN cd /app/example/client \
    && go build -o $APP_NAME -ldflags "-s -w" \
    && upx -9 $APP_NAME

FROM alpine:3.10

LABEL maintainer="9d77v@9d77v.me"

ENV APP_NAME=iec104
COPY --from=0  /app/example/client/$APP_NAME /app/$APP_NAME

RUN echo "http://mirrors.aliyun.com/alpine/v3.10/main/" > /etc/apk/repositories \
    && echo "http://mirrors.aliyun.com/alpine/v3.10/community/" >> /etc/apk/repositories \
    && apk add  --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" >  /etc/timezone \
    && apk del tzdata 

WORKDIR /app

CMD ["./iec104"]