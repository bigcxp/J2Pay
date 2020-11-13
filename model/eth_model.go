package model

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jmoiron/sqlx"
	"j2pay-server/hcommon"
	"j2pay-server/pkg/setting"
	"log"
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
	ID int64  `GetDb():"id" json:"id"`
	K  string `gorm:"unique;comment:'键名'";json:"k"` // 配置键名
	V  int64  `gorm:"comment:'键值'";json:"v"`        // 配置键值
}

//eth交易
type TTx struct {
	ID           int64
	UserId       int64  `gorm:"default:0;comment:'组织ID'";json:"user_id"`            //组织ID
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

//通知
type TUserNotify struct {
	ID           int64
	Nonce        string `gorm:"default:'0';comment:'唯一标识'";json:"nonce"`
	SystemID     string `gorm:"default:'';comment:'系统编号'";json:"system_id"` // 系统编号
	UserId       int64  `gorm:"default:0;comment:'组织ID'";json:"user_id"`    //组织ID
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
	RelatedType  int64  `gorm:"default:0;comment:'关联类型 1 零钱整理 2 提币'";json:"related_type"` // 关联类型 1 零钱整理 2 提币 3 代发
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

// GetDb()TTxErc20 t_tx_erc20 数据表 eth  erc20交易
type TTxErc20 struct {
	ID           int64
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
}

// 根据条件获取配置
func  SQLGetTAppConfigIntValueByK(k string) *int64 {
	var v *int64
	v = new(int64)
	GetDb().Table("t_app_config_int").Where("k = ?", k).Select("v").Row().Scan(&v)
	return v
}

// 根据条件获取block配置
func  SQLGetTAppStatusIntValueByK(k string) *int64 {
	var v *int64
	v = new(int64)
	GetDb().Table("t_app_status_int").Where("k = ?", k).Select("v").Row().Scan(&v)
	return v
}

// 根据条件获取str配置
func  SQLGetTAppConfigStrValueByK(k string) *string {
	var v *string
	v = new(string)
	GetDb().Table("t_app_config_str").Where("k = ?", k).Select("v").Row().Scan(&v)
	return v
}

//根据条件获取[]ttx
func  SQLSelectTTxColByOrgForUpdate(orgStatus int) ([]TTx, error) {
	var tx []TTx
	 err := GetDb().Model(&TTx{}).Where("org_status =?", orgStatus).Find(&tx).Error
	if err != nil {
		return nil, err
	}
	return tx, err
}

//根据handle_status获取[]ttx
func  SQLSelectTTxColByStatus(handleStatus int) ([]TTx, error) {
	var tx []TTx
	err := GetDb().Model(&TTx{}).Where("handle_status =?", handleStatus).Find(&tx).Error
	if err != nil {
		return nil, err
	}
	return tx, err
}

//根据handle_status获取[]TTxErc20
func  SQLSelectTTxErc20ColByStatus(handleStatus int64) ([]TTxErc20, error) {
	var te []TTxErc20
	err := GetDb().Model(&TTxErc20{}).Where("handle_status =?", handleStatus).Find(&te).Error
	if err != nil {
		return nil, err
	}
	return te, err
}

//根据org_status获取[]TTxErc20
func SQLSelectTTxErc20ColByOrgForUpdate(orgStatuses []int64) ([]TTxErc20, error) {
	var te []TTxErc20
	err := GetDb().Model(&TTxErc20{}).Where("org_status in (?)", orgStatuses).Find(&te).Error
	if err != nil {
		return nil, err
	}
	return te, err
}

//根据handle_status获取[]TSend
func  SQLSelectTSendColByStatus(handleStatus int) ([]TSend, error) {
	var ts []TSend
	 err := GetDb().Model(&TSend{}).Where("handle_status =?", handleStatus).Find(&ts).Error
	if err != nil {
		return nil, err
	}
	return ts, err
}

//根据address查询地址
func  SQLSelectTAddressKeyColByAddress(addresses []string) ([]Address, error) {
	var address []Address
	if len(addresses) == 0 {
		return address, nil
	}
	rows, err := GetDb().Model(&Address{}).Where("user_address in (?)", addresses).Rows()
	defer rows.Close()
	for rows.Next() {
		var add Address
		rows.Scan(add)
		address = append(address, add)
	}
	if err != nil {
		return nil, err
	}
	return address, err

}

//根据address查询地址
func  SQLGetTAddressKeyColByAddress(address string) *Address {
	var row *Address
	row = new(Address)
	GetDb().Table("address").Where("user_address = ?", address).Row().Scan(&row)
	return row

}

//根据handleStatus 、id查询*TWithdraw
func  SQLGetTWithdrawColForUpdate(id int64, handleStatus int) *TWithdraw {
	var row *TWithdraw
	row = new(TWithdraw)
	GetDb().Table("t_withdraw").Where("handle_status = ? and id = ?", handleStatus, id).Row().Scan(&row)
	return row
}
//根据handleStatus 查询[]*TWithdraw
func  SQLSelectTWithdrawColByStatus(handleStatus int) ([]TWithdraw, error) {
	var th []TWithdraw
	rows, err := GetDb().Model(&TWithdraw{}).Where("handle_status =?", handleStatus).Rows()
	defer rows.Close()
	for rows.Next() {
		var td TWithdraw
		rows.Scan(td)
		th = append(th, td)
	}
	if err != nil {
		return nil, err
	}
	return th, err
}

//获取token配置
func  SQLSelectTAppConfigTokenColAll() ([]TAppConfigToken, error) {
	var tact []TAppConfigToken
	rows, err := GetDb().Model(&TAppConfigToken{}).Select("*").Rows()
	defer rows.Close()
	for rows.Next() {
		var td TAppConfigToken
		rows.Scan(td)
		tact = append(tact, td)
	}
	if err != nil {
		return nil, err
	}
	return tact, err
}

//获取私钥map
func  SQLGetAddressKeyMap(addresses []string) (map[string]*Address, error) {
	itemMap := make(map[string]*Address)
	byAddress, err := SQLSelectTAddressKeyColByAddress(addresses)
	if err != nil {
		return nil, err
	}
	for _, itemRow := range byAddress {
		itemMap[itemRow.UserAddress] = &itemRow
	}
	return itemMap, nil
}

//获取私钥
func  GetPkOfAddress(address string) (*ecdsa.PrivateKey, error) {
	if address == ""{
		return nil,nil
	}
	var addr Address
	GetDb().Table("address").Where("user_address = ?", address).Select("pwd").Take(&addr)
	key := hcommon.AesDecrypt(addr.Pwd, fmt.Sprintf("%s", setting.AesConf.Key))
	if len(key) == 0 {
		log.Panicf("error key of: %s", address)
		return nil, fmt.Errorf("no key of: %s", address)
	}
	if strings.HasPrefix(key, "0x") {
		key = key[2:]
	}
	privateKey, err := crypto.HexToECDSA(key)


	if err != nil {
		log.Panicf("HexToECDSA err: [%T] %s", err, err.Error())
		return nil, err
	}
	return privateKey, nil

}

//获取提币map
func  SQLGetWithdrawMap(ids []int64) (map[int64]*TWithdraw, error) {
	itemMap := make(map[int64]*TWithdraw)
	var pick []TWithdraw
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := GetDb().Model(&TWithdraw{}).Where("id in (?)", ids).Rows()
	defer rows.Close()
	for rows.Next() {
		var add TWithdraw
		rows.Scan(add)
		pick = append(pick, add)
	}
	if err != nil {
		return nil, err
	}
	for _, itemRow := range pick {
		itemMap[itemRow.ID] = &itemRow
	}
	return itemMap, nil
}

//获取组织map
func  SQLGetUserMap(ids []int64) (map[int64]*AdminUser, error) {
	itemMap := make(map[int64]*AdminUser)
	var user []AdminUser
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := GetDb().Model(&AdminUser{}).Where("id in (?)", ids).Rows()
	defer rows.Close()
	for rows.Next() {
		var add AdminUser
		rows.Scan(add)
		user = append(user, add)
	}
	if err != nil {
		return nil, err
	}
	for _, itemRow := range user {
		itemMap[itemRow.ID] = &itemRow
	}
	return itemMap, nil
}

//  获取erc代币map
func SQLGetAppConfigTokenMap(ids []int64) (map[int64]*TAppConfigToken, error) {
	itemMap := make(map[int64]*TAppConfigToken)
	var token []TAppConfigToken
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := GetDb().Model(&TAppConfigToken{}).Where("id in (?)", ids).Rows()
	defer rows.Close()
	for rows.Next() {
		var add TAppConfigToken
		rows.Scan(add)
		token = append(token, add)
	}
	if err != nil {
		return nil, err
	}
	for _, itemRow := range token {
		itemMap[itemRow.ID] = &itemRow
	}
	return itemMap, nil
}

//获取nonce
func SQLGetTSendMaxNonce(address string) int64 {
	var i int64
	GetDb().Table("t_send").Where("from_address = ?", address).Select("IFNULL(MAX(nonce), -1)").Row().Scan(&i)
	return i + 1
}

// 获取地址的打包数额
func SQLGetTSendPendingBalanceReal(address string) string {
	var i string
	GetDb().Table("t_send").Where("from_address = ? and handle_status <? limit 1", address, 2).Select("IFNULL(SUM(CAST(balance_real as DECIMAL(65,18))), \"0\")").Row().Scan(&i)
	return i
}

//创建多个交易
func  SQLCreateIgnoreManyTTx(rows []*TTx) (int64, error) {
	if rows == nil {
		return 0, nil
	}
	for _, v := range rows {
		tx := GetDb().Begin()
		if err := tx.Model(&TTx{}).Create(v).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		tx.Commit()
	}
	return 0, nil
}

//创建发送数据
func SQLCreateTSend(rows *TSend) (int64, error) {
	tx := GetDb().Begin()
	if err := tx.Model(&TSend{}).Create(rows).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	return 0, nil

}

//创建多个发送数据
func  SQLCreateIgnoreManyTSend(rows []*TSend, isIgnore bool) (int64, error) {
	if len(rows) == 0 || rows == nil {
		return 0, nil
	}
	tx := GetDb().Begin()
	//需要做逻辑插入
	var args []interface{}
	if rows[0].ID > 0 {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ID,
					row.RelatedType,
					row.RelatedID,
					row.TokenID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.BalanceReal,
					row.Gas,
					row.GasPrice,
					row.Nonce,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.RelatedType,
					row.RelatedID,
					row.TokenID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.BalanceReal,
					row.Gas,
					row.GasPrice,
					row.Nonce,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	}
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_send ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`)
	var err error
	insertArgs := strings.Repeat("(?),", len(rows))
	insertArgs = strings.TrimSuffix(insertArgs, ",")
	query1 := fmt.Sprintf(query.String(), insertArgs)
	query1, args, err = sqlx.In(query1, args...)
	if err != nil {
		return 0, err
	}
	tx.Exec(query1, args)
	tx.Commit()
	return 0, nil

}

//创建通知
func  SQLCreateIgnoreManyTProductNotify(rows []*TUserNotify) (int64, error) {
	if len(rows) == 0 || rows == nil {
		return 0, nil
	}
	for _, v := range rows {
		tx := GetDb().Begin()
		if err := tx.Model(&TUserNotify{}).Create(v).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		tx.Commit()
	}
	return 0, nil
}

//创建多个TTxErc20对象
func  SQLCreateIgnoreManyTTxErc20(rows []*TTxErc20) (int64, error) {
	if rows == nil || len(rows) == 0 {
		return 0, nil
	}
	for _, v := range rows {
		tx := GetDb().Begin()
		if err := tx.Model(&TTxErc20{}).Create(v).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		tx.Commit()
	}
	return 0, nil
}

//更新区块
func  SQLUpdateTAppStatusIntByKGreater(row TAppStatusInt) (err error) {
	if row.K == "" {
		return
	}
	tx := GetDb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	//更新单个记录
	err = tx.Model(&TAppStatusInt{}).Where("k = ?", row.K).
		Update("v", row.V).Error
	return
}

//更新gas费用
func  SQLUpdateTAppStatusIntByK(row *TAppStatusInt) (err error) {
	if row.K == "" {
		return
	}
	tx := GetDb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Model(&TAppStatusInt{}).Where("k =?",row.K).
		Update("v",row.V).Error
	return
}

//更改ttx的org状态
func SQLUpdateTTxOrgStatusByIDs(ids []int64, row *TTx) (err error) {
	if ids == nil {
		return
	}
	tx := GetDb().Begin()
	for _, v := range ids {
		err = tx.Model(&TTx{}).Where("id = ?", v).
			Update("org_status,org_msg,org_time", row.OrgStatus, row.OrgMsg, row.OrgTime).Error
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
	if ids == nil{
		return
	}
	tx := GetDb().Begin()
	for _, v := range ids {
		var ttx TTx
		GetDb().First(&ttx, v)
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
func  SQLUpdateTWithdrawStatusByIDs(ids []int64, row *TWithdraw) (err error) {
	if ids == nil {
		return
	}
	tx := GetDb().Begin()
	for _, v := range ids {
		err = tx.Model(TWithdraw{}).Where("id = ?", v).
			Update("handle_status,handle_msg,handle_time", row.HandleStatus, row.HandleMsg, row.HandleTime).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()

		}
	}
	return
}

//更新erc20零钱整理状态
func  SQLUpdateTTxErc20OrgStatusByIDs(ids []int64, row *TTxErc20) (err error) {
	if ids == nil{
		return
	}
	tx := GetDb().Begin()
	for _, v := range ids {
		err = tx.Model(&TTxErc20{}).Where("id = ?", v).
			Update("org_status,org_msg,org_time", row.OrgStatus, row.OrgMsg, row.OrgTime).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return
}

//根据ids更新erc20处理整理状态
func  SQLUpdateTTxErc20StatusByIDs(ids []int64, row TTxErc20) (int64, error) {
	if ids == nil{
		return 0, nil
	}
	tx := GetDb().Begin()
	for _, v := range ids {
		err := tx.Model(&TTxErc20{}).Where("id = ?", v).
			Update("handle_status,handle_msg,handle_time", row.HandleStatus, row.HandleMsg, row.HandleTime).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()

		}
	}
	return 0, nil
}

//更新提领
func  SQLUpdateTWithdrawGenTx(row *TWithdraw) (err error) {
	if row.ID == 0 {
		return
	}
	tx := GetDb().Begin()
	err = tx.Model(&TWithdraw{}).Where("id = ?", row.ID).
		Update("tx_hash,handle_status,handle_msg,handle_time", row.TxHash, row.HandleStatus, row.HandleMsg, row.HandleTime).Error
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()

	}
	return err
}

//更新发送状态
func  SQLUpdateTSendStatusByIDs(ids []int64, row *TSend) (err error) {
	if ids == nil{
		return nil
	}
	tx := GetDb().Begin()
	for _, v := range ids {
		err = tx.Model(&TSend{}).Where("id = ?", v).
			Update("handle_status,handle_msg,handle_time", row.HandleStatus, row.HandleMsg, row.HandleTime).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()

		}
	}
	return
}

//将地址分配给用户
func  ToAddress(userId int, useTag int64, addr []Address) (err error) {
	if userId == 0{
		return
	}
	tx := GetDb().Begin()
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
func  SQLUpdateTProductNotifyStatusByID(row *TUserNotify) (err error) {
	if row.ID == 0{
		return
	}
	tx := GetDb().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	updateInfo := map[string]interface{}{
		"handle_status": row.HandleStatus,
		"handle_msg":    row.HandleMsg,
		"update_time":   row.UpdateTime,
	}
	err = tx.Model(&row).Updates(updateInfo).Error
	if err != nil {
		return err
	}
	return err
}

// SQLSelectTProductNotifyColByStatusAndTime 根据ids获取通知
func  SQLSelectTProductNotifyColByStatusAndTime(status int64, time int64) ([]TUserNotify, error) {
	var rows []TUserNotify
	err := GetDb().Select("id,url,msg").Where("handle_status = ? and update_time < ?", status, time).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
