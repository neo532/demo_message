FROM hub.com/standard/golang:1.16-debian

ADD ./configs /home/www/demo_message/configs
ADD ./bin/script /home/www/demo_message/bin/script
ADD ./cmd/script/supervisord.conf /home/www/demo_message/supervisord.conf

WORKDIR /home/www/demo_message

CMD ["supervisord","-c","/home/www/demo_message/supervisord.conf"]
