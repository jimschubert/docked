# Rules
*  [D2:single-cmd](#d2single-cmd)
*  [D3:avoid-copy-all](#d3avoid-copy-all)
*  [D5:no-debian-frontend](#d5no-debian-frontend)
*  [D5:secret-aws-access-key](#d5secret-aws-access-key)
*  [D5:secret-aws-secret-access-key](#d5secret-aws-secret-access-key)
*  [D6:questionable-expose](#d6questionable-expose)
*  [D7:tagged-latest](#d7tagged-latest)
*  [D7:tagged-latest-builder](#d7tagged-latest-builder)
*  [D9:oci-labels](#d9oci-labels)
*  [DA:maintainer-deprecated](#damaintainer-deprecated)
*  [DC:avoid-sudo](#dcavoid-sudo)
*  [DC:consider-multistage](#dcconsider-multistage)
*  [DC:curl-without-fail](#dccurl-without-fail)
*  [DC:gpg-without-batch](#dcgpg-without-batch)
*  [DC:layered-ownership-change](#dclayered-ownership-change)
*  [DC:minimize-layers](#dcminimize-layers)
*  [DC:sort-installer-args](#dcsort-installer-args)


## D2:single-cmd

> _Only a single CMD instruction is supported_

More than one CMD may indicate a programming error. Docker will run the last CMD instruction only, but this could be a security concern.

Priority: **Critical**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#cmd">CMD</a></kbd>

## D3:avoid-copy-all

> _Avoid copying entire source directory into image_

Explicitly copying sources helps avoid accidentally persisting secrets or other files that should not be shared.

Priority: **High**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#copy">COPY</a></kbd>

## D5:no-debian-frontend

> _Convert DEBIAN_FRONTEND to an ARG._

Avoid DEBIAN_FRONTEND, which affects derived images and docker run. Change this to an ARG.
This rule matches against the pattern `\bDEBIAN_FRONTEND\b`

Priority: **Critical**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#env">ENV</a></kbd>

## D5:secret-aws-access-key

> _Secrets shouldn&#39;t be hard-coded. You should remove and rotate any secrets._

This rule matches against the pattern `\bAK[A-Z0-9]{18}\b`

Priority: **Critical**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#env">ENV</a></kbd>

## D5:secret-aws-secret-access-key

> _Secrets shouldn&#39;t be hard-coded. You should remove and rotate any secrets._

This rule matches against the pattern `\b[A-Za-z0-9/+=]{40}\b`

Priority: **Critical**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#env">ENV</a></kbd>

## D6:questionable-expose

> _Avoid documenting EXPOSE with sensitive ports_

The EXPOSE command is metadata and does not actually open ports. Documenting the intention to expose sensitive ports poses a security concern.

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#expose">EXPOSE</a></kbd>

## D7:tagged-latest

> _Avoid using images tagged as Latest in production builds_

Docker best practices suggest avoiding `latest` images in production builds

Priority: **High**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#from">FROM</a></kbd>

## D7:tagged-latest-builder

> _Avoid using images tagged as Latest in builder stages_

Using `latest` images in builders is not recommended (builds are not repeatable).

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#from">FROM</a></kbd>

## D9:oci-labels

> _Consider using common annotations defined by Open Containers Initiative_

Open Containers Initiative defines a common set of annotations which expose as labels on containers

Priority: **Medium**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#label">LABEL</a></kbd>

## DA:maintainer-deprecated

> _MAINTAINER is deprecated_

MAINTAINER instruction is deprecated; Use LABEL instead, which can be queried via `docker inspect`.
This rule matches against the pattern `[[:graph:]]+`

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#maintainer">MAINTAINER</a></kbd>

## DC:avoid-sudo

> _Avoid running root elevation tasks like sudo/su_

Non-root users should avoid having sudo access in containers. Consider using gosu instead.

Priority: **Medium**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#run">RUN</a></kbd>

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

> _Change ownership in the same layer as file operation (RUN or COPY)_

In AUFS, ownership defined in an earlier layer can not be overridden by a broader mask in a later layer.
This rule matches against the pattern `[^ch(own|mod)\b]`

Priority: **Medium**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#run">RUN</a></kbd>

## DC:minimize-layers

> _Try to minimize the number of layers which increase image size_

RUN, ADD, and COPY create new layers which may increase the size of the final image. Consider condensing these to fewer than 7 combined layers or use multi-stage builds where possible.

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#run">RUN</a></kbd><kbd><a href="https://docs.docker.com/engine/reference/builder/#add">ADD</a></kbd><kbd><a href="https://docs.docker.com/engine/reference/builder/#copy">COPY</a></kbd>

## DC:sort-installer-args

> _Sort installed packages for package managers: apt-get, apk, npm, etc._

Sorting installed packages alphabetically prevents duplicates and simplifies maintainability.

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#run">RUN</a></kbd>

