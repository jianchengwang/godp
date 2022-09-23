package db

import (
	"encoding/json"
	"fmt"
	"godp/internal/api/types"
	"godp/internal/global"
	"godp/internal/pojo"
	"godp/internal/process"
)

type ProjectInfo struct {
	BaseModel
	ProjectApp  string        `json:"projectApp"`
	ProjectName string        `json:"projectName"`
	ClientName  string        `json:"clientName"`
	OnlineDate  pojo.JsonDate `json:"onlineDate"gorm:"type:date"`
	CreateUser  int           `json:"createUser"`
	EditorUser  int           `json:"editorUser"`
	Internal    string        `json:"internal" gorm:"default:0"`
	Tags        string        `json:"tags"`
	Config      string        `json:"config"`

	ProjectConfig pojo.ProjectConfig `json:"projectConfig"gorm:"-"`
}

type ProjectInfoDbStruct struct{}

var ProjectInfoDb = new(ProjectInfoDbStruct)

func (ProjectInfoDb *ProjectInfoDbStruct) Page(param types.ProjectQueryParam) (error, interface{}, int64) {
	limit := param.Limit
	offset := param.Limit * (param.CurPage - 1)
	db := global.DB.Model(&ProjectInfo{})

	var objList []ProjectInfo
	var total int64

	if param.Q != "" {
		db.Or(
			db.Where("project_name LIKE ?", "%"+param.Q+"%"),
			db.Where("project_app LIKE ?", "%"+param.Q+"%"),
			db.Where("project_name LIKE ?", "%"+param.Q+"%"),
			db.Where("client_name LIKE ?", "%"+param.Q+"%"),
		)
	}
	if param.Internal != "" {
		db.Where("internal = ?", param.Internal)
	}

	err := db.Count(&total).Error
	if err != nil {
		return err, objList, total
	} else {
		db = db.Limit(limit).Offset(offset).Order("id desc")
		err = db.Find(&objList).Error
		for i := 0; i < len(objList); i++ {
			objList[i].ProjectConfig = buildProjectConfig(&objList[i])
		}
	}
	return err, objList, total
}

func (ProjectInfoDb *ProjectInfoDbStruct) GetById(id uint) (error, *ProjectInfo) {
	var getObj = &ProjectInfo{}
	err := global.DB.First(getObj, id).Error
	getObj.ProjectConfig = buildProjectConfig(getObj)
	return err, getObj
}

func (ProjectInfoDb *ProjectInfoDbStruct) GetByProjectApp(projectApp string, id uint) (error, ProjectInfo) {
	var projectInfo ProjectInfo
	if id > 0 {
		err := global.DB.Where("project_app = ? and id != ?", projectApp, id).First(&projectInfo).Error
		return err, projectInfo
	} else {
		err := global.DB.Where("project_app = ?", projectApp).First(&projectInfo).Error
		return err, projectInfo
	}
}

func (ProjectInfoDb *ProjectInfoDbStruct) CreateOrUpdate(updateParam ProjectInfo) error {
	updateObj := &ProjectInfo{}
	if updateParam.ID == 0 {
		err := global.DB.Model(&updateObj).Create(
			&updateParam,
		).Error
		return err
	} else {
		updateObj.ID = updateParam.ID
		err := global.DB.Model(&updateObj).Updates(
			&ProjectInfo{
				ProjectName: updateParam.ProjectName,
				ProjectApp:  updateParam.ProjectApp,
				ClientName:  updateParam.ClientName,
				OnlineDate:  updateParam.OnlineDate,
				EditorUser:  updateParam.EditorUser,
				Internal:    updateParam.Internal,
				Tags:        updateParam.Tags,
			},
		).Error
		return err
	}
}

func (ProjectInfoDb *ProjectInfoDbStruct) UpdateConfig(id uint, config pojo.ProjectConfig) error {
	updateObj := &ProjectInfo{}
	updateObj.ID = id
	configJson, _ := json.Marshal(config)
	err := global.DB.Model(&updateObj).Updates(
		ProjectInfo{
			Config: string(configJson),
		},
	).Error
	return err
}

func (ProjectInfoDb *ProjectInfoDbStruct) DeleteById(id uint) error {
	return global.DB.Delete(&ProjectInfo{}, id).Error
}

func buildProjectConfig(projectInfo *ProjectInfo) pojo.ProjectConfig {
	ci := pojo.IpAddressStruct{
		IP:         global.CI_IP,
		Port:       uint16(global.CI_PORT),
		IntranetIP: global.CI_IntranetIP,
		User:       global.CI_USER,
		Password:   global.CI_PASSWORD,
		DeployDir:  global.CI_DEPLOYDIR,
	}
	if projectInfo.Config == "" {
		projectTemplate := process.LoadProjectConfigJsonTpl()
		projectTemplate.ProjectName = projectInfo.ProjectName
		projectTemplate.ProjectApp = projectInfo.ProjectApp
		projectTemplate.CI = ci
		return projectTemplate
	} else {
		projectConfig := pojo.ProjectConfig{}
		err := json.Unmarshal([]byte(projectInfo.Config), &projectConfig)
		if err != nil {
			fmt.Println(err)
		}
		projectConfig.ProjectName = projectInfo.ProjectName
		projectConfig.ProjectApp = projectInfo.ProjectApp
		if projectConfig.CI.IP == "" {
			projectConfig.CI = ci
		}
		return projectConfig
	}
}
