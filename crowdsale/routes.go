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
        sale.POST("/add", rest.SignRequired(), postSaleAddAction)
        sale.POST("/multiplier", rest.SignRequired(), postSetBonusMultiplierAction)
        sale.POST("/close", rest.SignRequired(), postCloseAction)
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
// Get status.
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
// Get token balance
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
// Get token balances
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

// swagger:route POST /crowdsale/add crowdsale addTokens
//
// Add tokens
//
// Add tokens for specified beneficiary (referral system tokens, for example).
//
// Consumes:
// - application/json
// Produces:
// - application/json
// Responses:
//   200: TxSuccessResponse
//   400: RestErrorResponse
//
func postSaleAddAction(c *gin.Context) {
    request := &models.AddressWithAmount{}
    err := c.BindJSON(request)
    if err != nil {
        rest.NewResponder(c).ErrorValidation(err.Error())
        return
    }

    tx, err := GetCrowdsale().Add(request.Address, request.Amount)
    if err != nil {
        rest.NewResponder(c).Error(err.Error())
        return
    }

    rest.NewResponder(c).Success(gin.H{
        "tx": tx.String(),
    })
}

// swagger:route POST /crowdsale/multiplier crowdsale setMultiplier
//
// Set multiplier
//
// Set Bonus tokens rate multiplier, x1000 (i.e. 1200 is 1.2 x 1000 = 120% x1000 = +20% bonus)
//
// Consumes:
// - application/json
// Produces:
// - application/json
// Responses:
//   200: TxSuccessResponse
//   400: RestErrorResponse
//
func postSetBonusMultiplierAction(c *gin.Context) {
    request := &models.SetMultiplierParams{}
    err := c.BindJSON(request)
    if err != nil {
        rest.NewResponder(c).ErrorValidation(err.Error())
        return
    }

    tx, err := GetCrowdsale().SetBonusMultiplier(request.BonusMultiplier)
    if err != nil {
        rest.NewResponder(c).Error(err.Error())
        return
    }

    rest.NewResponder(c).Success(gin.H{
        "tx": tx.String(),
    })
}

// swagger:route POST /crowdsale/close crowdsale closeSale
//
// Close crowdsale
//
// Set closed crowdsale status
//
// Consumes:
// - application/json
// Produces:
// - application/json
// Responses:
//   200: TxSuccessResponse
//   400: RestErrorResponse
//
func postCloseAction(c *gin.Context) {
    request := &models.CloseSaleParams{}
    err := c.BindJSON(request)
    if err != nil {
        rest.NewResponder(c).ErrorValidation(err.Error())
        return
    }

    tx, err := GetCrowdsale().Close(request.Close)
    if err != nil {
        rest.NewResponder(c).Error(err.Error())
        return
    }

    rest.NewResponder(c).Success(gin.H{
        "tx": tx.String(),
    })
}

// swagger:route POST /crowdsale/events crowdsale events
//
// Get events.
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

    events, err := GetCrowdsale().Events(request.Addresses, request.EventNames, request.StartBlock)
    if err != nil {
        rest.NewResponder(c).Error(err.Error())
        return
    }
    rest.NewResponder(c).Success(gin.H{
        "events": events,
    })
}
