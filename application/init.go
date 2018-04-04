package application

import (
    "github.com/spf13/viper"
    "ubex-api/common/ethereum"
    "ubex-api/token"
    "ubex-api/crowdsale"
)

func Init() error {
    err := ethereum.InitWallet(
        viper.GetString("ethereum.socket"),
        viper.GetString("ethereum.wallet.file"),
        viper.GetString("ethereum.wallet.pass"),
    )
    if err != nil {
        return err
    }

    if err := token.Init(); err != nil {
        return err
    }

    if err := crowdsale.Init(); err != nil {
        return err
    }

    return nil
}
