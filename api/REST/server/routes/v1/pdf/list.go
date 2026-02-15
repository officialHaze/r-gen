package pdf

import (
	"net/http"
	"path"
	"reportgenengine/constant"
	"reportgenengine/helper"

	"github.com/gin-gonic/gin"
)

// Handle listing all generated PDFs
func HandleListingPDF(c *gin.Context) {
	ip := c.ClientIP()

	dir := path.Join(constant.PDF_PATH, ip)

	infos := helper.ListFileInfo(dir)

	helper.SendResponse(c, http.StatusOK, "", infos)
}
