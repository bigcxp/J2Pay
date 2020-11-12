package cron

import (
	"github.com/robfig/cron/v3"
	"j2pay-server/hcommon"
	"j2pay-server/heth"
)

// 定时处理检测任务
func Cron() {
	go func() {
		c := cron.New(
			cron.WithSeconds(),
			cron.WithChain(
				cron.Recover(cron.DefaultLogger),
			),
		)
		var err error
		// --- common --
		// 检测 通知发送
		_, err = c.AddFunc("@every 1m", heth.CheckDoNotify)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// --- eth ---
		// 检测 eth 生成地址
		_, err = c.AddFunc("@every 1m", heth.CheckAddressFree)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// 检测 eth 冲币
		_, err = c.AddFunc("@every 30s", heth.CheckBlockSeek)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// 检测 eth 零钱整理=>充币地址整理到冷钱包
		_, err = c.AddFunc("@every 10m", heth.CheckAddressOrg)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// 检测 eth 提币
		_, err = c.AddFunc("@every 3m", heth.CheckWithdraw)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// 检测 eth 发送交易
		_, err = c.AddFunc("@every 1m", heth.CheckRawTxSend)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// 检测 eth 交易上链=》是否完成打包
		_, err = c.AddFunc("@every 15s", heth.CheckRawTxConfirm)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// 检测 eth 通知到账
		_, err = c.AddFunc("@every 15s", heth.CheckTxNotify)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// 检测 eth gas price
		_, err = c.AddFunc("@every 2m", heth.CheckGasPrice)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}

		// --- erc20 ---
		// 检测 erc20 冲币
		_, err = c.AddFunc("@every 30s", heth.CheckErc20BlockSeek)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// 检测 erc20 通知到账
		_, err = c.AddFunc("@every 20s", heth.CheckErc20TxNotify)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// 检测 erc20 零钱整理
		_, err = c.AddFunc("@every 10m", heth.CheckErc20TxOrg)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		// 检测 erc20 提币
		_, err = c.AddFunc("@every 3m", heth.CheckErc20Withdraw)
		if err != nil {
			hcommon.Log.Errorf("cron add func error: %#v", err)
		}
		c.Start()
	}()

}
