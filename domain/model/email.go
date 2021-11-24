package model

type Emails struct {
	EmailID         int      `gorm:"column:EmailID; primary_key; auto_increment"`
	EmailAddress    string    `gorm:"column:EmailAddress" json:"EmailAddress"`
	
}

func (Emails) TableName() string { return "Emails" }