FROM scratch

FROM scratch

COPY --from=0 --chown=12345 /bin/sh bin/sh
