package authentication

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// MyAuth 自定义 Auth 需要实现 credentials.PerRPCCredentials 接口
type MyAuth struct {
	Username string
	Password string
}

// GetRequestMetadata 定义授权信息的具体存放形式，最终会按这个格式存放到 metadata map 中。
func (a *MyAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"username": a.Username, "password": a.Password}, nil
}

// RequireTransportSecurity 是否需要基于 TLS 加密连接进行安全传输
func (a *MyAuth) RequireTransportSecurity() bool {
	return false
}

const (
	Admin = "admin"
	Root  = "root"
)

func NewMyAuth() *MyAuth {
	return &MyAuth{
		Username: "error",
		Password: Root,
	}
}

// IsValidAuth 具体的验证逻辑
func IsValidAuth(ctx context.Context) error {
	var (
		user     string
		password string
	)
	// 从 ctx 中获取 metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	// 从metadata中获取授权信息
	// 这里之所以通过md["username"]和md["password"] 可以取到对应的授权信息
	// 是因为我们自定义的 GetRequestMetadata 方法是按照这个格式返回的.
	if val, ok := md["username"]; ok {
		user = val[0]
	}
	if val, ok := md["password"]; ok {
		password = val[0]
	}
	// 简单校验一下 用户名密码是否正确.
	if user != Admin || password != Root {
		return status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	return nil
}
