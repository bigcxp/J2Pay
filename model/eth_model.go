package model

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"j2pay-server/hcommon"
	"j2pay-server/pkg/setting"
	"strings"
)

//eth整型参数
type TAppConfigInt struct {
	ID int64
	K  string `gorm:"unique;comment:'键名'";json:"k"` // 配置键名
	V  int64  `gorm:"comment:'键值'";json:"v"`        // 配置键值
}

//eth字符串类型参数
type TAppConfigStr struct {
	ID int64
	K  string `gorm:"unique;comment:'键名'";json:"k"` // 配置键名
	V  string `gorm:"comment:'键值'";json:"v"`        // 配置键值
}

//eth block_Number 当前处理区块数 状态表
type TAppStatusInt struct {
	ID int64  `db:"id" json:"id"`
	K  string `gorm:"unique;comment:'键名'";json:"k"` // 配置键名
	V  int64  `gorm:"comment:'键值'";json:"v"`        // 配置键值
}

//eth交易
type TTx struct {
	ID           int64
	SystemID     string `gorm:"default:'';comment:'系统编号'";json:"system_id"`         // 系统编号
	TxID         string `gorm:"unique;comment:'交易id'";json:"tx_id"`                 // 交易id
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
}

//erc20 代币 合约配置
type TAppConfigToken struct {
	ID            int64
	TokenAddress  string `gorm:"default:'';comment:'合约地址'";json:"token_address"`    //合约地址
	TokenDecimals int64  `gorm:"default:0;comment:'合约小数';" ;json:"token_decimals"`  //合约小数
	TokenSymbol   string `gorm:"default:'';comment:'代币单位'";json:"token_symbol"`     //代币单位
	ColdAddress   string `gorm:"default:'';comment:'冷钱包地址'" ;json:"cold_address"`   //冷钱包地址
	HotAddress    string `gorm:"default:'';comment:'热钱包地址'" json:"hot_address"`     //热钱包地址
	OrgMinBalance string `gorm:"default:'0';comment:'最小金额'";json:"org_min_balance"` //最小金额
	CreateTime    int64  `gorm:"default:0;comment:'创建时间'" ;json:"create_time"`      //创建时间
}

// const TProductNotify
const (
	DBColTProductNotifyID           = "t_product_notify.id"
	DBColTProductNotifyNonce        = "t_product_notify.nonce"
	DBColTProductNotifyProductID    = "t_product_notify.product_id"
	DBColTProductNotifyItemType     = "t_product_notify.item_type"
	DBColTProductNotifyItemID       = "t_product_notify.item_id"
	DBColTProductNotifyNotifyType   = "t_product_notify.notify_type"
	DBColTProductNotifyTokenSymbol  = "t_product_notify.token_symbol"
	DBColTProductNotifyURL          = "t_product_notify.url"
	DBColTProductNotifyMsg          = "t_product_notify.msg"
	DBColTProductNotifyHandleStatus = "t_product_notify.handle_status"
	DBColTProductNotifyHandleMsg    = "t_product_notify.handle_msg"
	DBColTProductNotifyCreateTime   = "t_product_notify.create_time"
	DBColTProductNotifyUpdateTime   = "t_product_notify.update_time"
)

//通知
type TProductNotify struct {
	ID           int64
	Nonce        string `gorm:"default:'0';comment:'唯一标识'";json:"nonce"`
	SystemID     string `gorm:"default:'';comment:'系统编号'";json:"system_id"` // 系统编号
	ItemType     int64  `gorm:"default:0;comment:'目标类型'";json:"item_type"`  //目标类型
	ItemID       int64  `gorm:"default:0;comment:'目标id'";json:"item_id"`
	NotifyType   int64  `gorm:"default:0;comment:'通知类型'";json:"notify_type"`       //通知类型 1 提领 2 充币
	TokenSymbol  string `gorm:"default:'USDT';comment:'代币类型'";json:"token_symbol"` //代币类型
	URL          string `gorm:"default:'';comment:'通知url'";json:"url"`             //通知url
	Msg          string `gorm:"default:'';comment:'通知消息'";json:"msg"`              //通知消息
	HandleStatus int64  `gorm:"default:0;comment:'处理状态'";json:"handle_status"`     //处理状态
	HandleMsg    string `gorm:"default:'';comment:'处理信息'";json:"handle_msg"`       //处理信息
	CreateTime   int64  `gorm:"default:0;comment:'创建时间戳'";json:"create_time"`      //创建时间戳
	UpdateTime   int64  `gorm:"default:0;comment:'更新时间戳'";json:"update_time"`      //更新时间戳
}

