// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.6.3
// - protoc             v4.23.4
// source: user.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationUserServiceGetUserInfo = "/user.service.v1.UserService/GetUserInfo"
const OperationUserServiceUpdateUserInfo = "/user.service.v1.UserService/UpdateUserInfo"
const OperationUserServiceUserLogin = "/user.service.v1.UserService/UserLogin"
const OperationUserServiceUserRegister = "/user.service.v1.UserService/UserRegister"

type UserServiceHTTPServer interface {
	GetUserInfo(context.Context, *UserInfoRequest) (*UserInfoReply, error)
	UpdateUserInfo(context.Context, *UpdateUserInfoRequest) (*UpdateUserInfoReply, error)
	UserLogin(context.Context, *UserLoginRequest) (*UserLoginReply, error)
	UserRegister(context.Context, *UserRegisterRequest) (*UserRegisterReply, error)
}

func RegisterUserServiceHTTPServer(s *http.Server, srv UserServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/douyin/user/register", _UserService_UserRegister0_HTTP_Handler(srv))
	r.POST("/douyin/user/login", _UserService_UserLogin0_HTTP_Handler(srv))
	r.GET("/douyin/user", _UserService_GetUserInfo0_HTTP_Handler(srv))
	r.POST("/douyin/user/update", _UserService_UpdateUserInfo0_HTTP_Handler(srv))
}

func _UserService_UserRegister0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UserRegisterRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUserRegister)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UserRegister(ctx, req.(*UserRegisterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UserRegisterReply)
		return ctx.Result(200, reply)
	}
}

func _UserService_UserLogin0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UserLoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUserLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UserLogin(ctx, req.(*UserLoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UserLoginReply)
		return ctx.Result(200, reply)
	}
}

func _UserService_GetUserInfo0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UserInfoRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceGetUserInfo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetUserInfo(ctx, req.(*UserInfoRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UserInfoReply)
		return ctx.Result(200, reply)
	}
}

func _UserService_UpdateUserInfo0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateUserInfoRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUpdateUserInfo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateUserInfo(ctx, req.(*UpdateUserInfoRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateUserInfoReply)
		return ctx.Result(200, reply)
	}
}

type UserServiceHTTPClient interface {
	GetUserInfo(ctx context.Context, req *UserInfoRequest, opts ...http.CallOption) (rsp *UserInfoReply, err error)
	UpdateUserInfo(ctx context.Context, req *UpdateUserInfoRequest, opts ...http.CallOption) (rsp *UpdateUserInfoReply, err error)
	UserLogin(ctx context.Context, req *UserLoginRequest, opts ...http.CallOption) (rsp *UserLoginReply, err error)
	UserRegister(ctx context.Context, req *UserRegisterRequest, opts ...http.CallOption) (rsp *UserRegisterReply, err error)
}

type UserServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewUserServiceHTTPClient(client *http.Client) UserServiceHTTPClient {
	return &UserServiceHTTPClientImpl{client}
}

func (c *UserServiceHTTPClientImpl) GetUserInfo(ctx context.Context, in *UserInfoRequest, opts ...http.CallOption) (*UserInfoReply, error) {
	var out UserInfoReply
	pattern := "/douyin/user"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceGetUserInfo))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserServiceHTTPClientImpl) UpdateUserInfo(ctx context.Context, in *UpdateUserInfoRequest, opts ...http.CallOption) (*UpdateUserInfoReply, error) {
	var out UpdateUserInfoReply
	pattern := "/douyin/user/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUpdateUserInfo))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserServiceHTTPClientImpl) UserLogin(ctx context.Context, in *UserLoginRequest, opts ...http.CallOption) (*UserLoginReply, error) {
	var out UserLoginReply
	pattern := "/douyin/user/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUserLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserServiceHTTPClientImpl) UserRegister(ctx context.Context, in *UserRegisterRequest, opts ...http.CallOption) (*UserRegisterReply, error) {
	var out UserRegisterReply
	pattern := "/douyin/user/register"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUserRegister))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
