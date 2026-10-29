package middleware

type Middleware struct {
	allowedOrigins []string
}

func InitMiddleware(allowedOrigins []string) *Middleware {
	return &Middleware{
		allowedOrigins: allowedOrigins,
	}
}
