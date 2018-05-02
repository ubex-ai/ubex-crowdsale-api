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
