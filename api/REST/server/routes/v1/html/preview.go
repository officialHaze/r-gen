package html

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"reportgenengine/constant"
	"reportgenengine/helper"
	"reportgenengine/log"
	"time"

	"github.com/gin-gonic/gin"
)

func PreviewTemplate(c *gin.Context) {
	templateId := c.Param("id")
	ip := c.ClientIP()

	filename := fmt.Sprintf("%s.html", templateId)
	path := path.Join(constant.TEMPLATE_PATH, ip, filename)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Errorf("File at %s does not exist! %v", path, err)
			helper.SendResponse(c, http.StatusInternalServerError, "File does not exist!", nil)
			c.Abort()
			return
		}

		// Some other error (perms, etc.)
		log.Errorf("%v", err)
		helper.SendResponse(c, http.StatusInternalServerError, "Error generating preview!", nil)
		c.Abort()
		return
	}

	c.File(path)
}

type TemplInfo struct {
	ID         string    `json:"id"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Size       int64     `json:"size"`
}

// Handle listing templates
func HandleListingTempl(c *gin.Context) {
	ip := c.ClientIP()

	dir := path.Join(constant.TEMPLATE_PATH, ip)

	infos := helper.ListFileInfo(dir)

	helper.SendResponse(c, http.StatusOK, "", infos)
}
