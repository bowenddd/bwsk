# Authenication

## Info

认证与授权由两部分组成，认证部分负责检查请求是否是由合法的用户发起，而授权部分则是判断用户是否拥有发起某项请求的权利。

## Authenication Design

认证部分使用JWT(Json Web Token)实现，认证部分的逻辑是：
用户首先登陆，比如向/user/login发送一个POST请求，请求中携带了用户名和密码。
服务器拿到用户名和密码之后与数据库中保存的密码进行对比，如果一致，就为用户生成一个token发送给用户。
之后用户每一次请求都在header的authenication字段中填入该token。
服务器在执行某个route的handler之前，首先经过一个Auth中间件(拦截器)，这个Auth中间价从header中拿到token，判断这个token是否合法，如果不合法就不会执行相应的请求，并返回错误。只有合法再继续执行

### Problem

使用JWT存在的问题是当一个用户注销后，使用token仍然可以访问请求，因为token未过期。
解决方案是在Redis中维护一个有效的token，过期后删除这个token。

## Authorize Design

授权部分使用RBAC模型
数据库新增角色表、权限表、角色权限表、用户角色表

```sql
CREATE TABLE `role`
(
    `id`         INT             PRIMARY KEY AUTO_INCREMENT,
    `name`       VARCHAR(20)     NOT NULL UNIQUE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `perm`
(
    `id`         INT             PRIMARY KEY AUTO_INCREMENT,
    `path`       VARCHAR(20)     NOT NULL UNIQUE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `user_role`
(
    `id`         INT             PRIMARY KEY AUTO_INCREMENT,
    `user_id`    INT             NOT NULL,
    `role_id`    INT             NOT NULL
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `role_perm`
(
    `id`         INT             PRIMARY KEY AUTO_INCREMENT,
    `role_id`    INT             NOT NULL,
    `perm_id`    INT             NOT NULL 
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

```

在RBAC模型中，一个用户可以拥有多个角色，一个角色可以拥有多个权限。
在BWSK中，权限就是访问某个URL PATH的权利。

在鉴权的过程中，我们通过数据库查询到某个用户拥有的全部权限， 然后将其存放到Redis缓存中，之后每一次鉴权，只要去Redis中查询即可。如果Redis中查询不到就去Mysql中查询，然后将结果写入Redis中。

拿到这个用户拥有的全部权限之后再判断用户是否可以访问这个URL即可。

鉴权部分也是在Auth中间件中实现的，实现位置是认证之后。


