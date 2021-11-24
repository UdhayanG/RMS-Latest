package model

type LoginTypes struct {
	LoginTypeID         int      `gorm:"column:LoginTypeID; primary_key; auto_increment"`
	LoginTypeDesc    	  string    `gorm:"column:LoginTypeDesc" json:"LoginTypeDesc"`
	
}

func (LoginTypes) TableName() string { return "LoginTypes" }