# 表示依赖 alpine 最新版
FROM alpine:latest

# 创建程序工作目录
WORKDIR /Neko

# # 挂载容器目录
VOLUME ["/Neko/data"]

# 拷贝编译出来的可执行执行文件
COPY Paimon ./

# 验证文件
RUN ls -R

# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# 设置编码
ENV LANG C.UTF-8

# 赋予权限
RUN chmod +x Paimon

# 安装运行环境
RUN apk --no-cache add libc6-compat libgcc libstdc++

# 运行程序
ENTRYPOINT ["./Paimon"]
