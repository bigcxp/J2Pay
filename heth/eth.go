package heth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	_ "encoding/hex"
	"encoding/json"
	_ "encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	_ "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	_ "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	_ "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	_ "github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/moremorefun/mcommon"
	_ "github.com/parnurzeal/gorequest"
	"j2pay-server/ethclient"
	_ "j2pay-server/ethclient"
	"j2pay-server/hcommon"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/pkg/setting"
	"j2pay-server/pkg/util"
	"log"
	"math"
	_ "math"
	"math/big"
	_ "math/big"
	_ "net/http"
	"strings"
	_ "strings"
	"time"
)

//获取地址和私钥
func genAddressAndAesKey() (string, string, error) {
	// 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyStr := hexutil.Encode(privateKeyBytes)
	// 加密密钥
	privateKeyStrEn := hcommon.AesEncrypt(privateKeyStr, fmt.Sprintf("%s", setting.AesConf.Key))
	// 获取地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("can't change public key")
	}
	// 地址全部储存为小写方便处理
	address := AddressBytesToStr(crypto.PubkeyToAddress(*publicKeyECDSA))
	return address, privateKeyStrEn, nil
}

// CreateHotAddress 创建自用地址 热钱包
func CreateHotAddress(addr request.AddressAdd) ([]string, error) {
	var rows []*model.Address
	// 当前时间
	now := time.Now().Unix()
	var userAddresses []string
	// 遍历差值次数
	for i := int64(0); i < addr.Num; i++ {
		address, privateKeyStrEn, err := genAddressAndAesKey()
		if err != nil {
			return nil, err
		}
		// 存入待添加队列
		rows = append(rows, &model.Address{
			Symbol:       CoinSymbol,
			UserAddress:  address,
			Pwd:          privateKeyStrEn,
			UseTag:       addr.UseTag,
			UserId:       addr.UserId,
			UsdtAmount:   0,
			EthAmount:    0,
			Status:       1,
			HandleStatus: addr.HandleStatus,
			CreateTime:   now,
			UpdateTime:   now,
		})
		userAddresses = append(userAddresses, address)
	}
	// 一次性将生成的地址存入数据库
	_, err := model.AddMoreAddress(rows)
	if err != nil {
		return nil, err
	}
	return userAddresses, nil
}

// CheckAddressFree 检测是否有充足的备用地址
func CheckAddressFree() {
	// 当前时间
	now := time.Now().Unix()
	//地址
	var userAddresses []string
	// 获取配置 允许的最小剩余地址数
	minAddress := model.SQLGetTAppConfigIntValueByK("min_free_address")
	// 获取当前剩余可用地址数
	count := model.GetAddressCount("use_tag = ?", 0)
	// 如果数据库中剩余可用地址小于最小允许可用地址
	if count < *minAddress {
		var rows []*model.Address
		// 遍历差值次数
		for i := int64(0); i < *minAddress-count; i++ {
			address, privateKeyStrEn, err := genAddressAndAesKey()
			if err != nil {
				return
			}
			// 存入待添加队列
			rows = append(rows, &model.Address{
				Symbol:      CoinSymbol,
				UserAddress: address,
				Pwd:         privateKeyStrEn,
				UseTag:      0,
				UserId:      0,
				UsdtAmount:  0,
				EthAmount:   0,
				Status:      0,
				CreateTime:  now,
				UpdateTime:  now,
			})
			userAddresses = append(userAddresses, address)
		}
		// 一次性将生成的地址存入数据库
		_, err := model.AddMoreAddress(rows)
		if err != nil {
			return
		}
	}

}

//将地址分配给商户
func ToMerchantAddress(addr request.AddressAdd) (err error) {
	//从钱包地址随机获取
	address, err := model.GetFreAddress(addr.Num)
	if err != nil {
		return err
	}
	err = model.ToAddress(addr.UserId, addr.UseTag, address)
	return err
}

// CheckBlockSeek 检测到账
func CheckBlockSeek() {
	// 获取配置 延迟确认数
	confirmValue := model.SQLGetTAppConfigIntValueByK("block_confirm_num")
	// 获取状态 当前处理完成的最新的block number
	seek := model.SQLGetTAppStatusIntValueByK("seek_num")
	// rpc 获取当前最新区块数
	rpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	startI := *seek + 1
	endI := rpcBlockNum - *confirmValue + 1
	if startI < endI {
		// 手续费钱包列表
		feeAddressValue := model.SQLGetTAppConfigStrValueByK("fee_wallet_address_list")
		addresses := strings.Split(*feeAddressValue, ",")
		var feeAddresses []string
		for _, address := range addresses {
			if address == "" {
				continue
			}
			feeAddresses = append(feeAddresses, address)
		}
		// 遍历获取需要查询的block信息
		for i := startI; i < endI; i++ {
			// rpc获取block信息
			rpcBlock, err := ethclient.RpcBlockByNum(context.Background(), i)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			// 接收地址列表
			var toAddresses []string
			// map[接收地址] => []交易信息
			toAddressTxMap := make(map[string][]*types.Transaction)
			// 遍历block中的tx
			for _, rpcTx := range rpcBlock.Transactions() {
				// 转账数额大于0 and 不是创建合约交易
				if rpcTx.Value().Int64() > 0 && rpcTx.To() != nil {
					msg, err := rpcTx.AsMessage(types.NewEIP155Signer(rpcTx.ChainId()))
					if err != nil {
						log.Panicf("AsMessage err: [%T] %s", err, err.Error())
						return
					}
					if IsStringInSlice(feeAddresses, AddressBytesToStr(msg.From())) {
						// 如果打币地址在手续费热钱包地址则不处理
						continue
					}
					toAddress := AddressBytesToStr(*(rpcTx.To()))
					toAddressTxMap[toAddress] = append(toAddressTxMap[toAddress], rpcTx)
					if !IsStringInSlice(toAddresses, toAddress) {
						toAddresses = append(toAddresses, toAddress)
					}
				}
			}
			//从db中查询这些地址是否是冲币地址中的地址
			dbAddressRows, err := model.SQLSelectTAddressKeyColByAddress(toAddresses)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			// 待插入数据
			var dbTxRows []*model.TTx
			// map[接收地址] => 系统id
			addressSystemMap := make(map[string]int64)
			for _, dbAddressRow := range dbAddressRows {
				addressSystemMap[dbAddressRow.UserAddress] = int64(dbAddressRow.UseTag)
			}
			// 时间
			now := time.Now().Unix()
			// 遍历数据库中有交易的地址
			for _, dbAddressRow := range dbAddressRows {
				if dbAddressRow.UseTag < 0 {
					continue
				}
				// 获取地址对应的交易列表
				txes := toAddressTxMap[dbAddressRow.UserAddress]
				for _, tx := range txes {
					msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
					if err != nil {
						log.Panicf("AsMessage err: [%T] %s", err, err.Error())
						return
					}
					fromAddress := AddressBytesToStr(msg.From())
					toAddress := AddressBytesToStr(*(tx.To()))
					balanceReal, err := WeiBigIntToEthStr(tx.Value())
					if err != nil {
						log.Panicf("err: [%T] %s", err, err.Error())
						return
					}
					dbTxRows = append(dbTxRows, &model.TTx{
						UserId:       addressSystemMap[toAddress],
						SystemID:     util.RandString(12),
						TxID:         tx.Hash().String(),
						FromAddress:  fromAddress,
						ToAddress:    toAddress,
						BalanceReal:  balanceReal,
						CreateTime:   now,
						HandleStatus: hcommon.TxStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
						OrgStatus:    hcommon.TxOrgStatusInit,
						OrgMsg:       "",
						OrgTime:      now,
					})
				}
			}
			// 插入交易数据
			_, err = model.SQLCreateIgnoreManyTTx(dbTxRows)
			if err != nil {
				return
			}
			// 更新检查到的最新区块数
			err = model.SQLUpdateTAppStatusIntByKGreater(
				model.TAppStatusInt{
					K: "seek_num",
					V: i,
				},
			)
			if err != nil {
				return
			}
		}
	}

}

