FROM busybox
MAINTAINER Rogier Lommers <rogier@lommers.org>

ADD https://github.com/rogierlommers/greedy/releases/download/1.0/greedy-1.0-linux-amd64.tar.bz2 /

EXPOSE 8080
ENTRYPOINT ["/greedy"]
