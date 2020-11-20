package model

import (
	"github.com/jinzhu/gorm"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
)

//汇率表
type Rate struct {
	gorm.Model
	Currency                 string  `gorm:"default:'';comment:'币别';"`
	OriginalRate             float64 `gorm:"default:0;comment:'原汇率';"`
	Collection               float64 `gorm:"default:0;comment:'代收加权';"`
	Payment                  float64 `gorm:"default:0;comment:'代发加权';"`
	ReceiveWeightType        int     `gorm:"default:0;comment:'代收加权类型：0：百分比，1：固定';"`
	PayWeightType            int     `gorm:"default:0;comment:'代发加权类型：0：百分比，1：固定';"`
	ReceiveWeightValue       float64 `gorm:"default:0;comment:'代收加权值';"`
	PayWeightValue           float64 `gorm:"default:0;comment:'代发加权值';"`
	PayWeightAddOrReduce     int     `gorm:"default:0;comment:'代发增加还是减少 0：增加 1：减少';"`
	ReceiveWeightAddOrReduce int     `gorm:"default:0;comment:'代收增加还是减少 0：增加 1：减少';"`
}

//查询记录
func (r *Rate) GetAllRate() (response.RatePage,error) {
	all := response.RatePage{
		Data: []response.Rate{},
	}
	err := DB.Find(&all.Data).Error
	if err != nil {
		return response.RatePage{}, err
	}
	return  all,nil
}

// 根据ID获取汇率详情
func (r *Rate) Detail(id ...int) (res response.Rate, err error) {
	searchId := r.ID
	if len(id) > 0 {
		searchId = uint(id[0])
	}
	err = DB.Table("rate").
		Where("id = ?", searchId).
		First(&res).
		Error
	return
}

// 根据ID获取汇率详情
func (r *Rate) TypeDetail(name ...string) (res response.Rate, err error) {
	err = DB.Table("rate").
		Where("currency = ?", name).
		First(&res).
		Error
	return
}

//修改代收、代发加权
func (r *Rate) Update(rate request.RateEdit) (err error) {
	tx := DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()

		}
	}()
	rates,err := GetRateByWhere("id = ?", rate.ID)
	if err != nil {
		return err
	}
	err = tx.Model(&rates).
		Updates(Rate{ReceiveWeightType: rate.ReceiveWeightType, PayWeightType: rate.PayWeightType,
			ReceiveWeightAddOrReduce: rate.ReceiveWeightAddOrReduce, PayWeightAddOrReduce: rate.PayWeightAddOrReduce,
			ReceiveWeightValue: rate.ReceiveWeightValue, PayWeightValue: rate.PayWeightValue}).Error
	return

}

// 根据条件获取详情
func GetRateByWhere(where ...interface{}) (ra Rate,err error) {
	err = DB.First(&ra, where...).Error
	return
}
