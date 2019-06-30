#!/bin/sh
# 没有使用最新的模块, 因为当前工程也在github,后面再想办法解决.
GOOS=linux GOARCH=amd64 go build -o eeblog ../main.go

cp -r ../web/ ./web/
cp -r ../templates/ ./templates/
cp /etc/localtime localtime