# 表示依赖 alpine 最新版
FROM alpine:latest

# 创建程序工作目录
WORKDIR /Neko/data

# # 挂载容器目录
# VOLUME ["/Neko/data"]

# 拷贝编译出来的可执行执行文件
COPY /home/runner/work/Paimon/Paimon /Neko/Paimon

# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# 设置编码
ENV LANG C.UTF-8

# 运行程序
ENTRYPOINT ["/Neko/Paimon"]