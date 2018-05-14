package models

type TokenDeployParams struct {
    TotalSupply  string `json:"totalSupply"`
}

// swagger:parameters getUbexBalance
type GetUbexBalanceParams struct {
    // Ethereum address
    // example: 0xFdb3Ae550c4f6a8FD170C3019c97D4F152b65373
    // in: query
    Address string `json:"address"`
}

// swagger:parameters getUbexBalances
type UbexAddressesParams struct {
    // in: body
    Body Addresses `json:"body"`
}