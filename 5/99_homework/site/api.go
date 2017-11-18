package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

var (
	statuses = map[string]int{
		"user":      0,
		"moderator": 10,
		"admin":     20,
	}
	users = map[string]*User{
		"rvasily": &User{
			ID:       42,
			Login:    "rvasily",
			FullName: "Vasily Romanov",
			Status:   statusAdmin,
		},
	}
	nextID  uint64 = 43
	muUsers        = &sync.RWMutex{}
)

const (
	statusUser      = 0
	statusModerator = 10
	statusAdmin     = 20
)

type MyApi struct{}

type ProfileParams struct {
	Login string `apivalidator:"required"`
}

type CreateParams struct {
	Login  string `apivalidator:"required,min=10"`
	Name   string `apivalidator:"paramname=full_name"`
	Status string `apivalidator:"enum=user|moderator|admin,default=user"`
	Age    int    `apivalidator:"min=0,max=128"`
}

type User struct {
	ID       uint64 `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
	Status   int    `json:"status"`
}

type NewUser struct {
	ID uint64 `json:"id"`
}

type ApiError struct {
	HTTPStatus int
	Err        error
}

func (ae ApiError) Error() string {
	return ae.Err.Error()
}

// apigen:api {"url": "/user/profile", "auth": false}
func (srv *MyApi) Profile(ctx context.Context, in ProfileParams) (*User, error) {

	if in.Login == "bad_user" {
		return nil, fmt.Errorf("bad user")
	}

	muUsers.RLock()
	user, exist := users[in.Login]
	muUsers.RUnlock()
	if !exist {
		return nil, ApiError{http.StatusNotFound, fmt.Errorf("user not exist")}
	}

	return user, nil
}

// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
func (srv *MyApi) Create(ctx context.Context, in CreateParams) (*NewUser, error) {
	if in.Login == "bad_username" {
		return nil, fmt.Errorf("bad user")
	}

	muUsers.Lock()
	defer muUsers.Unlock()

	_, exist := users[in.Login]
	if exist {
		return nil, ApiError{http.StatusConflict, fmt.Errorf("user %s exist", in.Login)}
	}

	id := nextID
	nextID++
	users[in.Login] = &User{
		ID:       id,
		Login:    in.Login,
		FullName: in.Name,
		Status:   statuses[in.Status],
	}

	return &NewUser{id}, nil
}
