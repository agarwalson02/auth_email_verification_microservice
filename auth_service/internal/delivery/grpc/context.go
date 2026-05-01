package grpc

type contextKey string

const (
	SessionKey contextKey = "sessionID"
	TokenKey   contextKey = "token"
)