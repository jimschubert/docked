FROM scratch
COPY . /

FROM ubuntu
COPY --from=0 README.md /README.md
ENTRYPOINT bash
