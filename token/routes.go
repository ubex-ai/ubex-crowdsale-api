package token

import (
    "github.com/gin-gonic/gin"
    "ubex-crowdsale-api/common/rest"
    "ubex-crowdsale-api/models"
)

func InitRoutes(router *gin.Engine) {
    ubex := router.Group("/ubex")
    {
        ubex.POST("/deploy", rest.SignRequired(), postDeployTokenAction)
        ubex.GET("/balance/:address", getUbexBalanceAction)
        ubex.POST("/balances", getUbexBalancesAction)
    }
}

func postDeployTokenAction(c *gin.Context) {
    request := &models.TokenDeployParams{}
    err := c.BindJSON(request)
    if err != nil {
        rest.NewResponder(c).ErrorValidation(err.Error())
        return
    }

    addr, tx, err := GetToken().Deploy(request.TotalSupply)
    if err != nil {
        rest.NewResponder(c).Error(err.Error())
        return
    }

    rest.NewResponder(c).Success(gin.H{
        "address": addr.String(),
        "tx":      tx.Hash().String(),
    })
}

// swagger:route GET /ubex/balance/:address token getUbexBalance
//
// Get balance
//
// Get UBEX token balance for particular Ethereum address.
//
// Consumes:
// - application/json
// Produces:
// - application/json
// Responses:
//   200: GetBalanceSuccessResponse
//   400: RestErrorResponse
//
func getUbexBalanceAction(c *gin.Context) {
    addr := c.Param("address")
    bal, err := GetToken().Balance(addr)
    if err != nil {
        rest.NewResponder(c).Error(err.Error())
        return
    }

    rest.NewResponder(c).Success(gin.H{
        "address": addr,
        "balance": bal.String(),
    })
}

// swagger:route POST /ubex/balances token getUbexBalances
//
// Get balances
//
// Get UBEX token balances for list of Ethereum addresses.
//
// Consumes:
// - application/json
// Produces:
// - application/json
// Responses:
//   200: GetBalancesSuccessResponse
//   400: RestErrorResponse
//
func getUbexBalancesAction(c *gin.Context) {
    request := &models.Addresses{}
    err := c.BindJSON(request)
    if err != nil {
        rest.NewResponder(c).ErrorValidation(err.Error())
        return
    }

    bals := map[string]string{}
    for _, addr := range request.Addresses {
        bal, err := GetToken().Balance(addr)
        if err != nil {
            rest.NewResponder(c).Error(err.Error())
            return
        }
        bals[addr] = bal.String()
    }

    rest.NewResponder(c).Success(gin.H{
        "balances": bals,
    })
}
