package router

import (
	"log"
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"dailybux/handler"
	"dailybux/model"
	"dailybux/router/result"
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
			c.JSON(http.StatusBadRequest, result.Response{Code: httpError.StatusCode, Msg: httpError.Message})
			return
		}
		c.Next()
	}
	r.Use(middleware)
	r.Use(Cors())

	r.POST("/login", login())
	r.POST("/userInfo", userInfo())
	r.GET("/crunch", crunchInfo())
	r.GET("/peanut", peanutInfo())
	r.GET("/dailyCheckIn", dailyCheckIn())

	err := r.Run(":8020")
	if err != nil {
		zap.S().Errorln("init client fail rpc url no available", err)
	}
}

func login() func(c *gin.Context) {
	return func(c *gin.Context) {
		var fetchErc20Req model.LoginReq
		err := c.BindJSON(&fetchErc20Req)
		if err != nil {
			zap.S().Errorf("[%v] BindJSON %v", fetchErc20Req, err.Error())
		}
		res, err := handler.Login(&fetchErc20Req)
		if err != nil {
			zap.S().Errorf("[%v] FetchErc20 %v", fetchErc20Req, err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, result.OK.WithData(res))
	}
}

func userInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req model.UserInfoReq
		err := c.BindJSON(&req)
		if err != nil {
			zap.S().Errorf("[%v] BindJSON %v", req, err.Error())
		}
		res, err := handler.UserInfo(&req)
		if err != nil {
			zap.S().Errorf("[%v] FetchErc20 %v", req, err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, result.OK.WithData(res))
	}
}

func dailyCheckIn() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req model.UserInfoReq
		err := c.BindJSON(&req)
		if err != nil {
			zap.S().Errorf("[%v] BindJSON %v", req, err.Error())
		}
		res, err := handler.DailyCheckIn(&req)
		if err != nil {
			zap.S().Errorf("[%v] FetchErc20 %v", req, err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, result.OK.WithData(res))
	}
}

func crunchInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		res, err := handler.Crunch()
		if err != nil {
			zap.S().Errorf("FetchErc20 %v", err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, result.OK.WithData(res))
	}
}

func peanutInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		res, err := handler.Peanut()
		if err != nil {
			zap.S().Errorf("FetchErc20 %v", err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, result.OK.WithData(res))
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