// CheckAddressOrg 零钱整理到冷钱包
func CheckAddressOrg() {
	// 获取冷钱包地址
	coldAddressValue := model.SQLGetTAppConfigStrValueByK("cold_wallet_address")
	coldAddress, err := StrToAddressBytes(*coldAddressValue)
	if err != nil {
		log.Panicf("eth organize cold address err: [%T] %s", err, err.Error())
		return
	}
	isComment := false
	//开启事务
	dbTx, err := model.GetDb().DB().BeginTx(context.Background(), nil)
	if err != nil {
		return
	}
	defer func() {
		if !isComment {
			_ = dbTx.Rollback()
		}
	}()

	// 获取待整理的交易列表
	txRows, err := model.SQLSelectTTxColByOrgForUpdate(hcommon.TxOrgStatusInit)
	if err != nil {
		log.Panicf("ttx transaction err: [%T] %s", err, err.Error())
		return
	}
	if len(txRows) <= 0 {
		// 没有要处理的信息
		return
	}
	// 获取gap price
	gasPriceValue := model.SQLGetTAppStatusIntValueByK("to_cold_gas_price")
	gasPrice := gasPriceValue
	gasLimit := int64(21000)
	feeValue := big.NewInt(gasLimit * *gasPrice)
	// chain id
	chainID, err := ethclient.RpcNetworkID(context.Background())
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	// 当前时间
	now := time.Now().Unix()
	// 将待整理地址按地址做归并处理
	type OrgInfo struct {
		RowIDs  []int64  // db t_tx.id
		Balance *big.Int // 金额
	}
	// addressMap map[地址] = []整理信息
	addressMap := make(map[string]*OrgInfo)
	// addresses 需要整理的地址列表
	var addresses []string
	for _, txRow := range txRows {
		info := addressMap[txRow.ToAddress]
		if info == nil {
			info = &OrgInfo{
				RowIDs:  []int64{},
				Balance: new(big.Int),
			}
			addressMap[txRow.ToAddress] = info
		}
		info.RowIDs = append(info.RowIDs, txRow.ID)
		txWei, err := EthStrToWeiBigInit(txRow.BalanceReal)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		info.Balance.Add(info.Balance, txWei)
		//将新的充币地址加入数组
		if !IsStringInSlice(addresses, txRow.ToAddress) {
			addresses = append(addresses, txRow.ToAddress)
		}
	}
	// 获取地址私钥
	addressPKMap, err := GetPKMapOfAddresses(addresses)
	if err != nil {
		log.Panicf("GetNonce err: [%T] %s", err, err.Error())
		return
	}
	for address, info := range addressMap {
		// 获取私钥
		privateKey, ok := addressPKMap[address]
		if !ok {
			log.Panicf("no key of: %s", address)
			continue
		}
		// 获取nonce值
		nonce, err := GetNonce(address)
		if err != nil {
			log.Panicf("GetNonce err: [%T] %s", err, err.Error())
			return
		}
		// 发送数量
		sendBalance := new(big.Int)
		sendBalance.Sub(info.Balance, feeValue)
		if sendBalance.Cmp(new(big.Int)) <= 0 {
			// 数额不足
			continue
		}
		sendBalanceReal, err := WeiBigIntToEthStr(sendBalance)
		if err != nil {
			log.Panicf("GetNonce err: [%T] %s", err, err.Error())
			return
		}
		// 生成tx
		var data []byte
		tx := types.NewTransaction(
			uint64(nonce),
			coldAddress,
			sendBalance,
			uint64(gasLimit),
			big.NewInt(*gasPrice),
			data,
		)
		// 签名
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
		if err != nil {
			log.Panicf("RpcNetworkID err: [%T] %s", err, err.Error())
			return
		}
		ts := types.Transactions{signedTx}
		rawTxBytes := ts.GetRlp(0)
		rawTxHex := hex.EncodeToString(rawTxBytes)
		txHash := strings.ToLower(signedTx.Hash().Hex())
		// 创建存入数据
		var sendRows []*model.TSend
		for rowIndex, rowID := range info.RowIDs {
			if rowIndex == 0 {
				// 只有第一条数据需要发送，其余数据为占位数据
				sendRows = append(sendRows, &model.TSend{
					RelatedType:  hcommon.SendRelationTypeTx,
					RelatedID:    rowID,
					TxID:         txHash,
					FromAddress:  address,
					ToAddress:    *coldAddressValue,
					BalanceReal:  sendBalanceReal,
					Gas:          gasLimit,
					GasPrice:     *gasPrice,
					Nonce:        nonce,
					Hex:          rawTxHex,
					CreateTime:   now,
					HandleStatus: hcommon.SendStatusInit,
					HandleMsg:    "",
					HandleTime:   now,
				})
			} else {
				// 占位数据
				sendRows = append(sendRows, &model.TSend{
					RelatedType:  hcommon.SendRelationTypeTx,
					RelatedID:    rowID,
					TxID:         txHash,
					FromAddress:  address,
					ToAddress:    *coldAddressValue,
					BalanceReal:  "0",
					Gas:          0,
					GasPrice:     0,
					Nonce:        -1,
					Hex:          "",
					CreateTime:   now,
					HandleStatus: hcommon.SendStatusInit,
					HandleMsg:    "",
					HandleTime:   now,
				})
			}
		}
		// 插入发送数据
		_, err = model.SQLCreateIgnoreManyTSend(sendRows, true)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 更改tx整理状态
		err = model.SQLUpdateTTxOrgStatusByIDs(
			info.RowIDs,
			&model.TTx{
				OrgStatus: hcommon.TxOrgStatusHex,
				OrgMsg:    "gen raw tx",
				OrgTime:   now,
			},
		)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 提交事物
		err = dbTx.Commit()
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		isComment = true
	}

}

