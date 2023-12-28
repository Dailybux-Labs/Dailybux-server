package router

import (
	"log"
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Make() {
	r := gin.Default()
	//Interface frequency limit
	limiter := tollbooth.NewLimiter(10, nil)
	limiter.SetIPLookups([]string{"X-Real-IP", "X-Forwarded-For", "RemoteAddr"})
	middleware := func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(limiter, c.Writer, c.Request)
		if httpError != nil {
			c.AbortWithStatus(httpError.StatusCode)
			c.JSON(http.StatusBadRequest, Response{Code: httpError.StatusCode, Msg: httpError.Message})
			return
		}
		c.Next()
	}
	r.Use(middleware)
	r.Use(Cors())

	r.POST("/test", test())

	err := r.Run(":8020")
	if err != nil {
		zap.S().Errorln("init client fail rpc url no available", err)
	}
}

func test() func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, Content-Type, X-CSRF-Token, Token,session")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
