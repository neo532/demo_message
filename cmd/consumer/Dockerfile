FROM hub.com/standard/golang:1.16-debian


ADD ./bin/consumer /home/www/demo_message/bin/consumer
ADD ./configs /home/www/demo_message/configs
ADD ./cmd/consumer/supervisord.conf /home/www/demo_message/supervisord.conf

WORKDIR /home/www/demo_message

CMD ["supervisord","-c","/home/www/demo_message/supervisord.conf"]
