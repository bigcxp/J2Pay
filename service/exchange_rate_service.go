package service

import (
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/pkg/casbin"
)

//获取所有汇率列表
func GetAllRate()  (response.RatePage,error){
	rate := model.Rate{}
	allRate,err := rate.GetAllRate()
	if err != nil {
		return response.RatePage{},err
	}
	return allRate,nil
}

// ID获取汇率详情
func RateDetail(id uint) (res response.Rate, err error) {
	rate := model.Rate{}
	rate.ID = id
	res, err = rate.Detail()
	if err != nil {
		return
	}
	return res, err
}

// 币别获取汇率详情
func TypeDetail(name string) (response.Rate,  error) {
	rate := model.Rate{}
	detail, _ := rate.TypeDetail("name")
	return detail, nil
}


//修改汇率
func UpdateRate(rate request.RateEdit)  error{
	defer casbin.ClearEnforcer()
	r := model.Rate{}
	r.ID = uint(rate.ID)
	//逻辑
	return r.Update(rate)
}



