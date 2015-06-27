FROM alpine:3.1

# install curl in the image so it is possible to get the runtime
# profiling information without any additional package installation.
RUN apk add --update-cache curl bash && rm -rf /var/cache/apk/*

ADD bin/publisher-gc /usr/local/bin/publisher-gc
ENTRYPOINT ["/usr/local/bin/publisher-gc"]

ENV DEIS_RELEASE 1.8.0-dev
