# Client Service

## Info

Client Service 用于处理用户的HTTP请求，并提供RESTFUL接口。用户请求到达Client Service中之后。
Client Service会通过Db Service以及Cache Service提供的RPC Client发送RPC请求，在Db Service
和 Cache Service中处理具体的业务逻辑。

## Seckill Method

在这一部分实现了四种解决超卖问题的方法：

- 数据库悲观锁：
    具体做法是，由于减库存和新增订单是一个原子操作。必须同时完成或者同时不完成。因此要开一个数据库事务
    首先做的是减库存：

    ```sql
    update product set stock = stock - num where id = pid and stock >= num
    ```

    最主要的是这个 stock >= num，由于mysql中update操作加了排他锁，是一个原子操作，只有当库存足够时这个操作才能成功，否则事务将会滚。

- 数据库乐观锁：

    具体做法是在Product表中加一个version版本号，开启事务后，首先select product表中的stock和版本号，然后进行更新时使用一下sql语句

    ```sql
    update product set stock = stock - num where id = pid and version = version
    ```

    只有当version与之前select的版本号一致时改库存才会成功。

    数据库乐观锁的缺点是在高并发场景下会有大量的请求失效，当高并发但是数量不足时会出现 少卖的情况。解决方案是在失败后增加重试。

- 服务器端加锁：
    通过全局mutex将订单事务加锁，使创建订单的事务串行执行

- 服务器使用channel：
    使用channel和服务器加锁的原理一样，都是使订单事务串行执行

上述两个服务器端方案的缺点是：当有多个服务器端进行负载均衡时无法解决超卖问题，仍然可以超卖。解决方案是使用分布式锁，比如zookeeper或者redis或者etcd的锁