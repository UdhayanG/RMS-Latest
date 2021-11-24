package form

type Login struct {
	UserName string `json:"UserName" form:"UserName" query:"UserName"`
	Password string `json:"Password" form:"Password" query:"Password"`
}
