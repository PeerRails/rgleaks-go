FROM ubuntu:14.04

ENV HOME /home/deploy

RUN useradd deply && echo 'deploy:docker' | chpasswd
RUN echo "deploy ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

COPY ./bin/rgleaks-64 /home/deploy
COPY .env /home/deploy
ENV http_proxy 'http:/172.17.42.1:8118/'
ENV https_proxy 'http:/172.17.42.1:8118/'
ENV HTTP_PROXY 'http:/172.17.42.1:8118/'
ENV HTTPS_PROXY 'http:/172.17.42.1:8118/'
CMD ["/home/deploy/rgleaks-64"]
