FROM gcr.io/distroless/static:nonroot
COPY /docked /
ENTRYPOINT ["/docked"]
