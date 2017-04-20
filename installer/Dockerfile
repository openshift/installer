FROM alpine:3.5
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
RUN apk --no-cache --update add ca-certificates
COPY bin/linux/installer /installer
EXPOSE 4444
ENV INSTALLER_OPEN_BROWSER=false
ENV INSTALLER_ADDRESS=0.0.0.0:4444
ENTRYPOINT ["/installer"]
