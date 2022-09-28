package service

var Instance Service

type Service interface {
	RegisterAuthMiddleware()
}
