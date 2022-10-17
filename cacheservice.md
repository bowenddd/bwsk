# Cache Service

## Info

Cache Service是Redis缓存层的服务，它向client service提供rpc服务供client service调用。

## Why Redis Cache

Q：为什么要使用Redis缓存层？
A：在高并发场景下，如果在秒杀时，client service创建订单直接调用db service层的服务，会给数据库带来巨大的压力。因为数据库中的数据是存放在磁盘中的，查询和写入操作可能
存在磁盘和内存之间数据交换，因此所需要的时间比较长

而Redis数据库是一个内存数据库，数据都存放在内存中，与Mysql相比查询速度更快。

使用缓存层，我们可以在秒杀之前，将商品的库存写入redis中，在实际秒杀过程中。先去Reids中查询是否有足够的库存，如果有的话，缓存层库存减少，再去调用db service，在mysql中创建订单，否则直接返回，不再对数据库进行操作。

## Problem

但是这引入了一个新的问题，就是缓存层和数据库层数据一致性问题，一次缓存层的超卖问题。

因为加入缓存层后创建订单的逻辑变成了：
1:查询缓存中的库存-->2:修改缓存中的库存-->3:修改mysql中的库存并创建订单。

第1步和第2步可能导致redis中的数据不一致性：
比如A和B两个请求，顺序为A查询库存B查询库存，A更新，B更新。就有可能导致应该减2但是结果减1，从而使得在mysql中库存变负。虽然redis是单线程，但依然会超卖。

还有一个问题就是Redis层和Mysql层数据一致性问题，即Mysql中对库存修改后如何保证Redis中的数据和Mysql中的数据一致呢？

可以使用「延时双删」方法，即在Mysql修改Stock之前，先删除Redis中的数据，然后修改Mysql中的stock，修改成功后，间隔500ms再删除一次。

而在Redis层中，如果查不到key，则从Mysql中取数据，当数据放到Redis中。

## Solution

解决方案是对缓存层加锁，加锁的方案同样有乐观锁和悲观锁两种

- 悲观锁：

    悲观锁使用的是Redis中自带的分布式锁。
    使用 SETNX设置key时，如果key已经存在，则返回false，否则设置成功。
    因此我们可以使用这个命令，设置一个key让他当锁，如果返回true，即设置
    key成功，则加锁成功，可以执行业务逻辑，否则返回false，表明锁已经在使用。
    用锁将redis中的查询和修改操作锁起来，就可以解决一致性问题。

- 乐观锁：

    乐观锁使用的是Redis中的MULTI和WATCH机制实现的。MULTI是Redis中的事务，
    但是它不保证事务的原子性。即事务中的操作可以有的成功有的失败；而WATCH表示
    监听一个KEY，只有当KEY没有发生改变时，执行的操作才可以成功。将这两个放到一起，即可以实现Redis层的乐观锁。
   
