FROM ubuntu
MAINTAINER Rogier Lommers <rogier@lommers.org>
LABEL description="Greedy, a personal readinglist"

# install dependencies
RUN apt-get update  
RUN apt-get install -y ca-certificates

# add binary
COPY bin/greedy /greedy

# change to data dir and run bianry
WORKDIR "/greedy-data"
CMD ["/greedy"]
