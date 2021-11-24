
package model

type Logins struct {
	LoginID             int      `gorm:"column:LoginID; primary_key; auto_increment"`
	UserName     	    string   `gorm:"column:UserName" json:"UserName"`
	LoginTypeID    	    int   	 `gorm:"column:LoginTypeID" json:"LoginTypeID"`
	UserNameVerified    int    	 `gorm:"column:UserNameVerified" json:"UserNameVerified"`
	LoginPassword 		string   `gorm:"column:LoginPassword" json:"LoginPassword"`
	LoginPasswordSalt 	string   `gorm:"column:LoginPasswordSalt" json:"LoginPasswordSalt"`
	
}

func (Logins) TableName() string { return "Logins" }