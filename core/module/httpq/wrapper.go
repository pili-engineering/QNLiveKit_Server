package httpq

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
