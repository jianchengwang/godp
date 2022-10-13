package db

import (
	"godp/internal/api/types"
	"godp/internal/global"
)

type BatchScript struct {
	BaseModel
	UserId  uint   `json:"userId"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Remark  string `json:"remark"`
}

type BatchScriptDbStruct struct{}

var BatchScriptDb = new(BatchScriptDbStruct)

func (BatchScriptDb *BatchScriptDbStruct) Page(param types.BatchScriptsQueryParam) (error, interface{}, int64) {
	limit := param.Limit
	offset := param.Limit * (param.CurPage - 1)
	db := global.DB.Model(&BatchScript{})

	var objList []BatchScript
	var total int64

	if param.Q != "" {
		db.Or(
			db.Where("name LIKE ?", "%"+param.Q+"%"),
		)
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

func (BatchScriptDb *BatchScriptDbStruct) GetById(id uint) (error, *BatchScript) {
	var getObj = &BatchScript{}
	err := global.DB.First(getObj, id).Error
	return err, getObj
}

func (BatchScriptDb *BatchScriptDbStruct) SaveOrUpdate(updateParam BatchScript) error {
	updateObj := &BatchScript{}
	if updateParam.ID == 0 {
		err := global.DB.Model(&updateObj).Create(
			&updateParam,
		).Error
		return err
	} else {
		updateObj.ID = updateParam.ID
		err := global.DB.Model(&updateObj).Updates(
			&BatchScript{
				Name:    updateParam.Name,
				Content: updateParam.Content,
				Remark:  updateParam.Remark,
			},
		).Error
		return err
	}
}

func (BatchScriptDb *BatchScriptDbStruct) DeleteById(id uint) error {
	return global.DB.Delete(&BatchScript{}, id).Error
}
