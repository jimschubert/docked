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
* LABEL: recommended open container labels
* LABEL: correct formatting for container labels
* ENV: recommend single-env formatting
* ENV: avoid mixing `key value` and `key=value` format
* COPY: recommend using `--chown`
* RUN: yum-clean or remove package list
* RUN: apt-clean or remove package list
* RUN: apk clean or remove package list
* RUN: yum-no-upgrades
* RUN: apt-no-upgrades
* RUN: apk-no-upgrades
* ~EXPOSE: valid port ranges~
* ~EXPOSE: avoid ssh et al. (low, since [EXPOSE is informational](https://docs.docker.com/engine/reference/builder/#expose))~
* ADD: warn on external files
* ADD: prefer copy for no tgz
* ADD: error for absolute paths
* ADD: Avoid fetching over HTTP(S)
* USER: require non-root user for "official" images (Docker official and Google Distro-less)
* USER: bind to username rather than UID (See [this](https://devopsbootcamp.org/dockerfile-security-best-practices/#1-2-don-t-bind-to-a-specific-uid))
* CMD/ENTRYPOINT scripts should be owned by root
* RUN: (need to research how to implement something like shellcheck)
