# 表示依赖 alpine 最新版
FROM alpine:latest

# 创建程序工作目录
WORKDIR /Neko

# 挂载容器目录
VOLUME ["/Neko/date"]

# 拷贝当前目录下 go_docker_demo1 可以执行文件
COPY /home/runner/work/Paimon/Paimon /Neko/Paimon

# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# 设置编码
ENV LANG C.UTF-8

# 运行程序
ENTRYPOINT ["/Neko/Paimon"]