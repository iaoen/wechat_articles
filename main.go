package main

import (
	"encoding/json"
	"mygo/wxapi/api"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	g.GET("/login", api_startLogin)
	g.POST("/search", api_search)
	g.POST("/appmsg", api_appmsg)
	g.Run("127.0.0.1:12312")
}

func api_startLogin(ctx *gin.Context) {
	res := api.StartLogin()
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": res,
	})
}

func api_search(ctx *gin.Context) {
	type MJSON struct {
		Cookie string `json:"cookie"`
		Query  string `json:"query"`
	}
	var mjson MJSON
	err := ctx.BindJSON(&mjson)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "error",
			"data": err.Error(),
		})
		return
	}
	res, err := api.Search(mjson.Cookie, mjson.Query)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "error",
			"data": err.Error(),
		})
	} else {
		var mmap []map[string]any
		err := json.Unmarshal([]byte(res), &mmap)
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": 500,
				"msg":  "error",
				"data": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "success",
			"data": mmap,
		})
	}
}

func api_appmsg(ctx *gin.Context) {
	type MJSON struct {
		Cookie string `json:"cookie"`
		Fakeid string `json:"fakeid"`
		Page   string `json:"page"`
	}
	var mjson MJSON
	err := ctx.BindJSON(&mjson)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "error",
			"data": err.Error(),
		})
		return
	}
	if mjson.Page == "" {
		mjson.Page = "1"
	}
	res, err := api.Appmsgpublish(mjson.Cookie, mjson.Fakeid, mjson.Page)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "error",
			"data": err.Error(),
		})
	} else {
		var mmap []map[string]any
		err := json.Unmarshal([]byte(res), &mmap)
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": 500,
				"msg":  "error",
				"data": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "success",
			"data": mmap,
		})
	}
}