//eth 发送交易
type TSend struct {
	ID           int64
	RelatedType  int64  `gorm:"default:0;comment:'关联类型 1 零钱整理 2 提币'";json:"related_type"` // 关联类型 1 零钱整理 2 提币
	RelatedID    int64  `gorm:"default:0;comment:'管理id'";json:"related_id"`               // 关联id
	TxID         string `gorm:"default:'';comment:'tx hash'";json:"tx_id"`                // tx hash
	TokenID      int64  `gorm:"default:0;comment:'合约id'";json:"token_id"`                 //合约id
	FromAddress  string `gorm:"default:'';comment:'打币地址'";json:"from_address"`            // 打币地址
	ToAddress    string `gorm:"default:'';comment:'收币地址'";json:"to_address"`              // 收币地址
	BalanceReal  string `gorm:"default:'';comment:'打币金额 ether'";json:"balance_real"`      // 打币金额 Ether
	Gas          int64  `gorm:"default:0;comment:'gas消耗'";json:"gas"`                     // gas消耗
	GasPrice     int64  `gorm:"default:0;comment:'gasPrice'";json:"gas_price"`            // gasPrice
	Nonce        int64  `gorm:"default:0;comment:'nonce'";json:"nonce"`                   // nonce
	Hex          string `gorm:"default:'';comment:'tx raw hex'";json:"hex"`               // tx raw hex
	CreateTime   int64  `gorm:"default:0;comment:'创建时间'";json:"create_time"`              // 创建时间
	HandleStatus int64  `gorm:"default:0;comment:'处理状态'";json:"handle_status"`            // 处理状态
	HandleMsg    string `gorm:"default:'';comment:'处理消息'";json:"handle_msg"`              // 处理消息
	HandleTime   int64  `gorm:"default:0;comment:'处理时间'" json:"handle_time"`              // 处理时间
}

// DBTTxErc20 t_tx_erc20 数据表 eth  erc20交易
type TTxErc20 struct {
	ID           int64
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
}

// DBTWithdraw t_withdraw 数据表  提币 代发
type TWithdraw struct {
	ID           int64
	WithdrawType int64  `gorm:"default:0;comment:' 类型 1 提币 2 代发'";json:"related_type"` // 类型 1 提币 2 代发
	SystemID     string `gorm:"default:'';comment:'系统编号'";json:"system_id"`            // 系统编号
	MerchantID   string `gorm:"default:'';comment:'商户订单编号'";json:"merchant_id"`        // 商户订单编号
	ToAddress    string `gorm:"default:'';comment:'提币地址'";json:"to_address"`           // 提币地址
	Symbol       string `gorm:"default:'';comment:'币种'";json:"symbol"`                 //币种
	BalanceReal  string `gorm:"default:'';comment:'金额'";json:"balance_real"`           // 提币金额
	TxHash       string `gorm:"default:'';comment:'提币tx hash'";json:"tx_hash"`         // 提币tx hash
	CreateTime   int64  `gorm:"default:0;comment:'创建时间'";json:"create_time"`           // 创建时间
	HandleStatus int64  `gorm:"default:0;comment:'处理状态'";json:"handle_status"`         // 处理状态
	HandleMsg    string `gorm:"default:'';comment:'处理消息'";json:"handle_msg"`           // 处理消息
	HandleTime   int64  `gorm:"default:'0';comment:'处理时间'";json:"handle_time"`         // 处理时间

}

