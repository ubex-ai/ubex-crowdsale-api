// Package classification UBEX Ethereum API.
//
// Simple HTTP API to the Ethereum blockchain smart contracts of the UBEX Exchange.
//
//     Schemes: http, https
//     BasePath: /
//     Version: 0.3.5
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Dan Gartman<dan.gartman@hotmail.com> https://github.com/dangartman
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - hash
//     - jwt
//
//     SecurityDefinitions:
//     hash:
//          type: apiKey
//          name: X-Sign
//          in: header
//     jwt:
//         type: apiKey
//         name: X-Authorization
//         in: header
//
// swagger:meta
package application

import (
    httpCommon "ubex-crowdsale-api/common/rest"
    "github.com/gin-gonic/gin"
    "sync"
    "ubex-crowdsale-api/common/rest"
    "ubex-crowdsale-api/wallet"
    "ubex-crowdsale-api/token"
    "ubex-crowdsale-api/crowdsale"
)

var o sync.Once

func Run() error {
    return httpCommon.Run(initRoutes)
}

func initRoutes(router *gin.Engine) {
    o.Do(func() {
        router.Use(rest.ExecutionTime())

        wallet.InitRoutes(router)
        token.InitRoutes(router)
        crowdsale.InitRoutes(router)
    })
}
