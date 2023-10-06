package models


type User struct {
	Id        string `json:"id"`
	Username  string `json:"username,omitempty"`
	CreatedOn string `json:"created_on"`
}
