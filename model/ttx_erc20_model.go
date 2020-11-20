package model

import (
	"j2pay-server/model/request"
	"j2pay-server/model/response"
	"time"
)

//erc20交易明细
type TTxErc20 struct {
	ID           int64
	OrderId      string `gorm:"default:0;comment:'订单编号'";json:"order_id"`           //订单编号
	UserId       int64  `gorm:"default:0;comment:'组织ID'";json:"user_id"`            //组织ID
	TokenID      int64  `gorm:"default:0;comment:'合约id'";json:"token_id"`           //合约id
	SystemID     string `gorm:"default:'';comment:'系统编号'";json:"system_id"`         // 系统编号
	TxID         string `gorm:"default:'';comment:'交易id'";json:"tx_id"`             // 交易id
	FromAddress  string `gorm:"defsult:'';comment:'来源地址'";json:"from_address"`      // 来源地址
	ToAddress    string `gorm:"default:'';comment:'目标地址'";json:"to_address"`        // 目标地址
	BalanceReal  string `gorm:"default:'';comment:'到账金额Ether'";json:"balance_real"` // 到账金额Ether
	CreateTime   int64  `gorm:"default:0;comment:'创建时间戳'";json:"create_time"`       // 创建时间戳
	HandleStatus int64  `gorm:"default:0;comment:'处理状态'";json:"handle_status"`      // 处理状态
	HandleMsg    string `gorm:"default:'';comment:'处理消息'";json:"handle_msg"`        // 处理消息
	HandleTime   int64  `gorm:"default:0;comment:'处理时间'";json:"handle_time"`        // 处理时间戳
	OrgStatus    int64  `gorm:"default:0;comment:'零钱整理状态'";json:"org_status"`       // 零钱整理状态
	OrgMsg       string `gorm:"default:'';comment:'零钱整理消息'";json:"org_msg"`         // 零钱整理消息
	OrgTime      int64  `gorm:"default:0;comment:'零钱整理时间'" json:"org_time"`         // 零钱整理时间
	Status       int    `gorm:"default:1;comment:'状态 1：未绑定，2：已绑定';";json:"status"`  //是否绑定订单
	Remark       string `gorm:"default:'';comment:'备注';";json:"remark"`             // 备注
}

//获取所有交易订单明细列表
func (t *TTxErc20) GetAllErc20Detail(page, pageSize int, where ...interface{}) (response.Erc20Page, error) {
	count, err2 := t.GetCount(where...)
	if err2 != nil {
		return response.Erc20Page{}, err2
	}
	all := response.Erc20Page{
		Total:       count,
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.Erc20List{},
	}
	//分页查询
	offset := GetOffset(page, pageSize)
	err := DB.Table("t_tx_erc20").Order("id desc").Limit(pageSize).Offset(offset).Find(&all.Data, where...).Error
	if err != nil {
		return response.Erc20Page{}, err
	}
	for index, v := range all.Data {
		//时间转换
		all.Data[index].Create = time.Unix(all.Data[index].CreateTime,0)
		if v.OrderId != "" && v.Status == 2 {
			//获取该交易订单
			orderByWhere,err := GetOrderByWhere("order_code = ?", v.OrderId)
			if err != nil {
				return response.Erc20Page{}, err
			}
			//获取组织
			userByWhere, _ := GetUserByWhere("id = ?", orderByWhere.UserId)
			all.Data[index].RealName = userByWhere.RealName
		}
	}
	return all, err
}

// 根据ID获取订单交易明细
func (t *TTxErc20) GetErc20Detail(id ...int) (res response.Erc20List, err error) {
	searchId := t.ID
	if len(id) > 0 {
		searchId = int64(id[0])
	}
	err = DB.Table("t_tx_erc20").
		Where("id = ?", searchId).
		First(&res).
		Error
	//时间转换
	res.Create = time.Unix(res.CreateTime,0)
	if res.OrderId != "" && res.Status == 2 {
		//获取该交易订单
		orderByWhere ,err:= GetOrderByWhere("order_code = ?", res.OrderId)
		if err != nil {
			return response.Erc20List{},err
		}
		//获取组织
		userByWhere, _ := GetUserByWhere("id = ?", orderByWhere.UserId)
		res.RealName = userByWhere.RealName
	}
	return
}

// 创建订单明细
func (t *TTxErc20) Create() error {
	tx := DB.Begin()
	if err := tx.Create(t).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

//修改交易订单明细
func (t *TTxErc20) BindOrder(order request.Erc20Edit) (err error) {
	tx := GetDb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	updateInfo := map[string]interface{}{
		"order_id" : order.OrderId,
		"status":order.Status,
	}
	err = tx.Model(&TTxErc20{}).Where("id = ?",order.ID).Updates(updateInfo).Error
	return
}

// 根据条件获取订单明细详情
func GetErc20ByWhere(where ...interface{}) (t TTxErc20,err error) {
	err = DB.First(&t, where...).Error
	return
}

// 获取所有交易订单数量
func (t *TTxErc20) GetCount(where ...interface{}) (count int,err error) {
	if len(where) == 0 {
		DB.Model(&t).Count(&count)
		return
	}
	err = DB.Model(&t).Where(where[0], where[1:]...).Count(&count).Error
	return
}
