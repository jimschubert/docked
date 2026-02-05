FROM gcr.io/distroless/static-debian12:nonroot
COPY /docked /
ENTRYPOINT ["/docked"]
