package entity

type User struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Sex     int    `json:"sex"`
	Phone   string `json:"phone"`
	Created string `json:"created"`
}

func (User) TableName() string {
	return "user"
}
