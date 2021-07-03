# Rules
*  [D5:no-debian-frontend](#d5:no-debian-frontend)
*  [D5:secret-aws-access-key](#d5:secret-aws-access-key)
*  [D5:secret-aws-secret-access-key](#d5:secret-aws-secret-access-key)
*  [D7:tagged-latest](#d7:tagged-latest)
*  [D7:tagged-latest-builder](#d7:tagged-latest-builder)
*  [DA:maintainer-deprecated](#da:maintainer-deprecated)
*  [DC:consider-multistage](#dc:consider-multistage)
*  [DC:curl-without-fail](#dc:curl-without-fail)
*  [DC:gpg-without-batch](#dc:gpg-without-batch)
*  [DC:layered-ownership-change](#dc:layered-ownership-change)


## D5:no-debian-frontend

> _Avoid DEBIAN_FRONTEND, which affects derived images and docker run. Change this to an ARG._

Found a string matching the Pattern `\bDEBIAN_FRONTEND\b`

Priority: **Critical**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#env">ENV</a></kbd>

## D5:secret-aws-access-key

> _Secrets should not be stored directly in the Dockerfile. You should remove and rotate any secrets used here._

Found a string matching the Pattern `\bAK[A-Z0-9]{18}\b`

Priority: **Critical**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#env">ENV</a></kbd>

## D5:secret-aws-secret-access-key

> _Secrets should not be stored directly in the Dockerfile. You should remove and rotate any secrets used here._

Found a string matching the Pattern `\b[A-Za-z0-9/+=]{40}\b`

Priority: **Critical**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#env">ENV</a></kbd>

## D7:tagged-latest

> _Avoid using images tagged as Latest in production builds_

Docker best practices suggest avoiding &#39;latest&#39; images in production builds

Priority: **High**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#from">FROM</a></kbd>

## D7:tagged-latest-builder

> _Avoid using images tagged as Latest in builder stages_

Using &#39;latest&#39; images in builders is not recommended.

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#from">FROM</a></kbd>

## DA:maintainer-deprecated

> _MAINTAINER instruction is deprecated; Use LABEL instead, which can be queried via &#39;docker inspect&#39;._

Found a string matching the Pattern `[[:graph:]]+`

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#maintainer">MAINTAINER</a></kbd>

## DC:consider-multistage

> _Consider using multi-stage builds for complex operations like building code._

A multi-stage build can reduce the final image size by building necessary components or downloading large archives in a separate build context. This can help keep your final image lean.

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#run">RUN</a></kbd><kbd><a href="https://docs.docker.com/engine/reference/builder/#from">FROM</a></kbd>

## DC:curl-without-fail

> _Avoid using curl without the silent failing option -f/--fail_

Invoking curl without -f/--fail may result in incorrect, missing or stale data, which is a security concern. Ignore this rule only if you&#39;re handling server errors or verifying file contents separately.

Priority: **Critical**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#run">RUN</a></kbd>

## DC:gpg-without-batch

> _GPG call without --batch (or --no-tty) may error._

Running GPG without --batch (or --no-tty) may cause GPG to fail opening /dev/tty, resulting in docker build failures.

Priority: **Medium**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#run">RUN</a></kbd>

## DC:layered-ownership-change

> _Some storage drivers may have issues with ownership changes in different layers. Move this to an earlier layer if possible._

Found a string matching the pattern `[^ch(own|mod)\b]`

Priority: **Medium**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#run">RUN</a></kbd>

