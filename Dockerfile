FROM alpine:latest

RUN echo http://repository.fit.cvut.cz/mirrors/alpine/v3.8/main > /etc/apk/repositories; \
    echo http://repository.fit.cvut.cz/mirrors/alpine/v3.8/community >> /etc/apk/repositories

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# install git - apt-get replace with apk
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh

RUN mkdir -p /etc/godp/frontend/dist
RUN mkdir -p /etc/godp/docs
WORKDIR /etc/godp

# 添加go二进制文件
ADD godp /etc/godp

# 添加frontend文件dist目录
COPY ./frontend/dist /etc/godp/frontend/dist
# 添加帮助文档目录
COPY ./docs /etc/godp/docs

RUN chmod 655 /etc/godp/godp

ENTRYPOINT ["/etc/godp/godp"]
EXPOSE 8081
