ARG ARCH="amd64"
ARG OS="linux"

FROM golang:1.13.4-buster AS build
WORKDIR /go/src/github.com/vglafirov/iexcloud_exporter/
COPY . /go/src/github.com/vglafirov/iexcloud_exporter/
RUN make

FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="Vladimir Glafirov <vglafirov@gmail.com>"

COPY --from=build /go/src/github.com/vglafirov/iexcloud_exporter/iexcloud_exporter   /bin/iexcloud_exporter
USER nobody
EXPOSE 9107
ENTRYPOINT [ "/bin/iexcloud_exporter" ]