package application

import (
    "os"
)

const (
    ENVNAME = "APPLICATION_ENV"
    MAINNET_ENV = "mainnet"
    RINKEBY_ENV = "rinkeby"
    DEFAULT_ENV = RINKEBY_ENV
)

var env string

func init() {
    env = getAppEnv(DEFAULT_ENV)
}

func Env() string {
    return env
}

func getAppEnv(def string) (env string) {
    env = os.Getenv(ENVNAME)

    if env = os.Getenv(ENVNAME); env == "" {
        env = def
    }

    return
}
