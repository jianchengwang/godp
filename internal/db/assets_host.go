package db

import (
	"godp/internal/api/types"
	"godp/internal/global"
)

type AssetsHost struct {
	BaseModel
	UserId      uint                  `json:"userId"`
	Name        string                `json:"name"`
	IP          string                `json:"ip"`
	IntranetIP  string                `json:"intranetIp"`
	Port        uint16                `json:"port"`
	User        string                `json:"user"`
	Password    string                `json:"password"`
	WorkDir     string                `json:"workDir"`
	Remark      string                `json:"remark"`
	CloudVendor types.CloudVendorType `json:"cloudVendor"`
	Verified    bool                  `json:"verified"`
}

type AssetsHostDbStruct struct{}

var AssetsHostDb = new(AssetsHostDbStruct)

func (AssetsHostDb *AssetsHostDbStruct) Page(param types.AssetsHostQueryParam) (error, interface{}, int64) {
	limit := param.Limit
	offset := param.Limit * (param.CurPage - 1)
	db := global.DB.Model(&AssetsHost{})

	var objList []AssetsHost
	var total int64

	if param.Q != "" {
		db.Or(
			db.Where("name LIKE ?", "%"+param.Q+"%"),
			db.Where("ip LIKE ?", "%"+param.Q+"%"),
		)
	}
	if param.Verified {
		db.Where("verified = ?", param.Verified)
	}

	err := db.Count(&total).Error
	if err != nil {
		return err, objList, total
	} else {
		db = db.Limit(limit).Offset(offset).Order("id desc")
		err = db.Find(&objList).Error
	}
	return err, objList, total
}

func (AssetsHostDb *AssetsHostDbStruct) GetByIP(ip string) (error, *AssetsHost) {
	var getObj = &AssetsHost{}
	err := global.DB.Where("ip = ?", ip).First(getObj).Error
	return err, getObj
}

func (AssetsHostDb *AssetsHostDbStruct) GetById(id uint) (error, *AssetsHost) {
	var getObj = &AssetsHost{}
	err := global.DB.First(getObj, id).Error
	return err, getObj
}

func (AssetsHostDb *AssetsHostDbStruct) SaveOrUpdate(updateParam AssetsHost) error {
	updateObj := &AssetsHost{}
	if updateParam.ID == 0 {
		err := global.DB.Model(&updateObj).Create(
			&updateParam,
		).Error
		return err
	} else {
		updateObj.ID = updateParam.ID
		err := global.DB.Model(&updateObj).Updates(
			&AssetsHost{
				Name:        updateParam.Name,
				IP:          updateParam.IP,
				Port:        updateParam.Port,
				User:        updateParam.User,
				Password:    updateParam.Password,
				WorkDir:     updateParam.WorkDir,
				Remark:      updateParam.Remark,
				CloudVendor: updateParam.CloudVendor,
			},
		).Error
		return err
	}
}
