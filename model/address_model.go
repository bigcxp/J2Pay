package model

import (
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"time"
)

type Address struct {
	ID           int64
	UserAddress  string  `gorm:"unique;comment:'钱包地址';"json:"user_address"`                           //地址
	EthAmount    float64 `gorm:"default:0;comment:'以太币余额';";json:"eth_amount"`                        //以太币余额
	UsdtAmount   float64 `gorm:"default:0;comment:'泰达币余额';";json:"usdt_amount"`                       //泰达币余额
	UserId       int     `gorm:"default:0;comment:'用户id';";json:"user_id"`                            //组织id
	Symbol       string  `gorm:"default:'eth';comment:'币种';"json:"symbol"`                            // 币种
	Pwd          string  `gorm:"default:'';comment:'加密私钥'";json:"pwd"`                                // 加密私钥
	Status       int     `gorm:"default:1;comment:'状态 0：所有，1：已完成，2：执行中，3：结账中';";json:"status"`        //状态 状态 0：所有，1：已完成，2：执行中，3：结账中
	HandleStatus int     `gorm:"default:1;comment:'指派状态 0：所有，1：启用，2：停用';";json:"status"`              //指派状态 0：所有，1：启用，2：停用
	UseTag       int64   `gorm:"default:0;comment:'-1：作为热钱包占用 ，0：未占用->其他 作为用户冲币地址占用'";json:"use_tag"` // HandleStatus
	CreateTime   int64   `gorm:"default:0;comment:'创建时间戳'";json:"create_time"`                        //创建时间戳
	UpdateTime   int64   //更新时间戳
}

//查询所有收款地址
func (a *Address) GetAllAddress(page, pageSize int, where ...interface{}) (response.AddressPage, error) {
	count, err2 := a.GetCount(where...)
	if err2 != nil {
		return response.AddressPage{}, err2
	}
	all := response.AddressPage{
		Total:       count,
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.Address{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := DB.Model(Address{}).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.AddressPage{}, err
	}
	for index, v := range all.Data {
		user, _ := GetUserByWhere("id = ?", v.UserId)
		all.Data[index].RealName = user.RealName
	}
	return all, err
}

//随机获取相对应数量的空闲地址
func GetFreAddress(num int64) ([]Address, error) {
	var rows []Address
	s := DB.Raw("SELECT * FROM address WHERE use_tag = 0 ORDER BY RAND() LIMIT ?", num).Scan(&rows).Error
	return rows, s
}

//随机获取商户不在收款中的充币地址
func GetAddress(id int) (Address, error) {
	var row Address
	err := GetDb().Model(&Address{}).Where("use_tag = ?", id).Order("RAND()").Limit(1).Take(&row).Error
	return row, err
}

//随机获取商户不在收款中的充币地址
func (a *Address) FindById(id int64) (Address, error) {
	var row Address
	err := DB.Model(row).Find(&row, "id = ? ", id).Error
	return row, err
}

//新增用户收款地址
func (a *Address) AddAddress() error {
	tx := DB.Begin()
	if err := tx.Create(a).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

//创建多个钱包地址
func AddMoreAddress(rows []*Address) (int64, error) {
	tx := DB.Begin()
	if len(rows) == 0 {
		return 0, nil
	}
	for _, v := range rows {
		if err := tx.Model(&Address{}).Create(v).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()

	return 0, nil
}

//启用 停用
func OpenOrStopAddress(handleStatus int, address []Address) (err error) {
	tx := DB.Begin()
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
func EditAddress(address []Address, userId int) error {
	tx := DB.Begin()
	for _, v := range address {
		if err := tx.Model(&v).
			Updates(Address{UserId: userId}).Error; err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()

	}
	return nil
}

//删除收款地址
func AddressDel(addr []Address) error {
	tx := DB.Begin()
	for _, v := range addr {
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
	tx := DB.Begin()
	addresss,err := GetAddressByWhere("id = ?", address.ID)
	if err != nil {
		return err
	}
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
	tx := DB.Begin()
	addresss ,err:= GetAddressByWhere("id = ?", address.ID)
	if err != nil {
		return err
	}
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
	tx := DB.Begin()
	now := time.Now().Unix()
	//查询出钱包地址
	address,err := GetAllAddress("id in (?)", ids.ID)
	if err != nil {
		return err
	}
	for _, v := range address {
		if err := tx.Model(&v).Updates(Address{UpdateTime: now}).Error; err != nil {
			return err
		}
	}
	tx.Commit()

	return
}

// 获取所有收款地址数量
func (a *Address) GetCount(where ...interface{}) (count int,err error) {
	if len(where) == 0 {
		DB.Model(&a).Count(&count)
		return
	}
	err = DB.Model(&a).Where(where[0], where[1:]...).Count(&count).Error
	return
}

// 获取所有收款地址数量
func GetAddressCount(where ...interface{}) (count int64,err error) {
	if len(where) == 0 {
		DB.Model(&Address{}).Count(&count)
		return
	}
	err = DB.Model(&Address{}).Where(where[0], where[1:]...).Count(&count).Error
	return
}

// 根据条件获取钱包地址
func GetAddressByWhere(where ...interface{}) (a Address,err error) {
	err = DB.First(&a, where...).Error
	return
}

// 根据条件获取钱包地址列表
func GetAllAddress(where ...interface{}) (res []Address,err error) {
	err = DB.Model(Address{}).
		Order("id asc").
		Find(&res, where...).Error
	return
}