// CheckRawTxSend 发送交易
func CheckRawTxSend() {
	// 获取待发送的数据
	sendRows, err := model.SQLSelectTSendColByStatus(hcommon.SendStatusInit)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	// 首先单独处理提领和代发，提取提币通知要使用的数据
	var withdrawIDs []int64
	for _, sendRow := range sendRows {
		switch sendRow.RelatedType {
		case hcommon.SendRelationTypeWithdraw:
			if !IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}

		case hcommon.SendRelationTypeSend:
			if !IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		}
	}
	withdrawMap, err := model.SQLGetWithdrawMap(withdrawIDs)
	var userIDs []int64
	for _, withdrawRow := range withdrawMap {
		if !mcommon.IsIntInSlice(userIDs, int64(withdrawRow.UserId)) {
			userIDs = append(userIDs, int64(withdrawRow.UserId))
		}
	}
	userMap, err := model.SQLGetUserMap(userIDs)
	// 执行发送
	var sendIDs []int64
	var txIDs []int64
	var erc20TxIDs []int64
	var erc20TxFeeIDs []int64
	withdrawIDs = []int64{}
	// 通知数据
	var notifyRows []*model.TUserNotify
	now := time.Now().Unix()
	var sendTxHashes []string
	onSendOk := func(sendRow *model.TSend) error {
		// 将发送成功和占位数据计入数组
		if !IsIntInSlice(sendIDs, sendRow.ID) {
			sendIDs = append(sendIDs, sendRow.ID)
		}
		switch sendRow.RelatedType {
		case hcommon.SendRelationTypeTx:
			if !IsIntInSlice(txIDs, sendRow.RelatedID) {
				txIDs = append(txIDs, sendRow.RelatedID)
			}
		case hcommon.SendRelationTypeWithdraw:
			if !IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		case hcommon.SendRelationTypeSend:
			if !IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		case hcommon.SendRelationTypeTxErc20:
			if !IsIntInSlice(erc20TxIDs, sendRow.RelatedID) {
				erc20TxIDs = append(erc20TxIDs, sendRow.RelatedID)
			}
		case hcommon.SendRelationTypeTxErc20Fee:
			if !IsIntInSlice(erc20TxFeeIDs, sendRow.RelatedID) {
				erc20TxFeeIDs = append(erc20TxFeeIDs, sendRow.RelatedID)
			}
		}
		// 如果是提币，创建通知信息
		if sendRow.RelatedType == hcommon.SendRelationTypeWithdraw ||sendRow.RelatedType == hcommon.SendRelationTypeSend {
			withdrawRow, ok := withdrawMap[sendRow.RelatedID]
			if !ok {
				log.Panicf("withdrawMap no: %d", sendRow.RelatedID)
				return nil
			}
			userRow, ok := userMap[int64(withdrawRow.UserId)]
			if !ok {
				mcommon.Log.Errorf("no userMap: %d", withdrawRow.UserId)
				return nil
			}
			nonce := GetUUIDStr()
			reqObj := gin.H{
				"tx_hash":     sendRow.TxID,
				"balance":     withdrawRow.BalanceReal,
				"real_name":   userRow.RealName,
				"system_id":   withdrawRow.SystemID,
				"address":     withdrawRow.ToAddress,
				"symbol":      withdrawRow.Symbol,
				"notify_type": hcommon.NotifyTypeWithdrawSend,
			}
			reqObj["sign"] = GetSign(userRow.UserName, reqObj)
			req, err := json.Marshal(reqObj)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return err
			}
			notifyRows = append(notifyRows, &model.TUserNotify{
				Nonce:        nonce,
				SystemID:     withdrawRow.SystemID,
				UserId:       userRow.ID,
				ItemType:     hcommon.SendRelationTypeWithdraw,
				ItemID:       withdrawRow.ID,
				NotifyType:   hcommon.NotifyTypeWithdrawSend,
				TokenSymbol:  withdrawRow.Symbol,
				Msg:          string(req),
				HandleStatus: hcommon.NotifyStatusInit,
				HandleMsg:    "",
				CreateTime:   now,
				UpdateTime:   now,
			})
			return nil
		}
		return err
	}
	for _, sendRow := range sendRows {
		// 发送数据中需要排除占位数据
		if sendRow.Hex != "" {
			rawTxBytes, err := hex.DecodeString(sendRow.Hex)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				continue
			}
			tx := new(types.Transaction)
			err = rlp.DecodeBytes(rawTxBytes, &tx)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				continue
			}
			err = ethclient.RpcSendTransaction(
				context.Background(),
				tx,
			)
			if err != nil {
				if !strings.Contains(err.Error(), "known transaction") {
					log.Panicf("err: [%T] %s", err, err.Error())
					continue
				}
			}
			sendTxHashes = append(sendTxHashes, sendRow.TxID)
			err = onSendOk(&sendRow)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
		} else if IsStringInSlice(sendTxHashes, sendRow.TxID) {
			err = onSendOk(&sendRow)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
		}
		// 插入通知
		_, err = model.SQLCreateIgnoreManyTProductNotify(notifyRows)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新提币状态
		err = model.SQLUpdateTWithdrawStatusByIDs(
			withdrawIDs,
			&model.TWithdraw{
				HandleStatus: hcommon.WithdrawStatusSend,
				HandleMsg:    "send",
				HandleTime:   now,
			},
		)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新eth零钱整理状态
		err = model.SQLUpdateTTxOrgStatusByIDs(
			txIDs,
			&model.TTx{
				OrgStatus: hcommon.TxOrgStatusSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20零钱整理状态
		err = model.SQLUpdateTTxErc20OrgStatusByIDs(
			erc20TxIDs,
			&model.TTxErc20{
				OrgStatus: hcommon.TxOrgStatusSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20手续费状态
		err = model.SQLUpdateTTxErc20OrgStatusByIDs(
			erc20TxFeeIDs,
			&model.TTxErc20{
				OrgStatus: hcommon.TxOrgStatusFeeSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新发送状态
		err = model.SQLUpdateTSendStatusByIDs(
			sendIDs,
			&model.TSend{
				HandleStatus: hcommon.SendStatusSend,
				HandleMsg:    "send",
				HandleTime:   now,
			},
		)
	}
}

// CheckRawTxConfirm 确认tx是否打包完成
func CheckRawTxConfirm() {
	sendRows, err := model.SQLSelectTSendColByStatus(hcommon.SendStatusSend)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	var withdrawIDs []int64
	for _, sendRow := range sendRows {
		if sendRow.RelatedType == hcommon.SendRelationTypeWithdraw {
			// 提币
			if !IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		}
	}
	withdrawMap, err := model.SQLGetWithdrawMap(withdrawIDs)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	var userIDs []int64
	for _, withdrawRow := range withdrawMap {
		if !mcommon.IsIntInSlice(userIDs, int64(withdrawRow.UserId)) {
			userIDs = append(userIDs, int64(withdrawRow.UserId))
		}
	}
	userMap, err := model.SQLGetUserMap(userIDs)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	now := time.Now().Unix()
	var notifyRows []*model.TUserNotify
	var sendIDs []int64
	var txIDs []int64
	var erc20TxIDs []int64
	var erc20TxFeeIDs []int64
	withdrawIDs = []int64{}
	var sendHashes []string
	for _, sendRow := range sendRows {
		if !IsStringInSlice(sendHashes, sendRow.TxID) {
			rpcTx, err := ethclient.RpcTransactionByHash(
				context.Background(),
				sendRow.TxID,
			)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				continue
			}
			if rpcTx == nil {
				continue
			}
			sendHashes = append(sendHashes, sendRow.TxID)
		}
		if sendRow.RelatedType == hcommon.SendRelationTypeWithdraw {
			// 提币
			withdrawRow, ok := withdrawMap[sendRow.RelatedID]
			if !ok {
				log.Panicf("no withdrawMap: %d", sendRow.RelatedID)
				return
			}
			userRow, ok := userMap[int64(withdrawRow.UserId)]
			if !ok {
				mcommon.Log.Errorf("no userMap: %d", withdrawRow.UserId)
				return
			}
			nonce := GetUUIDStr()
			reqObj := gin.H{
				"tx_hash":     sendRow.TxID,
				"balance":     withdrawRow.BalanceReal,
				"realName":    userRow.RealName,
				"system_id":   withdrawRow.SystemID,
				"address":     withdrawRow.ToAddress,
				"symbol":      withdrawRow.Symbol,
				"notify_type": hcommon.NotifyTypeWithdrawConfirm,
			}
			reqObj["sign"] = GetSign(userRow.UserName, reqObj)
			req, err := json.Marshal(reqObj)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			notifyRows = append(notifyRows, &model.TUserNotify{
				Nonce:        nonce,
				UserId:       userRow.ID,
				SystemID:     withdrawRow.SystemID,
				ItemType:     hcommon.SendRelationTypeWithdraw,
				ItemID:       withdrawRow.ID,
				NotifyType:   hcommon.NotifyTypeWithdrawConfirm,
				TokenSymbol:  withdrawRow.Symbol,
				Msg:          string(req),
				HandleStatus: hcommon.NotifyStatusInit,
				HandleMsg:    "",
				CreateTime:   now,
				UpdateTime:   now,
			})

		}
		// 将发送成功和占位数据计入数组
		if !IsIntInSlice(sendIDs, sendRow.ID) {
			sendIDs = append(sendIDs, sendRow.ID)
		}
		switch sendRow.RelatedType {
		case hcommon.SendRelationTypeTx:
			if !IsIntInSlice(txIDs, sendRow.RelatedID) {
				txIDs = append(txIDs, sendRow.RelatedID)
			}
		case hcommon.SendRelationTypeWithdraw:
			if !IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		case hcommon.SendRelationTypeTxErc20:
			if !IsIntInSlice(erc20TxIDs, sendRow.RelatedID) {
				erc20TxIDs = append(erc20TxIDs, sendRow.RelatedID)
			}
		case hcommon.SendRelationTypeTxErc20Fee:
			if !IsIntInSlice(erc20TxFeeIDs, sendRow.RelatedID) {
				erc20TxFeeIDs = append(erc20TxFeeIDs, sendRow.RelatedID)
			}
		}
	}
	// 添加通知信息
	_, err = model.SQLCreateIgnoreManyTProductNotify(notifyRows)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	// 更新提币状态
	err = model.SQLUpdateTWithdrawStatusByIDs(
		withdrawIDs,
		&model.TWithdraw{
			HandleStatus: hcommon.WithdrawStatusConfirm,
			HandleMsg:    "confirmed",
			HandleTime:   now,
		},
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	// 更新eth零钱整理状态
	err = model.SQLUpdateTTxOrgStatusByIDs(
		txIDs,
		&model.TTx{
			OrgStatus: hcommon.TxOrgStatusConfirm,
			OrgMsg:    "confirm",
			OrgTime:   now,
		},
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	// 更新erc20零钱整理状态
	err = model.SQLUpdateTTxErc20OrgStatusByIDs(
		erc20TxIDs,
		&model.TTxErc20{
			OrgStatus: hcommon.TxOrgStatusConfirm,
			OrgMsg:    "confirm",
			OrgTime:   now,
		},
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	// 更新erc20零钱整理eth手续费状态
	err = model.SQLUpdateTTxErc20OrgStatusByIDs(
		erc20TxFeeIDs,
		&model.TTxErc20{
			OrgStatus: hcommon.TxOrgStatusFeeConfirm,
			OrgMsg:    "eth fee confirmed",
			OrgTime:   now,
		},
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	// 更新发送状态
	err = model.SQLUpdateTSendStatusByIDs(
		sendIDs,
		&model.TSend{
			HandleStatus: hcommon.SendStatusConfirm,
			HandleMsg:    "confirmed",
			HandleTime:   now,
		},
	)
}

// CheckWithdraw 检测提现
func CheckWithdraw() {
	// 获取需要处理的提币数据
	withdrawRows, err := model.SQLSelectTWithdrawColByStatus(hcommon.WithdrawStatusInit)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	if len(withdrawRows) == 0 {
		// 没有要处理的提币
		return
	}
	// 获取热钱包地址
	hotAddressValue := model.SQLGetTAppConfigStrValueByK("hot_wallet_address")

	_, err = StrToAddressBytes(*hotAddressValue)
	if err != nil {
		log.Panicf("eth hot address err: [%T] %s", err, err.Error())
		return
	}
	// 获取私钥
	privateKey, err := model.GetPkOfAddress(*hotAddressValue)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	// 获取热钱包余额
	hotAddressBalance, err := ethclient.RpcBalanceAt(
		context.Background(),
		*hotAddressValue,
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	pendingBalanceRealStr := model.SQLGetTSendPendingBalanceReal(*hotAddressValue)
	pendingBalance, err := EthStrToWeiBigInit(pendingBalanceRealStr)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	hotAddressBalance.Sub(hotAddressBalance, pendingBalance)
	// 获取gap price
	gasPriceValue := model.SQLGetTAppStatusIntValueByK("to_user_gas_price")
	gasPrice := gasPriceValue
	gasLimit := int64(21000)
	feeValue := gasLimit * *gasPrice
	chainID, err := ethclient.RpcNetworkID(context.Background())
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	for _, withdrawRow := range withdrawRows {
		err = handleWithdraw(withdrawRow.ID, chainID, *hotAddressValue, privateKey, hotAddressBalance, gasLimit, *gasPrice, feeValue)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			continue
		}
	}

}

//提交eth提币
func handleWithdraw(withdrawID int64, chainID int64, hotAddress string, privateKey *ecdsa.PrivateKey, hotAddressBalance *big.Int, gasLimit, gasPrice, feeValue int64) error {
	isComment := false
	//开启事务
	dbTx, err := model.GetDb().DB().BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if !isComment {
			_ = dbTx.Rollback()
		}
	}()
	// 处理业务
	withdrawRow := model.SQLGetTWithdrawColForUpdate(
		withdrawID,
		hcommon.WithdrawStatusInit,
	)

	if withdrawRow == nil {
		return nil
	}
	balanceBigInt, err := EthStrToWeiBigInit(withdrawRow.BalanceReal)
	if err != nil {
		return err
	}
	hotAddressBalance.Sub(hotAddressBalance, balanceBigInt)
	hotAddressBalance.Sub(hotAddressBalance, big.NewInt(feeValue))
	if hotAddressBalance.Cmp(new(big.Int)) < 0 {
		log.Panicf("hot balance limit")
		hotAddressBalance.Add(hotAddressBalance, balanceBigInt)
		hotAddressBalance.Add(hotAddressBalance, big.NewInt(feeValue))
		return nil
	}
	// nonce
	nonce, err := GetNonce(
		hotAddress,
	)
	if err != nil {
		return err
	}
	// 创建交易
	var data []byte
	toAddress, err := StrToAddressBytes(withdrawRow.ToAddress)
	if err != nil {
		return err
	}
	tx := types.NewTransaction(
		uint64(nonce),
		toAddress,
		balanceBigInt,
		uint64(gasLimit),
		big.NewInt(gasPrice),
		data,
	)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return err
	}
	ts := types.Transactions{signedTx}
	rawTxBytes := ts.GetRlp(0)
	rawTxHex := hex.EncodeToString(rawTxBytes)
	txHash := strings.ToLower(signedTx.Hash().Hex())
	now := time.Now().Unix()
	err = model.SQLUpdateTWithdrawGenTx(
		&model.TWithdraw{
			ID:           withdrawID,
			TxHash:       txHash,
			HandleStatus: hcommon.WithdrawStatusHex,
			HandleMsg:    "gen tx hex",
			HandleTime:   now,
		},
	)
	if err != nil {
		return err
	}
	_, err = model.SQLCreateTSend(
		&model.TSend{
			RelatedType:  hcommon.SendRelationTypeWithdraw,
			RelatedID:    withdrawID,
			TxID:         txHash,
			FromAddress:  hotAddress,
			ToAddress:    withdrawRow.ToAddress,
			BalanceReal:  withdrawRow.BalanceReal,
			Gas:          gasLimit,
			GasPrice:     gasPrice,
			Nonce:        nonce,
			Hex:          rawTxHex,
			HandleStatus: hcommon.SendStatusInit,
			HandleMsg:    "init",
			HandleTime:   now,
		},
	)
	// 处理完成 提交事务
	err = dbTx.Commit()
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return err
	}
	isComment = true
	return nil
}

// CheckTxNotify 创建eth冲币通知
func CheckTxNotify() {
	txRows, err := model.SQLSelectTTxColByStatus(hcommon.TxStatusInit)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	var userIDs []int64
	for _, txRow := range txRows {
		if !mcommon.IsIntInSlice(userIDs, txRow.UserId) {
			userIDs = append(userIDs, txRow.UserId)
		}
	}
	userMap, err := model.SQLGetUserMap(userIDs)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	var notifyTxIDs []int64
	var notifyRows []*model.TUserNotify
	now := time.Now().Unix()
	for _, txRow := range txRows {
		userRow, ok := userMap[txRow.UserId]
		if !ok {
			mcommon.Log.Warnf("no userMap: %d", txRow.UserId)
			notifyTxIDs = append(notifyTxIDs, txRow.ID)
			continue
		}
		nonce := GetUUIDStr()
		reqObj := gin.H{
			"tx_hash":     txRow.TxID,
			"real_name":   userRow.RealName,
			"address":     txRow.ToAddress,
			"balance":     txRow.BalanceReal,
			"symbol":      CoinSymbol,
			"notify_type": hcommon.NotifyTypeTx,
		}
		reqObj["sign"] = GetSign(userRow.UserName, reqObj)
		req, err := json.Marshal(reqObj)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			continue
		}
		notifyRows = append(notifyRows, &model.TUserNotify{
			Nonce:        nonce,
			UserId:       userRow.ID,
			SystemID:     txRow.SystemID,
			ItemType:     hcommon.SendRelationTypeTx,
			ItemID:       txRow.ID,
			NotifyType:   hcommon.NotifyTypeTx,
			TokenSymbol:  CoinSymbol,
			Msg:          string(req),
			HandleStatus: hcommon.NotifyStatusInit,
			HandleMsg:    "",
			CreateTime:   now,
			UpdateTime:   now,
		})
		notifyTxIDs = append(notifyTxIDs, txRow.ID)
	}
	_, err1 := model.SQLCreateIgnoreManyTProductNotify(notifyRows)
	if err1 != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	err = model.SQLUpdateTTxStatusByIDs(
		notifyTxIDs,
		&model.TTx{
			HandleStatus: hcommon.TxStatusNotify,
			HandleMsg:    "notify",
			HandleTime:   now,
		},
	)
	if err != nil {
		return
	}

}

// CheckErc20BlockSeek 检测erc20到账
func CheckErc20BlockSeek() {
	// 获取配置 延迟确认数
	confirmValue := model.SQLGetTAppConfigIntValueByK("block_confirm_num")
	// 获取状态 当前处理完成的最新的block number
	seekValue := model.SQLGetTAppStatusIntValueByK("erc20_seek_num")
	// rpc 获取当前最新区块数
	rpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
	if err != nil {
		return
	}
	startI := *seekValue + 1
	endI := rpcBlockNum - *confirmValue + 1
	if startI < endI {
		// 读取abi
		type LogTransfer struct {
			From   string
			To     string
			Tokens *big.Int
		}
		contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取所有token
		var configTokenRowAddresses []string
		configTokenRowMap := make(map[string]*model.TAppConfigToken)
		configTokenRows, err := model.SQLSelectTAppConfigTokenColAll()
		if err != nil {
			return
		}
		for _, contractRow := range configTokenRows {
			configTokenRowAddresses = append(configTokenRowAddresses, contractRow.TokenAddress)
			configTokenRowMap[strings.ToLower(contractRow.TokenAddress)] = &contractRow
		}
		// 遍历获取需要查询的block信息
		for i := startI; i < endI; i++ {
			if len(configTokenRowAddresses) > 0 {
				// rpc获取block信息
				logs, err := ethclient.RpcFilterLogs(
					context.Background(),
					i,
					i,
					configTokenRowAddresses,
					contractAbi.Events["Transfer"],
				)
				if err != nil {
					log.Panicf("err: [%T] %s", err, err.Error())
					return
				}
				// 接收地址列表
				var toAddresses []string
				// map[接收地址] => []交易信息
				toAddressLogMap := make(map[string][]types.Log)
				for _, log := range logs {
					if log.Removed {
						continue
					}
					toAddress := AddressBytesToStr(common.HexToAddress(log.Topics[2].Hex()))
					if !IsStringInSlice(toAddresses, toAddress) {
						toAddresses = append(toAddresses, toAddress)
					}
					toAddressLogMap[toAddress] = append(toAddressLogMap[toAddress], log)
				}
				// 从db中查询这些地址是否是冲币地址中的地址
				dbAddressRows, err := model.SQLSelectTAddressKeyColByAddress(toAddresses)
				if err != nil {
					log.Panicf("err: [%T] %s", err, err.Error())
					return
				}
				// 时间
				now := time.Now().Unix()
				// 待添加数组
				var txErc20Rows []*model.TTxErc20
				// map[接收地址] => 系统id
				addressSystemMap := make(map[string]int64)
				for _, dbAddressRow := range dbAddressRows {
					addressSystemMap[dbAddressRow.UserAddress] = int64(dbAddressRow.UseTag)
				}
				// 遍历数据库中有交易的地址
				for _, dbAddressRow := range dbAddressRows {
					if dbAddressRow.UseTag < 0 {
						continue
					}
					// 获取地址对应的交易列表
					logs, ok := toAddressLogMap[dbAddressRow.UserAddress]
					if !ok {
						log.Panicf("toAddressLogMap no: %s", dbAddressRow.UserAddress)
						return
					}
					for _, log1 := range logs {
						var transferEvent LogTransfer
						err := contractAbi.Unpack(&transferEvent, "Transfer", log1.Data)
						if err != nil {
							log.Panicf("err: [%T] %s", err, err.Error())
							return
						}
						transferEvent.From = strings.ToLower(common.HexToAddress(log1.Topics[1].Hex()).Hex())
						transferEvent.To = strings.ToLower(common.HexToAddress(log1.Topics[2].Hex()).Hex())
						contractAddress := strings.ToLower(log1.Address.Hex())
						configTokenRow, ok := configTokenRowMap[contractAddress]
						if !ok {
							log.Panicf("no configTokenRowMap of: %s", contractAddress)
							return
						}
						rpcTxReceipt, err := ethclient.RpcTransactionReceipt(
							context.Background(),
							log1.TxHash.Hex(),
						)
						if err != nil {
							log.Panicf("err: [%T] %s", err, err.Error())
							return
						}
						if rpcTxReceipt.Status <= 0 {
							continue
						}
						rpcTx, err := ethclient.RpcTransactionByHash(
							context.Background(),
							log1.TxHash.Hex(),
						)
						if err != nil {
							log.Panicf("err: [%T] %s", err, err.Error())
							return
						}
						if strings.ToLower(rpcTx.To().Hex()) != contractAddress {
							// 合约地址和tx的to地址不匹配
							continue
						}
						// 检测input
						input, err := contractAbi.Pack(
							"transfer",
							common.HexToAddress(log1.Topics[2].Hex()),
							transferEvent.Tokens,
						)
						if err != nil {
							log.Panicf("err: [%T] %s", err, err.Error())
							return
						}
						if hexutil.Encode(input) != hexutil.Encode(rpcTx.Data()) {
							// input 不匹配
							continue
						}
						balanceReal, err := TokenWeiBigIntToEthStr(transferEvent.Tokens, configTokenRow.TokenDecimals)
						if err != nil {
							log.Panicf("err: [%T] %s", err, err.Error())
							return
						}
						// 放入待插入数组
						txErc20Rows = append(txErc20Rows, &model.TTxErc20{
							TokenID:      configTokenRow.ID,
							UserId:       addressSystemMap[transferEvent.To],
							SystemID:     util.RandString(12),
							TxID:         log1.TxHash.Hex(),
							FromAddress:  transferEvent.From,
							ToAddress:    transferEvent.To,
							BalanceReal:  balanceReal,
							CreateTime:   now,
							HandleStatus: hcommon.TxStatusInit,
							HandleMsg:    "",
							HandleTime:   now,
							OrgStatus:    hcommon.TxOrgStatusInit,
							OrgMsg:       "",
							OrgTime:      now,
						})
					}
				}
				_, err = model.SQLCreateIgnoreManyTTxErc20(txErc20Rows)
				if err != nil {
					log.Panicf("err: [%T] %s", err, err.Error())
					return
				}
			}
			// 更新检查到的最新区块数
			err = model.SQLUpdateTAppStatusIntByKGreater(
				model.TAppStatusInt{
					K: "erc20_seek_num",
					V: i,
				},
			)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
		}
	}
}

// CheckErc20TxNotify 创建erc20冲币通知
func CheckErc20TxNotify() {
	txRows, err := model.SQLSelectTTxErc20ColByStatus(hcommon.TxStatusInit)
	if err != nil {
		return
	}
	var userIds []int64
	var tokenIDs []int64
	for _, txRow := range txRows {
		if !IsIntInSlice(userIds, txRow.UserId) {
			userIds = append(userIds, txRow.UserId)
		}
		if !IsIntInSlice(tokenIDs, txRow.TokenID) {
			tokenIDs = append(tokenIDs, txRow.TokenID)
		}
	}
	tokenMap, err := model.SQLGetAppConfigTokenMap(tokenIDs)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	userMap, err := model.SQLGetUserMap(userIds)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	var notifyTxIDs []int64
	var notifyRows []*model.TUserNotify
	now := time.Now().Unix()
	for _, txRow := range txRows {
		tokenRow, ok := tokenMap[txRow.TokenID]
		if !ok {
			log.Panicf("tokenMap no: %d", txRow.TokenID)
			continue
		}
		userRow, ok := userMap[txRow.UserId]
		if !ok {
			mcommon.Log.Warnf("userMap no: %d", txRow.UserId)
			notifyTxIDs = append(notifyTxIDs, txRow.ID)
			continue
		}
		nonce := GetUUIDStr()
		reqObj := gin.H{
			"tx_hash":     txRow.TxID,
			"real_name":   userRow.RealName,
			"address":     txRow.ToAddress,
			"balance":     txRow.BalanceReal,
			"symbol":      tokenRow.TokenSymbol,
			"notify_type": hcommon.NotifyTypeTx,
		}
		reqObj["sign"] = GetSign(userRow.UserName, reqObj)
		req, err := json.Marshal(reqObj)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		notifyRows = append(notifyRows, &model.TUserNotify{
			Nonce:        nonce,
			UserId:       userRow.ID,
			SystemID:     txRow.SystemID,
			ItemType:     hcommon.SendRelationTypeTx,
			ItemID:       txRow.ID,
			NotifyType:   hcommon.NotifyTypeTx,
			TokenSymbol:  tokenRow.TokenSymbol,
			Msg:          string(req),
			HandleStatus: hcommon.NotifyStatusInit,
			HandleMsg:    "",
			CreateTime:   now,
			UpdateTime:   now,
		})
		notifyTxIDs = append(notifyTxIDs, txRow.ID)
	}
	_, err = model.SQLCreateIgnoreManyTProductNotify(notifyRows)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	_, err = model.SQLUpdateTTxErc20StatusByIDs(
		notifyTxIDs,
		model.TTxErc20{
			HandleStatus: hcommon.TxStatusNotify,
			HandleMsg:    "notify",
			HandleTime:   now,
		},
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
}

// CheckErc20TxOrg erc20零钱整理
func CheckErc20TxOrg() {
	// 计算转账token所需的手续费
	erc20GasUseValue := model.SQLGetTAppConfigIntValueByK("erc20_gas_use")
	gasPriceValue := model.SQLGetTAppStatusIntValueByK("to_cold_gas_price")
	erc20Fee := big.NewInt(*erc20GasUseValue * *gasPriceValue)
	ethGasUse := int64(21000)
	ethFee := big.NewInt(ethGasUse * *gasPriceValue)
	// chainID
	chainID, err := ethclient.RpcNetworkID(context.Background())
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	//开启事务
	isComment := false
	dbTx, err := model.GetDb().DB().BeginTx(context.Background(), nil)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	defer func() {
		if !isComment {
			_ = dbTx.Rollback()
		}
	}()
	// 查询需要处理的交易
	txRows, err := model.SQLSelectTTxErc20ColByOrgForUpdate(
		[]int64{hcommon.TxOrgStatusInit, hcommon.TxOrgStatusFeeConfirm},
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	if len(txRows) <= 0 {
		//没有需要整理的交易
		return
	}
	// 整理信息
	type StOrgInfo struct {
		TxIDs        []int64
		ToAddress    string
		TokenID      int64
		TokenBalance *big.Int
	}
	var tokenIDs []int64
	for _, txRow := range txRows {
		if !IsIntInSlice(tokenIDs, txRow.TokenID) {
			tokenIDs = append(tokenIDs, txRow.TokenID)
		}
	}
	tokenMap, err := model.SQLGetAppConfigTokenMap(tokenIDs)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	txMap := make(map[int64]*model.TTxErc20)
	// 地址eth余额
	addressEthBalanceMap := make(map[string]*big.Int)
	// 整理信息map
	orgMap := make(map[string]*StOrgInfo)
	// 整理地址
	var toAddresses []string
	for _, txRow := range txRows {
		tokenRow, ok := tokenMap[txRow.TokenID]
		if !ok {
			log.Panicf("no token of: %d", txRow.TokenID)
			return
		}
		// 转换为map
		txMap[txRow.ID] = &txRow
		// 读取eth余额
		_, ok = addressEthBalanceMap[txRow.ToAddress]
		if !ok {
			balance, err := ethclient.RpcBalanceAt(
				context.Background(),
				txRow.ToAddress,
			)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			addressEthBalanceMap[txRow.ToAddress] = balance
		}
		// 整理信息
		orgKey := fmt.Sprintf("%s-%d", txRow.ToAddress, txRow.TokenID)
		orgInfo, ok := orgMap[orgKey]
		if !ok {
			orgInfo = &StOrgInfo{
				TokenID:      txRow.TokenID,
				ToAddress:    txRow.ToAddress,
				TokenBalance: new(big.Int),
			}
			orgMap[orgKey] = orgInfo
		}
		orgInfo.TxIDs = append(orgInfo.TxIDs, txRow.ID)
		txBalance, err := TokenEthStrToWeiBigInit(txRow.BalanceReal, tokenRow.TokenDecimals)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		orgInfo.TokenBalance.Add(orgInfo.TokenBalance, txBalance)
		// 待查询id
		if !IsStringInSlice(toAddresses, txRow.ToAddress) {
			toAddresses = append(toAddresses, txRow.ToAddress)
		}
	}
	// 整理地址key
	addressPKMap, err := GetPKMapOfAddresses(toAddresses)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	// 需要手续费的整理信息
	now := time.Now().Unix()
	needEthFeeMap := make(map[string]*StOrgInfo)
	for k, orgInfo := range orgMap {
		toAddress := orgInfo.ToAddress
		// 计算eth费用
		addressEthBalanceMap[toAddress] = addressEthBalanceMap[toAddress].Sub(addressEthBalanceMap[toAddress], erc20Fee)
		if addressEthBalanceMap[toAddress].Cmp(new(big.Int)) < 0 {
			// eth手续费不足
			// 处理添加手续费
			needEthFeeMap[k] = orgInfo
			continue
		}
		tokenRow, ok := tokenMap[orgInfo.TokenID]
		if !ok {
			log.Panicf("no tokenMap: %d", orgInfo.TokenID)
			continue
		}
		orgMinBalance, err := TokenEthStrToWeiBigInit(tokenRow.OrgMinBalance, tokenRow.TokenDecimals)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			continue
		}
		if orgInfo.TokenBalance.Cmp(orgMinBalance) < 0 {
			log.Panicf("token balance < org min balance")
			continue
		}
		// 处理token转账
		privateKey, ok := addressPKMap[toAddress]
		if !ok {
			log.Panicf("addressMap no: %s", toAddress)
			continue
		}
		// 获取nonce值
		nonce, err := GetNonce(toAddress)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			continue
		}
		// 生成交易
		contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		input, err := contractAbi.Pack(
			"transfer",
			common.HexToAddress(tokenRow.ColdAddress),
			orgInfo.TokenBalance,
		)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		rpcTx := types.NewTransaction(
			uint64(nonce),
			common.HexToAddress(tokenRow.TokenAddress),
			big.NewInt(0),
			uint64(*erc20GasUseValue),
			big.NewInt(*gasPriceValue),
			input,
		)
		signedTx, err := types.SignTx(rpcTx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			continue
		}
		ts := types.Transactions{signedTx}
		rawTxBytes := ts.GetRlp(0)
		rawTxHex := hex.EncodeToString(rawTxBytes)
		txHash := strings.ToLower(signedTx.Hash().Hex())
		// 创建存入数据
		balanceReal, err := TokenWeiBigIntToEthStr(orgInfo.TokenBalance, tokenRow.TokenDecimals)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			continue
		}
		// 待插入数据
		var sendRows []*model.TSend
		for rowIndex, txID := range orgInfo.TxIDs {
			if rowIndex == 0 {
				sendRows = append(sendRows, &model.TSend{
					RelatedType:  hcommon.SendRelationTypeTxErc20,
					RelatedID:    txID,
					TokenID:      orgInfo.TokenID,
					TxID:         txHash,
					FromAddress:  toAddress,
					ToAddress:    tokenRow.ColdAddress,
					BalanceReal:  balanceReal,
					Gas:          *erc20GasUseValue,
					GasPrice:     *gasPriceValue,
					Nonce:        nonce,
					Hex:          rawTxHex,
					CreateTime:   now,
					HandleStatus: hcommon.SendStatusInit,
					HandleMsg:    "",
					HandleTime:   now,
				})
			} else {
				sendRows = append(sendRows, &model.TSend{
					RelatedType:  hcommon.SendRelationTypeTxErc20,
					RelatedID:    txID,
					TokenID:      orgInfo.TokenID,
					TxID:         txHash,
					FromAddress:  toAddress,
					ToAddress:    tokenRow.ColdAddress,
					BalanceReal:  "",
					Gas:          0,
					GasPrice:     0,
					Nonce:        -1,
					Hex:          "",
					CreateTime:   now,
					HandleStatus: hcommon.SendStatusInit,
					HandleMsg:    "",
					HandleTime:   now,
				})
			}
		}
		// 插入发送队列
		_, err = model.SQLCreateIgnoreManyTSend(sendRows, true)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新整理状态
		err = model.SQLUpdateTTxErc20OrgStatusByIDs(
			orgInfo.TxIDs,
			&model.TTxErc20{
				OrgStatus: hcommon.TxOrgStatusHex,
				OrgMsg:    "hex",
				OrgTime:   now,
			},
		)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
	}
	// 生成eth转账
	if len(needEthFeeMap) > 0 {
		// 获取热钱包地址
		feeAddressValue := model.SQLGetTAppConfigStrValueByK("fee_wallet_address")
		_, err = StrToAddressBytes(*feeAddressValue)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取私钥
		privateKey, err := model.GetPkOfAddress(*feeAddressValue)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		feeAddressBalance, err := ethclient.RpcBalanceAt(
			context.Background(),
			*feeAddressValue,
		)
		if err != nil {
			log.Panicf("RpcBalanceAt err: [%T] %s", err, err.Error())
			return
		}
		pendingBalanceReal := model.SQLGetTSendPendingBalanceReal(*feeAddressValue)
		pendingBalance, err := EthStrToWeiBigInit(pendingBalanceReal)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		feeAddressBalance.Sub(feeAddressBalance, pendingBalance)
		// 生成手续费交易
		for _, orgInfo := range needEthFeeMap {
			feeAddressBalance.Sub(feeAddressBalance, ethFee)
			feeAddressBalance.Sub(feeAddressBalance, erc20Fee)
			if feeAddressBalance.Cmp(new(big.Int)) < 0 {
				log.Panicf("eth fee balance limit")
				return
			}
			// nonce
			nonce, err := GetNonce(*feeAddressValue)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			// 创建交易
			var data []byte
			tx := types.NewTransaction(
				uint64(nonce),
				common.HexToAddress(orgInfo.ToAddress),
				erc20Fee,
				uint64(ethGasUse),
				big.NewInt(*gasPriceValue),
				data,
			)
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			ts := types.Transactions{signedTx}
			rawTxBytes := ts.GetRlp(0)
			rawTxHex := hex.EncodeToString(rawTxBytes)
			txHash := strings.ToLower(signedTx.Hash().Hex())
			now := time.Now().Unix()
			balanceReal, err := WeiBigIntToEthStr(erc20Fee)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			// 待插入数据
			var sendRows []*model.TSend
			for rowIndex, txID := range orgInfo.TxIDs {
				if rowIndex == 0 {
					sendRows = append(sendRows, &model.TSend{
						RelatedType:  hcommon.SendRelationTypeTxErc20Fee,
						RelatedID:    txID,
						TokenID:      0,
						TxID:         txHash,
						FromAddress:  *feeAddressValue,
						ToAddress:    orgInfo.ToAddress,
						BalanceReal:  balanceReal,
						Gas:          ethGasUse,
						GasPrice:     *gasPriceValue,
						Nonce:        nonce,
						Hex:          rawTxHex,
						CreateTime:   now,
						HandleStatus: hcommon.SendStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
					})
				} else {
					sendRows = append(sendRows, &model.TSend{
						RelatedType:  hcommon.SendRelationTypeTxErc20Fee,
						RelatedID:    txID,
						TokenID:      0,
						TxID:         txHash,
						FromAddress:  *feeAddressValue,
						ToAddress:    orgInfo.ToAddress,
						BalanceReal:  "",
						Gas:          0,
						GasPrice:     0,
						Nonce:        -1,
						Hex:          "",
						CreateTime:   now,
						HandleStatus: hcommon.SendStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
					})
				}
			}
			// 插入发送数据
			_, err = model.SQLCreateIgnoreManyTSend(sendRows, true)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			// 更新整理状态
			err = model.SQLUpdateTTxErc20OrgStatusByIDs(
				orgInfo.TxIDs,
				&model.TTxErc20{
					OrgStatus: hcommon.TxOrgStatusFeeHex,
					OrgMsg:    "fee hex",
					OrgTime:   now,
				},
			)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
		}
	}
	//提交事务
	err = dbTx.Commit()
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	isComment = true
}

//erc20提币
func CheckErc20Withdraw() {
	var tokenSymbols []string
	tokenMap := make(map[string]*model.TAppConfigToken)
	addressKeyMap := make(map[string]*ecdsa.PrivateKey)
	addressEthBalanceMap := make(map[string]*big.Int)
	addressTokenBalanceMap := make(map[string]*big.Int)
	tokenRows, err := model.SQLSelectTAppConfigTokenColAll()
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	for _, tokenRow := range tokenRows {
		tokenMap[tokenRow.TokenSymbol] = &tokenRow
		if !IsStringInSlice(tokenSymbols, tokenRow.TokenSymbol) {
			tokenSymbols = append(tokenSymbols, tokenRow.TokenSymbol)
		}
	}
	withdrawRows, err := model.SQLSelectTWithdrawColByStatus(hcommon.WithdrawStatusInit)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	if len(withdrawRows) == 0 {
		//没有提币请求
		return
	}
	for _, tokenRow := range tokenRows {
		tokenMap[tokenRow.TokenSymbol] = &tokenRow
		if !IsStringInSlice(tokenSymbols, tokenRow.TokenSymbol) {
			tokenSymbols = append(tokenSymbols, tokenRow.TokenSymbol)
		}
		// 获取私钥
		_, err = StrToAddressBytes(tokenRow.HotAddress)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		hotAddress := tokenRow.HotAddress
		_, ok := addressKeyMap[hotAddress]
		if !ok {
			// 获取私钥
			keyRow := model.SQLGetTAddressKeyColByAddress(hotAddress)
			if keyRow == nil {
				log.Panicf("no key of: %s", hotAddress)
				return
			}
			key := hcommon.AesDecrypt(keyRow.Pwd, fmt.Sprintf("%s", setting.AesConf.Key))
			if len(key) == 0 {
				log.Panicf("error key of: %s", hotAddress)
				return
			}
			if strings.HasPrefix(key, "0x") {
				key = key[2:]
			}
			privateKey, err := crypto.HexToECDSA(key)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			addressKeyMap[hotAddress] = privateKey
		}
		_, ok = addressEthBalanceMap[hotAddress]
		if !ok {
			hotAddressBalance, err := ethclient.RpcBalanceAt(
				context.Background(),
				hotAddress,
			)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			pendingBalanceReal := model.SQLGetTSendPendingBalanceReal(hotAddress)
			pendingBalance, err := EthStrToWeiBigInit(pendingBalanceReal)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			hotAddressBalance.Sub(hotAddressBalance, pendingBalance)
			addressEthBalanceMap[hotAddress] = hotAddressBalance
		}
		tokenBalanceKey := fmt.Sprintf("%s-%s", tokenRow.HotAddress, tokenRow.TokenSymbol)
		_, ok = addressTokenBalanceMap[tokenBalanceKey]
		if !ok {
			tokenBalance, err := ethclient.RpcTokenBalance(
				context.Background(),
				tokenRow.TokenAddress,
				tokenRow.HotAddress,
			)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				return
			}
			addressTokenBalanceMap[tokenBalanceKey] = tokenBalance
		}
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取gap price
		gasPriceValue := model.SQLGetTAppStatusIntValueByK("to_user_gas_price")
		gasPrice := gasPriceValue
		erc20GasUseValue := model.SQLGetTAppConfigIntValueByK("erc20_gas_use")
		gasLimit := erc20GasUseValue
		// eth fee
		feeValue := big.NewInt(*gasLimit * *gasPrice)
		chainID, err := ethclient.RpcNetworkID(context.Background())
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			return
		}
		for _, withdrawRow := range withdrawRows {
			err = handleErc20Withdraw(withdrawRow.ID, chainID, &tokenMap, &addressKeyMap, &addressEthBalanceMap, &addressTokenBalanceMap, *gasLimit, *gasPrice, feeValue)
			if err != nil {
				log.Panicf("err: [%T] %s", err, err.Error())
				continue
			}
		}
	}
}

