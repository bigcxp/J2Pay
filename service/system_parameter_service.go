package service

import (
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/pkg/casbin"
)

//系统参数详情
func GetDetail() (response.Parameter,error) {
	p := model.Parameter{}
	detail,err := p.GetDetail()
	return detail,err
}

//更新系统参数
func UpdateParameter(edit request.ParameterEdit)  error{
	defer casbin.ClearEnforcer()
	p := model.Parameter{}
	p.ID = uint(edit.ID)
	return p.UpdateParameter(edit)
}

//更新gasPrice
func UpdateGasPrice(edit request.ParameterEdit)  error{
	defer casbin.ClearEnforcer()
	p := model.Parameter{}
	p.ID = uint(edit.ID)
	return p.UpdateGasPrice(edit)
}
