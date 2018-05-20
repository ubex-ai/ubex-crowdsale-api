package crowdsale

import (
    "ubex-crowdsale-api/solidity/bindings/ubex_crowdsale"
    "github.com/ethereum/go-ethereum/common"
    "github.com/spf13/viper"
    "errors"
    "fmt"
    "math/big"
    "ubex-crowdsale-api/common/ethereum"
    "ubex-crowdsale-api/models"
    modelsCommon "ubex-crowdsale-api/common/models"
    "github.com/ethereum/go-ethereum/core/types"
    "ubex-crowdsale-api/token"
    "github.com/sirupsen/logrus"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
)

var cr *Crowdsale

type Crowdsale struct {
    *ethereum.Contract
    Crowdsale *ubex_crowdsale.UbexCrowdsale
}

func Init() error {
    c := ethereum.NewContract(viper.GetString("ethereum.address.crowdsale"))
    c.InitEvents(ubex_crowdsale.UbexCrowdsaleABI)

    s, err := ubex_crowdsale.NewUbexCrowdsale(c.Address, c.Wallet.Connection)
    if err != nil {
        return errors.New(fmt.Sprintf("Failed to instantiate a Crowdsale contract: %v", err))
    }

    cr = &Crowdsale{
        Contract: c,
        Crowdsale: s,
    }

    return nil
}

func GetCrowdsale() *Crowdsale {
    return cr
}

func (s *Crowdsale) Deploy(params *models.CrowdsaleDeployParams) (*common.Address, *types.Transaction, error) {
    tokenAddr := token.GetToken().Address

    tokenRate, ok := big.NewInt(0).SetString(params.TokenRate, 0)
    if !ok {
        return nil, nil, fmt.Errorf("wrong TokenRate provided: %s", params.TokenRate)
    }

    bonusMultiplier, ok := big.NewInt(0).SetString(params.BonusMultiplier, 0)
    if !ok {
        return nil, nil, fmt.Errorf("wrong BonusMultiplier provided: %s", params.TokenRate)
    }

    address, tx, _, err := ubex_crowdsale.DeployUbexCrowdsale(
        &bind.TransactOpts{
            From: s.Wallet.Account.From,
            Signer: s.Wallet.Account.Signer,
            Nonce: big.NewInt(int64(1)),
            GasPrice: big.NewInt(1).Mul(big.NewInt(10), ethereum.Gwei),
        },
        s.Wallet.Connection,
        tokenRate,
        common.HexToAddress(params.WalletAddress),
        tokenAddr,
        bonusMultiplier,
    )
    if err != nil {
        return nil, nil, fmt.Errorf("failed to deploy contract: %v", err)
    }
    return &address, tx, nil
}

func (s *Crowdsale) Balance(addr string) (*big.Int, error) {
    return s.Crowdsale.Balances(nil, common.HexToAddress(addr))
}

func (s *Crowdsale) Status() (*models.CrowdsaleStatus, error) {
    weiRaised, err := s.Crowdsale.WeiRaised(nil)
    if err != nil {
        return nil, err
    }

    rate, err := s.Crowdsale.Rate(nil)
    if err != nil {
        return nil, err
    }

    multiplier, err := s.Crowdsale.BonusMultiplier(nil)
    if err != nil {
        return nil, err
    }

    tokensIssued, err := s.Crowdsale.TokensIssued(nil)
    if err != nil {
        return nil, err
    }

    return &models.CrowdsaleStatus{
        Address: s.Address.String(),
        TokensIssued: tokensIssued.String(),
        Rate: rate.String(),
        BonusMultiplier: multiplier.String(),
        WeiRaised: weiRaised.String(),
    }, nil
}

func (s *Crowdsale) Add(addr string, amount string) (common.Hash, error) {
    tokenAmount, ok := big.NewInt(0).SetString(amount, 0)
    if !ok {
        return common.Hash{}, fmt.Errorf("wrong number provided: %s", amount)
    }

    opts, err := s.Wallet.GetTransactOpts()
    if err != nil {
        s.Wallet.OnFailTransaction(err)
        return common.Hash{}, err
    }

    tx, err := s.Crowdsale.AddTokens(opts, common.HexToAddress(addr), tokenAmount)
    if err != nil {
        logrus.Error("Add failed ", err.Error())
        s.Wallet.OnFailTransaction(err)

        if s.Wallet.ValidateRepeatableTransaction(err) {
            logrus.Warn("Repeat Add to ", addr)

            return s.Add(addr, amount)
        }

        return common.Hash{}, err
    }
    s.Wallet.OnSuccessTransaction()

    logrus.Info("Added ", amount, " tokens to ", addr, ", tx ", tx.Hash().String())

    return tx.Hash(), nil
}

