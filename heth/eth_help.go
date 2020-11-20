package heth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"j2pay-server/myerr"

	//"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/parnurzeal/gorequest"
	"j2pay-server/ethclient"
	"j2pay-server/hcommon"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/pkg/setting"
	"j2pay-server/pkg/util"
	"log"
	"math"
	"math/big"
	"net/http"
	"strings"
	"time"
)

//eth 交易中用到的属性
type ETHHelp struct {
	//数据库获取的交易数据
	GasData *GasData
	//发送链上的GasLimit
	SendGasLimit int64
	//发送链上的GasPrice
	SendGasPrice int64
	//交易发送金额
	SendBalance float64
	//密钥
	PrivateKey string
	//同地址交易索引
	Nonce int64
	//区块确认数
	ChainID int64
	//交易发起地址，合约交易不需要本参数
	FromAddress string
	//交易目标地址
	ToAddress string
	//交易合约地址
	ContractAddress string
	//代币位数
	Places int64
	// 提交交易后，返回的hash
	TxHash string
	//作用未明，后面补上
	RawTxHex string
}

//eth 获取的gas相关数据
type GasData struct {
	//eth打包区块数
	EthSeekNum int64
	//erc20打包区块数
	Erc20SeekNum int64
	//链上平均GasPrice数，每笔交易的平均值，一般在冷热钱包转账时使用
	AvgGasPrice int64
	//链上快速交易GasPrice数，交易成功比较快的数值，一般在用户提领转发的时候使用
	FastGasPrice int64
	//eth交易需要有矿工挖矿完成，一笔交易需要矿工挖矿多次完成，每次扣取GasLimit，共计完成GasPrice这么多次
	//eth交易一般默认为21000gwei
	DefaultGasLimit int64
	// 通过计算获取的gasLimit
	FastGasLimit int64
	//系统规定的最大gaslimit值
	MaxGasPrice int64
}

//gaslimit默认值21000gwei
const DefaultGasLimit int64 = 21000

//全零地址
const HAX_0 = "0x0000000000000000000000000000000000000000"

//获取chainID 区块确认数，作用还未明白
func (*ETHHelp) GetchainID() (chainID int64, err error) {
	chainID, err = ethclient.RpcNetworkID(context.Background())
	return
}

//获取默认契费gas
//汽费单位分为
// wei		1   					wei
//kwei		1000					wei
//mwei		1000000					wei
//gwei		1000000000				wei
//szabo		1000000000000			wei
//finey		1000000000000000		wei
//ether		1000000000000000000		wei   1eth=1亿亿wei
func (*ETHHelp) GetGas() *GasData {
	var gasData = GasData{}
	// 获取gap price
	gasPriceModes, err := model.SQLGetTAppConfigInt()
	if err != nil {
		log.Print("未获取到eth的gas数据,[%T] %s", err, err.Error())
	}
	for _, v := range gasPriceModes {
		if v.K == "seek_num" {
			gasData.EthSeekNum = v.V
		} else if v.K == "to_cold_gas_price" {
			gasData.AvgGasPrice = v.V
		} else if v.K == "erc20_seek_num" {
			gasData.Erc20SeekNum = v.V
		} else if v.K == "to_user_gas_price" {
			gasData.FastGasPrice = v.V
		} else if v.K == "default_gas_limit" {
			gasData.DefaultGasLimit = v.V
		} else if v.K == "max_gas_limit" {
			gasData.MaxGasPrice = v.V
		}
	}
	if gasData.FastGasLimit == 0 {
		gasData.FastGasLimit = DefaultGasLimit * 3
	}
	return &gasData
}

