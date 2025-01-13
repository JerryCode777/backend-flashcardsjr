package models

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password,omitempty"` // omitirá la contraseña al marshalling
}
