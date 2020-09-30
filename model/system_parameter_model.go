package model

import (
	"github.com/jinzhu/gorm"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"j2pay-server/validate"
)

//系统参数
type Parameter struct {
	gorm.Model
	Confirmation int     `gorm:"default:0;comment:'交易确认数';"`
	GasLimit     int     `gorm:"default:0;comment:'USDT gas limit';"`
	GasPrice     float64 `gorm:"default:0;comment:'GasPrice';"`
	EthFee       float64 `gorm:"default:0;comment:'ETH最小矿工费';"`
}

//查询系统参数数据
func (p *Parameter) GetDetail() response.Parameter {
	var parameter response.Parameter
	Db.First(&parameter)
	return parameter
}

//更新系统参数
func (p *Parameter) UpdateParameter(edit request.ParameterEdit) (err error) {
	tx := Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	parameter := GetParameterByWhere("id = ?", edit.Id)
	ethFee := validate.Decimal(validate.WrapToFloat64(edit.GasPrice, 10) * validate.Unwrap(int64(edit.GasLimit), 10) * 0.0000000001)
	err = tx.Model(&parameter).
		Updates(Parameter{Confirmation: edit.Confirmation, GasLimit: edit.GasLimit, GasPrice: edit.GasPrice, EthFee: ethFee}).Error
	return
}

//更新gasPrice
func (p *Parameter) UpdateGasPrice(edit request.ParameterEdit) (err error) {
	tx := Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	parameter := GetParameterByWhere("id = ?", edit.Id)
	ethFee := validate.Decimal(validate.WrapToFloat64(edit.GasPrice, 10) * validate.Unwrap(int64(parameter.GasLimit), 10) * 0.0000000001)
	err = tx.Model(&parameter).
		Updates(Parameter{GasPrice: edit.GasPrice, EthFee: ethFee}).Error
	return
}

// 根据条件获取详情
func GetParameterByWhere(where ...interface{}) (pa Parameter) {
	Db.First(&pa, where...)
	return
}

//查询GasFee
func GetGasFeeDetail() response.Parameter {
	var parameter response.Parameter
	Db.First(&parameter)
	return parameter
}
