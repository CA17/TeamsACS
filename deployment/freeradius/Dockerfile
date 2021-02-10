FROM freeradius/freeradius-server:latest

USER root
COPY freeradius/ /etc/raddb/
COPY lfreemate /usr/sbin/lfreemate
RUN chmod +x /usr/sbin/lfreemate
