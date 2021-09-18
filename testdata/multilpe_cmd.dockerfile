FROM gcr.io/distroless/base-debian10:nonroot
COPY --from=builder /go/bin/app /
ENTRYPOINT ["/app"]
CMD "first"
CMD "second"
