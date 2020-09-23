package util

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethrpc "j2pay-server/pkg/eth"
	"j2pay-server/pkg/setting"
	"log"
)

//连接以太坊   返回 以太坊web3对象
func EthClient() *ethrpc.EthRPC {

	client := ethrpc.New(fmt.Sprintf("%s", setting.EthConf.Url))
	version, err := client.Web3ClientVersion()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(version)
	return client
}

//创建用户并返回地址
func GetUserAddress(password string) string {
	ks := keystore.NewKeyStore(
		"./address",
		keystore.LightScryptN,
		keystore.LightScryptP)

	address, _ := ks.NewAccount(password)
	fmt.Println("Account address", address.Address.Hex())
	account, err := ks.Export(address, password, password)
	err = ks.Unlock(address, password)
	if err != nil {
		panic(err)
	}
	fmt.Println("private key",account)
	return address.Address.Hex()
}

