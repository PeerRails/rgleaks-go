FROM ubuntu:14.04

COPY bin/rgleaks-64 /bin/
COPY .env /bin/
ENV http_proxy 'http:/172.17.42.1:8118/'
ENV https_proxy 'http:/172.17.42.1:8118/'
ENV HTTP_PROXY 'http:/172.17.42.1:8118/'
ENV HTTPS_PROXY 'http:/172.17.42.1:8118/'
CMD ["/bin/rgleaks-64"]
