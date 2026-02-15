package middleware

import (
	"net/http"
	"os"
	"path"
	"reportgenengine/constant"
	"reportgenengine/helper"
	"reportgenengine/log"
	"reportgenengine/settings"
	"reportgenengine/util"

	"github.com/gin-gonic/gin"
)

// ========== Verification rules ============
// Rule: Check if file was sent with form data
func FilePresentRule(c *gin.Context) {
	_, err := c.FormFile(constant.TEMPLATE_UPLD_KEY)
	if err != nil {
		log.Errorf("No file named %s found in context!", constant.TEMPLATE_UPLD_KEY)
		helper.SendResponse(
			c,
			http.StatusBadRequest,
			"File is required!",
			nil,
		)
		c.Abort()
		return
	}
	log.Successf("File named %s present in context!", constant.TEMPLATE_UPLD_KEY)
	c.Next()
}

// Rule: Check if file size is a valid as per defined threshold
func ValidSizeRule(c *gin.Context) {
	f, _ := c.FormFile(constant.TEMPLATE_UPLD_KEY)
	filesize := f.Size
	thres := util.DecodeSize(settings.MySettings.UPLOAD_MAX_SIZE)
	if thres < 0 {
		log.Errorln("Invalid Threshold size. Please check app settings.")
		helper.SendResponse(
			c,
			http.StatusInternalServerError,
			"Invalid file upload threshold size!",
			nil,
		)
		c.Abort()
		return
	}
	log.Printf("Upload Threshold %s", settings.MySettings.UPLOAD_MAX_SIZE)

	if !isValidSize(thres, filesize) {
		log.Errorf("File size exceeds threshold! Thresh: %d; Filesize: %d", thres, filesize)
		helper.SendResponse(
			c,
			http.StatusInternalServerError,
			"File size exceeds threshold!",
			nil,
		)
		c.Abort()
		return
	}
	log.Successf("File size within threshold! Thresh: %d; Filesize: %d", thres, filesize)
	c.Next()
}

// Rule: Check if max file upload reached for client
func UploadLimitRule(c *gin.Context) {
	// Every client will be limited to a fixed number of template uploads
	limit := settings.MySettings.UPLOAD_LIMIT
	ip := c.ClientIP()
	if ip == "" {
		log.Errorf("Empty IP detected!")
		helper.SendResponse(
			c,
			http.StatusBadRequest,
			"Unrecognized IP!",
			nil,
		)
		c.Abort()
		return
	}
	log.Printf("Client IP - %s", ip)

	dst := path.Join(constant.TEMPLATE_PATH, ip)

	// Create client's upload dir if not already created (just to make sure)
	if err := os.MkdirAll(dst, 0777); err != nil {
		log.Errorf("Failed to create dst at %s: %v", dst, err)
		helper.SendResponse(
			c,
			http.StatusInternalServerError,
			"Internal server error! Please try after sometime.",
			nil,
		)
		c.Abort()
		return
	}

	diritems := dirItems(dst)
	if uploadLimitReached(limit, diritems) {
		log.Errorf("Upload limit reached - %s! Exists: %d", ip, diritems)
		helper.SendResponse(c, http.StatusForbidden, "Upload limit reached!", diritems)
		c.Abort()
		return
	}

	log.Successf("Within upload limit - %s. Exists: %d", ip, diritems)
	c.Next()
}

// ============= Helper methods =================
func uploadLimitReached(limit, entries int) bool {
	return entries >= limit
}

func dirItems(path string) int {
	entries, _ := os.ReadDir(path)
	return len(entries)
}

func isValidSize(thresh, filesize int64) bool {
	return filesize <= thresh
}
