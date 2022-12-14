package entity

type Order struct {
	ID        uint    `json:"id"`
	UserId    uint    `json:"user_id"`
	ProductId uint    `json:"product_id"`
	Num       int     `json:"num"`
	Price     float64 `json:"price"`
	Created   string  `json:"created"`
}

func (Order) TableName() string {
	return "orders"
}
