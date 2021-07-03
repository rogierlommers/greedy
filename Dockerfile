FROM alpine
LABEL description="Greedy, a personal readinglist"

# add binary
COPY bin/greedy /greedy

RUN apk update

# needed to do GETs through https
RUN apk --no-cache add tzdata zip ca-certificates && update-ca-certificates

# change to data dir and run bianry
WORKDIR "/greedy-data"
CMD ["/greedy"]
