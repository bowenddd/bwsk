package store

import (
	"fmt"
	"seckill/seetings"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

var Mysql *datastore

func init() {
	Mysql = &datastore{}
	setting, err := seetings.GetSetting()
	mysqlSetting := setting.Mysql
	if err != nil {
		fmt.Println("init mysql datastore error, caused by read mysql setting!")
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%ds",
		mysqlSetting.Username,
		mysqlSetting.Password,
		mysqlSetting.Host,
		mysqlSetting.Port,
		mysqlSetting.Dbname,
		mysqlSetting.Timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("connect to database error!")
		return
	}
	Mysql.db = db
}

func (ds *datastore) NewUserStore() UserStore {
	return &UserOp{db: ds.db}
}

func (ds *datastore) NewOrderStore() OrderStore {
	orderOp := &OrderOp{
		db: ds.db,
		mu: &sync.Mutex{},
		ch: make(chan orderChan, 100),
	}
	go orderOp.HandleOrderChan()
	return orderOp
}

func (ds *datastore) NewProductStore() ProductStore {
	return &ProductOp{db: ds.db}
}

func(ds *datastore) NewPermStore() PermStore {
	return &PermOp{db: ds.db}
}

func (ds *datastore) DB() *gorm.DB {
	return ds.db
}
