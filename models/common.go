package models

type Address struct {
    Address string `json:"address"`
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
    // get events from specified block
    StartBlock int64 `json:"startBlock"`
}

// Success
//
// swagger:response
type TxSuccessResponse struct {
    // in: body
    Body struct {
        Data struct {
            // Ethereum transaction hash
            // example: 0x83422ef776bfb465433e8f6a6e84eab71f039f039e86933aeca20ee46d01d251
            Tx string `json:"tx"`
        } `json:"data"`
    }
}