// 根据条件获取配置
func SQLGetTAppConfigIntValueByK(where ...interface{}) (ci TAppConfigInt) {
	Getdb().First(&ci, where...)
	return
}

// 根据条件获取block配置
func SQLGetTAppStatusIntValueByK(where ...interface{}) (si TAppStatusInt) {
	Getdb().First(&si, where...)
	return
}

// 根据条件获取str配置
func SQLGetTAppConfigStrValueByK(where ...interface{}) (cf TAppConfigStr) {
	Getdb().First(&cf, where...)
	return
}

// 根据条件获取ttx
func SQLGetTTx(where ...interface{}) (tx TTx) {
	Getdb().First(&tx, where...)
	return
}

//根据条件获取[]ttx
func SQLSelectTTxColByOrgForUpdate(where ...interface{}) (tx []*TTx) {
	Getdb().First(&tx, where...)
	return
}

//根据条件获取[]TSend
func SQLSelectTSendColByStatus(where ...interface{}) (ts []*TSend) {
	Getdb().First(&ts, where...)
	return
}

//根据ids查询地址
func SQLSelectTAddressKeyColByAddress(addresses []string) ([]*Address, error) {
	var rows []*Address
	err := Getdb().Model(&Address{}).Find(&rows, "id in (?)", strings.Split(strings.Join(addresses, ","), ",")).Error
	return rows, err
}

//根据address查询地址
func SQLGetTAddressKeyColByAddress(address string) (*Address, error) {
	var row Address
	err := Getdb().Model(&Address{}).Find(&row, "user_address = ? ", address).Error
	return &row, err
}

//根据status 、id查询
func SQLGetTWithdrawColForUpdate(id int64, status int) (*TWithdraw, error) {
	var rows *TWithdraw
	err := Getdb().Model(&TWithdraw{}).Find(&rows, "handle_status = ? and id = ?", status, id).Error
	return rows, err
}

//根据status 查询需要处理的提币数据
func SQLSelectTWithdrawColByStatus(twithdraws int) ([]*TWithdraw, error) {
	var rows []*TWithdraw
	err := Getdb().Model(&TWithdraw{}).Find(&rows, "handle_status = ?", twithdraws).Error
	return rows, err
}

//获取token配置
func SQLSelectTAppConfigTokenColAll() ([]*TAppConfigToken, error) {
	var rows []*TAppConfigToken
	err := Getdb().Model(&TAppConfigToken{}).Find(&rows).Error
	return rows, err
}

//获取私钥map
func SQLGetAddressKeyMap(addresses []string) (map[string]*Address, error) {
	itemMap := make(map[string]*Address)
	err := Getdb().Where("user_address in (?)", addresses).Find(&itemMap).Error
	if err != nil {
		return nil, err
	}
	for _, itemRow := range itemMap {
		itemMap[itemRow.UserAddress] = itemRow
	}
	return itemMap, nil
}

//获取私钥
func GetPkOfAddress(address string) (*ecdsa.PrivateKey, error) {
	var addr Address
	err := Getdb().Select("pwd").Where("user_address = ?", address).Find(&addr).Error
	if err != nil {
		return nil, err
	}
	key := hcommon.AesDecrypt(addr.Pwd, fmt.Sprintf("%s", setting.AesConf))
	if len(key) == 0 {
		hcommon.Log.Errorf("error key of: %s", address)
		return nil, fmt.Errorf("no key of: %s", address)
	}
	if strings.HasPrefix(key, "0x") {
		key = key[2:]
	}
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		hcommon.Log.Errorf("HexToECDSA err: [%T] %s", err, err.Error())
		return nil, err
	}
	return privateKey, nil

}

//获取提币map
func SQLGetWithdrawMap(ids []int64) (map[int64]*TWithdraw, error) {

	itemMap := make(map[int64]*TWithdraw)
	err := Getdb().Where("id in (?)", ids).Find(&itemMap).Error
	if err != nil {
		return nil, err
	}
	for _, itemRow := range itemMap {
		itemMap[itemRow.ID] = itemRow
	}
	return itemMap, nil
}

