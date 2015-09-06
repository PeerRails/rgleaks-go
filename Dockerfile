FROM ubuntu:14.04
MAINTAINER me@gmail.com
COPY bin/rgleaks-64 /bin/
COPY .env /bin/
CMD ["/bin/rgleaks-64"]
