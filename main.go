package main

import (
    appCommon "ubex-api/common/application"
    "ubex-api/application"
    "github.com/sirupsen/logrus"
)

func init() {
    appCommon.Init()
    if err := application.Init(); err != nil {
        logrus.Fatal(err)
    }
}

func main() {
    appChannel := appCommon.GetNewChannel()

    go application.Run()

    appCommon.StartLoop(appChannel)
}
