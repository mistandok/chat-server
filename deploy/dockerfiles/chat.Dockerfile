FROM golang:alpine3.19 AS builder

WORKDIR /build

COPY go.mod .

COPY . .

RUN go build -o chat_service ./cmd/service/main.go

FROM ubuntu:20.04

ENV APP_DIR /build

WORKDIR $APP_DIR

RUN apt-get update \
    && groupadd -r web \
    && useradd -d $APP_DIR -r -g web web \
    && chown web:web -R $APP_DIR \
    && apt-get install -y netcat-traditional \
    && apt-get install -y acl

COPY --from=builder /build/chat_service $APP_DIR/chat_service
COPY --from=builder /build/deploy/scripts/prod-chat-service-start.sh $APP_DIR/prod-chat-service-start.sh
COPY --from=builder /build/deploy/env/.env.prod $APP_DIR//deploy/env/.env.prod

RUN setfacl -R -m u:web:rwx $APP_DIR/prod-chat-service-start.sh

USER web

ENTRYPOINT ["bash", "prod-chat-service-start.sh"]
