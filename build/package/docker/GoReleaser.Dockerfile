FROM gcr.io/distroless/base-debian12

ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/lexware-office-mcp /usr/local/bin/lexware-office-mcp

ENTRYPOINT ["/usr/local/bin/lexware-office-mcp"]
