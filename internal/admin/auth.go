package admin

import (
	"bic-cd/internal/model"
	"bic-cd/pkg/db"
	"bic-cd/pkg/gen/api"
	"bic-cd/pkg/jwt"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type Auth struct {
}

func checkAuth(username, password string) (*model.User, error) {
	user := &model.User{}
	if err := db.DB().Where("username = ?", username).First(user).Error; err != nil {
		return nil, fmt.Errorf("no this user(%s)", username)
	}
	if user.Forbidden {
		return nil, errors.New("用户已被禁用")
	}
	return user, nil
}

func (a Auth) PostAuth(ctx *gin.Context, in *api.RequestAuth) (out *api.ResponseAuth, code api.ErrCode) {
	code = api.ECSuccess
	user, err := checkAuth(in.Username, in.Password)
	if err != nil {
		code = api.ECAuthError
		return
	}
	token, err := jwt.GenerateToken(user.ID, user.Username, time.Hour*12)
	if err != nil {
		code = api.ECAuthError
		return
	}
	out = &api.ResponseAuth{Token: token, Username: user.Username}
	return
}

func (a Auth) PostAdmin(ctx *gin.Context, in *api.RequestPostAdmin) (out *api.ResponsePostAdmin, code api.ErrCode) {
	code = api.ECSuccess
	var count int64
	if err := db.DB().Model(&model.User{}).Count(&count).Error; err != nil {
		return nil, api.ECAuthError.Wrap(err)
	}
	u := model.User{Username: in.Username}
	u.SetPassword(in.Password)
	if err := db.DB().Create(&u); err != nil {
		return nil, api.ECDbCreateError.Wrap(err)
	}
	return
}
