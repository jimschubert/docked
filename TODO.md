# TODOs

## Features

* ~User-facing configuration~
* ~docked initialization (similar to logrus NewLog)~ 
* ~HTML Reporting~
* ~JSON Reporting (junit style?)~
* Concurrent evaluation of rules
* ~Testing~

## Commands

* ~RUN: avoid running su/sudo~
* ~COPY: avoid copying entire context (`.`)~. See [this](https://devopsbootcamp.org/dockerfile-security-best-practices/#3-3-build-context-and-dockerignore).
* ~LABEL: recommended open container labels~
* LABEL: correct formatting for container labels. See [this](https://docs.docker.com/config/labels-custom-metadata/)
  * ~`com.docker.*`, `io.docker.*`, and `org.dockerproject.*` namespaces are reserved by Docker for internal use~
  * ~Label keys should begin and end with a lower-case letter and should only contain lower-case alphanumeric characters, the period character (.), and the hyphen character (-). Consecutive periods or hyphens are not allowed.~
  * The period character (.) separates namespace “fields”. Label keys without namespaces are reserved for CLI use, allowing users of the CLI to interactively label Docker objects using shorter typing-friendly strings.
* ENV: recommend single-env formatting
* ENV: avoid mixing `key value` and `key=value` format
* RUN: unsetting environment variable set by ENV. See [this](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#env) 
* RUN: include `--no-log-init` to useradd. See [this](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#user)
* COPY: recommend using `--chown`
* RUN: yum-clean or remove package list
* RUN: apt-clean or remove package list
* RUN: apk clean or remove package list
* ~RUN: yum-no-upgrades, apt-no-upgrades, apk-no-upgrades~ this advice was removed in [docker docs](https://github.com/docker/docker.github.io/pull/12571) and [owasp](https://github.com/OWASP/CheatSheetSeries/pull/614) in March 2021.
* ~EXPOSE: valid port ranges~
* ~EXPOSE: avoid ssh et al. (low, since [EXPOSE is informational](https://docs.docker.com/engine/reference/builder/#expose))~
* ~ADD: warn on external files~
* ~ADD: prefer copy for no tgz~
* ~ADD: error for absolute paths~
* ~ADD: Avoid fetching over HTTP(S), at least in final build context; consider using multi-stage build.~
* USER: require non-root user for "official" images (Docker official and Google Distro-less)
* USER: bind to username rather than UID (See [this](https://devopsbootcamp.org/dockerfile-security-best-practices/#1-2-don-t-bind-to-a-specific-uid))
* CMD/ENTRYPOINT scripts should be owned by root
* RUN: (need to research how to implement something like shellcheck)
