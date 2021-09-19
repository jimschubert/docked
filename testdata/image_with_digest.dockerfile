FROM gcr.io/distroless/base-debian10@sha256:a74f307185001c69bc362a40dbab7b67d410a872678132b187774fa21718fa13
COPY --from=builder /go/bin/app /
ENTRYPOINT ["/app"]
