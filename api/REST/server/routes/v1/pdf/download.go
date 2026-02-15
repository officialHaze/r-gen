package pdf

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"reportgenengine/constant"
	"reportgenengine/helper"
	"reportgenengine/log"

	"github.com/gin-gonic/gin"
)

// Handle PDF download
func HandlePDFDownload(c *gin.Context) {
	sendPDF(c, true)
}

// Handle PDF preview
func HandlePDFPreview(c *gin.Context) {
	sendPDF(c, false)
}

// Send PDF file in response
func sendPDF(c *gin.Context, download bool) {
	name := c.Query("filename")
	id := c.Param("id") // PDF ID
	ip := c.ClientIP()

	filename := fmt.Sprintf("%s.pdf", id)
	filepath := path.Join(constant.PDF_PATH, ip, filename)

	// Verify if file exists
	info, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			helper.SendResponse(c, http.StatusBadRequest, "PDF for the given ID does not exist!", nil)
			return
		}
		// Some other kind or error (perms, etc.)
		// In case of unexpected error, log it
		log.Errorf("%v", err)
		helper.SendResponse(c, http.StatusInternalServerError, "", nil)
		return
	}

	if download {
		// Send file as an attachment
		if name == "" {
			name = info.Name()
		} else {
			name = fmt.Sprintf("%s.pdf", name)
		}
		c.FileAttachment(filepath, name)
		// Once downloaded, remove the file
		err = helper.ReportGenEngineInit().RemoveFile(filepath)
		if err != nil {
			log.Errorf("Failed to remove %s: %v", filepath, err)
		}
		return
	}

	// Else, send the file but do not remove it
	c.File(filepath)
}
