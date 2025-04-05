FROM registry.access.redhat.com/ubi9/ubi-minimal

COPY bin/app1 .

WORKDIR /templates

COPY templates .

WORKDIR /

CMD ["/app1"]