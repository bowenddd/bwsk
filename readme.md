# BWSK
## Info
BWSK(bowen seckill)是一个秒杀系统，采用三层设计，分别为 
client service、cache service、 db service三层
- client service，用于处理用户请求，向用户提供RESTFUL API
- cache service，操作redis缓存，client service通过RPC调用
本层的服务，对缓存数据库进行写入和查询
- db service，用于对数据库进行增删改查操作，cache service中的
操作完成之后需要将信息持久化到数据库中。
## Preparation
- 服务器中安装Mysql5.7
- 服务器中安装Redis6.2.6
- 数据库增加相应的表
 
```sql
create database seckill;

use seckill;

drop table if exists orders;

CREATE TABLE `orders`
(
    `id` INT     PRIMARY KEY AUTO_INCREMENT,
    `user_id`    INT NOT NULL,
    `product_id` INT NOT NULL,
    `num`        INT NOT NULL,
    `price`      DECIMAL(11,2) NOT NULL,
    `created`    DATETIME
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `user`
(
    `id`         INT PRIMARY KEY AUTO_INCREMENT,
    `name`       VARCHAR(20) NOT NULL,
    `sex`        INT NOT NULL,
    `phone`      VARCHAR(20) NOT NULL,
    `created`    DATETIME
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `product`
(
    `id`          INT PRIMARY KEY AUTO_INCREMENT,
    `name`        VARCHAR(20) NOT NULL,
    `price`       DECIMAL(11,2) NOT NULL,
    `description` VARCHAR(20) NOT NULL,
    `stock`       INT NOT NULL,
    `created`     DATETIME
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

## DB Service

[DB Service说明](./dbservice.md)

[Client Service说明](./clientservice.md)

[Cache Service说明](./cacheservice.md)