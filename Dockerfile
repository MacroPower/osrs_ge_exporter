FROM golang:1.19

COPY osrs_ge_exporter /usr/local/bin/

ENTRYPOINT ["osrs_ge_exporter"]