//  获取erc代币map
func SQLGetAppConfigTokenMap(ids []int64) (map[int64]*TAppConfigToken, error) {
	itemMap := make(map[int64]*TAppConfigToken)
	//var token []TAppConfigToken
	err := Getdb().Where("id in (?)", ids).Scan(&itemMap).Error
	if err != nil {
		return nil, err
	}
	for _, itemRow := range itemMap {
		itemMap[itemRow.ID] = itemRow
	}
	return itemMap, nil
}

//获取nonce
func SQLGetTSendMaxNonce(address string) (int64, error) {
	var i int64
	err := Getdb().Select("IFNULL(MAX(nonce), -1)").Where("user_address = ?", address).Find(&i).Error
	if err != nil {
		return 0, err
	}
	return i + 1, nil
}

// 获取地址的打包数额
func SQLGetTSendPendingBalanceReal(address string) (string, error) {
	var i string
	err := Getdb().Select("IFNULL(SUM(CAST(balance_real as DECIMAL(65,18))), \"0\")").Where("from_address = ? and handle_status < ?", address, 2).Limit(1).Find(&i).Error
	if err != nil {
		return "", err
	}
	return i, nil
}

//根据状态获取Erc20条约
func SQLSelectTTxErc20ColByStatus(status int64) ([]*TTxErc20, error) {
	var rows []*TTxErc20
	err := Getdb().Model(&TTxErc20{}).Find(&rows, "handle_status = ?", status).Error
	return rows, err
}

func SQLSelectTTxErc20ColByOrgForUpdate(orgStatuses []int64) ([]*TTxErc20, error) {
	var rows []*TTxErc20
	err := Getdb().Model(&TTxErc20{}).Find(&rows, "org_status = ?", orgStatuses).Error
	return rows, err
}

