FROM scratch
COPY osrs_ge_exporter /usr/local/bin/
ENTRYPOINT ["osrs_ge_exporter"]
