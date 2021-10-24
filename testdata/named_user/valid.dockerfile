FROM scratch

FROM scratch

COPY --from=0 --chown=jim /bin/sh bin/sh

USER jim
