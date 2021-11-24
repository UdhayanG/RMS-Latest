package model

type UserPhoneNumbers struct {
	UserID         int      `gorm:"column:UserID; primary_key"`
	PhoneNoID    	  int    `gorm:"column:PhoneNoID; primary_key"`
	
}

func (UserPhoneNumbers) TableName() string { return "User_PhoneNumbers" }