package rest

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "ubex-crowdsale-api/common/models"
    "github.com/sirupsen/logrus"
)

const (
    ErrorNoError       = "NO_ERROR"
    ErrorDefault       = "ERROR"
    ErrorValidation    = "ERROR_VALIDATION"
    ErrorAuthorization = "ERROR_AUTHORIZATION"
)

type Responder struct {
    context *gin.Context
}

func (resp Responder) Success(content interface{}) {
    data := &models.RestResponse{
        Data: content,
        Error: models.RestError{
            Code: ErrorNoError,
        },
    }
    resp.context.JSON(http.StatusOK, data)
    resp.Done()
}

func (resp Responder) Error(message interface{}) {
    resp.ErrorWithCode(http.StatusBadRequest, ErrorDefault, message)
}

func (resp Responder) ErrorValidation(message interface{}) {
    resp.ErrorWithCode(http.StatusPreconditionFailed, ErrorValidation, message)
}

func (resp Responder) ErrorAuthorization() {
    resp.ErrorWithCode(http.StatusUnauthorized, ErrorAuthorization, "Auth required")
}
func (resp Responder) ErrorWithCode(httpCode int, restCode string, message interface{}) {
    defer logrus.Errorf("Responding with error: %v", message)

    data := &models.RestResponse{
        Error: models.RestError{
            Code:    restCode,
            Message: message,
        },
    }
    resp.context.JSON(httpCode, data)
    resp.context.Abort()
    resp.Done()
}

func (resp Responder) Done() {
    if done, exists := resp.context.Get("done"); exists {
        if doneChan, ok := done.(chan bool); ok {
            close(doneChan)
        }
    }
}

func NewResponder(c *gin.Context) *Responder {
    return &Responder{c}
}
