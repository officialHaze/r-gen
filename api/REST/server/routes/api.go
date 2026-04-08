package routes

import (
	v1 "reportgenengine/api/REST/server/routes/v1"
	"reportgenengine/api/REST/server/routes/v1/middleware"

	"github.com/gin-gonic/gin"
)

func MapRoutes(r *gin.Engine) {
	api := r.Group("/api",middleware.RateLimmiter())

	// ======= v1 API ========
	v1.MapRoutes(api.Group("/v1"))
}
