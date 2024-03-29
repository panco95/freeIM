## FreeIM-Server

A free open source im server used by golang.

## Features

- [x] Free open source
- [x] Relation chain
- [x] Friend social
- [x] Group socail
- [x] Neighborhood social (contains chatGroup)
- [x] Distributed connection
- [x] Websocket network protocol
- [x] Json/Protobuf data protocol
- [x] Multiple message types
- [x] High performance
- [ ] Large number of users chat group
- [ ] Real time audio and video (RTC)
- [ ] Multiple platform client
- [ ] Manager web system

## Doc

[API doc](docs/api.md)

## Component

- [x] Mysql
- [x] Redis
- [x] Etcd
- [ ] Clickhouse

## Thirdparty

- [x] Qiniu-cloud
- [x] Smsbao
- [x] Smtp(email)

## How to run

cp development.yml.example development.yml

go run main.go server

## How to build

make build

## Docker

cp development.yml.example development.yml

docker-compose up