FROM ubuntu
MAINTAINER Rogier Lommers <rogier@lommers.org>
LABEL description="Greedy, a personal readinglist"
COPY bin/greedy /greedy
WORKDIR "/greedy-data"
CMD ["/greedy"]
