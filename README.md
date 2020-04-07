# filestore
分布式文件存储客户端，支持文件分块上传，断点续传(TODO)等功能
### Getting Started
进入db.conf配置postgres与redis
启动redis server
```
redis-server
```
创建数据库（保证postgreSQL service已经启动）
```
CREATE DATABASE filestore;
```
建表（参考db/potgres/tables.txt）
```
CREATE TABLE tbl_file (
    filehash varchar(100),
    filename varchar(20),
    filesize integer,
    location varchar(30),
    uploadtime timestamp
);
```
运行项目
```
go run main.go
```
