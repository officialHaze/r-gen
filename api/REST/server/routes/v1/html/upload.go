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
	"reportgenengine/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handle fresh template upload
func HandleTemplateUpload(c *gin.Context) {
	uuid := util.UUIDGen() // Every file will be assigned a UUID
	ip := c.ClientIP()

	f, _ := c.FormFile(constant.TEMPLATE_UPLD_KEY)
	filename := f.Filename

	filename = createUUIDFilename(filename, uuid)               // update the filename with uuid
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

	helper.SendResponse(c, http.StatusOK, "Uploaded!", uuid)
}

// ========= Helper methods ============
func createUUIDFilename(name, uuid string) string {
	parts := strings.Split(name, ".")
	ext := parts[len(parts)-1] // Extension of the file

	return fmt.Sprintf("%s.%s", uuid, ext)
}
