FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY osrs_ge_exporter /usr/local/bin/
ENTRYPOINT ["osrs_ge_exporter"]
