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
	SendRelationTypeTx         = 1 //eth充币
	SendRelationTypeWithdraw   = 2 //eth提币
	SendRelationTypeTxErc20    = 3 //erc20充币
	SendRelationTypeTxErc20Fee = 4 //ercx20手续费
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
)

// 提领状态
const (
	PickStatusWait    = 0 //等待中
	PickStatusDo      = 1 //执行中
	PickStatusSuccess = 2 //成功
	PickStatusCancel  = 3 //已取消
	PickStatusFail    = 4 //失败
)
