FROM scratch
MAINTAINER Rogier Lommers <rogier@lommers.org>
LABEL description="Greedy, a personal readinglist"

# add binary
COPY bin/greedy /greedy

# change to data dir and run bianry
WORKDIR "/greedy-data"
CMD ["/greedy"]
