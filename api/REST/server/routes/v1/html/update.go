package html

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"reportgenengine/constant"
	"reportgenengine/helper"
	"reportgenengine/log"

	"github.com/gin-gonic/gin"
)

// Handle old template update
func HandleTemplateUpdate(c *gin.Context) {
	templId := c.Param("id") // provided template ID
	ip := c.ClientIP()

	f, _ := c.FormFile(constant.TEMPLATE_UPLD_KEY)
	filename := f.Filename

	filename = fmt.Sprintf("%s.html", templId)
	filepath := path.Join(constant.TEMPLATE_PATH, ip, filename) // Dir creation and all is handled in middleware rules

	file, _ := f.Open()
	defer file.Close()

	b, _ := io.ReadAll(file)

	if err := os.WriteFile(filepath, b, 0666); err != nil {
		log.Errorf("Failed to save uploaded file at %s", filepath)
		helper.SendResponse(c, http.StatusInternalServerError, "Upload failed!", nil)
		c.Abort()
		return
	}

	helper.SendResponse(c, http.StatusOK, "Old template updated with new.", nil)
}
