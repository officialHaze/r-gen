package v1

import (
	"reportgenengine/api/REST/server/routes/v1/controllers"
	"reportgenengine/api/REST/server/routes/v1/html"
	"reportgenengine/api/REST/server/routes/v1/pdf"

	"github.com/gin-gonic/gin"
)

func MapRoutes(v1 *gin.RouterGroup) {
	// ===== HOME =====
	v1.GET("/", controllers.Home)

	// ======= HTML related API routes =======
	html.MapRoutes(v1.Group("/html"))

	// ======= PDF related API routes =======
	pdf.MapRoutes(v1.Group("/pdf"))
}
