FROM ubuntu:16.04
RUN  mkdir -p /etc/ssl/certs && mkdir /waimai
COPY cacert.pem /etc/ssl/certs
COPY waimai /waimai/waimai
COPY look /waimai/look
RUN chmod +x /waimai/look
VOLUME /waimai/runtime
VOLUME /waimai/files
VOLUME /waimai/config
VOLUME /waimai/logs
EXPOSE 80
WORKDIR /waimai
CMD ./look ./waimai