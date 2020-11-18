package hcommon

//账户状态
const (
	Open = 1 //开启
	Shut = 2 //关闭
)

//是否开启google验证
const (
	Yes = 1 //是
	No  = 2 //否
)

// 交处理易状态
const (
	TxStatusInit   = 0 //待处理
	TxStatusNotify = 1 //已处理
)

// 零钱整理状态
const (
	TxOrgStatusInit       = 0 //待整理
	TxOrgStatusHex        = 1 //生成交易
	TxOrgStatusSend       = 2 //已发送交易
	TxOrgStatusConfirm    = 3 //已完成
	TxOrgStatusFeeHex     = 4 //生成手续费交易
	TxOrgStatusFeeSend    = 5 //发送手续费交易
	TxOrgStatusFeeConfirm = 6 //手续费加以已完成
)

// 发送状态
const (
	SendStatusInit    = 0 //待发送
	SendStatusSend    = 1 //已发送
	SendStatusConfirm = 2 //成功
)

// 发送类型
const (
	SendRelationTypeTx         = 1 //零钱整理
	SendRelationTypeWithdraw   = 2 //提币
	SendRelationTypeSend       = 3 //代发
	SendRelationTypeTxErc20    = 4 //erc20零钱整理
	SendRelationTypeTxErc20Fee = 5 //erc20手续费
)

// 通知状态
const (
	NotifyStatusInit = 0 //待通知
	NotifyStatusFail = 1 //失败
	NotifyStatusPass = 2 //成功
)

// 通知类型
const (
	NotifyTypeTx              = 1 //充币
	NotifyTypeWithdrawSend    = 2 //提币
	NotifyTypeWithdrawConfirm = 3 //交易确认
)

// 提币状态
const (
	WithdrawStatusInit    = 0 //待处理
	WithdrawStatusHex     = 1 //已生成交易
	WithdrawStatusSend    = 2 //已发送
	WithdrawStatusConfirm = 3 //交易成功
	WithdrawStatusCancel  = 4 //已取消
	WithdrawStatusFail    = 5 //失败
)

//提币类型
const (
	Deal     = 0 //零钱整理
	WithDraw = 1 //提领
	Send     = 2 // 代发

)
