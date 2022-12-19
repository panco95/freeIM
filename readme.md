## FreeIM-Server

A free open source im server used by golang.

## Features

- [x] Free open source
- [x] Username/email/mobile social
- [x] Relation chain
- [x] Websocket protocol
- [ ] Large number of users chat group
- [x] Distributed connection
- [ ] Real time audio and video (RTC)
- [x] Multiple message types
- [ ] Multiple platform client

## Doc

[API doc](docs/api.md)

## Environment

- mysql
- redis
- etcd
- thirdparty: smtp, qiniu, smsbao

## How to run

cp development.yml.example development.yml

go run main.go server

## How to build

make build

## Docker

cp development.yml.example development.yml

docker-compose up