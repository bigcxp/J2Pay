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
	"github.com/parnurzeal/gorequest"
	_ "github.com/parnurzeal/gorequest"
	"j2pay-server/ethclient"
	_ "j2pay-server/ethclient"
	"j2pay-server/hcommon"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/pkg/setting"
	"j2pay-server/pkg/util"
	"math"
	_ "math"
	"math/big"
	_ "math/big"
	"net/http"
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
			Symbol:      CoinSymbol,
			UserAddress: address,
			Pwd:         privateKeyStrEn,
			UseTag:      addr.UseTag,
			UserId:      addr.UserId,
			UsdtAmount:  0,
			EthAmount:   0,
			Status:      1,
			HandleStatus: addr.HandleStatus,
			CreateTime:  now,
			UpdateTime:  now,
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
func CheckAddressFree() ([]string, error) {
	// 当前时间
	now := time.Now().Unix()
	//地址
	var userAddresses []string
	// 获取配置 允许的最小剩余地址数
	appConfig := model.SQLGetTAppConfigIntValueByK("k = ?", "min_free_address")
	// 获取当前剩余可用地址数
	count := model.GetAddressCount("use_tag = ?", 0)
	// 如果数据库中剩余可用地址小于最小允许可用地址
	if count < appConfig.V {
		var rows []*model.Address
		// 遍历差值次数
		for i := int64(0); i < appConfig.V-count; i++ {
			address, privateKeyStrEn, err := genAddressAndAesKey()
			if err != nil {
				return nil, err
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
			return nil, err
		}
	}
	return userAddresses, nil
}

//将地址分配给商户
func ToMerchantAddress(addr request.AddressAdd)(err error)  {
	//从钱包地址随机获取
	address, err := model.GetFreAddress(addr.Num)
	if err != nil {
		return err
	}
	err = model.ToAddress(addr.UserId,addr.UseTag, address)
	return err
}

// CheckBlockSeek 检测到账
func CheckBlockSeek() {
	// 获取配置 延迟确认数
	confirmValue := model.SQLGetTAppConfigIntValueByK("k = ?", "block_confirm_num")

	// 获取状态 当前处理完成的最新的block number
	seek := model.SQLGetTAppStatusIntValueByK("k = ?", "seek_num")

	// rpc 获取当前最新区块数
	rpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	startI := seek.V + 1
	endI := rpcBlockNum - confirmValue.V + 1
	if startI < endI {
		// 手续费钱包列表
		feeAddressValue := model.SQLGetTAppConfigStrValueByK("k = ?", "fee_wallet_address_list")
		addresses := strings.Split(feeAddressValue.V, ",")
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
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
						hcommon.Log.Errorf("AsMessage err: [%T] %s", err, err.Error())
						return
					}
					if hcommon.IsStringInSlice(feeAddresses, AddressBytesToStr(msg.From())) {
						// 如果打币地址在手续费热钱包地址则不处理
						continue
					}
					toAddress := AddressBytesToStr(*(rpcTx.To()))
					toAddressTxMap[toAddress] = append(toAddressTxMap[toAddress], rpcTx)
					if !hcommon.IsStringInSlice(toAddresses, toAddress) {
						toAddresses = append(toAddresses, toAddress)
					}
				}
			}
			//从db中查询这些地址是否是冲币地址中的地址
			dbAddressRows, err := model.SQLSelectTAddressKeyColByAddress(toAddresses)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 待插入数据
			var dbTxRows []*model.TTx
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
						hcommon.Log.Errorf("AsMessage err: [%T] %s", err, err.Error())
						return
					}
					fromAddress := AddressBytesToStr(msg.From())
					toAddress := AddressBytesToStr(*(tx.To()))
					balanceReal, err := WeiBigIntToEthStr(tx.Value())
					if err != nil {
						hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
						return
					}
					dbTxRows = append(dbTxRows, &model.TTx{
						SystemID:     util.RandString(12),
						TxID:         tx.Hash().String(),
						FromAddress:  fromAddress,
						ToAddress:    toAddress,
						BalanceReal:  balanceReal,
						CreateTime:   now,
						HandleStatus: model.TxStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
						OrgStatus:    model.TxOrgStatusInit,
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
				&model.TAppStatusInt{
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
	coldAddressValue := model.SQLGetTAppConfigStrValueByK("k = ?", "cold_wallet_address")
	coldAddress, err := StrToAddressBytes(coldAddressValue.V)
	if err != nil {
		hcommon.Log.Errorf("eth organize cold address err: [%T] %s", err, err.Error())
		return
	}
	// 获取待整理的交易列表
	txRows := model.SQLSelectTTxColByOrgForUpdate("org_status = ?", model.TxOrgStatusInit)
	if len(txRows) <= 0 {
		// 没有要处理的信息
		return
	}
	// 获取gap price
	gasPriceValue := model.SQLGetTAppStatusIntValueByK("k = ?", "to_cold_gas_price")
	gasPrice := gasPriceValue.V
	gasLimit := int64(21000)
	feeValue := big.NewInt(gasLimit * gasPrice)
	// chain id
	chainID, err := ethclient.RpcNetworkID(context.Background())
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		info.Balance.Add(info.Balance, txWei)

		if !hcommon.IsStringInSlice(addresses, txRow.ToAddress) {
			addresses = append(addresses, txRow.ToAddress)
		}
	}
	// 获取地址私钥
	addressPKMap, err := GetPKMapOfAddresses(addresses)
	if err != nil {
		hcommon.Log.Errorf("GetNonce err: [%T] %s", err, err.Error())
		return
	}
	for address, info := range addressMap {
		// 获取私钥
		privateKey, ok := addressPKMap[address]
		if !ok {
			hcommon.Log.Errorf("no key of: %s", address)
			continue
		}
		// 获取nonce值
		nonce, err := GetNonce(address)
		if err != nil {
			hcommon.Log.Errorf("GetNonce err: [%T] %s", err, err.Error())
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
			hcommon.Log.Errorf("GetNonce err: [%T] %s", err, err.Error())
			return
		}
		// 生成tx
		var data []byte
		tx := types.NewTransaction(
			uint64(nonce),
			coldAddress,
			sendBalance,
			uint64(gasLimit),
			big.NewInt(gasPrice),
			data,
		)
		// 签名
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
		if err != nil {
			hcommon.Log.Warnf("RpcNetworkID err: [%T] %s", err, err.Error())
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
					RelatedType:  model.SendRelationTypeTx,
					RelatedID:    rowID,
					TxID:         txHash,
					FromAddress:  address,
					ToAddress:    coldAddressValue.V,
					BalanceReal:  sendBalanceReal,
					Gas:          gasLimit,
					GasPrice:     gasPrice,
					Nonce:        nonce,
					Hex:          rawTxHex,
					CreateTime:   now,
					HandleStatus: model.SendStatusInit,
					HandleMsg:    "",
					HandleTime:   now,
				})
			} else {
				// 占位数据
				sendRows = append(sendRows, &model.TSend{
					RelatedType:  model.SendRelationTypeTx,
					RelatedID:    rowID,
					TxID:         txHash,
					FromAddress:  address,
					ToAddress:    coldAddressValue.V,
					BalanceReal:  "0",
					Gas:          0,
					GasPrice:     0,
					Nonce:        -1,
					Hex:          "",
					CreateTime:   now,
					HandleStatus: model.SendStatusInit,
					HandleMsg:    "",
					HandleTime:   now,
				})
			}
		}
		// 插入发送数据
		_, err = model.SQLCreateIgnoreManyTSend(sendRows)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更改tx整理状态
		err = model.SQLUpdateTTxOrgStatusByIDs(
			info.RowIDs,
			&model.TTx{
				OrgStatus: model.TxOrgStatusHex,
				OrgMsg:    "gen raw tx",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	}
}

