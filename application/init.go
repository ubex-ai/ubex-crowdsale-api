package application

import (
    "github.com/spf13/viper"
    "ubex-crowdsale-api/common/ethereum"
    "ubex-crowdsale-api/token"
    "ubex-crowdsale-api/crowdsale"
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
