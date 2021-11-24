package model

type Addresses struct {
	AddressID  int      `gorm:"column:AddressID; primary_key; auto_increment"`
	Address    string   `gorm:"column:Address" json:"Address"`
	Town       string   `gorm:"column:Town" json:"Town"`
	State      string   `gorm:"column:State" json:"State"`
	CountryID  int      `gorm:"column:CountryID" json:"CountryID"`
	PostCode   int      `gorm:"column:PostCode" json:"PostCode"`
	
}

func (Addresses) TableName() string { return "Addresses" }