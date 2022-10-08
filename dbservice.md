# DB Service

## Info
DB Service是对数据库中的User、Orser、Product表进行相应的增删改查。
它采用GORM来实现对数据库的操作，并向其他Service提供RPC服务。其他Service
通过RPC调用DB Service层提供的方法来实现对数据库表的操作

## Implement

### entity

首先新建一个entity包，在包中新建User、Product、Order三个结构体，结构体中
的字段与数据库中对应表的属性一致

⚠️注意：由于 order在mysql中是关键字，不能用做表名。因此订单表的表名设为orders
而用户表和产品表的表名都为单数，分别为user和product。

### store

在dbservice中新建一个store包，用于使用GORM实现对数据库的相关操作。

首先gorm需要连接MySql数据库，我们在mysql.go中新建一个datastore结构体
结构体中设置一个gorm.DB用于数据库的操作。

然后我们在init()方法中初始化一个datastore结构体实例，命名为Mysql，并将其的
db设置为我们要连接的数据库的gorm.DB()

(在init()方法中初始化一个全局的datastore使用的是设计模式中的单例模式，由于我们连接
数据库使用一个db即可，因此使用单例模式可以减少资源的消耗。这里具体细分的话是单例模式中
的饿汉模式，即在初始化时就初始化了这个datastore，而不是在使用时再初始化。)

我们将数据库连接的相关参数(如host,username,password,database)保存在cofigs
文件夹下的conf.ymal中,并使用Viper库读取相应的设置，并连接数据库。

由于对用户表、订单表以及产品表的操作并不相同，我们要分别为订单操作，用户操作，产品操作定义
相应的结构体并实现相应的操作。

除了定义结构体之外，我们为用户操作、订单操作和产品操作定义了相应的接口，并让相对应的结构体实现
对应的接口。

在我们通过datastore获得相应的操作时返回的是接口而不是结构体。这样的好处是解耦合。

如果我们直接返回结构体，我们在做单元测试的时候必须要连接数据库进行测试，但是如果我们在获得操作时
返回相应的接口，那我们在进行Mock测试时可以定义自己的Mock结构体实现相应的方法，这样就可以不用连接
数据库进行测试。

而且也方便对相应的操作进行切换，只要实现了接口中的方法即可。

我们通过在datastore结构体定义相应的New方法来获得相应的数据库操作，每一个
数据库操作中也包含一个gorm.DB,这个gorm.DB与datastore中的db相同，在New方法
中直接复制即可。

使用New方法获得相应的操作实例是设计模式中的工厂模式。采用工厂模式的好处一是将对象的
创建和使用进一步解耦合，二是减小代码量，便于维护。

### service

service层对store层的方法进行了一层封装。同样分别对user、product和order的service定义
接口，然后定义serv结构体实现相应的接口，在实现方法中调用相应的store的方法即可。

在service 中获得相应的Service采用的仍然是设计模式中的单例模式，
首先定义一个全局的service变量，然后在Get方法中使用锁或者sync.Once()确保
这个全局的service只被初始化一次，然后韩惠这个service的接口即可。

这里的单例模式是懒汉模式，与在init()方法中初始化相比，只有使用到service并且service还未
被初始化时才会对servi进行初始化，并且仅初始化一次。使用懒汉模式可以保证服务在启动的时候可以
拥有更快的速度，并且减少资源的浪费。

### rpc

db service 中的服务是被外部的服务通过rpc调用的，因此还需要实现rpc service 和rpc client
这里我们使用google的grpc

首先新建一个rpc文件夹，然后新建一个dbservice.proto文件，在里边儿定义相应的rpc服务和消息，
然后使用

```shell
 protoc -I. --go_out=plugins=grpc:. dbservice.proto
```

生成相应的rpc代码。

然后我们在dbservice文件夹下新建rpc包，在rpc包下实现rpcserv和rpccli

rpcserv是rpc服务器，我们要实现rpc定义中的相应的服务，在服务中调用service中相应的
方法，即可实现对数据库的操作。

同时我们还要实现一个StartRpcServServer()方法，启动一个rpcservServer
在方法中我们新建一个tcp连接，监听的端口号同样定义在配置文件conf.ymal中，并通过
Viper读取配置。同时我们还要在这个方法中注册相应的rpc 服务。

rpccli中我们要对客户端的相应方法进行封装。这里主要是实现类型的转换。
因为在entity中定义的user、order、以及product中结构体中属性的类型和在rpc通信中
相应的结构体的类型并不完全相同，比如时间类型等。
因此在Add方法中要将entity中的类型转化成rpc中的类型，在Get方法中要将rpc得到的reply中
的类型转化为相应的entity中的类型。这就是rpcclient中对客户端中方法进行封装的主要目的。

这样我们在其他服务器中，通过dbserivce 中提供的rpc的client拿到对应服务的client即可实现
远程的rpc调用。

