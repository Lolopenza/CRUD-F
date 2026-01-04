package models

import "time"

type User struct {
	User_ID    int       `json:"usr_id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Created_At time.Time `json:"created_at,omitempty"`
	Updated_At time.Time `json:"updated_at,omitempty"`
}
