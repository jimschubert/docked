FROM golang:1.22-alpine as builder
ENV GOOS=linux \
    GOARCH=386 \
    CGO_ENABLED=0

WORKDIR /go/src/app
ADD . /go/src/app

# Install git and deps
RUN apk --no-cache add gcc g++ make ca-certificates && apk add git

RUN go mod download && go build -o /go/bin/app

FROM gcr.io/distroless/base-debian10:nonroot
COPY --from=builder /go/bin/app /
ENTRYPOINT ["/app"]