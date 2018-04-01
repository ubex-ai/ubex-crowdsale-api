package models

type RestResponse struct {
    Data  interface{} `json:"data"`
    Error RestError   `json:"error"`
}

type RestError struct {
    Code    string      `json:"code"`
    Message interface{} `json:"message"`
}

// Error
//
// swagger:response
type RestErrorResponse struct {
    // in: body
    Body struct {
        Error struct {
            // error code
            // example: VALIDATION_ERROR
            Code string `json:"code"`
            // error description
            // example: id should be in UUID format
            Message string `json:"message"`
        } `json:"error"`
    }
}
