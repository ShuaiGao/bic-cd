package manager

import (
	"bic-cd/internal/manager/nginx"
	"bic-cd/internal/manager/service"
	"bic-cd/internal/model"
	"bic-cd/internal/util"
	"bic-cd/pkg/db"
	"bic-cd/pkg/gen/api"
	"bic-cd/pkg/log"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"path"
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
			Version:     v.Version,
		})
	}
	return out, api.ECSuccess
}

func (m Manager) PostServices(ctx *gin.Context, in *api.RequestPostService) (out *api.ResponsePostService, code api.ErrCode) {
	data := &model.Service{
		Domain:      in.Domain,
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

func (m Manager) DeleteServices(ctx *gin.Context, id uint) (out *api.CommonNil, code api.ErrCode) {
	var instanceCount int64
	if err := db.DB().Model(&model.ServiceInstance{}).Where("service_id = ?", id).Count(&instanceCount).Error; err != nil {
		return nil, api.ECDbFindError.Wrap(err)
	}
	if instanceCount > 0 {
		return nil, api.ECServiceHasInstance
	}
	if err := db.DB().Delete(&model.Service{}, id).Error; err != nil {
		return nil, api.ECDbDeleteError.Wrap(err)
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
	exePath := path.Join(data.WorkingDir, in.Version)
	if _, err := os.Stat(exePath); err != nil {
		code = api.ECExeFileError.Wrap(err)
		return
	}
	var minPort uint16
	if err := db.DB().Model(&model.ServiceInstance{}).
		Select("max(port) as port").
		Pluck("port", &minPort).Error; err != nil {
		code = api.ECDbFindError.Wrap(err)
		return
	}
	port, err := service.GetAvailablePort(minPort, 2048)
	if err != nil {
		code = api.ECServerError.Wrap(err)
		return
	}
	instance := model.ServiceInstance{ServiceID: data.ID, Service: data, Version: in.Version}
	instance.SetExecStart(in.Version, port)
	if err = service.BuildService(service.Config{Instance: instance}); err != nil {
		code = api.ECServerError
		return
	}
	if err = db.DB().Create(&instance).Error; err != nil {
		code = api.ECDbCreateError.Wrap(err)
		return
	}
	// TODO test status
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
	var data []*model.Service
	var total int64
	if err := db.DB().Model(&model.Service{}).
		Preload("Instances").
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
	versionMap := make(map[string]bool)
	for _, v := range data {
		ServiceItem := &api.ServiceItem{
			Id:          uint32(v.ID),
			Name:        v.Name,
			Description: v.Description,
			WorkingDir:  v.WorkingDir,
			User:        v.User,
			PortMin:     int32(v.PortMin),
			PortMax:     int32(v.PortMax),
			Config:      v.Config,
			Version:     v.Version,
		}
		var instances []*api.Instance
		for _, inst := range v.Instances {
			versionMap[inst.Version] = true
			instances = append(instances, &api.Instance{
				Id:           uint32(inst.ID),
				Port:         uint32(inst.Port),
				ExecStart:    inst.ExecStart,
				Version:      inst.Version,
				CreateAt:     inst.CreatedAt.Format(time.DateTime),
				InstanceName: v.GetInstanceName(inst),
			})
		}
		if entries, err := os.ReadDir(v.WorkingDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				if !util.IsValidTag(entry.Name()) {
					continue
				}
				if _, ok := versionMap[entry.Name()]; ok {
					continue
				}
				instances = append(instances, &api.Instance{
					Version: entry.Name(),
				})
			}
		}
		out.Items = append(out.Items, &api.ServiceInstanceItem{ServiceItem: ServiceItem, Instances: instances})
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
	if data.Service.Version == data.Version {
		code = api.ECInstanceRunning.Wrap(data.Version)
		return
	}
	if err := service.RemoveService(data); err != nil {
		code = api.ECServerError.Wrap(err)
		return
	}
	if err := db.DB().Delete(&data).Error; err != nil {
		code = api.ECDbDeleteError.Wrap(err)
		return
	}
	return
}

func (m Manager) PostServiceVersion(ctx *gin.Context, in *api.RequestPostServiceVersion, id uint) (out *api.CommonNil, code api.ErrCode) {
	code = api.ECSuccess
	var data model.ServiceInstance
	log.XInfo("nginx ", id)
	if err := db.DB().Preload("Service").First(&data, id).Error; err != nil {
		code = api.ECDbFindError.Wrap(err)
		return
	}
	log.XInfo("nginx 2222")
	if data.Service.Version == in.Version {
		code = api.ECSuccess
		return
	}
	log.XInfo("nginx 3333")
	config := nginx.NginxConfig{
		Domain: data.Service.Domain,
		Port:   data.Port,
	}
	log.XInfo("nginx 4444")
	if err := nginx.CreateNginxConfig(config); err != nil {
		code = api.ECNginxConfig.Wrap(err)
		return
	}
	log.XInfo("nginx 5555")
	if err := nginx.ExecuteNginxTest(); err != nil {
		code = api.ECNginxTest.Wrap(err)
		return
	}
	if err := db.DB().Model(&model.Service{}).Where("id = ?", data.ServiceID).Update("version", in.Version).Error; err != nil {
		code = api.ECDbUpdate.Wrap(err)
		return
	}
	log.XInfo("nginx 6666")
	if err := nginx.ExecuteNginxReload(); err != nil {
		code = api.ECNginxReload.Wrap(err)
		return
	}
	return
}
