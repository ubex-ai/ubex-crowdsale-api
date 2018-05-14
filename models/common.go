package models

// swagger:parameters getBalanceByAddress
type GetBalanceParams struct {
    // Ethereum address
    // example: 0xED2F74E1fb73b775E6e35720869Ae7A7f4D755aD
    // in: query
    Address string `json:"address"`
}

type Address struct {
    Address string `json:"address"`
}

// swagger:parameters getBalancesByAddress
type GetBalanccesByAddressesParams struct {
    // in: body
    Body Addresses `json:"body"`
}

type Addresses struct {
    Addresses []string `json:"addresses"`
}

type AddressWithAmount struct {
    Address string `json:"address"`
    Amount  string `json:"amount"`
}

// swagger:parameters events
type EventsParams struct {
    // in: body
    Body Events `json:"body"`
}

type Events struct {
    // filter events by list of Ethereum addresses
    Addresses []string `json:"addresses"`
    // filter events by list of event names
    EventNames []string `json:"eventNames"`
    // get events of only latest N blocks
    Latest int64 `json:"latest"`
}
