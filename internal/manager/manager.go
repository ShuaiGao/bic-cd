package manager

import (
	"bic-cd/internal/manager/service"
	"bic-cd/internal/model"
	"bic-cd/pkg/db"
	"bic-cd/pkg/gen/api"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Manager struct{}

func (m Manager) GetServices(ctx *gin.Context, in *api.RequestGetService) (out *api.ResponseGetService, code api.ErrCode) {
	var data []*model.Service
	var total int64
	if err := db.DB().Model(&model.Service{}).
		Count(&total).
		Offset(int((in.Page - 1) * in.PageSize)).
		Limit(int(in.PageSize)).
		Find(&data).Error; err != nil {
		code = api.ECDbFindError.Wrap(err)
		return
	}
	out = &api.ResponseGetService{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint32(total),
	}
	for _, v := range data {
		out.Items = append(out.Items, &api.ServiceItem{
			Id:          uint32(v.ID),
			Name:        v.Name,
			Description: v.Description,
			ExecStart:   v.ExecStart,
			WorkingDir:  v.WorkingDir,
			User:        v.User,
			PortMin:     int32(v.PortMin),
			PortMax:     int32(v.PortMax),
			Config:      v.Config,
		})
	}
	return out, api.ECSuccess
}

func (m Manager) PostServices(ctx *gin.Context, in *api.RequestPostService) (out *api.ResponsePostService, code api.ErrCode) {
	data := &model.Service{
		Name:        in.Name,
		Description: in.Description,
		ExecStart:   in.ExecStart,
		WorkingDir:  in.WorkingDir,
		User:        in.User,
		PortMin:     uint16(in.PortMin),
		PortMax:     uint16(in.PortMax),
		Config:      in.Config,
	}
	if err := db.DB().Create(data).Error; err != nil {
		code = api.ECDbCreateError.Wrap(err)
		return
	}
	return out, api.ECSuccess
}

func (m Manager) PostServiceDeploy(ctx *gin.Context, id uint) (out *api.ResponsePostServiceDeploy, code api.ErrCode) {
	var data model.Service
	if err := db.DB().First(&data, id).Error; err != nil {
		code = api.ECDbFindError.Wrap(err)
		return
	}
	f, err := service.CreateService(service.Config{Service: data})
	if err != nil {
		code = api.ECServerError
		return
	}
	out = &api.ResponsePostServiceDeploy{
		Id:          uint32(id),
		FileContent: f,
	}
	fmt.Println(f)
	return out, api.ECSuccess
}