func (s *Crowdsale) SetBonusMultiplier(multiplier string) (common.Hash, error) {
    multiplierCoeff, ok := big.NewInt(0).SetString(multiplier, 0)
    if !ok {
        return common.Hash{}, fmt.Errorf("wrong number provided: %s", multiplier)
    }

    opts, err := s.Wallet.GetTransactOpts()
    if err != nil {
        s.Wallet.OnFailTransaction(err)
        return common.Hash{}, err
    }

    tx, err := s.Crowdsale.SetBonusMultiplier(opts, multiplierCoeff)
    if err != nil {
        logrus.Error("SetBonusMultiplier failed ", err.Error())
        s.Wallet.OnFailTransaction(err)

        if s.Wallet.ValidateRepeatableTransaction(err) {
            logrus.Warn("Repeat SetBonusMultiplier to ", multiplier)

            return s.SetBonusMultiplier(multiplier)
        }

        return common.Hash{}, err
    }
    s.Wallet.OnSuccessTransaction()

    logrus.Info("Multiplier is set to ", multiplier, ", tx ", tx.Hash().String())

    return tx.Hash(), nil
}

func (s *Crowdsale) Close(close bool) (common.Hash, error) {
    opts, err := s.Wallet.GetTransactOpts()
    if err != nil {
        s.Wallet.OnFailTransaction(err)
        return common.Hash{}, err
    }

    tx, err := s.Crowdsale.CloseCrowdsale(opts, close)
    if err != nil {
        logrus.Error("Close failed ", err.Error())
        s.Wallet.OnFailTransaction(err)

        if s.Wallet.ValidateRepeatableTransaction(err) {
            logrus.Warn("Repeat Close")

            return s.Close(close)
        }

        return common.Hash{}, err
    }
    s.Wallet.OnSuccessTransaction()

    logrus.Info("Closed status of sale is set to ", close, ", tx ", tx.Hash().String())

    return tx.Hash(), nil
}

func (s *Crowdsale) Events(addrs []string, eventNames []string, startBlock int64) ([]modelsCommon.ContractEvent, error) {
    hashAddrs := make([]common.Hash, len(addrs))
    for _, addr := range addrs {
        hashAddrs = append(hashAddrs, common.HexToHash(addr))
    }

    hashEventNames := make([]common.Hash, len(eventNames))
    for _, eventName := range eventNames {
        name, ok := s.Contract.EventHashes[eventName]
        if !ok {
            return nil, fmt.Errorf("unknown event name provided: %s", eventName)
        }
        hashEventNames = append(hashEventNames, common.HexToHash(name))
    }

    from := big.NewInt(viper.GetInt64("ethereum.start_block.crowdsale"))
    if startBlock != 0 {
        from = big.NewInt(startBlock)
    }

    events, err := s.GetEventsByTopics(
        [][]common.Hash{hashEventNames, hashAddrs},
        from,
    )
    if err != nil {
        return nil, err
    }

    resEvents := make([]modelsCommon.ContractEvent, 0)

    for _, event := range events {
        switch {
        case event.Name == "TokenDelivered" || event.Name == "TokenAdded":
            event.Args = models.TokenWithAddressEventArgs{
                Address: common.BytesToAddress(event.RawArgs[0]).String(),
                TokensAmount: common.BytesToHash(event.RawArgs[1]).Big().String(),
            }
        case event.Name == "TokenPurchase":
            event.Args = models.TokenPurchaseEventArgs{
                Purchaser: common.BytesToAddress(event.RawArgs[0]).String(),
                Beneficiary: common.BytesToAddress(event.RawArgs[1]).String(),
                WeiAmount: common.BytesToHash(event.RawArgs[2]).Big().String(),
                TokensAmount: common.BytesToHash(event.RawArgs[3]).Big().String(),
            }
        default:
            return nil, fmt.Errorf("unknown event type: %s", event.Name)
        }

        resEvents = append(resEvents, event)
    }

    return resEvents, nil
}
