package pdf

import (
	"fmt"
	"net/http"
	"path"
	"reportgenengine/constant"
	"reportgenengine/helper"
	"reportgenengine/log"

	"github.com/gin-gonic/gin"
)

// Handle PDF deletion
func HandlePDFDelete(c *gin.Context) {
	id := c.Param("id") // Template ID by which the PDF was generated
	ip := c.ClientIP()

	filename := fmt.Sprintf("%s.pdf", id)
	filepath := path.Join(constant.PDF_PATH, ip, filename)

	if err := helper.ReportGenEngineInit().RemoveFile(filepath); err != nil {
		log.Errorf("Failed to remove %s: %v", filepath, err)
		helper.SendResponse(c, http.StatusInternalServerError, "", nil)
		return
	}

	helper.SendResponse(c, http.StatusOK, "PDF removed successfully!", nil)
}
