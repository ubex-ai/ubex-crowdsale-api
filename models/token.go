package models

type TokenDeployParams struct {
    TotalSupply  string `json:"totalSupply"`
}

// swagger:parameters getUbexBalance
type GetUbexBalanceParams struct {
    // Ethereum address
    // example: 0xED2F74E1fb73b775E6e35720869Ae7A7f4D755aD
    // in: query
    Address string `json:"address"`
}

// swagger:parameters getUbexBalances
type UbexAddressesParams struct {
    // in: body
    Body Addresses `json:"body"`
}