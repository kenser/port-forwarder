package adminuserdao

type AdminUser struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	LastLogin string `json:"last_login"`
}
