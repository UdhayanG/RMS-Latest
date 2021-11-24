package model

type User struct {
	UserID           int      `gorm:"column:UserID; primary_key; auto_increment"`
	FirstName     	 string    `gorm:"column:FirstName" json:"FirstName"`
	MiddleName    	 string    `gorm:"column:MiddleName" json:"MiddleName"`
	LastName   		 string    `gorm:"column:LastName" json:"LastName"`
	DefaultAddressID int      `gorm:"column:DefaultAddressID" json:"DefaultAddressID"`
	DefaultPhoneID 	 int      `gorm:"column:DefaultPhoneID" json:"DefaultPhoneID"`
	
}

func (User) TableName() string { return "Users" }
