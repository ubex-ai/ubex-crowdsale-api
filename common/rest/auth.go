package rest

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "github.com/spf13/viper"
    "fmt"
)

const SIGN_HEADER = "X-Authorization"

func SignRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        signHeader := c.GetHeader(SIGN_HEADER)
        if signHeader == "" {
            NewResponder(c).ErrorAuthorization()
            return
        }

        if signHeader != viper.GetString("auth.secret") {
            NewResponder(c).ErrorWithCode(
                http.StatusUnauthorized,
                ErrorAuthorization,
                fmt.Sprintf(
                    "Incorrect signature, provided %v, expected %v",
                    signHeader,
                    "xxx",
                ),
            )
            return
        }
    }
}
