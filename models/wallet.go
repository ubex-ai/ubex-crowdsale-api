package models

import "github.com/ethereum/go-ethereum/core/types"

// Success
//
// swagger:response
type GetBalanceSuccessResponse struct {
    // in: body
    Body struct {
        Data struct {
            // Ethereum address
            // example: 0xED2F74E1fb73b775E6e35720869Ae7A7f4D755aD
            Address string `json:"address"`
            // Ethereum address balance
            // example: 42340000000000000000
            Balance string `json:"balance"`
        } `json:"data"`
    }
}

// Success
//
// swagger:response
type GetBalancesSuccessResponse struct {
    // in: body
    Body struct {
        Data struct {
            // Ethereum address balances list
            // example: {"0xED2F74E1fb73b775E6e35720869Ae7A7f4D755aD": "42340000000000000000"}
            Balances map[string]string `json:"balances"`
        } `json:"data"`
    }
}

// Success
//
// swagger:response
type GetBlockSuccessResponse struct {
    // in: body
    Body struct {
        Data struct {
            // Block number
            // example: 5349009
            Number string `json:"number"`
            // Block hash
            // example: 0x43dbc6e17450307c226e78cc1146a15662fb2ff16369d0e5c24652576b9d8e5e
            Hash string `json:"hash"`
            // Block creation timestamp
            // example: 1522415005
            Timestamp string `json:"timestamp"`
            // Block raw data
            Raw *types.Header `json:"raw"`
        } `json:"data"`
    }
}