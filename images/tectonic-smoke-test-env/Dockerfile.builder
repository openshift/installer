FROM alpine:edge
RUN apk update && apk upgrade && apk add alpine-sdk linux-headers bash curl openssl-dev readline-dev zlib-dev && \
    adduser -s /bin/bash -D rspec 
VOLUME /home/rspec/.rbenv
WORKDIR /home/rspec
USER rspec
RUN echo "export PATH=\$HOME/.rbenv/bin:\$PATH ; eval \"\$(rbenv init -)\"" >> .bashrc
ENTRYPOINT ["/bin/bash", "-c"]
