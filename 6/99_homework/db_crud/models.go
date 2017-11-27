package main

type Item struct {
	ID          *int    `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Updated     *string `json:"updated"`
}

type User struct {
	ID       *int    `json:"id"`
	Login    *string `json:"login"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Info     *string `json:"info"`
	Updated  *string `json:"updated"`
}
