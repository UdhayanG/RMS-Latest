
package model

type Countries struct {
	CountryID         int      `gorm:"column:CountryID; primary_key; auto_increment"`
	CountryName    	  string    `gorm:"column:CountryName" json:"CountryName"`
	PhonePrefix 	  int      `gorm:"column:PhonePrefix" json:"PhonePrefix"`
	
}

func (Countries) TableName() string { return "Countries" }