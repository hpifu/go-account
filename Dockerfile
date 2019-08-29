FROM centos:centos7
COPY docker/account /var/docker/account
RUN mkdir -p /var/docker/account/log
EXPOSE 6060
WORKDIR /var/docker/account
CMD [ "bin/account", "-c", "configs/account.json" ]
