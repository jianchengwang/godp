#!/bin/bash
# Get real path
BASEDIR=$(cd `dirname $0` && pwd)
cd ${BASEDIR}

# Log Location on Server.
LOG_LOCATION=${BASEDIR}
exec > >(tee -i $LOG_LOCATION/build.`date +%Y%m%d%H%M%S`.log)
exec 2>&1

cd ../
# 拉取代码
git pull origin main

# 打包frontend
cd ./frontend
date "+%Y-%m-%d %H:%M:%S"
docker run --net=host --rm -it -e TZ=Asia/Shanghai -v npm-repo:/root/.npm -v yarn-cache:/usr/local/share/.cache/yarn -v "$(pwd)":/src -w /src node:16-alpine sh -c '
    echo -e "https://mirrors.aliyun.com/alpine/v3.11/main/\nhttps://mirrors.aliyun.com/alpine/v3.11/community/" > /etc/apk/repositories
    apk add curl
    curl -f https://get.pnpm.io/v6.16.js | node - add --global pnpm
    rm -rf node_modules dist
    pnpm install
    npm run build:stage
    git rev-parse HEAD > dist/git-commit-id.txt
    echo "frontend.git.build.time="$(date +"%Y-%m-%dT%H\:%M\:%S%z") > dist/git-commit.txt
    echo "frontend.git.commit.id.abbrev="$(git rev-parse --short HEAD) >> dist/git-commit.txt
    echo "frontend.git.commit.id.full="$(git rev-parse HEAD) >> dist/git-commit.txt '\
    || exit 0
date "+%Y-%m-%d %H:%M:%S"

# 编译go
cd ../
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/
go env -w GO111MODULE=on
go env -w CGO_ENABLED=0
go env -w GOARCH=amd64
go env -w GOOS=linux
go mod vendor
go build -tags netgo -o godp main.go

# 打包docker镜像
docker build -t godp:v0.0.1 .
# docker tag godp:v0.0.1 jianchengwang/godp
# docker login
# docker push jianchengwang/godp

# 部署
cd /root/godp/docker-compose-godp
kill -9 `netstat -nlp | grep :8081 | awk '{print $7}' | awk -F"/" '{ print $1 }'`
docker pull jianchengwang/godp:latest
docker-compose up -d
docker rmi $(docker images | grep "none" | awk '{print $3}')
