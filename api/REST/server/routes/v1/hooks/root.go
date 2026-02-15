package hooks

import (
	"reportgenengine/api/REST/server/routes/v1/middleware"

	"github.com/gin-gonic/gin"
)

func FileUploadVerify(r *gin.RouterGroup) {
	r.Use(middleware.FilePresentRule)
	r.Use(middleware.ValidSizeRule)
	r.Use(middleware.UploadLimitRule)
}

func PDFGenVerify(r *gin.RouterGroup) {
	r.Use(middleware.GenQuotaRule)
}
