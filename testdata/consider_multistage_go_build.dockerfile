FROM golang:1.22-alpine
ENV GOOS=linux \
    GOARCH=386 \
    CGO_ENABLED=0

RUN go build -o docked cmd/docked/
CMD ["docked"]
