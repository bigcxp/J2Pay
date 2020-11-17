package heth

import (
	"context"
	"crypto/ecdsa"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/parnurzeal/gorequest"
	"github.com/shopspring/decimal"
	"io"
	"j2pay-server/ethclient"
	"j2pay-server/hcommon"
	"j2pay-server/model"
	"j2pay-server/pkg/setting"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v8"
)

// IsStringInSlice 字符串是否在数组中
func IsStringInSlice(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// IsIntInSlice 数字是否在数组中
func IsIntInSlice(arr []int64, str int64) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// GinFillBindError 检测gin输入绑定错误
func GinFillBindError(c *gin.Context, err error) {
	validatorError, ok := err.(validator.ValidationErrors)
	if ok {
		errMsgList := make([]string, 0, 16)
		for _, v := range validatorError {
			errMsgList = append(errMsgList, fmt.Sprintf("[%s] is %s", strcase.ToSnake(v.Field), v.ActualTag))
		}
		c.JSON(http.StatusOK, gin.H{"error": hcommon.ErrorBind, "err_msg": strings.Join(errMsgList, ", ")})
		return
	}
	unmarshalError, ok := err.(*json.UnmarshalTypeError)
	if ok {
		c.JSON(http.StatusOK, gin.H{"error": hcommon.ErrorBind, "err_msg": fmt.Sprintf("[%s] type error", unmarshalError.Field)})
		return
	}
	if err == io.EOF {
		c.JSON(http.StatusOK, gin.H{"error": hcommon.ErrorBind, "err_msg": fmt.Sprintf("empty body")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": hcommon.ErrorInternal})
}

// GetSign 获取签名
func GetSign(appSecret string, paramsMap gin.H) string {
	var args []string
	var keys []string
	for k := range paramsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := fmt.Sprintf("%s=%v", k, paramsMap[k])
		args = append(args, v)
	}
	baseString := strings.Join(args, "&")
	baseString += fmt.Sprintf("&key=%s", appSecret)
	data := []byte(baseString)
	r := md5.Sum(data)
	signedString := hex.EncodeToString(r[:])
	return strings.ToUpper(signedString)
}

// GetUUIDStr 获取唯一字符串
func GetUUIDStr() string {
	u1 := uuid.NewV4()
	return strings.Replace(u1.String(), "-", "", -1)
}

const (
	// EthToWei 数据单位
	EthToWei = 1e18
	// CoinSymbol 单位标志
	CoinSymbol = "eth"
)

// ethToWeiDecimal 转换单位
var ethToWeiDecimal decimal.Decimal

func init() {
	ethToWeiDecimal = decimal.NewFromInt(EthToWei)
}

// GetNonce 获取nonce值
func GetNonce(address string) (int64, error) {
	// 通过rpc获取
	rpcNonce, err := ethclient.RpcNonceAt(
		context.Background(),
		address,
	)
	if nil != err {
		return 0, err
	}
	// 获取db nonce
	dbNonce:= model.SQLGetTSendMaxNonce(address)
	if dbNonce > rpcNonce {
		rpcNonce = dbNonce
	}
	return rpcNonce, nil
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// AddressBytesToStr 地址转化为字符串
func AddressBytesToStr(addressBytes common.Address) string {
	return strings.ToLower(addressBytes.Hex())
}

// StrToAddressBytes 字符串转化为地址
func StrToAddressBytes(str string) (common.Address, error) {
	if !IsValidAddress(str) {
		return common.HexToAddress("0x0"), fmt.Errorf("str not address: %s", str)
	}
	return common.HexToAddress(str), nil
}

// EthStrToWeiBigInit 转换金额 eth to wei
func EthStrToWeiBigInit(balanceRealStr string) (*big.Int, error) {
	balanceReal, err := decimal.NewFromString(balanceRealStr)
	if err != nil {
		return nil, err
	}
	balanceStr := balanceReal.Mul(ethToWeiDecimal).StringFixed(0)
	b := new(big.Int)
	_, ok := b.SetString(balanceStr, 10)
	if !ok {
		return nil, errors.New("error str to bigint")
	}
	return b, nil
}

// WeiBigIntToEthStr 转换金额 wei to eth
func WeiBigIntToEthStr(wei *big.Int) (string, error) {
	balance, err := decimal.NewFromString(wei.String())
	if err != nil {
		return "0", err
	}
	balanceStr := balance.Div(ethToWeiDecimal).StringFixed(18)
	return balanceStr, nil
}

// TokenEthStrToWeiBigInit 转换金额 eth to wei
func TokenEthStrToWeiBigInit(balanceRealStr string, tokenDecimals int64) (*big.Int, error) {
	balanceReal, err := decimal.NewFromString(balanceRealStr)
	if err != nil {
		return nil, err
	}
	balanceStr := balanceReal.Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(tokenDecimals))).StringFixed(0)
	b := new(big.Int)
	_, ok := b.SetString(balanceStr, 10)
	if !ok {
		return nil, errors.New("error str to bigint")
	}
	return b, nil
}

// TokenWeiBigIntToEthStr 转换金额 wei to eth
func TokenWeiBigIntToEthStr(wei *big.Int, tokenDecimals int64) (string, error) {
	balance, err := decimal.NewFromString(wei.String())
	if err != nil {
		return "0", err
	}
	balanceStr := balance.Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(tokenDecimals))).StringFixed(int32(tokenDecimals))
	return balanceStr, nil
}

// GetPKMapOfAddresses 获取地址私钥
func GetPKMapOfAddresses(addresses []string) (map[string]*ecdsa.PrivateKey, error) {
	addressPKMap := make(map[string]*ecdsa.PrivateKey)
	addressKeyMap, err := model.SQLGetAddressKeyMap(addresses)
	if err != nil {
		return nil, err
	}
	for k, v := range addressKeyMap {
		key := hcommon.AesDecrypt(v.Pwd, fmt.Sprintf("%s", setting.AesConf.Key))
		if len(key) == 0 {
			log.Panicf("error key of: %s", k)
			continue
		}
		if strings.HasPrefix(key, "0x") {
			key = key[2:]
		}
		privateKey, err := crypto.HexToECDSA(key)
		if err != nil {
			log.Panicf("err: [%T] %s", err, err.Error())
			continue
		}
		addressPKMap[k] = privateKey
	}
	return addressPKMap, nil
}

//resp结构体
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
//获取最新gas
func GetGas() (StRespGasPrice,error){
	gresp, body, errs := gorequest.New().
		Get("https://ethgasstation.info/api/ethgasAPI.json").
		Timeout(time.Second * 120).
		End()
	if errs != nil {
		log.Print("err: %s", errs[0], errs[0].Error())
		return StRespGasPrice{},errs[0]
	}
	if gresp.StatusCode != http.StatusOK {
		// 状态错误
		log.Print("req status error: %d", gresp.StatusCode)
		return StRespGasPrice{},errs[0]
	}
	var resp StRespGasPrice
	err := json.Unmarshal([]byte(body), &resp)
	if err != nil {
		log.Print("err: [%T] %s", err, err.Error())
		return StRespGasPrice{},errs[0]
	}
	return resp,nil
}

