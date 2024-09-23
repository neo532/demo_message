FROM hub.com/standard/golang:1.21-debian


ADD ./bin/api /home/www/demo_message/bin/api
ADD ./configs /home/www/demo_message/configs
ADD ./cmd/api/supervisord.conf /home/www/demo_message/supervisord.conf

WORKDIR /home/www/demo_message

EXPOSE 8501
EXPOSE 9501

CMD ["bin/api","-conf","configs/config.yaml"]
