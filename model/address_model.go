package model

import (
	"github.com/jinzhu/gorm"
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"time"
)

type Address struct {
	gorm.Model
	UserAddress string    `gorm:"default:'';comment:'钱包地址';"json:"user_address"`                //地址
	EthAmount   float64   `gorm:"default:0;comment:'以太币余额';";json:"eth_amount"`                 //以太币余额
	UsdtAmount  float64   `gorm:"default:0;comment:'泰达币余额';";json:"usdt_amount"`                //泰达币余额
	UserId      int       `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"user_id"`                   //组织id
	AdminUser   AdminUser `json:"admin_user";gorm:"foreignkey:UserId"`                          //指定关联外键
	SearchTime  time.Time `gorm:"comment:'完成时间';";json:"search_time"`                           //最后查询余额时间
	Status      int       `gorm:"default:1;comment:'状态 0：所有，1：已完成，2：执行中，3：结账中';";json:"status"` //状态
	HandStatus  int       `gorm:"default:1;comment:'指派状态 0：所有，1：启用，2：停用';";json:"status"`       //指派状态
}

//查询所有收款地址
func (a *Address) GetAllAddress(page, pageSize int, where ...interface{}) (response.AddressPage, error) {
	all := response.AddressPage{
		Total:       a.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.AddressList{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := Db.Model(&a).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.AddressPage{}, err
	}
	for index, v := range all.Data {
		all.Data[index].RealName = GetUserByWhere("id = ?", v.UserId).RealName
	}
	return all, err
}

//新增用户收款地址
func (a *Address) AddAddress() error {
	tx := Db.Begin()
	if err := tx.Create(a).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//编辑用户收款地址
func (a *Address) EditAddress(address request.AddressEdit) error {
	tx := Db.Begin()
	addresss := GetAddressByWhere("id = ?", address.Id)
	if err := tx.Model(&addresss).
		Updates(Address{UserId: address.UserId}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//启用 停用
func (a *Address) OpenOrStopAddress(address request.OpenOrStopAddress) error {
	tx := Db.Begin()
	addresss := GetAddressByWhere("id = ?", address.Id)
	if err := tx.Model(&addresss).
		Updates(Address{HandStatus: address.HandStatus}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//储值
func (a *Address) Save(address request.SaveAmount) error {
	tx := Db.Begin()
	addresss := GetAddressByWhere("id = ?", address.Id)
	//判断余额是否足够

	if err := tx.Model(&addresss).
		Updates(Address{EthAmount: address.EthAmount}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//删除收款地址
func (a *Address) Del() error {
	tx := Db.Begin()
	if err := tx.Delete(a, "id = ?", a.ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//结账 只有关闭指派的情况下才能结账
func (a *Address) Col(address request.Math) error {
	tx := Db.Begin()
	addresss := GetAddressByWhere("id = ?", address.Id)
	if err := tx.Model(&addresss).
		Updates(Address{Status: address.Status}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//更新余额
func Update(ids request.UpdateAmount) (addresses []Address, err error) {
	for _, v := range ids.Id {
		address := GetAddressByWhere("id = ?", v)
		if err := Db.Model(&address).Updates(Address{SearchTime: time.Now()}).Error; err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// 获取所有收款地址数量
func (a *Address) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Db.Model(&a).Count(&count)
		return
	}
	Db.Model(&a).Where(where[0], where[1:]...).Count(&count)
	return
}

// 根据条件获取收款地址看详情
func GetAddressByWhere(where ...interface{}) (a Address) {
	Db.First(&a, where...)
	return
}
