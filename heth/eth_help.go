package heth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"j2pay-server/ethclient"
	"j2pay-server/hcommon"
	"j2pay-server/model"
	"j2pay-server/model/request"
	"j2pay-server/pkg/setting"
	"log"
	"math"
	"math/big"
	"strings"
	"time"
)
//eth 交易中用到的属性
type ETHHelp struct {
	//打包一次交易单位
	GasPrice *int64
	//本次交易最大汽费,go-ethereum提供了计算方式，未找到
	GasLimit *int64
	//交易发送金额
	SendBalance  float64
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
}

//获取chainID 区块确认数，作用还未明白
func (*ETHHelp) GetchainID()(chainID int64, err error){
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
func ( *ETHHelp) GetGas()( *int64, *int64){
	// 获取gap price
	gasPriceValue:= model.SQLGetTAppStatusIntValueByK("to_cold_gas_price")
	var gasLimit int64= 21000// gasPriceValue.V
	return gasPriceValue,&gasLimit
}
//通过交易发起地址获取nonce,
//return nonce交易编号每次加一
func ( *ETHHelp) GetNONCE(fromAddress string)(nonce int64, err error){

	// 通过rpc获取
	rpcNonce, err := ethclient.RpcNonceAt(
		context.Background(),
		fromAddress,
	)
	if nil != err {
		return 0, err
	}
	// 获取db nonce
	nonce=model.SQLGetTSendMaxNonce(fromAddress)
	if(err!=nil){
		log.Panicf("GetNonce err: [%T] %s", err, err.Error())
		return 0, err
	}
	if nonce==0{
		nonce=-2
	}else{
		nonce=nonce+1
	}
	if nonce > rpcNonce {
		rpcNonce = nonce+1
	}
	return rpcNonce+1, nil
}

//获取密钥，获取发起交易地址的密钥
//param 发起交易的地址
//return 地址对应的密钥
func ( *ETHHelp) GetPrivateKey(fromAddress string) (privateKey string,err error){
	//通过热钱包地址查询pwd
	keyRow := model.SQLGetTAddressKeyColByAddress(fromAddress)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	if keyRow == nil {
		log.Panicf("no key of: %s", fromAddress)
		return
	}
	privateKey = hcommon.AesDecrypt(keyRow.Pwd, fmt.Sprintf("%s", setting.AesConf.Key))
	if len(privateKey) == 0 {
		log.Panicf("error key of: %s", fromAddress)
		return
	}
	if strings.HasPrefix(privateKey, "0x") {
		privateKey = privateKey[2:]
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


//检查余额
func ( *ETHHelp) CheckBalance(){

}

//eth交易转账
//param ethTransactionMode 包含交易得所有参数
func (e *ETHHelp) ETHTransaction(ethTransactionMode ETHHelp)(success bool,err error){
	return e.transaction(ethTransactionMode,1)
}

//desc erc20交易转账
//param ethTransactionMode 包含交易得所有参数

func (e *ETHHelp) ERC20Transaction(ethTransactionMode ETHHelp)(success bool,err error){

	return e.transaction(ethTransactionMode,2)
}

//发起交易
//eth交易可以不适用合约，目前没有研究，逻辑上也不需要eth交易，不做处理
//param transactionType 1=eth,2=代币交易
func (e *ETHHelp) transaction(ethTransactionMode ETHHelp,transactionType int)(success bool,err error){
	success=false
	if strings.TrimSpace( ethTransactionMode.ContractAddress)=="" && transactionType==2 {
		log.Panicf("请求发送链上交易参数中没有合约地址！",)
		return
	}
	if strings.TrimSpace( ethTransactionMode.FromAddress)=="" && transactionType==1 {
		log.Panicf("请求发送链上交易参数中没有发起交易地址！",)
		return
	}
	if strings.TrimSpace( ethTransactionMode.ToAddress)=="" {
		log.Panicf("请求发送链上交易参数中没有目标地址！",)
		return
	}
	if ethTransactionMode.SendBalance<=0 {
		log.Panicf("请求发送链上交易金额小于等于0！",)
		return
	}
	if  ethTransactionMode.Places<=8 {
		log.Panicf("请求发送链上交易的代币位数不合理！",)
		return
	}
	// 加载合约信息
	contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	//本次eth的交易数量
	var ethTransactionNum =big.NewInt(0)
	//本次代币的交易数量
	var erc20TransactionNum =big.NewInt(0)
	var fromAddress string
	var typeStr string
	if transactionType==1{
		//eth交易
		ethTransactionNum= e.BalanceToVirtualCurrencyWei(ethTransactionMode.SendBalance,18)
		fromAddress=ethTransactionMode.FromAddress
		typeStr="eth"
	}else if transactionType==2{
		//erc20 代币交易
		erc20TransactionNum=e.BalanceToVirtualCurrencyWei(ethTransactionMode.SendBalance,18)
		fromAddress=ethTransactionMode.ContractAddress
		typeStr="usdt"
	}

	input, err := contractAbi.Pack(
		"transfer",
		//目标地址
		common.HexToAddress(ethTransactionMode.ToAddress),
		erc20TransactionNum,
	)
	if err != nil {
		log.Panicf("err: [%T] %s", err, err.Error())
		return
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	//委托合约发起交易
	rpcTx := types.NewTransaction(
		//big.NewInt(int64(ethTransactionMode.Nonce)),
		uint64(ethTransactionMode.Nonce),
		// 交易发起地址，eth是发起地址，erc20交易是合约地址
		common.HexToAddress(fromAddress),
		//ecr20转账不需要设置eth
		ethTransactionNum,
		// gasPrice为最低费用，gasLimit为限额费用，理解为首次给予gasPrice最低费用，如果失败，会在这个基础上增加费用，费用会一直增加，超过gaslimit就会结束
		uint64(*ethTransactionMode.GasLimit),
		big.NewInt(*ethTransactionMode.GasPrice),
		input,
	)
	key, err := crypto.HexToECDSA(ethTransactionMode.PrivateKey)
	if err != nil {
		log.Panicf("签名验证失败err: [%T] %s", err, err.Error())
		return
	}
	//验证签名 ，获取链id
	signedTx, err := types.SignTx(rpcTx, types.NewEIP155Signer(big.NewInt(ethTransactionMode.ChainID)), key)
	if err != nil {
		log.Panicf("获取链上id失败，err: [%T] %s", err, err.Error())
		return
	}

	// 交易发送
	serr := ethclient.RpcSendTransaction(context.Background(),signedTx)
	if serr != nil {
		log.Panicf("发送交易链上失败，err: [%T] %s",serr, serr.Error())
		return
	}

	var dbi =DeployBackendImpl{}
	// 等待挖矿完成
	bind.WaitMined(context.Background(), &dbi,signedTx)
	log.Print("一笔%s支付交易发送，费用%s",typeStr,ethTransactionMode.SendBalance)
	success=true
	return
}



type DeployBackendImpl struct {

}
func (d *DeployBackendImpl) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error){
	//println("TransactionReceipt")

	//println(txHash)
	return nil,nil
}
func (d *DeployBackendImpl) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error){
	println("CodeAt")
	//println(blockNumber)
	return nil,nil
}


// 将交易的数字转换为该虚拟币对应的位数
func ( *ETHHelp) BalanceToVirtualCurrencyWei(balanceReal float64,wei float64) ( *big.Int) {
	val:= int64(float64(balanceReal)*math.Pow(10,float64(wei)))
	b := new(big.Int)
	return b.SetInt64(val)
}