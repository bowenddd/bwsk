package consts

const (
	DBPESSIMISTICLOCK = "DBPESSIMISTICLOCK"
	DBOPTIMISTICLOCK  = "DBOPTIMISTICLOCK"
	SERVICELOCK       = "SERVICELOCK"
	SERVICECHANNEL    = "SERVICECHANNEL"
)

var MethodSet map[string]struct{}

func init() {
	MethodSet = make(map[string]struct{})
	MethodSet[DBPESSIMISTICLOCK] = struct{}{}
	MethodSet[DBOPTIMISTICLOCK] = struct{}{}
	MethodSet[SERVICELOCK] = struct{}{}
	MethodSet[SERVICECHANNEL] = struct{}{}
}