// CheckRawTxSend 发送交易
func CheckRawTxSend() {
	// 获取待发送的数据
	sendRows := model.SQLSelectTSendColByStatus("handle_status = ?", model.SendStatusInit)
	// 首先单独处理提币，提取提币通知要使用的数据
	var withdrawIDs []int64
	for _, sendRow := range sendRows {
		switch sendRow.RelatedType {
		case model.SendRelationTypeWithdraw:
			if !hcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		}
	}
	withdrawMap, err := model.SQLGetWithdrawMap(withdrawIDs)
	// 执行发送
	var sendIDs []int64
	var txIDs []int64
	var erc20TxIDs []int64
	var erc20TxFeeIDs []int64
	withdrawIDs = []int64{}
	// 通知数据
	var notifyRows []*model.TProductNotify
	now := time.Now().Unix()
	var sendTxHashes []string
	onSendOk := func(sendRow *model.TSend) error {
		// 将发送成功和占位数据计入数组
		if !hcommon.IsIntInSlice(sendIDs, sendRow.ID) {
			sendIDs = append(sendIDs, sendRow.ID)
		}
		switch sendRow.RelatedType {
		case model.SendRelationTypeTx:
			if !hcommon.IsIntInSlice(txIDs, sendRow.RelatedID) {
				txIDs = append(txIDs, sendRow.RelatedID)
			}
		case model.SendRelationTypeWithdraw:
			if !hcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		case model.SendRelationTypeTxErc20:
			if !hcommon.IsIntInSlice(erc20TxIDs, sendRow.RelatedID) {
				erc20TxIDs = append(erc20TxIDs, sendRow.RelatedID)
			}
		case model.SendRelationTypeTxErc20Fee:
			if !hcommon.IsIntInSlice(erc20TxFeeIDs, sendRow.RelatedID) {
				erc20TxFeeIDs = append(erc20TxFeeIDs, sendRow.RelatedID)
			}
		}
		// 如果是提币，创建通知信息
		if sendRow.RelatedType == model.SendRelationTypeWithdraw {
			withdrawRow, ok := withdrawMap[sendRow.RelatedID]
			if !ok {
				hcommon.Log.Errorf("withdrawMap no: %d", sendRow.RelatedID)
				return nil
			}
			nonce := hcommon.GetUUIDStr()
			reqObj := gin.H{
				"tx_hash":     sendRow.TxID,
				"balance":     withdrawRow.BalanceReal,
				"system_id":   withdrawRow.SystemID,
				"address":     withdrawRow.ToAddress,
				"symbol":      withdrawRow.Symbol,
				"notify_type": model.NotifyTypeWithdrawSend,
			}
			reqObj["sign"] = hcommon.GetSign("j2pay", reqObj)
			req, err := json.Marshal(reqObj)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return err
			}
			notifyRows = append(notifyRows, &model.TProductNotify{
				Nonce:        nonce,
				SystemID:     withdrawRow.SystemID,
				ItemType:     model.SendRelationTypeWithdraw,
				ItemID:       withdrawRow.ID,
				NotifyType:   model.NotifyTypeWithdrawSend,
				TokenSymbol:  withdrawRow.Symbol,
				Msg:          string(req),
				HandleStatus: model.NotifyStatusInit,
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
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			tx := new(types.Transaction)
			err = rlp.DecodeBytes(rawTxBytes, &tx)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			err = ethclient.RpcSendTransaction(
				context.Background(),
				tx,
			)
			if err != nil {
				if !strings.Contains(err.Error(), "known transaction") {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					continue
				}
			}
			sendTxHashes = append(sendTxHashes, sendRow.TxID)
			err = onSendOk(sendRow)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		} else if hcommon.IsStringInSlice(sendTxHashes, sendRow.TxID) {
			err = onSendOk(sendRow)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
		// 插入通知
		_, err = model.SQLCreateIgnoreManyTProductNotify(notifyRows)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新提币状态
		err = model.SQLUpdateTWithdrawStatusByIDs(
			withdrawIDs,
			&model.TWithdraw{
				HandleStatus: model.WithdrawStatusSend,
				HandleMsg:    "send",
				HandleTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新eth零钱整理状态
		err = model.SQLUpdateTTxOrgStatusByIDs(
			txIDs,
			&model.TTx{
				OrgStatus: model.TxOrgStatusSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20零钱整理状态
		err = model.SQLUpdateTTxErc20OrgStatusByIDs(
			erc20TxIDs,
			&model.TTxErc20{
				OrgStatus: model.TxOrgStatusSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20手续费状态
		err = model.SQLUpdateTTxErc20OrgStatusByIDs(

			erc20TxFeeIDs,
			&model.TTxErc20{
				OrgStatus: model.TxOrgStatusFeeSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新发送状态
		err = model.SQLUpdateTSendStatusByIDs(
			sendIDs,
			&model.TSend{
				HandleStatus: model.SendStatusSend,
				HandleMsg:    "send",
				HandleTime:   now,
			},
		)
	}
}

// CheckRawTxConfirm 确认tx是否打包完成
func CheckRawTxConfirm() {
	sendRows := model.SQLSelectTSendColByStatus("handle_status = ?", model.SendStatusSend)
	var withdrawIDs []int64
	for _, sendRow := range sendRows {
		if sendRow.RelatedType == model.SendRelationTypeWithdraw {
			// 提币
			if !hcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		}
	}
	withdrawMap, err := model.SQLGetWithdrawMap(withdrawIDs)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	now := time.Now().Unix()
	var notifyRows []*model.TProductNotify
	var sendIDs []int64
	var txIDs []int64
	var erc20TxIDs []int64
	var erc20TxFeeIDs []int64
	withdrawIDs = []int64{}
	var sendHashes []string
	for _, sendRow := range sendRows {
		if !hcommon.IsStringInSlice(sendHashes, sendRow.TxID) {
			rpcTx, err := ethclient.RpcTransactionByHash(
				context.Background(),
				sendRow.TxID,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			if rpcTx == nil {
				continue
			}
			sendHashes = append(sendHashes, sendRow.TxID)
		}
		if sendRow.RelatedType == model.SendRelationTypeWithdraw {
			// 提币
			withdrawRow, ok := withdrawMap[sendRow.RelatedID]
			if !ok {
				hcommon.Log.Errorf("no withdrawMap: %d", sendRow.RelatedID)
				return
			}
			nonce := hcommon.GetUUIDStr()
			reqObj := gin.H{
				"tx_hash":     sendRow.TxID,
				"balance":     withdrawRow.BalanceReal,
				"system_id":   withdrawRow.SystemID,
				"address":     withdrawRow.ToAddress,
				"symbol":      withdrawRow.Symbol,
				"notify_type": model.NotifyTypeWithdrawConfirm,
			}
			reqObj["sign"] = hcommon.GetSign("j2pay", reqObj)
			req, err := json.Marshal(reqObj)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			notifyRows = append(notifyRows, &model.TProductNotify{
				Nonce:        nonce,
				SystemID:     withdrawRow.SystemID,
				ItemType:     model.SendRelationTypeWithdraw,
				ItemID:       withdrawRow.ID,
				NotifyType:   model.NotifyTypeWithdrawConfirm,
				TokenSymbol:  withdrawRow.Symbol,
				Msg:          string(req),
				HandleStatus: model.NotifyStatusInit,
				HandleMsg:    "",
				CreateTime:   now,
				UpdateTime:   now,
			})

		}
		// 将发送成功和占位数据计入数组
		if !hcommon.IsIntInSlice(sendIDs, sendRow.ID) {
			sendIDs = append(sendIDs, sendRow.ID)
		}
		switch sendRow.RelatedType {
		case model.SendRelationTypeTx:
			if !hcommon.IsIntInSlice(txIDs, sendRow.RelatedID) {
				txIDs = append(txIDs, sendRow.RelatedID)
			}
		case model.SendRelationTypeWithdraw:
			if !hcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedID) {
				withdrawIDs = append(withdrawIDs, sendRow.RelatedID)
			}
		case model.SendRelationTypeTxErc20:
			if !hcommon.IsIntInSlice(erc20TxIDs, sendRow.RelatedID) {
				erc20TxIDs = append(erc20TxIDs, sendRow.RelatedID)
			}
		case model.SendRelationTypeTxErc20Fee:
			if !hcommon.IsIntInSlice(erc20TxFeeIDs, sendRow.RelatedID) {
				erc20TxFeeIDs = append(erc20TxFeeIDs, sendRow.RelatedID)
			}
		}
	}
	// 添加通知信息
	_, err = model.SQLCreateIgnoreManyTProductNotify(notifyRows)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	// 更新提币状态
	err = model.SQLUpdateTWithdrawStatusByIDs(
		withdrawIDs,
		&model.TWithdraw{
			HandleStatus: model.WithdrawStatusConfirm,
			HandleMsg:    "confirmed",
			HandleTime:   now,
		},
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	// 更新eth零钱整理状态
	err = model.SQLUpdateTTxOrgStatusByIDs(
		txIDs,
		&model.TTx{
			OrgStatus: model.TxOrgStatusConfirm,
			OrgMsg:    "confirm",
			OrgTime:   now,
		},
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	// 更新erc20零钱整理状态
	err = model.SQLUpdateTTxErc20OrgStatusByIDs(
		erc20TxIDs,
		&model.TTxErc20{
			OrgStatus: model.TxOrgStatusConfirm,
			OrgMsg:    "confirm",
			OrgTime:   now,
		},
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	// 更新erc20零钱整理eth手续费状态
	err = model.SQLUpdateTTxErc20OrgStatusByIDs(
		erc20TxFeeIDs,
		&model.TTxErc20{
			OrgStatus: model.TxOrgStatusFeeConfirm,
			OrgMsg:    "eth fee confirmed",
			OrgTime:   now,
		},
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	// 更新发送状态
	err = model.SQLUpdateTSendStatusByIDs(
		sendIDs,
		&model.TSend{
			HandleStatus: model.SendStatusConfirm,
			HandleMsg:    "confirmed",
			HandleTime:   now,
		},
	)
}

// CheckWithdraw 检测提现
func CheckWithdraw() {
	// 获取需要处理的提币数据
	withdrawRows, err := model.SQLSelectTWithdrawColByStatus(model.WithdrawStatusInit)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	if len(withdrawRows) == 0 {
		// 没有要处理的提币
		return
	}
	// 获取热钱包地址
	hotAddressValue := model.SQLGetTAppConfigStrValueByK("k = ?", "hot_wallet_address")

	_, err = StrToAddressBytes(hotAddressValue.V)
	if err != nil {
		hcommon.Log.Errorf("eth hot address err: [%T] %s", err, err.Error())
		return
	}
	// 获取私钥
	privateKey, err := model.GetPkOfAddress(hotAddressValue.V)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	// 获取热钱包余额
	hotAddressBalance, err := ethclient.RpcBalanceAt(
		context.Background(),
		hotAddressValue.V,
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	pendingBalanceRealStr, err := model.SQLGetTSendPendingBalanceReal(hotAddressValue.V)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	pendingBalance, err := EthStrToWeiBigInit(pendingBalanceRealStr)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	hotAddressBalance.Sub(hotAddressBalance, pendingBalance)
	// 获取gap price
	gasPriceValue := model.SQLGetTAppStatusIntValueByK("k = ?", "to_user_gas_price")

	gasPrice := gasPriceValue.V
	gasLimit := int64(21000)
	feeValue := gasLimit * gasPrice
	chainID, err := ethclient.RpcNetworkID(context.Background())
	if err != nil {
		hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
		return
	}
	for _, withdrawRow := range withdrawRows {
		err = handleWithdraw(withdrawRow.ID, chainID, hotAddressValue.V, privateKey, hotAddressBalance, gasLimit, gasPrice, feeValue)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			continue
		}
	}
}

//提交提币
func handleWithdraw(withdrawID int64, chainID int64, hotAddress string, privateKey *ecdsa.PrivateKey, hotAddressBalance *big.Int, gasLimit, gasPrice, feeValue int64) error {

	// 处理业务
	withdrawRow, err := model.SQLGetTWithdrawColForUpdate(
		withdrawID,
		model.WithdrawStatusInit,
	)
	if err != nil {
		return err
	}
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
		hcommon.Log.Errorf("hot balance limit")
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
			HandleStatus: model.WithdrawStatusHex,
			HandleMsg:    "gen tx hex",
			HandleTime:   now,
		},
	)
	if err != nil {
		return err
	}
	_, err = model.SQLCreateTSend(
		&model.TSend{
			RelatedType:  model.SendRelationTypeWithdraw,
			RelatedID:    withdrawID,
			TxID:         txHash,
			FromAddress:  hotAddress,
			ToAddress:    withdrawRow.ToAddress,
			BalanceReal:  withdrawRow.BalanceReal,
			Gas:          gasLimit,
			GasPrice:     gasPrice,
			Nonce:        nonce,
			Hex:          rawTxHex,
			HandleStatus: model.SendStatusInit,
			HandleMsg:    "init",
			HandleTime:   now,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// CheckTxNotify 创建eth冲币通知
func CheckTxNotify() {
	txRows := model.SQLSelectTTxColByOrgForUpdate("hand_status", model.TxStatusInit)

	var notifyTxIDs []int64
	var notifyRows []*model.TProductNotify
	now := time.Now().Unix()
	for _, txRow := range txRows {
		nonce := hcommon.GetUUIDStr()
		reqObj := gin.H{
			"tx_hash":     txRow.TxID,
			"system_id":   txRow.SystemID,
			"address":     txRow.ToAddress,
			"balance":     txRow.BalanceReal,
			"symbol":      CoinSymbol,
			"notify_type": model.NotifyTypeTx,
		}
		reqObj["sign"] = hcommon.GetSign("j2pay", reqObj)
		req, err := json.Marshal(reqObj)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			continue
		}
		notifyRows = append(notifyRows, &model.TProductNotify{
			Nonce:        nonce,
			SystemID:     txRow.SystemID,
			ItemType:     model.SendRelationTypeTx,
			ItemID:       txRow.ID,
			NotifyType:   model.NotifyTypeTx,
			TokenSymbol:  CoinSymbol,
			Msg:          string(req),
			HandleStatus: model.NotifyStatusInit,
			HandleMsg:    "",
			CreateTime:   now,
			UpdateTime:   now,
		})
		notifyTxIDs = append(notifyTxIDs, txRow.ID)
	}
	_, err := model.SQLCreateIgnoreManyTProductNotify(notifyRows)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}

	err = model.SQLUpdateTTxStatusByIDs(
		notifyTxIDs,
		&model.TTx{
			HandleStatus: model.TxStatusNotify,
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
	confirmValue := model.SQLGetTAppConfigIntValueByK("k = ? ", "block_confirm_num")

	// 获取状态 当前处理完成的最新的block number
	seekValue := model.SQLGetTAppStatusIntValueByK("k = ?", "erc20_seek_num")
	// rpc 获取当前最新区块数
	rpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
	if err != nil {
		return
	}
	startI := seekValue.V + 1
	endI := rpcBlockNum - confirmValue.V + 1
	if startI < endI {
		// 读取abi
		type LogTransfer struct {
			From   string
			To     string
			Tokens *big.Int
		}
		contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
			configTokenRowMap[contractRow.TokenAddress] = contractRow
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
					hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
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
					if !hcommon.IsStringInSlice(toAddresses, toAddress) {
						toAddresses = append(toAddresses, toAddress)
					}
					toAddressLogMap[toAddress] = append(toAddressLogMap[toAddress], log)
				}
				// 从db中查询这些地址是否是冲币地址中的地址
				dbAddressRows, err := model.SQLSelectTAddressKeyColByAddress(toAddresses)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// map[接收地址] => 系统id
				addressSystemMap := make(map[string]string)
				for _, dbAddressRow := range dbAddressRows {
					addressSystemMap[dbAddressRow.UserAddress] = string(dbAddressRow.UseTag)
				}
				// 时间
				now := time.Now().Unix()
				// 待添加数组
				var txErc20Rows []*model.TTxErc20
				// 遍历数据库中有交易的地址
				for _, dbAddressRow := range dbAddressRows {
					if dbAddressRow.UseTag < 0 {
						continue
					}
					// 获取地址对应的交易列表
					logs, ok := toAddressLogMap[dbAddressRow.UserAddress]
					if !ok {
						hcommon.Log.Errorf("toAddressLogMap no: %s", dbAddressRow.UserAddress)
						return
					}
					for _, log := range logs {
						var transferEvent LogTransfer
						err := contractAbi.Unpack(&transferEvent, "Transfer", log.Data)
						if err != nil {
							hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
							return
						}
						transferEvent.From = strings.ToLower(common.HexToAddress(log.Topics[1].Hex()).Hex())
						transferEvent.To = strings.ToLower(common.HexToAddress(log.Topics[2].Hex()).Hex())
						contractAddress := strings.ToLower(log.Address.Hex())
						configTokenRow, ok := configTokenRowMap[contractAddress]
						if !ok {
							hcommon.Log.Errorf("no configTokenRowMap of: %s", contractAddress)
							return
						}
						rpcTxReceipt, err := ethclient.RpcTransactionReceipt(
							context.Background(),
							log.TxHash.Hex(),
						)
						if err != nil {
							hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
							return
						}
						if rpcTxReceipt.Status <= 0 {
							continue
						}
						rpcTx, err := ethclient.RpcTransactionByHash(
							context.Background(),
							log.TxHash.Hex(),
						)
						if err != nil {
							hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
							return
						}
						if strings.ToLower(rpcTx.To().Hex()) != contractAddress {
							// 合约地址和tx的to地址不匹配
							continue
						}
						// 检测input
						input, err := contractAbi.Pack(
							"transfer",
							common.HexToAddress(log.Topics[2].Hex()),
							transferEvent.Tokens,
						)
						if err != nil {
							hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
							return
						}
						if hexutil.Encode(input) != hexutil.Encode(rpcTx.Data()) {
							// input 不匹配
							continue
						}
						balanceReal, err := TokenWeiBigIntToEthStr(transferEvent.Tokens, configTokenRow.TokenDecimals)
						if err != nil {
							hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
							return
						}
						// 放入待插入数组
						txErc20Rows = append(txErc20Rows, &model.TTxErc20{
							TokenID:      configTokenRow.ID,
							SystemID:     addressSystemMap[transferEvent.To],
							TxID:         log.TxHash.Hex(),
							FromAddress:  transferEvent.From,
							ToAddress:    transferEvent.To,
							BalanceReal:  balanceReal,
							CreateTime:   now,
							HandleStatus: model.TxStatusInit,
							HandleMsg:    "",
							HandleTime:   now,
							OrgStatus:    model.TxOrgStatusInit,
							OrgMsg:       "",
							OrgTime:      now,
						})
					}
				}
				_, err = model.SQLCreateIgnoreManyTTxErc20(txErc20Rows)
				if err != nil {
					hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
			}
			// 更新检查到的最新区块数
			err = model.SQLUpdateTAppStatusIntByKGreater(
				&model.TAppStatusInt{
					K: "erc20_seek_num",
					V: i,
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
	}
}

// CheckErc20TxNotify 创建erc20冲币通知
func CheckErc20TxNotify() {
	txRows, err := model.SQLSelectTTxErc20ColByStatus(model.TxStatusInit)
	if err != nil {
		return
	}
	var systemIds []string
	var tokenIDs []int64
	for _, txRow := range txRows {
		if !hcommon.IsStringInSlice(systemIds, txRow.SystemID) {
			systemIds = append(systemIds, txRow.SystemID)
		}
		if !hcommon.IsIntInSlice(tokenIDs, txRow.TokenID) {
			tokenIDs = append(tokenIDs, txRow.TokenID)
		}
	}
	tokenMap, err := model.SQLGetAppConfigTokenMap(tokenIDs)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	var notifyTxIDs []int64
	var notifyRows []*model.TProductNotify
	now := time.Now().Unix()
	for _, txRow := range txRows {
		tokenRow, ok := tokenMap[txRow.TokenID]
		if !ok {
			hcommon.Log.Errorf("tokenMap no: %d", txRow.TokenID)
			continue
		}
		nonce := hcommon.GetUUIDStr()
		reqObj := gin.H{
			"tx_hash":     txRow.TxID,
			"system_id":   txRow.SystemID,
			"address":     txRow.ToAddress,
			"balance":     txRow.BalanceReal,
			"symbol":      tokenRow.TokenSymbol,
			"notify_type": model.NotifyTypeTx,
		}
		reqObj["sign"] = hcommon.GetSign("j2pay", reqObj)
		req, err := json.Marshal(reqObj)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		notifyRows = append(notifyRows, &model.TProductNotify{
			Nonce:        nonce,
			SystemID:     txRow.SystemID,
			ItemType:     model.SendRelationTypeTx,
			ItemID:       txRow.ID,
			NotifyType:   model.NotifyTypeTx,
			TokenSymbol:  tokenRow.TokenSymbol,
			Msg:          string(req),
			HandleStatus: model.NotifyStatusInit,
			HandleMsg:    "",
			CreateTime:   now,
			UpdateTime:   now,
		})
		notifyTxIDs = append(notifyTxIDs, txRow.ID)
	}
	_, err = model.SQLCreateIgnoreManyTProductNotify(notifyRows)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	_, err = model.SQLUpdateTTxErc20StatusByIDs(
		notifyTxIDs,
		model.TTxErc20{
			HandleStatus: model.TxStatusNotify,
			HandleMsg:    "notify",
			HandleTime:   now,
		},
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}

}

// CheckErc20TxOrg erc20零钱整理
func CheckErc20TxOrg() {
	// 计算转账token所需的手续费
	erc20GasUseValue := model.SQLGetTAppConfigIntValueByK("k = ?", "erc20_gas_use")
	gasPriceValue := model.SQLGetTAppStatusIntValueByK("k = ?", "to_cold_gas_price")

	erc20Fee := big.NewInt(erc20GasUseValue.V * gasPriceValue.V)
	ethGasUse := int64(21000)
	ethFee := big.NewInt(ethGasUse * gasPriceValue.V)
	// chainID
	chainID, err := ethclient.RpcNetworkID(context.Background())
	if err != nil {
		hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
		return
	}
	// 查询需要处理的交易
	txRows, err := model.SQLSelectTTxErc20ColByOrgForUpdate(
		[]int64{model.TxOrgStatusInit, model.TxOrgStatusFeeConfirm},
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	if len(txRows) <= 0 {
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
		if !hcommon.IsIntInSlice(tokenIDs, txRow.TokenID) {
			tokenIDs = append(tokenIDs, txRow.TokenID)
		}
	}
	tokenMap, err := model.SQLGetAppConfigTokenMap(tokenIDs)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
			hcommon.Log.Errorf("no token of: %d", txRow.TokenID)
			return
		}
		// 转换为map
		txMap[txRow.ID] = txRow
		// 读取eth余额
		_, ok = addressEthBalanceMap[txRow.ToAddress]
		if !ok {
			balance, err := ethclient.RpcBalanceAt(
				context.Background(),
				txRow.ToAddress,
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		orgInfo.TokenBalance.Add(orgInfo.TokenBalance, txBalance)
		// 待查询id
		if !hcommon.IsStringInSlice(toAddresses, txRow.ToAddress) {
			toAddresses = append(toAddresses, txRow.ToAddress)
		}
	}
	// 整理地址key
	addressPKMap, err := GetPKMapOfAddresses(toAddresses)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
			hcommon.Log.Errorf("no tokenMap: %d", orgInfo.TokenID)
			continue
		}

		orgMinBalance, err := TokenEthStrToWeiBigInit(tokenRow.OrgMinBalance, tokenRow.TokenDecimals)
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			continue
		}
		if orgInfo.TokenBalance.Cmp(orgMinBalance) < 0 {
			hcommon.Log.Errorf("token balance < org min balance")
			continue
		}
		// 处理token转账
		privateKey, ok := addressPKMap[toAddress]
		if !ok {
			hcommon.Log.Errorf("addressMap no: %s", toAddress)
			continue
		}
		// 获取nonce值
		nonce, err := GetNonce(toAddress)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			continue
		}
		// 生成交易
		contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		input, err := contractAbi.Pack(
			"transfer",
			common.HexToAddress(tokenRow.ColdAddress),
			orgInfo.TokenBalance,
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		rpcTx := types.NewTransaction(
			uint64(nonce),
			common.HexToAddress(tokenRow.TokenAddress),
			big.NewInt(0),
			uint64(erc20GasUseValue.V),
			big.NewInt(gasPriceValue.V),
			input,
		)
		signedTx, err := types.SignTx(rpcTx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			continue
		}
		ts := types.Transactions{signedTx}
		rawTxBytes := ts.GetRlp(0)
		rawTxHex := hex.EncodeToString(rawTxBytes)
		txHash := strings.ToLower(signedTx.Hash().Hex())
		// 创建存入数据
		balanceReal, err := TokenWeiBigIntToEthStr(orgInfo.TokenBalance, tokenRow.TokenDecimals)
		if err != nil {
			hcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			continue
		}
		// 待插入数据
		var sendRows []*model.TSend
		for rowIndex, txID := range orgInfo.TxIDs {
			if rowIndex == 0 {
				sendRows = append(sendRows, &model.TSend{
					RelatedType:  model.SendRelationTypeTxErc20,
					RelatedID:    txID,
					TokenID:      orgInfo.TokenID,
					TxID:         txHash,
					FromAddress:  toAddress,
					ToAddress:    tokenRow.ColdAddress,
					BalanceReal:  balanceReal,
					Gas:          erc20GasUseValue.V,
					GasPrice:     gasPriceValue.V,
					Nonce:        nonce,
					Hex:          rawTxHex,
					CreateTime:   now,
					HandleStatus: model.SendStatusInit,
					HandleMsg:    "",
					HandleTime:   now,
				})
			} else {
				sendRows = append(sendRows, &model.TSend{
					RelatedType:  model.SendRelationTypeTxErc20,
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
					HandleStatus: model.SendStatusInit,
					HandleMsg:    "",
					HandleTime:   now,
				})
			}
		}
		// 插入发送队列
		_, err = model.SQLCreateIgnoreManyTSend(sendRows)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新整理状态
		err = model.SQLUpdateTTxErc20OrgStatusByIDs(
			orgInfo.TxIDs,
			&model.TTxErc20{
				OrgStatus: model.TxOrgStatusHex,
				OrgMsg:    "hex",
				OrgTime:   now,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	}
	// 生成eth转账
	if len(needEthFeeMap) > 0 {
		// 获取热钱包地址
		feeAddressValue := model.SQLGetTAppConfigStrValueByK("k = ?", "fee_wallet_address")
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		_, err = StrToAddressBytes(feeAddressValue.V)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取私钥
		// 获取私钥
		privateKey, err := model.GetPkOfAddress(feeAddressValue.V)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		feeAddressBalance, err := ethclient.RpcBalanceAt(
			context.Background(),
			feeAddressValue.V,
		)
		if err != nil {
			hcommon.Log.Errorf("RpcBalanceAt err: [%T] %s", err, err.Error())
			return
		}
		pendingBalanceReal, err := model.SQLGetTSendPendingBalanceReal(feeAddressValue.V)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		pendingBalance, err := EthStrToWeiBigInit(pendingBalanceReal)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		feeAddressBalance.Sub(feeAddressBalance, pendingBalance)
		// 生成手续费交易
		for _, orgInfo := range needEthFeeMap {
			feeAddressBalance.Sub(feeAddressBalance, ethFee)
			feeAddressBalance.Sub(feeAddressBalance, erc20Fee)
			if feeAddressBalance.Cmp(new(big.Int)) < 0 {
				hcommon.Log.Errorf("eth fee balance limit")
				return
			}
			// nonce
			nonce, err := GetNonce(feeAddressValue.V)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 创建交易
			var data []byte
			tx := types.NewTransaction(
				uint64(nonce),
				common.HexToAddress(orgInfo.ToAddress),
				erc20Fee,
				uint64(ethGasUse),
				big.NewInt(gasPriceValue.V),
				data,
			)
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			ts := types.Transactions{signedTx}
			rawTxBytes := ts.GetRlp(0)
			rawTxHex := hex.EncodeToString(rawTxBytes)
			txHash := strings.ToLower(signedTx.Hash().Hex())
			now := time.Now().Unix()
			balanceReal, err := WeiBigIntToEthStr(erc20Fee)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 待插入数据
			var sendRows []*model.TSend
			for rowIndex, txID := range orgInfo.TxIDs {
				if rowIndex == 0 {
					sendRows = append(sendRows, &model.TSend{
						RelatedType:  model.SendRelationTypeTxErc20Fee,
						RelatedID:    txID,
						TokenID:      0,
						TxID:         txHash,
						FromAddress:  feeAddressValue.V,
						ToAddress:    orgInfo.ToAddress,
						BalanceReal:  balanceReal,
						Gas:          ethGasUse,
						GasPrice:     gasPriceValue.V,
						Nonce:        nonce,
						Hex:          rawTxHex,
						CreateTime:   now,
						HandleStatus: model.SendStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
					})
				} else {
					sendRows = append(sendRows, &model.TSend{
						RelatedType:  model.SendRelationTypeTxErc20Fee,
						RelatedID:    txID,
						TokenID:      0,
						TxID:         txHash,
						FromAddress:  feeAddressValue.V,
						ToAddress:    orgInfo.ToAddress,
						BalanceReal:  "",
						Gas:          0,
						GasPrice:     0,
						Nonce:        -1,
						Hex:          "",
						CreateTime:   now,
						HandleStatus: model.SendStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
					})
				}
			}
			// 插入发送数据
			_, err = model.SQLCreateIgnoreManyTSend(sendRows)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 更新整理状态
			err = model.SQLUpdateTTxErc20OrgStatusByIDs(
				orgInfo.TxIDs,
				&model.TTxErc20{
					OrgStatus: model.TxOrgStatusFeeHex,
					OrgMsg:    "fee hex",
					OrgTime:   now,
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
	}
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
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
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	for _, tokenRow := range tokenRows {
		tokenMap[tokenRow.TokenSymbol] = tokenRow
		if !hcommon.IsStringInSlice(tokenSymbols, tokenRow.TokenSymbol) {
			tokenSymbols = append(tokenSymbols, tokenRow.TokenSymbol)
		}
		// 获取私钥
		_, err = StrToAddressBytes(tokenRow.HotAddress)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		hotAddress := tokenRow.HotAddress
		_, ok := addressKeyMap[hotAddress]
		if !ok {
			// 获取私钥
			keyRow, err := model.SQLGetTAddressKeyColByAddress(hotAddress)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			if keyRow == nil {
				hcommon.Log.Errorf("no key of: %s", hotAddress)
				return
			}
			key := hcommon.AesDecrypt(keyRow.Pwd, fmt.Sprintf("%s", setting.AesConf.Key))
			if len(key) == 0 {
				hcommon.Log.Errorf("error key of: %s", hotAddress)
				return
			}
			if strings.HasPrefix(key, "0x") {
				key = key[2:]
			}
			privateKey, err := crypto.HexToECDSA(key)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			pendingBalanceReal, err := model.SQLGetTSendPendingBalanceReal(hotAddress)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			pendingBalance, err := EthStrToWeiBigInit(pendingBalanceReal)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			addressTokenBalanceMap[tokenBalanceKey] = tokenBalance
		}
	}
	withdrawRows, err := model.SQLSelectTWithdrawColByStatus(model.WithdrawStatusInit)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	if len(withdrawRows) == 0 {
		return
	}
	// 获取gap price
	gasPriceValue := model.SQLGetTAppStatusIntValueByK("k = ? ", "to_user_gas_price")
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	gasPrice := gasPriceValue
	erc20GasUseValue := model.SQLGetTAppConfigIntValueByK("k = ?", "erc20_gas_use")
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	gasLimit := erc20GasUseValue
	// eth fee
	feeValue := big.NewInt(gasLimit.V * gasPrice.V)
	chainID, err := ethclient.RpcNetworkID(context.Background())
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	for _, withdrawRow := range withdrawRows {
		err = handleErc20Withdraw(withdrawRow.ID, chainID, &tokenMap, &addressKeyMap, &addressEthBalanceMap, &addressTokenBalanceMap, gasLimit.V, gasPrice.V, feeValue)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			continue
		}
	}
}

//提币
func handleErc20Withdraw(withdrawID int64, chainID int64, tokenMap *map[string]*model.TAppConfigToken, addressKeyMap *map[string]*ecdsa.PrivateKey, addressEthBalanceMap *map[string]*big.Int, addressTokenBalanceMap *map[string]*big.Int, gasLimit, gasPrice int64, feeValue *big.Int) error {
	// 处理业务
	withdrawRow, err := model.SQLGetTWithdrawColForUpdate(
		withdrawID,
		model.WithdrawStatusInit,
	)
	if err != nil {
		return err
	}
	if withdrawRow == nil {
		return nil
	}
	tokenRow, ok := (*tokenMap)[withdrawRow.Symbol]
	if !ok {
		hcommon.Log.Errorf("no tokenMap: %s", withdrawRow.Symbol)
		return nil
	}
	hotAddress := tokenRow.HotAddress
	key, ok := (*addressKeyMap)[hotAddress]
	if !ok {
		hcommon.Log.Errorf("no addressKeyMap: %s", hotAddress)
		return nil
	}
	(*addressEthBalanceMap)[hotAddress] = (*addressEthBalanceMap)[hotAddress].Sub(
		(*addressEthBalanceMap)[hotAddress],
		feeValue,
	)
	if (*addressEthBalanceMap)[hotAddress].Cmp(new(big.Int)) < 0 {
		hcommon.Log.Errorf("%s eth limit", hotAddress)
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
		hcommon.Log.Errorf("%s token limit", tokenBalanceKey)
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
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	input, err := contractAbi.Pack(
		"transfer",
		common.HexToAddress(withdrawRow.ToAddress),
		tokenBalance,
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
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
			HandleStatus: model.WithdrawStatusHex,
			HandleMsg:    "gen tx hex",
			HandleTime:   now,
		},
	)
	if err != nil {
		return err
	}
	_, err = model.SQLCreateTSend(
		&model.TSend{
			RelatedType:  model.SendRelationTypeWithdraw,
			RelatedID:    withdrawID,
			TxID:         txHash,
			FromAddress:  hotAddress,
			ToAddress:    withdrawRow.ToAddress,
			BalanceReal:  withdrawRow.BalanceReal,
			Gas:          gasLimit,
			GasPrice:     gasPrice,
			Nonce:        nonce,
			Hex:          rawTxHex,
			HandleStatus: model.SendStatusInit,
			HandleMsg:    "init",
			HandleTime:   now,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// 检测gas price
func CheckGasPrice() {
		type StRespGasPrice struct {
			Fast        int64   `json:"fast"`
			Fastest     int64   `json:"fastest"`
			SafeLow     int64   `json:"safeLow"`
			Average     int64   `json:"average"`
			BlockTime   float64 `json:"block_time"`
			BlockNum    int64   `json:"blockNum"`
			Speed       float64 `json:"speed"`
			SafeLowWait float64 `json:"safeLowWait"`
			AvgWait     float64 `json:"avgWait"`
			FastWait    float64 `json:"fastWait"`
			FastestWait float64 `json:"fastestWait"`
		}
		gresp, body, errs := gorequest.New().
			Get("https://ethgasstation.info/api/ethgasAPI.json").
			Timeout(time.Second * 120).
			End()
		if errs != nil {
			hcommon.Log.Errorf("err: [%T] %s", errs[0], errs[0].Error())
			return
		}
		if gresp.StatusCode != http.StatusOK {
			// 状态错误
			hcommon.Log.Errorf("req status error: %d", gresp.StatusCode)
			return
		}
		var resp StRespGasPrice
		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		toUserGasPrice := resp.Fast * int64(math.Pow10(8))
		toColdGasPrice := resp.Average * int64(math.Pow10(8))
		err = model.SQLUpdateTAppStatusIntByK(
			&model.TAppStatusInt{
				K: "to_user_gas_price",
				V: toUserGasPrice,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		 err = model.SQLUpdateTAppStatusIntByK(
			&model.TAppStatusInt{
				K: "to_cold_gas_price",
				V: toColdGasPrice,
			},
		)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
}
