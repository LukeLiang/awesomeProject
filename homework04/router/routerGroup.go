package router

import "awesomeProject/homework04/api"

type RoutersGroup struct {
	UserRouter
	PostRouter
}

var (
	userApi = api.UserApi{}
	postApi = api.PostApi{}
)