//通过交易发起地址获取nonce,
//return nonce交易编号每次加一
func (*ETHHelp) GetNONCE(fromAddress string) (nonce int64, err error) {
	// 通过rpc获取nonce交易数量，例如nonce=5 count=6
	rpcNonceCount, err := ethclient.RpcNonceAt(
		context.Background(),
		fromAddress,
	)
	if nil != err {
		return -1, err
	}
	rpcNonce := rpcNonceCount - 1
	// 获取db nonce
	nonceRecording := model.NonceRecording{}
	maxAndMin, err := nonceRecording.SQLGetTSendMaxAndMinNonceByFrom(fromAddress, rpcNonce)
	if err != nil {
		log.Print("GetNonce err:  %s", err, err.Error())
		return -1, err
	}
	if maxAndMin.MinNonce < 0 && rpcNonceCount <= 0 {
		//说明该地址是第一次交易，nooce起始就应该是0
		return 0, nil
	}
	if maxAndMin.MaxNonce <= rpcNonce {
		//说明链上的nonce大于数据库的nonce,中途有第三方提交交易，这个时候使用rpcNonce
		return rpcNonce + 1, nil
	} else if maxAndMin.MinNonce > rpcNonce+1 {
		//说明链上的nonce与数据库的nonce有一个以上数字空缺，需要补上
		return rpcNonce + 1, nil
	} else if maxAndMin.MaxNonce > rpcNonce {
		//说明我们提交的交易，有一个以上被卡住了，这里不适合继续交易，需要盘查错误，找出nonce
		return -1, errors.New("该账号的交易nonce被卡住,请管理员盘查!")
	}

	//TODO 其他情况暂时无法判断,就以链上的nonce为准
	return rpcNonce + 1, nil
}

//获取密钥，获取发起交易地址的密钥
//param 发起交易的地址
//return 地址对应的密钥
func (*ETHHelp) GetPrivateKey(fromAddress string) (privateKey string, err error) {
	if strings.TrimSpace(fromAddress) == "" {
		err = errors.New("fromAddress 不可以为空！")
		return
	}
	//通过热钱包地址查询pwd
	keyRow := model.SQLGetTAddressKeyColByAddress(fromAddress)
	if err != nil {
		log.Print("err: [%T] %s", err, err.Error())
		return
	}
	if keyRow == nil || keyRow.Pwd == "" {
		log.Print("未查询到地址对应的密码!")
		return
	}
	privateKey = hcommon.AesDecrypt(keyRow.Pwd, fmt.Sprintf("%s", setting.AesConf.Key))
	if len(privateKey) == 0 {
		log.Print("error key of: %s", fromAddress)
		return
	}
	if strings.HasPrefix(privateKey, "0x") {
		privateKey = privateKey[2:]
	}

	return
}

//获取eth上交易的各种数据
func (*ETHHelp) GetETHData() {
	// 获取gas price
	GetETHGasPrice()
}

// 获取gas price
func GetETHGasPrice() (err error) {
	gresp, body, errs := gorequest.New().
		Get("https://ethgasstation.info/api/ethgasAPI.json").
		Timeout(time.Second * 120).
		End()
	if errs != nil {
		log.Print("err: %s", errs[0], errs[0].Error())
	}
	if gresp.StatusCode != http.StatusOK {
		// 状态错误
		log.Print("req status error: %d", gresp.StatusCode)
	}
	var resp StRespGasPrice
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		log.Print("err: [%T] %s", err, err.Error())
	}
	//最快值gasPricew
	fastGasPrice := resp.Fast * int64(math.Pow10(8))
	//平均值gasPrice
	avgGasPrice := resp.Average * int64(math.Pow10(8))
	var appStatusModes = []*model.TAppStatusInt{
		&model.TAppStatusInt{K: "to_user_gas_price", V: fastGasPrice},
		&model.TAppStatusInt{K: "to_cold_gas_price", V: avgGasPrice},
	}
	err = model.SQLUpdateTAppStatusInt(appStatusModes)
	if err != nil {
		log.Print("err: [%T] %s", err, err.Error())
	}
	return
}

