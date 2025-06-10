package admin

import (
	"bic-cd/internal/model"
	"bic-cd/pkg/db"
	"bic-cd/pkg/gen/api"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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
	token := ""
	out = &api.ResponseAuth{Token: token, Username: user.Username}
	return
}
