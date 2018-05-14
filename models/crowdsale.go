package models

type CrowdsaleDeployParams struct {
    WalletAddress string `json:"wallet"`
    TokenRate string `json:"rate"`
}

type CrowdsaleStatus struct {
    Address string `json:"address"`
    Rate string `json:"rate"`
    WeiRaised string `json:"weiRaised"`
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
    // example: 0xED2F74E1fb73b775E6e35720869Ae7A7f4D755aD
    // in: query
    Address string `json:"address"`
}


// swagger:parameters getCrowdsaleBalances
type CrowdsaleAddressesParams struct {
    // in: body
    Body Addresses `json:"body"`
}