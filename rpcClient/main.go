package rpcClient

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	CtxTimeout = 100 * time.Second
)

type Chain struct {
	ChainName       string
	ChainID         *big.Int
	SuggestGasPrice *big.Int
	RpcCli          *rpc.Client
	ethClient       *ethclient.Client
}

func newChain(rpcUrl string) (*Chain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("rpcClient err: %v", err)
	}
	defer rpcClient.Close()

	ethClient := ethclient.NewClient(rpcClient)
	// ethClient, err := ethclient.Dial(rpcUrl)
	// if err != nil {
	// 	return nil, err
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	// defer cancel()
	// gasPrice, err := ethClient.SuggestGasPrice(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// chainID, err := ethClient.NetworkID(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return &Chain{
		ChainID:         big.NewInt(7001),
		SuggestGasPrice: big.NewInt(7),
		// RpcCli:          rpcClient,
		ethClient: ethClient,
	}, nil
}

func NewChain(rpcUrl string) (*Chain, error) {
	return newChain(rpcUrl)
}

func (c *Chain) GetNonce(fromAddress common.Address) (nonce uint64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	nonce, err = c.ethClient.PendingNonceAt(ctx, fromAddress)
	return
}

func (c *Chain) MakeTxn(privKeyStr string, nonce uint64) (*types.Transaction, error) {

	// 处理私钥
	privKey, err := crypto.HexToECDSA(privKeyStr)
	if err != nil {
		return nil, err
	}

	// 处理fromaddress
	to := common.HexToAddress("0x387F83710c848Ead3047B2cDF85Ad87127309A49")

	// The signer activates the 1559 features even before the fork,
	// so the new 1559 txs can be created with this signer.
	signer := types.LatestSignerForChainID(c.ChainID)

	return types.MustSignNewTx(privKey, signer, &types.DynamicFeeTx{
		ChainID:    signer.ChainID(),
		Nonce:      nonce,
		GasTipCap:  big.NewInt(1.5 * 1e9),
		GasFeeCap:  big.NewInt(1.5*1e9 + 16),
		Gas:        21000,
		To:         &to,
		Value:      big.NewInt(100),
		Data:       nil,
		AccessList: nil,
	}), nil
}

func (c *Chain) EnRawTxn(signedTx *types.Transaction) string {

	buf := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf.Bytes())
}

func (c *Chain) SendRawTxn(rawTx string) (string, uint64, error) {

	tx := new(types.Transaction)
	txBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return "", 0, fmt.Errorf("DecodeString err: %v", err)
	}
	rlpStream := rlp.NewStream(bytes.NewBuffer(txBytes), 0)
	if err := tx.DecodeRLP(rlpStream); err != nil {
		return "", 0, fmt.Errorf("DecodeRLP err: %v", err)
	}

	err = c.ethClient.SendTransaction(context.Background(), tx)
	if err != nil {
		// 异常处理
		// tx already in mempool
		if strings.Contains(fmt.Sprintf("%v", err), "tx already in mempool") {
			err = fmt.Errorf("tx already in mempool, it's nonce: %v", tx.Nonce())
		}
		return "", 0, err
	}

	return tx.Hash().Hex(), tx.Nonce(), nil
}
