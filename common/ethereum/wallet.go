package ethereum

import (
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "fmt"
    "errors"
    "strings"
    "io/ioutil"
    "math/big"
    "github.com/ethereum/go-ethereum/common"
    "context"
    "github.com/sirupsen/logrus"
    "sync"
    "time"
    "github.com/spf13/viper"
    "math/rand"
    "github.com/ethereum/go-ethereum/core/types"
)

var Kilo *big.Int = big.NewInt(1e3)
var Gwei *big.Int = big.NewInt(1e9)
var Szabo *big.Int = big.NewInt(1).Mul(Gwei, Kilo)
var Finney *big.Int = big.NewInt(1).Mul(Szabo, Kilo)
var Ether *big.Int = big.NewInt(1).Mul(Finney, Kilo)

var wallet *Wallet

type Wallet struct {
    Connection *ethclient.Client
    Account    *bind.TransactOpts
    nonce uint64
    nonceGapReached time.Time
    nextTxDelay time.Duration
    mutex *sync.Mutex
}

func InitWallet(endpoint string, keyFile string, keyPass string) error {
    conn, err := ethclient.Dial(endpoint)
    if err != nil {
        return errors.New(fmt.Sprintf("Failed to connect to the Ethereum network: %v", err))
    }

    keyJson, err := ioutil.ReadFile(keyFile)
    if err != nil {
        return err
    }

    acc, err := bind.NewTransactor(strings.NewReader(string(keyJson)), keyPass)
    if err != nil {
        return errors.New(fmt.Sprintf("Failed to authorize main account: %s", err.Error()))
    }

    n, err := conn.PendingNonceAt(context.TODO(), acc.From)
    if err != nil {
        return errors.New(fmt.Sprintf("Failed to get account Nonce: %s", err.Error()))
    }

    if n == 0 {
        logrus.Warn("Account Nonce = 0, something maybe wrong!!!")
    }

    wallet = &Wallet{
        Connection: conn,
        Account:    acc,
        nonce: n,
        nonceGapReached: time.Now().Add(24 * 365 * time.Hour),
        mutex: &sync.Mutex{},
    }

    return nil
}

func GetWallet() *Wallet {
    return wallet
}

func (w *Wallet) BalanceAt(addr string) (*big.Int, error) {
    return w.Connection.BalanceAt(context.TODO(), common.HexToAddress(addr), nil)
}

func (w *Wallet) GetBlockHeaderByNumber(number *big.Int) (*types.Header, error) {
    return w.Connection.HeaderByNumber(context.TODO(), number)
}

func (w *Wallet) GetBlockHeaderByHash(hash common.Hash) (*types.Header, error) {
    return w.Connection.HeaderByHash(context.TODO(), hash)
}

// WARNING!!!
//
// Every call to that method should end up with call to one of OnSuccessTransaction or OnFailTransaction
// Otherwise it leads to DeadLock
func (w *Wallet) GetTransactOpts() (*bind.TransactOpts, error) {
    w.mutex.Lock()

    logrus.Info("Delaying for ", w.nextTxDelay)
    time.Sleep(w.nextTxDelay)
    w.nextTxDelay = time.Duration(viper.GetInt64("geth.tx_delay")) * time.Second

    gas, err := w.Connection.SuggestGasPrice(context.TODO())
    if err != nil {
        return nil, errors.New("cannot get gas price")
    }

    if gas.Cmp(
        big.NewInt(1).Mul(big.NewInt(viper.GetInt64("geth.gas_price")), Gwei),
    ) == 1 {
        logrus.Warn("gas price ", gas, " too high, reducing...")
        gas = big.NewInt(1).Mul(
            big.NewInt(viper.GetInt64("geth.gas_price") + int64(rand.Intn(6))),
            Gwei,
        )
    }

    if n, err := w.Connection.PendingNonceAt(context.TODO(), w.Account.From); err == nil {
        if w.nonceGapReached.Add(
            time.Duration(viper.GetInt64("geth.nonce_replacement_timeout")),
        ).Before(time.Now()) {
            logrus.Warn("local nonce ", w.nonce, " was replaced by chain nonce ", n)
            w.nonce = n
            w.nonceGapReached = time.Now().Add(24 * 365 * time.Hour)
            gas = gas.Add(
                gas,
                big.NewInt(1).Div(
                    gas,
                    big.NewInt(2),
                ),
            )
            w.nextTxDelay = time.Duration(viper.GetInt64("geth.tx_delay")) * 10 * time.Second
        }

        if w.nonce > n {
            logrus.Warn("local nonce ", w.nonce, " > chain nonce ", n, " use local nonce")
        }

        if w.nonce > n + uint64(viper.GetInt64("geth.nonce_gap")) {
            if w.nonceGapReached.After(time.Now()) {
                w.nonceGapReached = time.Now()
            }
            d := time.Duration((w.nonce - n) * 5) * time.Second
            logrus.Warn("Delay response for ", d)
            time.Sleep(d)
            return nil, errors.New(fmt.Sprintf("nonce is too far from chain: %d > %d", w.nonce, n))
        } else {
            w.nonceGapReached = time.Now().Add(24 * 365 * time.Hour)
        }
    }

    logrus.Info("nonce ", w.nonce, " gas ", gas)

    opts := &bind.TransactOpts{
        From: w.Account.From,
        Signer: w.Account.Signer,
        Nonce: big.NewInt(int64(w.nonce)),
        GasPrice: gas,
    }

    return opts, nil
}

func (w *Wallet) OnSuccessTransaction() {
    w.nonce++
    w.mutex.Unlock()
}

func (w *Wallet) OnFailTransaction(err error) {
    if w.ValidateWrongNonce(err) {
        w.nonce++
    }
    w.mutex.Unlock()
}

func (w *Wallet) ValidateWrongNonce(err error) bool {
    return strings.Contains(err.Error(), "replacement") ||
        strings.Contains(err.Error(), "nonce too low") ||
        strings.Contains(err.Error(), "known transaction")
}

func (w *Wallet) ValidateRepeatableTransaction(err error) bool {
    return w.ValidateWrongNonce(err)
}
