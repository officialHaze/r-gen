package html

import (
	"fmt"
	"net/http"
	"path"
	"reportgenengine/constant"
	"reportgenengine/helper"
	"reportgenengine/log"

	"github.com/gin-gonic/gin"
)

// Delete a template by ID
func HandleTemplDelete(c *gin.Context) {
	templId := c.Param("id")
	ip := c.ClientIP()

	filename := fmt.Sprintf("%s.html", templId)
	filepath := path.Join(constant.TEMPLATE_PATH, ip, filename)

	if err := helper.ReportGenEngineInit().RemoveFile(filepath); err != nil {
		log.Errorf("Failed to remove %s: %v", filepath, err)
		helper.SendResponse(c, http.StatusInternalServerError, "", nil)
		c.Abort()
		return
	}

	helper.SendResponse(c, http.StatusOK, "Template deleted successfully!", nil)
}
