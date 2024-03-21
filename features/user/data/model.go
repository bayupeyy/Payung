package data

type User struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	Email    string
	Password string `json:"password" form:"password" validate:"required"`
	Hp       string `gorm:"type:varchar(13);uniqueIndex;not null" json:"hp" form:"hp" validate:"required,max=13,min=10"`
}
