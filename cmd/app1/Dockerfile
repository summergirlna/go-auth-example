FROM registry.access.redhat.com/ubi9/ubi-minimal

COPY bin/app1 .
COPY cmd/app1/wait-for-it.sh .
RUN chmod +x /wait-for-it.sh

WORKDIR /templates

COPY templates .

WORKDIR /

CMD ["/app1"]