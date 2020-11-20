package model

import (
	"j2pay-server/model/response"
)

//eth 钱包交易明细实体
type EthTransaction struct {
	ID             int64
	From           string  `gorm:"default:'';comment:'打币地址'";json:"from"`
	To             string  `gorm:"default:'';comment:'收币地址'";json:"to"`
	Balance        float64 `gorm:"default:0;comment:'金额';";json:"balance"`
	ScheduleStatus int     `gorm:"default:1;comment:'排程状态：1：等待中，:成功,2：失败,3:执行中'"`
	TXID           string  `gorm:"default:'';comment:'交易哈希';";json:"txid"`
	ChainStatus    int     `gorm:"default:1;comment:'链上状态:1：none,2:等待中,3：失败,4:dropped,5：成功'"`
	CreateTime     int64   `gorm:"default:0;comment:'创建时间戳'";json:"create_time"`
}

//hot 钱包交易明细实体
type HotTransaction struct {
	ID             int64
	SystemCode     string  `gorm:"default:'';comment:'系统编号'";json:"system_code"`
	From           string  `gorm:"default:'';comment:'打币地址'";json:"from"`
	To             string  `gorm:"default:'';comment:'收币地址'";json:"to"`
	Balance        float64 `gorm:"default:0;comment:'金额';";json:"balance"`
	Type           int     `gorm:"default:0;comment:'类型:1:代发,2:排程结账,3:手动结账';"`
	ScheduleStatus int     `gorm:"default:1;comment:'排程状态：1：等待中，:成功,2：失败,3:执行中'"`
	TXID           string  `gorm:"default:'';comment:'交易哈希';";json:"txid"`
	GasFee         int64   `gorm:"default:0;comment:'gas费'";json:"gas_fee"`
	ChainStatus    int     `gorm:"default:1;comment:'链上状态:1：none,2:等待中,3：失败,4:dropped,5：成功'"`
	CreateTime     int64   `gorm:"default:0;comment:'创建时间戳'";json:"create_time"`
}

//查询所有eth交易
func (e *EthTransaction) GetAll(page, pageSize int, where ...interface{}) (response.EthTransactionPage, error) {
	count, err2 := e.GetCount(where...)
	if err2 != nil {
		return response.EthTransactionPage{}, err2
	}
	all := response.EthTransactionPage{
		Total:       count,
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.EthTransaction{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := DB.Model(EthTransaction{}).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.EthTransactionPage{}, err
	}
	return all, err
}

//查询所有hot交易
func (h *HotTransaction) GetAll(page, pageSize int, where ...interface{}) (response.HotTransactionPage, error) {
	count, err2 := h.GetCount(where...)
	if err2 != nil {
		return response.HotTransactionPage{}, err2
	}
	all := response.HotTransactionPage{
		Total:       count,
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.HotTransaction{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := DB.Model(EthTransaction{}).Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.HotTransactionPage{}, err
	}
	return all, err
}

//创建eth交易
func (e *EthTransaction) AddEthTx() error {
	tx := DB.Begin()
	if err := tx.Create(e).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

//创建hot交易
func (h *HotTransaction) AddHotTx() error {
	tx := DB.Begin()
	if err := tx.Create(h).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

//修改eth交易 暂时没有此功能
//修改hot钱包交易 只有在等待着才能修改
func UpdateAmount(handleStatus int, address []Address) (err error) {
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

// 获取所有ETH交易数量
func (e *EthTransaction) GetCount(where ...interface{}) (count int,err error) {
	if len(where) == 0 {
		err = DB.Model(&e).Count(&count).Error
		return
	}
	err = DB.Model(&e).Where(where[0], where[1:]...).Count(&count).Error
	return
}

// 获取所有Hot交易数量
func (h *HotTransaction) GetCount(where ...interface{}) (count int,err error) {
	if len(where) == 0 {
		err = DB.Model(&h).Count(&count).Error
		return
	}
	err = DB.Model(&h).Where(where[0], where[1:]...).Count(&count).Error
	return
}
