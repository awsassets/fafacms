#!/bin/bash

mkdir /opt
chmod 777 /opt
mkdir /opt/mydocker
chmod 777 /opt/mydocker
mkdir -p /opt/mydocker/redis/data
mkdir -p /opt/mydocker/redis/conf
mkdir -p /opt/mydocker/mysql/data
mkdir -p /opt/mydocker/mysql/conf
mkdir -p /opt/mydocker/fafacms
mkdir -p /opt/mydocker/fafacms/storage
mkdir -p /opt/mydocker/fafacms/storage_x
mkdir -p /opt/mydocker/fafacms/log
cp my.cnf /opt/mydocker/mysql/conf/my.cnf
chmod 644 /opt/mydocker/mysql/conf/my.cnf
cp redis.conf /opt/mydocker/redis/conf/redis.conf
cp config.yaml /opt/mydocker/fafacms/config.yaml
docker-compose stop
docker-compose rm -f
docker-compose up -d