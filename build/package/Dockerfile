FROM 9d77v/iec104:base-0.0.2

ENV APP_NAME=iec104

COPY . /app

RUN cd /app/example/client \
    && CGO_ENABLED=0 go build -o $APP_NAME -ldflags "-s -w" \
    && upx -9 $APP_NAME

FROM 9d77v/iec104:base-deploy-0.0.1

ENV APP_NAME=iec104
COPY --from=0  /app/example/client/$APP_NAME /app/$APP_NAME

WORKDIR /app

CMD ["./iec104"]