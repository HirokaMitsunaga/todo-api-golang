package model

type Todo struct {
	ID       string `gorm:"type:char(26);primaryKey"`
	Title    string `gorm:"type:varchar(100);not null"`
	Status   string `gorm:"type:varchar(20);not null"`
	UserID   string `gorm:"type:char(26);not null;index"`
	Priority uint16 `gorm:"type:smallint;not null"`

	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Todo) TableName() string {
	return "todos"
}
