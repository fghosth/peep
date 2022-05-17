#!/bin/bash
set -e


if [ "${1:0:1}" = '-' ]; then
    set -- app "$@" #如果第一个参数的第一个字符是【-】,在所有参数前添加segment 以空格分割
fi

if [ "$1" = 'app' ]; then
    mkdir -p /app/logs
    mkdir -p /app/profile
    touch /app/logs/app.log
    touch /app/logs/costEngin.log
	mkdir -p /app/logs/
	/app/server up --conf /app/config/config.yaml >> /app/logs/demo.log 2>&1 &
	sleep 1
	tail -qf /app/logs/*.log
fi