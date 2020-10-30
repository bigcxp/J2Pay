package model

import (
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"time"
)

type Address struct {
	ID           int64
	UserAddress  string  `gorm:"unique;comment:'钱包地址';"json:"user_address"`                                              //地址
	EthAmount    float64 `gorm:"default:0;comment:'以太币余额';";json:"eth_amount"`                                           //以太币余额
	UsdtAmount   float64 `gorm:"default:0;comment:'泰达币余额';";json:"usdt_amount"`                                          //泰达币余额
	UserId       int     `gorm:"default:0;comment:'用户id';";json:"user_id"`                                               //组织id
	Symbol       string  `gorm:"default:'eth';comment:'币种';"json:"symbol"`                                               // 币种
	Pwd          string  `gorm:"default:'';comment:'加密私钥'";json:"pwd"`                                                   // 加密私钥
	Status       int     `gorm:"default:1;comment:'状态 0：所有，1：已完成，2：执行中，3：结账中';";json:"status"`                           //状态 状态 0：所有，1：已完成，2：执行中，3：结账中
	HandleStatus int     `gorm:"default:1;comment:'指派状态 0：所有，1：启用，2：停用';";json:"status"`                                 //指派状态 0：所有，1：启用，2：停用
	UseTag       int64   `gorm:"default:0;comment:'占用标志 -2：作为eth钱包占用， -1：作为热钱包占用 ，0：未占用->其他 作为用户冲币地址占用'";json:"use_tag"` // HandleStatus
	CreateTime   int64   `gorm:"default:0;comment:'创建时间戳'";json:"create_time"`                                           //创建时间戳
	UpdateTime   int64   `gorm:"default:0;comment:'更新时间戳'";json:"update_time"`                                           //更新时间戳
}

//查询所有收款地址
func (a *Address) GetAllAddress(page, pageSize int, where ...interface{}) (response.AddressPage, error) {
	all := response.AddressPage{
		Total:       a.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.Address{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := Getdb().Model(Address{}).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.AddressPage{}, err
	}
	for index, v := range all.Data {
		all.Data[index].RealName = GetUserByWhere("id = ?", v.UserId).RealName
	}
	return all, err
}

//随机获取相对应数量的空闲地址
func GetFreAddress(num int64) ([]Address, error) {
	var rows []Address
	s := Getdb().Raw("SELECT * FROM address WHERE use_tag = 0 ORDER BY RAND() LIMIT ?", num).Scan(&rows).Error
	return rows, s
}

//随机获取商户不在收款中的充币地址
func GetAddress(id int) (Address, error) {
	var row Address
	s := Getdb().Raw("SELECT * FROM address WHERE use_tag = ? ORDER BY RAND() LIMIT 1", id).Scan(&row).Error
	return row, s
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

//启用 停用
func OpenOrStopAddress(handleStatus int, address []Address) (err error) {
	tx := Getdb().Begin()
	for _, v := range address {
		if err = tx.Model(&v).
			Updates(Address{HandleStatus: handleStatus}).Error; err != nil {
			tx.Rollback()
			return err
			tx.Commit()
		}
	}
	return err
}

//编辑用户收款地址 地址在停用状态下
func EditAddress(address []Address,userId int) error {
	tx := Getdb().Begin()
	for _, v := range address {
		if err := tx.Model(&v).
			Updates(Address{UserId:userId}).Error; err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	return nil
}


//删除收款地址
func  AddressDel(addr []Address) error {
	tx := Getdb().Begin()
	for _,v := range addr {
		if err := tx.Delete(&v).Error; err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	return nil
}

//储值  eth钱包=>用户收款地址 eth余额必须足够 =>生成eth交易
func (a *Address) Save(address request.SaveAmount) error {
	tx := Getdb().Begin()
	addresss := GetAddressByWhere("id = ?", address.ID)
	if err := tx.Model(&addresss).
		Updates(Address{EthAmount: address.EthAmount}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//结账 指派状态为停用 eth金额不能小于最小矿工费 =》生成热钱包交易
func (a *Address) Col(address request.Math) error {
	tx := Getdb().Begin()
	addresss := GetAddressByWhere("id = ?", address.ID)
	if err := tx.Model(&addresss).
		Updates(Address{Status: address.Status}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//更新余额
func UpdateBalance(ids request.UpdateAmount) (err error) {
	tx := Getdb().Begin()
	now := time.Now().Unix()
	//查询出钱包地址
	address := GetAllAddress("id in (?)", ids.ID)
	for _, v := range address {
		if err := tx.Model(&v).Updates(Address{UpdateTime: now}).Error; err != nil {
			return err
		}
	}
	tx.Commit()
	return
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

// 根据条件获取钱包地址列表
func GetAllAddress(where ...interface{}) (res []Address) {
	Getdb().Model(Address{}).
		Order("id asc").
		Find(&res, where...)
	return
}
