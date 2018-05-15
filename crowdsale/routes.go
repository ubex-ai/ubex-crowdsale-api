package crowdsale

import (
    "github.com/gin-gonic/gin"
    "ubex-crowdsale-api/common/rest"
    "ubex-crowdsale-api/models"
)

func InitRoutes(router *gin.Engine) {
    sale := router.Group("/crowdsale")
    {
        sale.POST("/deploy", rest.SignRequired(), postDeploySaleAction)
        sale.GET("/status", getSaleStatusAction)
        sale.GET("/balance/:address", getSaleTokensBalanceAction)
        sale.POST("/balances", postSaleTokensBalancesAction)
        sale.POST("/events", postSaleEventsAction)
    }
}

func postDeploySaleAction(c *gin.Context) {
    request := &models.CrowdsaleDeployParams{}
    err := c.BindJSON(request)
    if err != nil {
        rest.NewResponder(c).ErrorValidation(err.Error())
        return
    }

    addr, tx, err := GetCrowdsale().Deploy(request)
    if err != nil {
        rest.NewResponder(c).Error(err.Error())
        return
    }

    rest.NewResponder(c).Success(gin.H{
        "address": addr.String(),
        "tx":      tx.Hash().String(),
    })
}

// swagger:route GET /crowdsale/status crowdsale getStatus
//
// Get UBEX Crowdsale status.
//
// Consumes:
// - application/json
// Produces:
// - application/json
// Responses:
//   200: CrowdsaleStatusResponse
//   400: RestErrorResponse
//
func getSaleStatusAction(c *gin.Context) {
    st, err := GetCrowdsale().Status()
    if err != nil {
        rest.NewResponder(c).Error(err.Error())
        return
    }

    rest.NewResponder(c).Success(gin.H{
        "status": st,
    })
}

// swagger:route GET /crowdsale/balance/:address crowdsale getCrowdsaleBalance
//
// Get UBEX Crowdsale token balance
//
// Get UBEX Crowdsale issued (but not paid) token balance for particular Ethereum address.
//
// Consumes:
// - application/json
// Produces:
// - application/json
// Responses:
//   200: GetBalanceSuccessResponse
//   400: RestErrorResponse
//
func getSaleTokensBalanceAction(c *gin.Context) {
    addr := c.Param("address")
    bal, err := GetCrowdsale().Balance(addr)
    if err != nil {
        rest.NewResponder(c).Error(err.Error())
        return
    }

    rest.NewResponder(c).Success(gin.H{
        "balance": bal.String(),
    })
}

// swagger:route POST /crowdsale/balances crowdsale getCrowdsaleBalances
//
// Get UBEX Crowdsale token balances
//
// Get UBEX Crowdsale issued (but not paid) token balances for list of Ethereum addresses.
//
// Consumes:
// - application/json
// Produces:
// - application/json
// Responses:
//   200: GetBalancesSuccessResponse
//   400: RestErrorResponse
//
func postSaleTokensBalancesAction(c *gin.Context) {
    request := &models.Addresses{}
    err := c.BindJSON(request)
    if err != nil {
        rest.NewResponder(c).ErrorValidation(err.Error())
        return
    }

    bals := map[string]string{}
    for _, addr := range request.Addresses {
        bal, err := GetCrowdsale().Balance(addr)
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

// swagger:route POST /crowdsale/events crowdsale events
//
// Get UBEX Crowdsale events.
//
// Consumes:
// - application/json
// Produces:
// - application/json
// Responses:
//   200: ContractEventResponse
//   400: RestErrorResponse
//
func postSaleEventsAction(c *gin.Context) {
    request := &models.Events{}
    err := c.BindJSON(request)
    if err != nil {
        rest.NewResponder(c).ErrorValidation(err.Error())
        return
    }

    events, err := GetCrowdsale().Events(request.Addresses, request.EventNames, request.Latest)
    if err != nil {
        rest.NewResponder(c).Error(err.Error())
        return
    }
    rest.NewResponder(c).Success(gin.H{
        "events": events,
    })
}
