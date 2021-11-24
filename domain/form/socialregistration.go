package form

type SocialRegistration struct {
	Id        string `json:"id" form:"id" query:"id"`
	Name      string `json:"name" form:"name" query:"name"`
	Email     string `json:"email" form:"email" query:"email"`
	PhotoUrl  string `json:"photoUrl" form:"photoUrl" query:"photoUrl"`
	FirstName string `json:"firstName" form:"firstName" query:"firstName"`
	LastName  string `json:"lastName" form:"lastName" query:"lastName"`
	AuthToken string `json:"authToken" form:"authToken" query:"authToken"`
	IdToken   string `json:"idToken" form:"idToken" query:"idToken"`
	Response  struct {
		Token_type    string `json:"token_type" form:"token_type" query:"token_type"`
		Access_token  string `json:"access_token" form:"access_token" query:"access_token"`
		Scope         string `json:"scope" form:"scope" query:"scope"`
		Login_hint    string `json:"login_hint" form:"login_hint" query:"login_hint"`
		Expires_in    int    `json:"expires_in" form:"expires_in" query:"expires_in"`
		Id_token      string `json:"id_token" form:"id_token" query:"id_token"`
		Session_state struct {
			ExtraQueryParams struct {
				Authuser string `json:"authuser" form:"authuser" query:"authuser"`
			}
		}
		First_issued_at int    `json:"first_issued_at" form:"first_issued_at" query:"first_issued_at"`
		Expires_at      int    `json:"expires_at" form:"expires_at" query:"expires_at"`
		IdpId           string `json:"idpId" form:"idpId" query:"idpId"`
	}
	Provider  string `json:"provider" form:"provider" query:"provider"`
	Err       string `json:"Err" form:"Err" query:"Err"`
	LoginType string `json:"LoginType" form:"LoginType" query:"LoginType" value:"Social"`
}
