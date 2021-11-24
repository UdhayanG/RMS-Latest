package model

type UserAddresses struct {
	UserID   int      `gorm:"column:UserID" primary_key; "`
	AddressID  int      `gorm:"column:AddressID; primary_key;"`
	
}

func (UserAddresses) TableName() string { return "User_Addresses" }