//创建多个交易
func SQLCreateIgnoreManyTTx(rows []*TTx) (int64, error) {
	tx := Getdb().Begin()
	if len(rows) == 0 {
		return 0, nil
	}
	for _, v := range rows {
		if err := Getdb().Model(&TTx{}).Create(v).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return 0, nil
}

//创建发送数据
func SQLCreateTSend(rows *TSend) (int64, error) {
	tx := Getdb().Begin()
	if err := Getdb().Model(&TSend{}).Create(rows).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return 0, nil

}

//创建多个发送数据
func SQLCreateIgnoreManyTSend(rows []*TSend) (int64, error) {
	tx := Getdb().Begin()
	if len(rows) == 0 {
		return 0, nil
	}
	for _, v := range rows {
		if err := Getdb().Model(&TSend{}).Create(v).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return 0, nil

}

//创建通知
func SQLCreateIgnoreManyTProductNotify(rows []*TProductNotify) (int64, error) {
	tx := Getdb().Begin()
	if len(rows) == 0 {
		return 0, nil
	}
	for _, v := range rows {
		if err := Getdb().Model(&TProductNotify{}).Create(v).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return 0, nil
}

//创建多个TTxErc20对象
func SQLCreateIgnoreManyTTxErc20(rows []*TTxErc20) (int64, error) {
	tx := Getdb().Begin()
	if len(rows) == 0 {
		return 0, nil
	}
	for _, v := range rows {
		if err := Getdb().Model(&TTxErc20{}).Create(v).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return 0, nil
}

//更新区块
func SQLUpdateTAppStatusIntByKGreater(row *TAppStatusInt) (err error) {
	tx := Getdb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	status := SQLGetTAppStatusIntValueByK("k = ?", row.K)
	err = tx.Model(&status).
		Updates(TAppStatusInt{V: row.V}).Error
	return
}

//更新gas费用
func SQLUpdateTAppStatusIntByK(row *TAppStatusInt) (err error) {
	tx := Getdb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	status := SQLGetTAppStatusIntValueByK("k = ?", row.K)
	err = tx.Model(&status).
		Updates(TAppStatusInt{V: row.V}).Error
	return
}

//更改ttx的org状态
func SQLUpdateTTxOrgStatusByIDs(ids []int64, row *TTx) (err error) {
	tx := Getdb().Begin()
	for _, v := range ids {
		tTx := SQLGetTTx("id = ?", v)
		err = tx.Model(&tTx).
			Updates(TTx{OrgStatus: row.OrgStatus, OrgMsg: row.OrgMsg, OrgTime: row.OrgTime}).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return
}

//更改ttx的handle状态
func SQLUpdateTTxStatusByIDs(ids []int64, row *TTx) (err error) {
	tx := Getdb().Begin()
	for _, v := range ids {
		var ttx TTx
		Getdb().First(&ttx, v)
		err = tx.Model(&ttx).
			Updates(TTx{HandleStatus: row.HandleStatus, HandleMsg: row.HandleMsg, HandleTime: row.HandleTime}).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return
}

//更新提币状态
func SQLUpdateTWithdrawStatusByIDs(ids []int64, row *TWithdraw) (err error) {
	tx := Getdb().Begin()
	for _, v := range ids {
		var tw TWithdraw
		Getdb().First(&tw, v)
		err = tx.Model(&tw).
			Updates(TWithdraw{HandleStatus: row.HandleStatus, HandleMsg: row.HandleMsg, HandleTime: row.HandleTime}).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return
}

//更新erc20零钱整理状态
func SQLUpdateTTxErc20OrgStatusByIDs(ids []int64, row *TTxErc20) (err error) {
	tx := Getdb().Begin()
	for _, v := range ids {
		var trc TTxErc20
		Getdb().First(&trc, v)
		err = tx.Model(&trc).
			Updates(TTxErc20{OrgStatus: row.OrgStatus, OrgMsg: row.OrgMsg, OrgTime: row.OrgTime}).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return
}

//根据ids更新erc20处理整理状态
func SQLUpdateTTxErc20StatusByIDs(ids []int64, row TTxErc20) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tx := Getdb().Begin()
	for _, v := range ids {
		var trc TTxErc20
		Getdb().First(&trc, v)
		err := tx.Model(&trc).
			Updates(TTxErc20{HandleStatus: row.HandleStatus, HandleMsg: row.HandleMsg, HandleTime: row.HandleTime}).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return 0, nil
}

//更新Twitndraw
func SQLUpdateTWithdrawGenTx(row *TWithdraw) (err error) {
	tx := Getdb().Begin()
	err = tx.Model(&row).
		Updates(TWithdraw{TxHash: row.TxHash, HandleStatus: row.HandleStatus, HandleMsg: row.HandleMsg, HandleTime: row.HandleTime}).Error
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

//更新发送状态
func SQLUpdateTSendStatusByIDs(ids []int64, row *TSend) (err error) {
	tx := Getdb().Begin()
	for _, v := range ids {
		var ts TSend
		Getdb().First(&ts, v)
		err = tx.Model(&ts).
			Updates(TWithdraw{HandleStatus: row.HandleStatus, HandleMsg: row.HandleMsg, HandleTime: row.HandleTime}).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return
}

//将地址分配给用户
func ToAddress(userId int, useTag int64, addr []Address) (err error) {
	tx := Getdb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	for _, v := range addr {
		err = tx.Model(&v).
			Updates(Address{UserId: userId, UseTag: useTag}).Error
	}
	return
}

// SQLUpdateTProductNotifyStatusByID 更新
func SQLUpdateTProductNotifyStatusByID(row *TProductNotify) (err error) {
	tx := Getdb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Model(&row).Updates(TProductNotify{HandleStatus: row.HandleStatus, HandleMsg: row.HandleMsg, UpdateTime: row.UpdateTime}).Error
	if err != nil {
		return err
	}
	return err
}

// SQLSelectTProductNotifyColByStatusAndTime 根据ids获取通知
func SQLSelectTProductNotifyColByStatusAndTime(cols []string, status int64, time int64) ([]*TProductNotify, error) {
	var rows []*TProductNotify
	err := Getdb().Model(&TProductNotify{}).Find(&rows, "handle_status = ? and update_time < ?", status, time).Error

	if err != nil {
		return nil, err
	}
	return rows, nil
}