//erc20提币
func handleErc20Withdraw(withdrawID int64, chainID int64, tokenMap *map[string]*model.TAppConfigToken, addressKeyMap *map[string]*ecdsa.PrivateKey, addressEthBalanceMap *map[string]*big.Int, addressTokenBalanceMap *map[string]*big.Int, gasLimit, gasPrice int64, feeValue *big.Int) error {
	isComment := false
	//开启事务
	dbTx, err := model.GetDb().DB().BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if !isComment {
			_ = dbTx.Rollback()
		}
	}()
	// 处理业务
	withdrawRow := model.SQLGetTWithdrawColForUpdate(
		withdrawID,
		hcommon.WithdrawStatusInit,
	)
	if withdrawRow == nil {
		return nil
	}
	tokenRow, ok := (*tokenMap)[withdrawRow.Symbol]
	if !ok {
		log.Panicf("no tokenMap: %s", withdrawRow.Symbol)
		return nil
	}
	hotAddress := tokenRow.HotAddress
	key, ok := (*addressKeyMap)[hotAddress]
	if !ok {
		log.Panicf("no addressKeyMap: %s", hotAddress)
		return nil
	}
	(*addressEthBalanceMap)[hotAddress] = (*addressEthBalanceMap)[hotAddress].Sub(
		(*addressEthBalanceMap)[hotAddress],
		feeValue,
	)
	if (*addressEthBalanceMap)[hotAddress].Cmp(new(big.Int)) < 0 {
		log.Panicf("%s eth limit", hotAddress)
		return nil
	}
	tokenBalanceKey := fmt.Sprintf("%s-%s", tokenRow.HotAddress, tokenRow.TokenSymbol)
	tokenBalance, err := TokenEthStrToWeiBigInit(withdrawRow.BalanceReal, tokenRow.TokenDecimals)
	if err != nil {
		return err
	}
	(*addressTokenBalanceMap)[tokenBalanceKey] = (*addressTokenBalanceMap)[tokenBalanceKey].Sub(
		(*addressTokenBalanceMap)[tokenBalanceKey],
		tokenBalance,
	)
	if (*addressTokenBalanceMap)[tokenBalanceKey].Cmp(new(big.Int)) < 0 {
		log.Panicf("%s token limit", tokenBalanceKey)
		return nil
	}
	// 获取nonce值
	nonce, err := GetNonce(hotAddress)
	if err != nil {
		return err
	}
	// 生成交易
	contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return err
	}
	input, err := contractAbi.Pack(
		"transfer",
		common.HexToAddress(withdrawRow.ToAddress),
		tokenBalance,
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return err
	}
	rpcTx := types.NewTransaction(
		uint64(nonce),
		common.HexToAddress(tokenRow.TokenAddress),
		big.NewInt(0),
		uint64(gasLimit),
		big.NewInt(gasPrice),
		input,
	)
	signedTx, err := types.SignTx(rpcTx, types.NewEIP155Signer(big.NewInt(chainID)), key)
	if err != nil {
		return err
	}
	ts := types.Transactions{signedTx}
	rawTxBytes := ts.GetRlp(0)
	rawTxHex := hex.EncodeToString(rawTxBytes)
	txHash := strings.ToLower(signedTx.Hash().Hex())
	now := time.Now().Unix()
	err = model.SQLUpdateTWithdrawGenTx(
		&model.TWithdraw{
			ID:           withdrawID,
			TxHash:       txHash,
			HandleStatus: hcommon.WithdrawStatusHex,
			HandleMsg:    "gen tx hex",
			HandleTime:   now,
		},
	)
	if err != nil {
		return err
	}
	_, err = model.SQLCreateTSend(
		&model.TSend{
			RelatedType:  hcommon.SendRelationTypeWithdraw,
			RelatedID:    withdrawID,
			TxID:         txHash,
			FromAddress:  hotAddress,
			ToAddress:    withdrawRow.ToAddress,
			BalanceReal:  withdrawRow.BalanceReal,
			Gas:          gasLimit,
			GasPrice:     gasPrice,
			Nonce:        nonce,
			Hex:          rawTxHex,
			HandleStatus: hcommon.SendStatusInit,
			HandleMsg:    "init",
			HandleTime:   now,
		},
	)
	if err != nil {
		return err
	}
	// 处理完成，提交事务
	err = dbTx.Commit()
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return err
	}
	isComment = true
	return nil
}

// 检测gas price
func CheckGasPrice() {
	// 获取最高单价
	maxValue := model.SQLGetTAppStatusIntValueByK("max_gas_price_eth")
	resp := GetGas()
	toUserGasPrice := resp.Fast * int64(math.Pow10(8))
	toColdGasPrice := resp.Average * int64(math.Pow10(8))
	if toUserGasPrice > *maxValue {
		toUserGasPrice = *maxValue
	}
	if toColdGasPrice > *maxValue {
		toColdGasPrice = *maxValue
	}
	err := model.SQLUpdateTAppStatusIntByK(
		&model.TAppStatusInt{
			K: "to_user_gas_price",
			V: toUserGasPrice,
		},
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	err = model.SQLUpdateTAppStatusIntByK(
		&model.TAppStatusInt{
			K: "to_cold_gas_price",
			V: toColdGasPrice,
		},
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
}
