package manager

import (
	"bic-cd/internal/manager/service"
	"bic-cd/internal/model"
	"bic-cd/pkg/db"
	"bic-cd/pkg/gen/api"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
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

func (m Manager) PostServiceDeploy(ctx *gin.Context, in *api.RequestPostServiceDeploy, id uint) (out *api.ResponsePostServiceDeploy, code api.ErrCode) {
	var count int64
	if err := db.DB().Model(&model.ServiceInstance{}).
		Where("service_id = ?", id).
		Where("version = ?", in.Version).Count(&count).Error; err != nil {
		code = api.ECDbFindError.Wrap(err)
		return
	}
	if count > 0 {
		code = api.ECRepeatedVersion
		return
	}
	var data model.Service
	if err := db.DB().First(&data, id).Error; err != nil {
		code = api.ECDbFindError.Wrap(err)
		return
	}
	port, err := service.GetAvailablePort(1024, 2048)
	if err != nil {
		code = api.ECServerError.Wrap(err)
		return
	}
	instance := model.ServiceInstance{ServiceID: data.ID, Service: data, Version: in.Version}
	instance.SetExecStart(port)
	if err = service.BuildService(service.Config{Instance: instance}); err != nil {
		code = api.ECServerError
		return
	}
	if err = db.DB().Create(&instance).Error; err != nil {
		code = api.ECDbCreateError.Wrap(err)
		return
	}
	out = &api.ResponsePostServiceDeploy{
		Id: uint32(id),
	}
	return out, api.ECSuccess
}

func (m Manager) PostServiceStart(ctx *gin.Context, id uint) (out *api.ResponsePostServiceDeploy, code api.ErrCode) {
	code = api.ECSuccess
	var data model.ServiceInstance
	if err := db.DB().Preload("Service").First(&data, id).Error; err != nil {
		code = api.ECDbFindError.Wrap(err)
		return
	}
	err := service.StartService(data)
	if err != nil {
		code = api.ECServerError
		return
	}
	out = &api.ResponsePostServiceDeploy{
		Id: uint32(id),
	}
	return
}

func (m Manager) PostServiceStop(ctx *gin.Context, id uint) (out *api.ResponsePostServiceDeploy, code api.ErrCode) {
	code = api.ECSuccess
	var data model.ServiceInstance
	if err := db.DB().Preload("Service").First(&data, id).Error; err != nil {
		code = api.ECDbFindError.Wrap(err)
		return
	}
	err := service.StopService(data)
	if err != nil {
		code = api.ECServerError
		return
	}
	out = &api.ResponsePostServiceDeploy{
		Id: uint32(id),
	}
	return
}

func (m Manager) PostServiceStatus(ctx *gin.Context, id uint) (out *api.ResponsePostServiceStatus, code api.ErrCode) {
	code = api.ECSuccess
	var data model.ServiceInstance
	if err := db.DB().Model(&data).Preload("Service").First(&data, id).Error; err != nil {
		fmt.Println(err)
		code = api.ECDbFindError.Wrap(err)
		return
	}
	fmt.Println("service status 111")
	stdout, err := service.StatusService(data)
	fmt.Println("service status 222 ", stdout)
	fmt.Println("service status 333 ", stdout)
	if err != nil {
		code = api.ECServerError.Wrap(err)
		return
	}
	out = &api.ResponsePostServiceStatus{
		Id:     uint32(id),
		Stdout: stdout,
	}
	return
}

func (m Manager) GetServiceInstances(ctx *gin.Context, in *api.RequestGetServiceInstance) (out *api.ResponseGetServiceInstance, code api.ErrCode) {
	code = api.ECSuccess
	var data []*model.ServiceInstance
	var total int64
	if err := db.DB().Model(&model.ServiceInstance{}).Debug().
		Preload("Service").
		Count(&total).
		Offset(int((in.Page - 1) * in.PageSize)).
		Limit(int(in.PageSize)).
		Find(&data).Error; err != nil {
		code = api.ECDbFindError.Wrap(err)
		return
	}
	out = &api.ResponseGetServiceInstance{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint32(total),
	}
	for _, v := range data {
		out.Items = append(out.Items, &api.ServiceInstanceItem{
			ServiceItem: &api.ServiceItem{
				Id:          uint32(v.ID),
				Name:        v.Service.Name,
				Description: v.Service.Description,
				WorkingDir:  v.Service.WorkingDir,
				User:        v.Service.User,
				PortMin:     int32(v.Service.PortMin),
				PortMax:     int32(v.Service.PortMax),
				Config:      v.Service.Config,
			},
			Id:           uint32(v.ID),
			Port:         uint32(v.Port),
			ExecStart:    v.ExecStart,
			Version:      v.Version,
			CreateAt:     v.CreatedAt.Format(time.DateTime),
			InstanceName: v.GetService(),
		})
	}
	return
}

func (m Manager) DeleteServiceInstance(ctx *gin.Context, id uint) (out *api.CommonNil, code api.ErrCode) {
	code = api.ECSuccess
	var data model.ServiceInstance
	if err := db.DB().Preload("Service").First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = api.ECSuccess
			return
		}
		code = api.ECDbFindError.Wrap(err)
		return
	}
	if err := service.RemoveService(data); err != nil {
		code = api.ECServerError
		return
	}
	if err := db.DB().Delete(&data).Error; err != nil {
		code = api.ECDbDeleteError.Wrap(err)
		return
	}
	return
}
