package entity

type User struct {
	Id       string `gorm:"primarykey" json:"id"`
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
}
