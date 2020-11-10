package hcommon

//账户状态
const (
	Open = 1
	Shut = 2
)

//是否开启google验证
const (
	Yes = 1
	No  = 2
)

// 交易状态
const (
	TxStatusInit   = 0
	TxStatusNotify = 1
)

// 零钱整理状态
const (
	TxOrgStatusInit       = 0
	TxOrgStatusHex        = 1
	TxOrgStatusSend       = 2
	TxOrgStatusConfirm    = 3
	TxOrgStatusFeeHex     = 4
	TxOrgStatusFeeSend    = 5
	TxOrgStatusFeeConfirm = 6
)

// 发送状态
const (
	SendStatusInit    = 0
	SendStatusSend    = 1
	SendStatusConfirm = 2
)

// 发送类型
const (
	SendRelationTypeTx         = 1 //eth充币
	SendRelationTypeWithdraw   = 2 //eth提币
	SendRelationTypeTxErc20    = 3 //erc20充币
	SendRelationTypeTxErc20Fee = 4 //ercx20提币
)

// 通知状态
const (
	NotifyStatusInit = 0
	NotifyStatusFail = 1
	NotifyStatusPass = 2
)

// 通知类型
const (
	NotifyTypeTx              = 1
	NotifyTypeWithdrawSend    = 2
	NotifyTypeWithdrawConfirm = 3
)

// 提币状态
const (
	WithdrawStatusInit    = 0
	WithdrawStatusHex     = 1
	WithdrawStatusSend    = 2
	WithdrawStatusConfirm = 3
)

// 提领状态
const (
	PickStatusWait    = 0 //等待中
	PickStatusDo      = 1 //执行中
	PickStatusSuccess = 2 //成功
	PickStatusCancel  = 3 //已取消
	PickStatusFail    = 4 //失败
)
