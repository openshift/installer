FROM alpine:3.6
VOLUME /tests
ENTRYPOINT ["/bin/bash", "-lc", "bundle", "exec"]
RUN apk add --no-cache bash curl openssl readline zlib && \
    adduser -s /bin/bash -D rspec && \
    echo "export PATH=\$HOME/.rbenv/bin:\$PATH ; eval \"\$(rbenv init -)\"" >> /home/rspec/.bashrc && \
    echo ". ~/.bashrc" >> /home/rspec/.bash_profile
ADD rbenv /home/rspec/.rbenv
RUN chown -R rspec:rspec /home/rspec/.rbenv
USER rspec
SHELL ["/bin/bash", "-c"]
WORKDIR /tests/rspec
# RUN "source ~/.bashrc && bundle install"
