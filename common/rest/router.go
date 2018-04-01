package rest

import (
    "fmt"
    "sync"
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
)

var router *gin.Engine
var o sync.Once

func getRouter() *gin.Engine {
    o.Do(func() {
        router = gin.New()
    })
    return router
}

func initHttp(initRouters func(engine *gin.Engine)) {
    router = getRouter()

    initRouters(router)
}

func Run(initRouters func(engine *gin.Engine)) error {
    initHttp(initRouters)

    addr := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))
    return router.Run(addr)
}

