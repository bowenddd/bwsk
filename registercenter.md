# Register Center

## Info

### What's Register Center?

注册中心可以理解为一个中间商，被调用的服务将自己的信息注册到注册中心，而服务的调用者去注册中心拿到被调用的服务的信息，调用该服务。

### Why Register Center?

在我们最初的实现中是没有注册中心的，比如db service这个服务，它通过RPC调用的方式向外提供服务，它本身是一个RPC Server，有自己的地址（端口号）；当需要调用db service提供的服务时，我们需要一个RPC Client，这个Client需要与RPC Server建立连接，因此需要知道这个 db service的RPC Service的地址。

在没有注册中心的情况下，我们只能通过硬编码的方法实现，比如在配置文件中写入RPC Server的地址。

当这个RPC Server的地址（端口）改变时我们也要对应改变配置文件中的地址。

而当有多个服务时，当服务改变后通过硬编码的方式会很难维护。

而引入注册中心之后，RPC Server向注册中心注册自己的信息，而调用者
到注册中心根据服务名拿到对应的RPC Client即可。

因为当RPC Server的地址变化时，注册中心中的信息会变化，而调用者是通过服务名去注册中心拿信息，服务名不变，所以调用者不用做任何修改。

注册中心的另一个用处是实现负载均衡。比如一个服务有多个实例提供服务，他们都注册到注册中心，而调用者到注册中心拿Client时可以根据流量等信息选择一个最优的服务实例调用服务。

## Design

在这里我们通过etcd实现注册中心，etcd是一个分布式高性能KV服务器，通过Raft协议实现分布式一致性。

我们的注册中心主要实现两大功能： 服务注册 和服务发现。

### Service Register

服务注册很简单，就是将服务的名称作为key，服务的地址作为value，写入到etcd中。

由于一个服务可能有多个实例，因此服务名采用scheme+service name + port的形式
e.g. 一个dbservice 的服务的地址是 www.dbservice.com:5678
则
key   : /bwsk/dbserbice/www.dbservice.com:5678
value : www.dbservice.com:5678

由于服务可能发生中断，比如宕机或网络故障，因此在写入时加入一个lease，即过期时间，然后让他自动续租。即当lease过期但是服务仍然正常运行时，延长这个lease。

这样当服务中断后，不会续租，因此一段时间后这个服务的地址会自动删除

### Service Discovery

服务发现是获得一个服务的所有实例的信息。
比如dbservice有三个实例，他们在注册中心注册的信息为：

key   : /bwsk/dbserbice/www.dbservice1.com:5678
value : www.dbservice3.com:5678

key   : /bwsk/dbserbice/www.dbservice2.com:5678
value : www.dbservice3.com:5678

key   : /bwsk/dbserbice/www.dbservice3.com:5678
value : www.dbservice3.com:5678

我们可以根据key的前缀/bwsk/dbserbice/，拿到这个服务所有实例的地址。

然后我们可以采用负载均衡的算法选择一个服务的实例连接。
常用的负载均衡算法有
Round Robin（轮询）
Random      (随机)
或根据流量等信息选择

⚠️注意：由于service在注册时key可能过期，因此每间隔一段时间就要重新从etcd中获取service的信息，保证service的信息不过时。

