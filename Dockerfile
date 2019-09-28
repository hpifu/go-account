FROM centos:centos7
COPY docker/go-account /var/docker/go-account
RUN mkdir -p /var/docker/go-account/log
EXPOSE 6060
WORKDIR /var/docker/go-account
CMD [ "bin/account", "-c", "configs/account.json" ]
