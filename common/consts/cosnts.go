package consts

import "time"

const (
	DBPESSIMISTICLOCK = "DBPESSIMISTICLOCK"
	DBOPTIMISTICLOCK  = "DBOPTIMISTICLOCK"
	SERVICELOCK       = "SERVICELOCK"
	SERVICECHANNEL    = "SERVICECHANNEL"
	NOMEASURE         = "NOMEASURE"
	CACHEPESSIMISTICLOCK = "CACHEPESSIMISTICLOCK"
	CACHEOPTIMISTICLOCK = "CACHEOPTIMISTICLOCK"
)

const (
	REDIS_CREATE_ORDER_MAX_RETRY = 20
)

const (
	TokenExpireTime = 2*time.Hour
)

var MethodSet map[string]struct{}

func init() {
	MethodSet = make(map[string]struct{})
	MethodSet[DBPESSIMISTICLOCK] = struct{}{}
	MethodSet[DBOPTIMISTICLOCK] = struct{}{}
	MethodSet[SERVICELOCK] = struct{}{}
	MethodSet[SERVICECHANNEL] = struct{}{}
	MethodSet[NOMEASURE] = struct{}{}
	MethodSet[CACHEPESSIMISTICLOCK] = struct{}{}
	MethodSet[CACHEOPTIMISTICLOCK] = struct{}{}
}
