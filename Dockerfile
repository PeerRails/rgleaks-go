FROM ubuntu:14.04

ENV HOME /home/deploy

RUN useradd deploy
RUN adduser deploy www-data
RUN echo "deploy ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

COPY bin/rgleaks-64 /home/deploy/rgleaks-64
COPY .env /home/deploy/.env
RUN chmod +x /home/deploy/rgleaks-64
RUN chown -R deploy:deploy $HOME
RUN echo 'Acquire::http::Proxy "http://172.17.42.1:8118/";' > /etc/apt/apt.conf.d/proxy
USER deploy
RUN cd /home/deploy/
CMD ["./rgleaks-64"]
