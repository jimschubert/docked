FROM gcr.io/distroless/base-debian10:nonroot
COPY /docked /
ENTRYPOINT ["/docked"]
