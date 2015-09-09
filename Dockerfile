FROM ubuntu:14.04

ENV HOME /home/deploy

RUN useradd deploy
RUN adduser deploy www-data
RUN echo "deploy ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

COPY bin/rgleaks-64 /home/deploy/rgleaks-64
COPY .env /home/deploy/.env
RUN chmod +x /home/deploy/rgleaks-64
RUN chown -R deploy:deploy $HOME
ENV http_proxy "http://tor:8118/"
ENV HTTP_PROXY "http://tor:8118/"
ENV https_proxy "http://tor:8118/"
ENV HTTPS_PROXY "http://tor:8118/"
USER deploy
WORKDIR /home/deploy/
CMD ["./rgleaks-64"]
