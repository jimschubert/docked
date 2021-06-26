# Rules
*  [D5:secret-aws-access-key](#d5:secret-aws-access-key)
*  [D5:secret-aws-secret-access-key](#d5:secret-aws-secret-access-key)
*  [D7:tagged-latest](#d7:tagged-latest)
*  [D7:tagged-latest-builder](#d7:tagged-latest-builder)
*  [DA:maintainer-deprecated](#da:maintainer-deprecated)
*  [DC:consider-multistage](#dc:consider-multistage)


## D5:secret-aws-access-key

> _Secrets should not be stored directly in the Dockerfile_

Found a string matching the pattern \bAK[A-Z0-9]{18}\b

Priority: **High**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#env">ENV</a></kbd>

## D5:secret-aws-secret-access-key

> _Secrets should not be stored directly in the Dockerfile_

Found a string matching the pattern \b[A-Za-z0-9/&#43;=]{40}\b

Priority: **High**  
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

Found a string matching the pattern [[:graph:]]&#43;

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#maintainer">MAINTAINER</a></kbd>

## DC:consider-multistage

> _Consider using multi-stage builds for complex operations like building code._

A multi-stage build can reduce the final image size by building necessary components or downloading large archives in a separate build context. This can help keep your final image lean.

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#run">RUN</a></kbd><kbd><a href="https://docs.docker.com/engine/reference/builder/#from">FROM</a></kbd>

