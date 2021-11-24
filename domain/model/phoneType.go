package model

type PhoneNoTypes struct {
	PhoneNoTypeID         int      `gorm:"column:PhoneNoTypeID; primary_key; auto_increment"`
	PhoneNoTypeDesc    	  string    `gorm:"column:PhoneNoTypeDesc" json:"PhoneNoTypeDesc"`
	
}

func (PhoneNoTypes) TableName() string { return "PhoneNoTypes" }