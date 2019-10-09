FROM centos:centos7

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >> /etc/timezone

COPY docker/go-account /var/docker/go-account
RUN mkdir -p /var/docker/go-account/log

EXPOSE 6060

WORKDIR /var/docker/go-account
CMD [ "bin/account", "-c", "configs/account.json" ]
