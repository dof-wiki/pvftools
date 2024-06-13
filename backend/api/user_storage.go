package api

import (
	"pvftools/backend/dao"
	"pvftools/backend/model"
)

func (a *App) GetCustomAttrTmpls() ([]*model.CustomAttrTmpl, error) {
	return dao.GetCustomAttrTmpls()
}

func (a *App) AddCustomAttrTmpl(tmpl *model.CustomAttrTmpl) error {
	return dao.AddCustomAttrTmpl(tmpl)
}

func (a *App) UpdateCustomAttrTmpl(tmpl *model.CustomAttrTmpl) error {
	return dao.UpdateCustomAttrTmpl(tmpl)
}

func (a *App) DeleteCustomAttrTmpl(id int) error {
	return dao.DeleteCustomAttrTmpl(id)
}
