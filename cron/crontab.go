package cron

import (
	"github.com/robfig/cron/v3"
	"j2pay-server/heth"
	"j2pay-server/model"
	"log"
)

// 定时处理检测任务
func Cron() {
	go func() {
		//初始化数据库
		model.Setup()
		c := cron.New(
			cron.WithSeconds(),
			cron.WithChain(
				cron.Recover(cron.DefaultLogger),
			),
		)
		var err error
		// --- common --
		// 检测 通知发送
		_, err = c.AddFunc("@every 5m", heth.CheckDoNotify)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// --- eth ---
		// 检测 eth 生成地址
		_, err = c.AddFunc("@every 10m", heth.CheckAddressFree)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// 检测 eth 冲币
		_, err = c.AddFunc("@every 20m", heth.CheckBlockSeek)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// 检测 eth 零钱整理=>充币地址整理到冷钱包
		_, err = c.AddFunc("@every 30m", heth.CheckAddressOrg)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// 检测 eth 提币
		_, err = c.AddFunc("@every 25m", heth.CheckWithdraw)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// 检测 eth 发送交易
		_, err = c.AddFunc("@every 15m", heth.CheckRawTxSend)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// 检测 eth 交易上链=》是否完成打包
		_, err = c.AddFunc("@every 18m", heth.CheckRawTxConfirm)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// 检测 eth 通知到账
		_, err = c.AddFunc("@every 22m", heth.CheckTxNotify)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// 检测 eth gas price
		_, err = c.AddFunc("@every 28m", heth.CheckGasPrice)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}

		// --- erc20 ---
		// 检测 erc20 冲币
		_, err = c.AddFunc("@every 20m", heth.CheckErc20BlockSeek)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// 检测 erc20 通知到账
		_, err = c.AddFunc("@every 30m", heth.CheckErc20TxNotify)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// 检测 erc20 零钱整理
		_, err = c.AddFunc("@every 40m", heth.CheckErc20TxOrg)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		// 检测 erc20 提币
		_, err = c.AddFunc("@every 10m", heth.CheckErc20Withdraw)
		if err != nil {
			log.Panicf("cron add func error: %#v", err)
		}
		c.Start()
	}()

}
