FROM alpine:3.14.2

RUN apk --no-cache update && apk --no-cache add curl
