package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"tast1/contract"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// client
	client, err := getClient()
	if err != nil {
		fmt.Errorf("failed to connect", err)
	}
	// privatekey

	privateKeyHex, err := getprivateKey()
	if err != nil {
		log.Fatal(err)
	}

	//info
	// BlockByNumber := big.NewInt(5671744)
	// getBlockMsg(client, BlockByNumber)

	// send
	// sendTransAuction(client, privateKeyHex)

	// depliy
	// deployContrcat(client, privateKeyHex)

	// run
	testContrcat(client, privateKeyHex)

}

func getClient() (*ethclient.Client, error) {
	err := godotenv.Load("message.env")
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	infuraURL := os.Getenv("INFURA_URL")
	if infuraURL == "" {
		return nil, fmt.Errorf("INFURA_URL not found in .env file")
	}

	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}

	return client, nil

}

func getBlockMsg(client *ethclient.Client, BlockByNumber *big.Int) {
	// header
	header, err := client.HeaderByNumber(context.Background(), BlockByNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(header.Number.Uint64())
	fmt.Println(header.Time)
	fmt.Println(header.Difficulty.Uint64())
	fmt.Println(header.Hash().Hex())
	// 5671744
	// 1712798400
	// 0
	// 0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5

	// block
	block, err := client.BlockByNumber(context.Background(), BlockByNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(block.Number().Uint64())
	fmt.Println(block.Time())
	fmt.Println(block.Difficulty().Uint64())
	fmt.Println(block.Hash().Hex())
	fmt.Println(len(block.Transactions()))

	// 5671744
	// 1712798400
	// 0
	// 0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5
	// 70

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)
}

func getprivateKey() (string, error) {

	err := godotenv.Load("message.env")
	if err != nil {
		return "", fmt.Errorf("failed to load .env file: %v", err)
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		return "", fmt.Errorf("INFURA_URL not found in .env file")
	}
	// delete 0x
	privateKey = strings.TrimPrefix(privateKey, "0x")
	if len(privateKey) != 64 {
		return "", fmt.Errorf("invalid private key length, expected 64 characters")
	}

	// uncode
	_, err = hex.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("解码失败: %v\n", err)

	}
	return privateKey, nil
}

func sendTransAuction(client *ethclient.Client, privateKeyHex string) {

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("From Address: %s\n", fromAddress.Hex())

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Nonce: %d\n", nonce)

	toAddress := common.HexToAddress("0x1263c125e12538bae30f6fb488049a0b91bf7dcb")
	fmt.Printf("To Address: %s\n", toAddress.Hex())

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	ethValue := new(big.Float).SetInt(value)
	ethValue.Quo(ethValue, big.NewFloat(1e18))
	fmt.Printf("ethValue: %s ETH\n", ethValue.Text('f', 6))

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("gas: %d\n", gasPrice)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Value: value,
		Data:  nil,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("gasLimit: %d\n", gasLimit)

	tx := types.NewTransaction(nonce, toAddress, value, gasPrice.Uint64(), big.NewInt(int64(gasLimit)), nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Chain ID: %d\n", chainID)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Signed Transaction: %s\n", signedTx.Hash().String())

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

}

func deployContrcat(client *ethclient.Client, privateKeyHex string) {

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("From Address: %s\n", fromAddress.Hex())

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Nonce: %d\n", nonce)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Chain ID: %d\n", chainID)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("gas: %d\n", gasPrice)

	// 创建交易对象
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal("创建交易对象失败:", err)
	}
	fmt.Printf("auth: ", auth)
	// 计算合适的gas limit
	bin := contract.ContractMetaData.Bin

	fmt.Printf("bin: ", bin)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		Data: common.FromHex(bin),
	})
	if err != nil {
		log.Println("获取gas limit失败:", err)
		gasLimit = uint64(600000)
	}
	fmt.Printf("gasLimit: %d\n", gasLimit)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = gasLimit

	address, tx, _, err := contract.DeployContract(auth, client, big.NewInt(0))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("合约部署中，地址:", address.Hex())
	fmt.Println("部署交易哈希:", tx.Hash().Hex())
	bind.WaitMined(context.Background(), client, tx)
	fmt.Println("合约部署完成!")

}
func testContrcat(client *ethclient.Client, privateKeyHex string) {

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("From Address: %s\n", fromAddress.Hex())

	// 这里要替换为部署成功的合约地址
	instance, err := contract.NewContract(common.HexToAddress("0xC708d216DE3886D4CB93841eDbf168086cCebB13"), client)
	if err != nil {
		log.Fatal("实例化合约失败:", err)
	}

	count, err := instance.GetCount(nil)
	if err != nil {
		log.Fatal("获取合约count失败:", err)
	}
	fmt.Println("合约初始计数:", count)

	//
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Nonce: %d\n", nonce)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Chain ID: %d\n", chainID)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("gas: %d\n", gasPrice)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal("创建交易对象失败:", err)
	}
	fmt.Printf("auth: ", auth)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(300000)
	auth.Value = big.NewInt(0)

	// tx, err := instance.AddOne(auth)
	// if err != nil {
	// 	log.Fatal("调用addOne失败:", err)
	// }
	// fmt.Println("调用addOne方法交易哈希:", tx.Hash().Hex())
	// bind.WaitMined(context.Background(), client, tx) // 等待交易完成

	count, err = instance.GetCount(nil)
	if err != nil {
		log.Fatal("获取合约count失败:", err)
	}
	fmt.Println("second计数:", count)

	// tx1, err := instance.SetCount(auth, big.NewInt(100))
	// if err != nil {
	// 	log.Fatal("调用setCount失败:", err)
	// }
	// fmt.Println("setCount交易哈希:", tx1.Hash().Hex())

	// 3. 调用 reset 方法（需要 owner 权限）
	tx3, err := instance.Reset(auth)
	if err != nil {
		log.Fatal("调用reset失败:", err)
	}
	fmt.Println("reset交易哈希:", tx3.Hash().Hex())

	count, err = instance.GetCount(nil)
	if err != nil {
		log.Fatal("获取合约count失败:", err)
	}
	fmt.Println("set计数:", count)

}
