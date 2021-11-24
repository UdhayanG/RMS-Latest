package model

type UserEmails struct {
	UserID         int      `gorm:"column:UserID; primary_key;"`
	EmailID    int      `gorm:"column:EmailID; primary_key; "`
	
}

func (UserEmails) TableName() string { return "User_Emails" }