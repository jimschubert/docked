FROM ubuntu:xenial
RUN gpg --batch --verbose --keyserver hkp://keyserver.ubuntu.com:80 --keyserver-options timeout=5 --recv-keys ABAF11C65A2970B130ABE3C479BE3E4300411886

