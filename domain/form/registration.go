package form

type Registration struct {
	Command   string `json:"Command" form:"Command" query:"Command"`
	AddressID int    `json:"AddressID,omitempty" form:"AddressID" query:"AddressID"`
	Address   string `json:"Address" form:"Address" query:"Address"`
	Town      string `json:"Town" form:"Town" query:"Town"`
	State     string `json:"State" form:"State" query:"State"`
	CountryID int    `json:"CountryID" form:"CountryID" query:"CountryID"`
	PostCode  int    `json:"PostCode" form:"PostCode" query:"PostCode"`

	CountryName string `json:"CountryName" form:"CountryName" query:"CountryName"`
	PhonePrefix int    `json:"PhonePrefix" form:"PhonePrefix" query:"PhonePrefix"`

	EmailID      int    `json:"EmailID" form:"EmailID" query:"EmailID"`
	EmailAddress string `json:"EmailAddress" form:"EmailAddress" query:"EmailAddress"`

	LoginID          int    `json:"LoginID" form:"LoginID" query:"LoginID"`
	UserName         string `json:"UserName" form:"UserName" query:"UserName"`
	LoginTypeID      int    `json:"LoginTypeID" form:"LoginTypeID" query:"LoginTypeID"`
	UserNameVerified int    `json:"UserNameVerified" form:"UserNameVerified" query:"UserNameVerified"`
	LoginPassword    string `json:"LoginPassword" form:"LoginPassword" query:"LoginPassword"`
	Password         string `json:"Password" form:"Password" query:"Password" validate:"required"`
	ConfirmPassword  string `json:"ConfirmPassword" form:"ConfirmPassword" query:"ConfirmPassword" validate:"required"`

	LoginTypeDesc string `json:"LoginTypeDesc" form:"LoginTypeDesc" query:"LoginTypeDesc"`

	PhoneNoID int `json:"PhoneNoID" form:"PhoneNoID" query:"PhoneNoID"`
	//PhoneNumber       string `json:"PhoneNumber" form:"PhoneNumber" query:"PhoneNumber" validate:"required"`
	PhoneNumber       string `json:"PhoneNumber" form:"PhoneNumber" query:"PhoneNumber"`
	NumberinInterForm int    `json:"NumberinInterForm" form:"NumberinInterForm" query:"NumberinInterForm"`
	PhoneNoTypeID     int    `json:"PhoneNoTypeID" form:"PhoneNoTypeID" query:"PhoneNoTypeID"`

	PhoneNoTypeDesc string `json:"PhoneNoTypeDesc" form:"PhoneNoTypeDesc" query:"PhoneNoTypeDesc"`

	UserID           int    `json:"UserID" form:"UserID" query:"UserID"`
	FirstName        string `json:"FirstName" form:"FirstName" query:"FirstName" validate:"required"`
	MiddleName       string `json:"MiddleName" form:"MiddleName" query:"MiddleName"`
	LastName         string `json:"LastName" form:"LastName" query:"LastName"`
	DefaultAddressID int    `json:"DefaultAddressID" form:"DefaultAddressID" query:"DefaultAddressID"`
	DefaultPhoneID   int    `json:"DefaultPhoneID" form:"DefaultPhoneID" query:"DefaultPhoneID"`
	RequestKey       string `json:"RequestKey" form:"RequestKey" query:"RequestKey"`
	ErrStatus        string `json:"ErrStatus" form:"ErrStatus" query:"ErrStatus"`
	RequestID        string `json:"RequestID" form:"RequestID" query:"RequestID"`
	OTP              string `json:"OTP" form:"OTP" query:"OTP"`
}
