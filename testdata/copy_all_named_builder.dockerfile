FROM scratch AS builder
COPY . /

FROM ubuntu
COPY --from=builder README.md /README.md
ENTRYPOINT bash
