package model

type UserLogins struct {
	UserID         int      `gorm:"column:UserID; primary_key;"`
	LoginID    int      `gorm:"column:LoginID; primary_key; "`
	
}

func (UserLogins) TableName() string { return "User_Logins" }