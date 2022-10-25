package httpq

import (
	"github.com/gin-gonic/gin"
)

func Handle(httpMethod, relativePath string, handler HandlerFunc) {
	instance.Handle(httpMethod, relativePath, handler)
}

func ClientHandle(httpMethod, relativePath string, handler HandlerFunc) {
	instance.ClientHandle(httpMethod, relativePath, handler)
}

func ServerHandle(httpMethod, relativePath string, handler HandlerFunc) {
	instance.ServerHandle(httpMethod, relativePath, handler)
}

func AdminHandle(httpMethod, relativePath string, handler HandlerFunc) {
	instance.AdminHandle(httpMethod, relativePath, handler)
}

func CallbackHandle(httpMethod, relativePath string, handler HandlerFunc) {
	instance.CallbackHandle(httpMethod, relativePath, handler)
}

func SetClientAuth(handler gin.HandlerFunc) {
	instance.clientAuthHandle = handler
}

func SetServerAuth(handler gin.HandlerFunc) {
	instance.serverAuthHandle = handler
}

func SetAdminAuth(handler gin.HandlerFunc) {
	instance.adminAuthHandle = handler
}
