FROM alpine:3.7 
RUN apk --no-cache add ca-certificates shadow && \
    groupadd -r avdb && useradd --no-log-init -r -g avdb avdb
ADD release/avdb_linux_amd64 /bin/avdb
USER avdb
ENTRYPOINT ["/bin/avdb"]
