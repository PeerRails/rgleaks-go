FROM ubuntu:14.04

ENV HOME /home/deploy

RUN useradd deploy
RUN echo "deploy ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

COPY bin/rgleaks-64 /home/deploy/rgleaks-64
COPY .env /home/deploy/.env
RUN chmod +x /home/deploy/rgleaks-64
RUN chown -R deploy:deploy $HOME
USER deploy
ENV http_proxy 'http:/172.17.42.1:8118/'
ENV https_proxy 'http:/172.17.42.1:8118/'
ENV HTTP_PROXY 'http:/172.17.42.1:8118/'
ENV HTTPS_PROXY 'http:/172.17.42.1:8118/'

CMD ["./home/deploy/rgleaks-64"]
