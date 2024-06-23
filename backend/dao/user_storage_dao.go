package dao

import (
	"pvftools/backend/common"
	"pvftools/backend/model"
)

func GetCustomAttrTmpls() ([]*model.CustomAttrTmpl, error) {
	list := make([]*model.CustomAttrTmpl, 0)
	if err := common.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func AddCustomAttrTmpl(tmpl *model.CustomAttrTmpl) error {
	return common.DB.Create(tmpl).Error
}

func UpdateCustomAttrTmpl(tmpl *model.CustomAttrTmpl) error {
	return common.DB.Model(tmpl).Updates(tmpl).Error
}

func DeleteCustomAttrTmpl(id int) error {
	return common.DB.Where("id = ?", id).Delete(new(model.CustomAttrTmpl)).Error
}
