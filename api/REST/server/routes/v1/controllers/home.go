package controllers

import (
	"net/http"
	"reportgenengine/helper"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.IndentedJSON(
		http.StatusOK,
		helper.ResponseGen(
			helper.RespTypes.SUCCESS,
			"You have reached the home of V1 API!",
			nil,
		),
	)
}
