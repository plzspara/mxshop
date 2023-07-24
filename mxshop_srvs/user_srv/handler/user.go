package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	. "mxshop_srvs/user_srv/global"
	. "mxshop_srvs/user_srv/model"
	. "mxshop_srvs/user_srv/proto"
	"strings"
	"time"
)

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type UserServers struct {
}

func ModelToResponse(user User) UserInfoResponse {
	userInfoRsp := UserInfoResponse{
		Id:       user.Id,
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

func (u *UserServers) GetUserList(_ context.Context, req *PageInfo) (*UserListResponse, error) {
	var users []User
	result := DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &UserListResponse{}
	rsp.Total = int32(result.RowsAffected)
	DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	for i := range users {
		userInfoRsp := ModelToResponse(users[i])
		rsp.UserInfo = append(rsp.UserInfo, &userInfoRsp)
	}
	return rsp, nil
}
func (u *UserServers) GetUserByMobile(_ context.Context, req *MobileRequest) (*UserInfoResponse, error) {
	var userInfo User
	result := DB.Where(&User{Mobile: req.Mobile}).First(&userInfo)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	response := ModelToResponse(userInfo)
	return &response, nil
}
func (u *UserServers) GetUserById(_ context.Context, req *UserIdRequest) (*UserInfoResponse, error) {
	var userInfo User
	result := DB.First(&userInfo, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	response := ModelToResponse(userInfo)
	return &response, nil
}
func (u *UserServers) CreateUser(_ context.Context, req *CreateUserInfo) (*UserInfoResponse, error) {
	var userInfo User
	result := DB.Where(&User{Mobile: req.Mobile}).First(&userInfo)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已经存在")
	}
	userInfo.Nickname = req.Nickname
	userInfo.Mobile = req.Mobile
	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	salt, encodedPwd := password.Encode(req.Password, options)
	userInfo.Password = fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
	result = DB.Create(&userInfo)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	response := ModelToResponse(userInfo)
	return &response, nil
}
func (u *UserServers) UpdateUser(_ context.Context, req *UpdateUserInfo) (*empty.Empty, error) {
	var userInfo User
	result := DB.First(&userInfo, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	birthday := time.Unix(int64(req.Birthday), 0)
	userInfo.Birthday = &birthday
	userInfo.Nickname = req.Nickname
	userInfo.Gender = req.Gender
	result = DB.Save(userInfo)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}
func (u *UserServers) CheckPassword(_ context.Context, req *CheckPasswordInfo) (*CheckResponse, error) {
	strs := strings.Split(req.EncryptedPassword, "$")
	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	result := password.Verify(req.Password, strs[1], strs[2], options)
	return &CheckResponse{Success: result}, nil
}
