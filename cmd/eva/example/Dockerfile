
FROM alpine:3.10

RUN apk --update add tzdata \
	&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
	&& echo "Asia/Shanghai" > /etc/timezone \
	&& apk del tzdata

RUN mkdir -p /app
# 以上部分可以弄成一个通用基础镜像

WORKDIR /app

##暴露端口
EXPOSE 8080

COPY conf/config.*.json /app/conf/
COPY example /app/example

#最终运行docker的命令
ENTRYPOINT  ["./example"]
