package model

type User struct {
	ID       string `gorm:"type:char(26);primaryKey"`
	Name     string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password string `gorm:"type:text;not null"`
}

func (User) TableName() string {
	return "users"
}
