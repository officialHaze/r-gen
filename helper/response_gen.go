/*
Helper methods to create response bodies
*/

package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ====== Main Response Schema =======
type Response struct {
	Type    RespType `json:"type"`
	Message string   `json:"message"`
	Data    any      `json:"data,omitempty"`
}

// ======= Response Type Slug - Error, Sucess, Warn =========
type RespType string
type respTypeSchema struct {
	ERROR   RespType
	SUCCESS RespType
	WARN    RespType
}

var RespTypes respTypeSchema = respTypeSchema{
	ERROR:   "ERROR",
	SUCCESS: "SUCCESS",
	WARN:    "WARNING",
}

// ===== Main Response Body Generator ====
func ResponseGen(resptype RespType, message string, data any) Response {
	return Response{
		Type:    resptype,
		Message: message,
		Data:    data,
	}
}

// ======== Send a response =========
type respstatusmapschema struct {
	resptype   RespType
	defaultmsg string
}

var respStatusMap map[int]respstatusmapschema = map[int]respstatusmapschema{
	http.StatusOK: {
		resptype:   RespTypes.SUCCESS,
		defaultmsg: "Success!",
	},
	http.StatusCreated: {
		resptype:   RespTypes.SUCCESS,
		defaultmsg: "Created!",
	},
	http.StatusBadRequest: {
		resptype:   RespTypes.ERROR,
		defaultmsg: "Check your request header and body!",
	},
	http.StatusInternalServerError: {
		resptype:   RespTypes.ERROR,
		defaultmsg: "Internal server error!",
	},
	http.StatusForbidden: {
		resptype:   RespTypes.WARN,
		defaultmsg: "Not allowed!",
	},
	http.StatusNotImplemented: {
		resptype:   RespTypes.WARN,
		defaultmsg: "Not implemented yet!",
	},
}

func SendResponse(c *gin.Context, status int, msg string, data any) {
	if msg == "" {
		msg = respStatusMap[status].defaultmsg
	}

	c.IndentedJSON(
		status,
		ResponseGen(
			respStatusMap[status].resptype,
			msg,
			data,
		),
	)
}
