package model

type PhoneNumbers struct {
	PhoneNoID         int      `gorm:"column:PhoneNoID; primary_key; auto_increment"`
	CountryID     	  int    `gorm:"column:CountryID" json:"CountryID"`
	PhoneNumber    	  string    `gorm:"column:PhoneNumber" json:"PhoneNumber"`
	NumberinInterForm int      `gorm:"column:NumberinInterForm" json:"NumberinInterForm"`
	PhoneNoTypeID 	  int      `gorm:"column:PhoneNoTypeID" json:"PhoneNoTypeID"`
	
}

func (PhoneNumbers) TableName() string { return "PhoneNumbers" }