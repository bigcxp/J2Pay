package model

import "strings"

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
	V  int64 `gorm:"comment:'键值'";json:"v"`        // 配置键值
}

//eth交易
type TTx struct {
	ID           int64
	SystemID     string `gorm:"default:'';comment:'系统编号'";json:"system_id"`         // 系统编号
	TxID         string `gorm:"unique;comment:'交易id'";json:"tx_id"`             // 交易id
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
	TokenAddress  string `gorm:"default:'';comment:'合约地址'";json:"token_address"`
	TokenDecimals int64  `gorm:"default:0;comment:'合约小数';" ;json:"token_decimals"`
	TokenSymbol   string `gorm:"default:'erc20_usdt';comment:'代币单位'";json:"token_symbol"`
	ColdAddress   string `gorm:"default:'';comment:'冷钱包地址'" ;json:"cold_address"`
	HotAddress    string `gorm:"default:'';comment:'热钱包地址'" json:"hot_address"`
	OrgMinBalance string `gorm:"default:'0';comment:'最小金额'";json:"org_min_balance"`
	CreateTime    int64  `gorm:"default:0;comment:'创建时间'" ;json:"create_time"`
}

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
	TokenID      int64  `gorm:"default:0;comment:'合约id'";json:"token_id"`                 //合约id
	TxID         string `gorm:"default:'';comment:'tx hash'";json:"tx_id"`                // tx hash
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

//根据ids查询
func SQLSelectTAddressKeyColByAddress(addresses []string)  ([]*Address, error) {
	var rows []*Address
	err := Getdb().Find(&rows, "id in (?)", strings.Split(strings.Join(addresses, ","), ",")).Error
	return rows,err
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

//更新区块
func SQLUpdateTAppStatusIntByKGreater(row *TAppStatusInt)(err error)  {
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

