package model

import (
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"time"
)

type Address struct {
	ID          int64
	UserAddress string    `gorm:"unique;comment:'钱包地址';"json:"user_address"`                               //地址
	EthAmount   float64   `gorm:"default:0;comment:'以太币余额';";json:"eth_amount"`                            //以太币余额
	UsdtAmount  float64   `gorm:"default:0;comment:'泰达币余额';";json:"usdt_amount"`                           //泰达币余额
	UserId      int       `gorm:"TYPE:int(11);NOT NULL;INDEX";json:"user_id"`                              //组织id
	AdminUser   AdminUser `json:"admin_user";gorm:"foreignkey:UserId"`                                     //指定关联外键
	Symbol      string    `gorm:"default:'eth';comment:'币种';"json:"symbol"`                                // 币种
	Pwd         string    `gorm:"default:'';comment:'加密私钥'";json:"pwd"`                                    // 加密私钥
	Status      int       `gorm:"default:1;comment:'状态 0：所有，1：已完成，2：执行中，3：结账中';";json:"status"`            //状态 状态 0：所有，1：已完成，2：执行中，3：结账中
	UseTag      int64     `gorm:"default:0;comment:'占用标志 -1 作为热钱包占用-0 未占用->其他 作为用户冲币地址占用'";json:"use_tag"` // 占用标志 -1 作为热钱包占用-0 未占用->0 作为用户冲币地址占用
	CreateTime  int64     `gorm:"default:0;comment:'创建时间戳'";json:"create_time"`                            //创建时间戳
	UpdateTime  int64     `gorm:"default:0;comment:'更新时间戳'";json:"update_time"`                            //更新时间戳
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
	err := Getdb().Model(&a).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
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
	tx := Getdb().Begin()
	if err := tx.Create(a).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//创建多个钱包地址
func AddMoreAddress(rows []*Address) (int64, error) {
	tx := Getdb().Begin()
	if len(rows) == 0 {
		return 0, nil
	}
	for _, v := range rows {
		if err := Getdb().Model(&Address{}).Create(v).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return 0, nil
}

//编辑用户收款地址
func (a *Address) EditAddress(address request.AddressEdit) error {
	tx := Getdb().Begin()
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
	tx := Getdb().Begin()
	addresss := GetAddressByWhere("id = ?", address.Id)
	if err := tx.Model(&addresss).
		Updates(Address{UseTag: address.UseTag}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//储值
func (a *Address) Save(address request.SaveAmount) error {
	tx := Getdb().Begin()
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
	tx := Getdb().Begin()
	if err := tx.Delete(a, "id = ?", a.ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//结账 只有关闭指派的情况下才能结账
func (a *Address) Col(address request.Math) error {
	tx := Getdb().Begin()
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
	now := time.Now().Unix()
	for _, v := range ids.Id {
		address := GetAddressByWhere("id = ?", v)
		if err := Getdb().Model(&address).Updates(Address{UpdateTime: now}).Error; err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// 获取所有收款地址数量
func (a *Address) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Getdb().Model(&a).Count(&count)
		return
	}
	Getdb().Model(&a).Where(where[0], where[1:]...).Count(&count)
	return
}

// 获取所有收款地址数量
func GetAddressCount(where ...interface{}) (count int64) {
	if len(where) == 0 {
		Getdb().Model(&Address{}).Count(&count)
		return
	}
	Getdb().Model(&Address{}).Where(where[0], where[1:]...).Count(&count)
	return
}

// 根据条件获取钱包地址
func GetAddressByWhere(where ...interface{}) (a Address) {
	Getdb().First(&a, where...)
	return
}
