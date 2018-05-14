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
    // get events of only latest N blocks
    Latest int64 `json:"latest"`
}
