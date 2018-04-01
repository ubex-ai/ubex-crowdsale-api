package application

import (
    "os"
    "time"
    "github.com/sirupsen/logrus"
)

func GetNewChannel() chan os.Signal {
    return make(chan os.Signal, 1)
}

func StartLoop(interrupt chan os.Signal) {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case int := <-interrupt:
            if int == os.Interrupt {
                logrus.Info("Close application...")
                select {
                case <-time.After(time.Second):
                }

                return
            }
        }
    }
}