//创建地址和私钥
//return 地址
//return 加密密钥
func (e *ETHHelp) CreateAddressAndAesKey() (string, string, error) {
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

// CreateHotAddress 创建自用地址
func (e *ETHHelp) Createddress(addr request.AddressAdd) ([]string, error) {
	var rows []*model.Address
	// 当前时间
	now := time.Now().Unix()
	var userAddresses []string
	// 遍历差值次数
	for i := int64(0); i < addr.Num; i++ {
		address, privateKeyStrEn, err := e.CreateAddressAndAesKey()
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

//检查erc20交易
//param 需要检查的地址
func (e *ETHHelp) CheckTransaction22(addresss []string) {
	// 获取配置 延迟确认数
	confirmValue := model.SQLGetTAppConfigIntValueByK("block_confirm_num")
	// 获取状态 erc20最后一次获取链上的块数量
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
			log.Print("err: [%T] %s", err, err.Error())
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
		if len(configTokenRowAddresses) > 0 {
			// rpc获取block信息
			logs, err := ethclient.RpcFilterLogs(
				context.Background(),
				startI,
				endI,
				configTokenRowAddresses,
				contractAbi.Events["Transfer"],
			)
			if err != nil {
				log.Print("err: [%T] %s", err, err.Error())
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
				log.Print("err: [%T] %s", err, err.Error())
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
					log.Print("toAddressLogMap no: %s", dbAddressRow.UserAddress)
					return
				}
				for _, log1 := range logs {
					var transferEvent LogTransfer
					err := contractAbi.Unpack(&transferEvent, "Transfer", log1.Data)
					if err != nil {
						log.Print("err: [%T] %s", err, err.Error())
						return
					}
					transferEvent.From = strings.ToLower(common.HexToAddress(log1.Topics[1].Hex()).Hex())
					transferEvent.To = strings.ToLower(common.HexToAddress(log1.Topics[2].Hex()).Hex())
					contractAddress := strings.ToLower(log1.Address.Hex())
					configTokenRow, ok := configTokenRowMap[contractAddress]
					if !ok {
						log.Print("no configTokenRowMap of: %s", contractAddress)
						return
					}
					rpcTxReceipt, err := ethclient.RpcTransactionReceipt(
						context.Background(),
						log1.TxHash.Hex(),
					)
					if err != nil {
						log.Print("err: [%T] %s", err, err.Error())
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
						log.Print("err: [%T] %s", err, err.Error())
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
						log.Print("err: [%T] %s", err, err.Error())
						return
					}
					if hexutil.Encode(input) != hexutil.Encode(rpcTx.Data()) {
						// input 不匹配
						continue
					}
					balanceReal, err := TokenWeiBigIntToEthStr(transferEvent.Tokens, configTokenRow.TokenDecimals)
					if err != nil {
						log.Print("err: [%T] %s", err, err.Error())
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
				log.Print("err: [%T] %s", err, err.Error())
				return
			}
		}
		// 更新检查到的最新区块数
		err = model.SQLUpdateTAppStatusIntByKGreater(
			model.TAppStatusInt{
				K: "erc20_seek_num",
				V: endI,
			},
		)
		if err != nil {
			log.Print("err: [%T] %s", err, err.Error())
			return
		}
	}
}

//检查erc20交易
//param 需要检查的地址
func (e *ETHHelp) CheckTransaction(addresss []string) {
	// 获取配置 延迟确认数
	confirmValue := model.SQLGetTAppConfigIntValueByK("block_confirm_num")
	// 获取状态 erc20最后一次获取链上的块数量
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
			log.Print("err: [%T] %s", err, err.Error())
			return
		}
		// 获取所有token
		var configTokenRowAddresses []string
		configTokenRowMap := make(map[string]*model.TAppConfigToken)
		var aa = model.TAppConfigToken{}
		configTokenRows, err := aa.SQLSelectTAppConfigTokenColAll()
		if err != nil || configTokenRows[0].ID == 0 {
			return
		}
		for _, contractRow := range configTokenRows {
			configTokenRowAddresses = append(configTokenRowAddresses, contractRow.TokenAddress)

			configTokenRowMap[strings.ToLower(contractRow.TokenAddress)] = contractRow
		}
		// 遍历获取需要查询的block信息
		//for i := startI; i < endI; i++ {
		if len(configTokenRowAddresses) > 0 {
			// rpc获取block信息
			logs, err := ethclient.RpcFilterLogs(
				context.Background(),
				startI,
				endI,
				configTokenRowAddresses,
				contractAbi.Events["Transfer"],
			)
			if err != nil {
				log.Print("err: [%T] %s", err, err.Error())
				return
			}
			// 接收地址列表
			//var toAddresses []string
			//// map[接收地址] => []交易信息
			//toAddressLogMap := make(map[string][]types.Log)
			for _, logMode := range logs {
				if logMode.Removed {
					continue
				}
				var transferEvent LogTransfer
				err := contractAbi.Unpack(&transferEvent, "Transfer", logMode.Data)
				if err != nil {
					log.Print("err: [%T] %s", err, err.Error())
					return
				}
				transferEvent.From = strings.ToLower(common.HexToAddress(logMode.Topics[1].Hex()).Hex())
				println("transferEvent.From==", transferEvent.From)
				transferEvent.To = strings.ToLower(common.HexToAddress(logMode.Topics[2].Hex()).Hex())
				println("transferEvent.To==", transferEvent.To)
				hax := logMode.TxHash.Hex()
				println("txhash=======", hax)
				var BlockHash = logMode.BlockHash.Hex()
				println("BlockHash=======", BlockHash)
				//if BlockHash == HAX_0 && logMode.BlockNumber == 0 {
				//	//判断为未打包状态
				//	continue
				//}

				//未打包前blockHash为0，blockNumber为null
				//打包后blockHash块哈希，blockNumber有块号
				//未打包无交易凭据，打包后有交易凭据（不管调用是否成功）
				//合约调用成功后交易凭据中logs会有数据（可能是事件数据，未再确定），但普通转帐交易凭据中logs无数据
				//合约调用成功后交易凭据中gasUsed可能小于gas，如果gas指定不足调用后gasUsed等于gas
				//合约调用input为调用数据，而普通转帐交易input为"0x"

				//hax:= strings.ToLower(common.HexToAddress(logMode.).Hex())
				//println(hax)
				contractAddress := strings.ToLower(logMode.Address.Hex())
				//确认交易是否打包完成
				rpcTxReceipt, err := ethclient.RpcTransactionReceipt(
					context.Background(),
					logMode.TxHash.Hex(),
				)
				if err != nil {
					log.Print("err: [%T] %s", err, err.Error())
					return
				}
				println(strings.ToLower(common.HexToAddress(rpcTxReceipt.TxHash.Hex()).Hex()))
				println("CumulativeGasUsed======", rpcTxReceipt.CumulativeGasUsed)
				println("GasUsed======", rpcTxReceipt.GasUsed)
				//TODO 1=成功，0=失败，是否还有其他状态
				if rpcTxReceipt.Status != 1 {
					continue
				}
				rpcTx, err := ethclient.RpcTransactionByHash(
					context.Background(),
					logMode.TxHash.Hex(),
				)
				if err != nil {
					log.Print("err: [%T] %s", err, err.Error())
					return
				}
				println(contractAddress)
				//if strings.ToLower(rpcTx.To().Hex()) != contractAddress {
				//	// 合约地址和tx的to地址不匹配
				//	continue
				//}
				// 检测input
				input, err := contractAbi.Pack(
					"transfer",
					common.HexToAddress(logMode.Topics[2].Hex()),
					transferEvent.Tokens,
				)
				if err != nil {
					log.Print("err: [%T] %s", err, err.Error())
					return
				}
				if hexutil.Encode(input) != hexutil.Encode(rpcTx.Data()) {
					// input 不匹配
					continue
				}
				println(logMode.BlockNumber)
				balanceReal, err := TokenWeiBigIntToEthStr(transferEvent.Tokens, 18)
				if err != nil {
					log.Print("err: [%T] %s", err, err.Error())
					return
				}
				println("实际支付金额", balanceReal)
			}
			//// 从db中查询这些地址是否是冲币地址中的地址
			//dbAddressRows, err := model.SQLSelectTAddressKeyColByAddress(toAddresses)
			//if err != nil {
			//	log.Print("err: [%T] %s", err, err.Error())
			//	return
			//}
			//// 时间
			//now := time.Now().Unix()
			//// 待添加数组
			//var txErc20Rows []*model.TTxErc20
			//// map[接收地址] => 系统id
			//addressSystemMap := make(map[string]int64)
			//for _, dbAddressRow := range dbAddressRows {
			//	addressSystemMap[dbAddressRow.UserAddress] = int64(dbAddressRow.UseTag)
			//}
			//// 遍历数据库中有交易的地址

		}
		// 更新检查到的最新区块数
		//err = model.SQLUpdateTAppStatusIntByKGreater(
		//	model.TAppStatusInt{
		//		K: "erc20_seek_num",
		//		V: endI,
		//	},
		//)
		//if err != nil {
		//	log.Print("err: [%T] %s", err, err.Error())
		//	return
		//}
	}
	//}
}

//检查余额
func (*ETHHelp) CheckBalance() {

}

//eth交易转账
//param ethTransactionMode 包含交易得所有参数
func (e *ETHHelp) ETHTransaction(ethTransactionMode ETHHelp) (ethHelp *ETHHelp, ethErr *myerr.EthError) {
	return e.transaction(ethTransactionMode, 1)
}

//desc erc20交易转账
//param ethTransactionMode 包含交易得所有参数
func (e *ETHHelp) ERC20Transaction(ethTransactionMode ETHHelp) (ethHelp *ETHHelp, ethErr *myerr.EthError) {
	return e.transaction(ethTransactionMode, 2)
}

//发起交易
//eth交易可以不适用合约，目前没有研究，逻辑上也不需要eth交易，不做处理
//param transactionType 1=eth,2=代币交易
func (e *ETHHelp) transaction(ethTransactionMode ETHHelp, transactionType int) (ethHelp *ETHHelp, ethErr *myerr.EthError) {
	if len(strings.TrimSpace(ethTransactionMode.ToAddress)) != 42 && transactionType == 2 {
		log.Print("请求发送链上交易参数中目标地址不合法！")
		return
	}
	if len(strings.TrimSpace(ethTransactionMode.ContractAddress)) != 42 && transactionType == 2 {
		log.Print("请求发送链上交易参数中合约地址不合法！")
		return
	}
	if len(strings.TrimSpace(ethTransactionMode.FromAddress)) != 42 && transactionType == 1 {
		log.Print("请求发送链上交易参数中没有发起交易地址！")
		return
	}
	if strings.TrimSpace(ethTransactionMode.ToAddress) == "" {
		log.Print("请求发送链上交易参数中没有目标地址！")
		return
	}
	if ethTransactionMode.SendBalance <= 0 {
		log.Print("请求发送链上交易金额小于等于0！")
		return
	}
	if ethTransactionMode.Places <= 8 {
		log.Print("请求发送链上交易的代币位数不合理！")
		return
	}
	if ethTransactionMode.Nonce < 0 {
		log.Print("请求的nonce不可以小于0！")
		return
	}
	// 加载合约信息
	contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
	if err != nil {
		log.Print("err: [%T] %s", err, err.Error())
		return
	}
	//本次eth的交易数量
	var ethTransactionNum = big.NewInt(0)
	//本次代币的交易数量
	var erc20TransactionNum = big.NewInt(0)
	var fromAddress string
	var typeStr string
	if transactionType == 1 {
		//eth交易
		ethTransactionNum = e.BalanceToVirtualCurrencyWei(ethTransactionMode.SendBalance, 18)
		fromAddress = ethTransactionMode.FromAddress
		typeStr = "eth"
		//eth交易我认为就是，冷热钱包的转账交易，不需要太快
		ethTransactionMode.SendGasPrice = ethTransactionMode.GasData.AvgGasPrice
		//一般eth交易默认limit就够了
		ethTransactionMode.SendGasLimit = DefaultGasLimit
	} else if transactionType == 2 {
		//erc20 代币交易
		erc20TransactionNum = e.BalanceToVirtualCurrencyWei(ethTransactionMode.SendBalance, 18)
		fromAddress = ethTransactionMode.ContractAddress
		typeStr = "usdt"
		//TODO erc20交易暂时写死
		ethTransactionMode.SendGasPrice = ethTransactionMode.GasData.FastGasPrice
		if ethTransactionMode.GasData.FastGasLimit < 21000 {
			log.Print("请求发送链上交易的GasLimit不可以小于21000！")
			return
		}
		//TODO erc20支付需要计算，暂时不知道怎么计算
		ethTransactionMode.SendGasLimit = ethTransactionMode.GasData.FastGasLimit
	}

	input, err := contractAbi.Pack(
		"transfer",
		//目标地址
		common.HexToAddress(ethTransactionMode.ToAddress),
		erc20TransactionNum,
	)
	if err != nil {
		log.Print("err: [%T] %s", err, err.Error())
		return
	}
	//组装交易信息
	rpcTx := types.NewTransaction(
		//big.NewInt(int64(ethTransactionMode.Nonce)),
		uint64(ethTransactionMode.Nonce),
		// 交易发起地址，eth是发起地址，erc20交易是合约地址
		common.HexToAddress(fromAddress),
		//ecr20转账不需要设置eth
		ethTransactionNum,
		uint64(ethTransactionMode.SendGasLimit),
		// gasPrice为最低费用，gasLimit为限额费用，理解为首次给予gasPrice最低费用，如果失败，会在这个基础上增加费用，费用会一直增加，超过gaslimit就会结束
		big.NewInt(ethTransactionMode.SendGasPrice),
		input,
	)
	key, err := crypto.HexToECDSA(ethTransactionMode.PrivateKey)
	if err != nil {
		log.Print("签名验证失败err: [%T] %s", err, err.Error())
		return
	}
	//验证签名 ，获取链id
	signedTx, err := types.SignTx(rpcTx, types.NewEIP155Signer(big.NewInt(ethTransactionMode.ChainID)), key)
	if err != nil {
		log.Print("获取链上id失败，err: [%T] %s", err, err.Error())
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 交易发送
	var err2 = ethclient.RpcSendTransaction(ctx, signedTx)
	ethErr = e.transactionErrHander(err2)
	if ethErr != nil {
		return
	}
	ts := types.Transactions{signedTx}
	rawTxBytes := ts.GetRlp(0)
	//TODO 补充交易信息，供后面流程使用，作用不明
	ethTransactionMode.RawTxHex = hex.EncodeToString(rawTxBytes)
	//补充交易信息，供后面流程使用，用来标记交易记录
	ethTransactionMode.TxHash = strings.ToLower(signedTx.Hash().Hex())
	log.Print("一笔%s支付交易发送，费用%s", typeStr, ethTransactionMode.SendBalance)
	ethHelp = &ethTransactionMode
	return
}

//发送交易错误处理
func (*ETHHelp) transactionErrHander(err error) (ethErr *myerr.EthError) {
	if err == nil {
		return nil
	} else {
		println("发送交易链上失败，err: [%T] %s", err, err.Error())
	}
	return ethErr.ErrorType(err.Error())
}

//通过txhash查询交易结果
//return result 交易结果 1=成功，2=未获取结果，0=失败
//return tType 1=eth加以，2=erc20交易
//return  gasUsed gas使用的数量
func (*ETHHelp) CheckTransactionByTxHash(txHash string) (result int, tType int, gasUsed uint64, err error) {
	result = 2
	if len(txHash) <= 30 {
		log.Print("txHash 参数值异常")
		return
	}
	c := ethclient.GetClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var index = 10
	var rpcTxReceipt *types.Receipt
	for {

		rpcTxReceipt, err = c.TransactionReceipt(ctx, common.HexToHash(txHash))
		if err != nil {
			if strings.Contains(err.Error(), "\"code\":-32005,\"message\"") {
				println("eth查询交易结果请求过多，被屏蔽了！", "err", err.Error())
				return
			}
			if strings.Contains(err.Error(), "Receipt retrieval failed err 401 Unauthorized: invalid project id") ||
				strings.Contains(err.Error(), "401 Unauthorized: invalid project id") {
				println("eth访问地址出现异常！", "err", err.Error())
				return
			}
			println("Receipt retrieval failed", "err", err.Error())
			index = index - 1
			if index == 0 {
				break
			}
		}
		//TODO 1=成功，0=失败，是否还有其他状态
		if rpcTxReceipt != nil && rpcTxReceipt.Status != 1 {
			println("交易失败，状态=", rpcTxReceipt.Status)
			return
		} else if rpcTxReceipt != nil && rpcTxReceipt.Status == 1 {
			break
		}

	}
	if rpcTxReceipt == nil {
		println("未获取到交易状态")
		return
	}
	println(rpcTxReceipt.TxHash.Hex())
	println("CumulativeGasUsed======", rpcTxReceipt.CumulativeGasUsed)
	println("GasUsed======", rpcTxReceipt.GasUsed)
	println(rpcTxReceipt.ContractAddress.Hex())
	//交易块 BlockNumber
	println(rpcTxReceipt.BlockNumber)
	//合约地址为空，说明是eth交易
	if rpcTxReceipt.ContractAddress.Hex() == HAX_0 {
		tType = 1
	} else {
		//erc20交易
		tType = 2
	}
	//已经计算出使用gas费用，说明交易成功
	if rpcTxReceipt.CumulativeGasUsed > 0 {
		result = 1
		gasUsed = rpcTxReceipt.GasUsed
	}
	return result, tType, gasUsed, err
}

// 将交易的数字转换为该虚拟币对应的位数
func (*ETHHelp) BalanceToVirtualCurrencyWei(balanceReal float64, wei float64) *big.Int {
	val := int64(float64(balanceReal) * math.Pow(10, float64(wei)))
	b := new(big.Int)
	return b.SetInt64(val)
}
