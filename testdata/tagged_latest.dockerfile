FROM gcr.io/distroless/base-debian10:latest
COPY --from=builder /go/bin/app /
ENTRYPOINT ["/app"]
