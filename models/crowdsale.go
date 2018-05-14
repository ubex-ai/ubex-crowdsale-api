package models

type CrowdsaleDeployParams struct {
    WalletAddress string `json:"wallet"`
    TokenRate string `json:"rate"`
}

type CrowdsaleStatus struct {
    // Ethereum address
    // example: 0xFdb3Ae550c4f6a8FD170C3019c97D4F152b65373
    Address string `json:"address"`
    // How many token units a buyer gets per wei
    // example: 10000
    Rate string `json:"rate"`
    // Total amount of wei raised during crowdsale
    // example: 234123000000000000000
    WeiRaised string `json:"weiRaised"`
    // Total amount of issued tokens during crowdsale
    // example: 2341230000000000000000000
    TokensIssued string `json:"tokensIssued"`
}

type TokenPaidEventArgs struct {
    Purchaser string `json:"purchaser"`
    Beneficiary string `json:"beneficiary"`
    WeiAmount string `json:"weiAmount"`
    Created string `json:"created"`
}

type TokenPurchaseEventArgs struct {
    Purchaser string `json:"purchaser"`
    Beneficiary string `json:"beneficiary"`
    WeiAmount string `json:"weiAmount"`
    TokensAmount string `json:"tokensAmount"`
}

// Success
//
// swagger:response
type CrowdsaleStatusResponse struct {
    // in: body
    Body struct {
        Data struct {
            Status CrowdsaleStatus `json:"status"`
        } `json:"data"`
    }
}

// swagger:parameters getCrowdsaleBalance
type GetCrowdsaleBalanceParams struct {
    // Ethereum address
    // example: 0xFdb3Ae550c4f6a8FD170C3019c97D4F152b65373
    // in: query
    Address string `json:"address"`
}


// swagger:parameters getCrowdsaleBalances
type CrowdsaleAddressesParams struct {
    // in: body
    Body Addresses `json:"body"`
}