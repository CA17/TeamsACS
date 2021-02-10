# GenieACS v1.2 Dockerfile #
############################

FROM ubuntu:18.04

USER root

RUN rm -f /etc/apt/sources.list
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb-src http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb-src http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb-src http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb-src http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb-src http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse" >> /etc/apt/sources.list

RUN apt-get update && apt-get install -y sudo supervisor git
RUN mkdir -p /var/log/supervisor

RUN apt-get install -y curl
RUN curl -sL https://deb.nodesource.com/setup_12.x | sudo -E bash -
RUN apt-get install -y gcc g++ make
RUN apt-get install -y nodejs

#sudo npm install -g --unsafe-perm genieacs@1.2.2
RUN npm install -g --unsafe-perm genieacs@1.2.3

RUN useradd --system --no-create-home --user-group genieacs
#RUN mkdir /opt/genieacs
RUN mkdir -p /opt/genieacs/ext
RUN chown genieacs:genieacs /opt/genieacs/ext

RUN mkdir -p /var/log/genieacs
RUN chown genieacs:genieacs /var/log/genieacs

ADD genieacs.logrotate /etc/logrotate.d/genieacs

WORKDIR /opt
RUN mkdir -p /opt/genieacs/dist/bin \
&& ln -s /usr/bin/genieacs-fs /opt/genieacs/dist/bin/genieacs-fs \
&& ln -s /usr/bin/genieacs-ui /opt/genieacs/dist/bin/genieacs-ui \
&& ln -s /usr/bin/genieacs-nbi /opt/genieacs/dist/bin/genieacs-nbi \
&& ln -s /usr/bin/genieacs-cwmp /opt/genieacs/dist/bin/genieacs-cwmp

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY run_with_env.sh /usr/bin/run_with_env.sh
RUN chmod +x /usr/bin/run_with_env.sh

COPY setup_ssl.sh /setup_ssl.sh
RUN chmod +x /setup_ssl.sh
RUN /setup_ssl.sh
RUN chown -R genieacs:genieacs /etc/ssl/genieacs

RUN rm -rf /var/lib/apt/lists/*
RUN apt-get purge -y && apt-get autoremove && apt-get clean

WORKDIR /var/log/genieacs

CMD ["/usr/bin/supervisord","-c","/etc/supervisor/conf.d/supervisord.conf"]
