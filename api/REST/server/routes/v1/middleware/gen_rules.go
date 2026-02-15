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

// ========= PDF Gen Rules ===========
// Rule: Check if max gen quota reached
func GenQuotaRule(c *gin.Context) {
	ip := c.ClientIP()
	if ip == "" {
		log.Errorln("Empty IP detected!")
		helper.SendResponse(c, http.StatusBadRequest, "Cannot detect any IP", nil)
		c.Abort()
		return
	}

	dir := path.Join(constant.PDF_PATH, ip)
	totalsize := calcTotalSize(dir)

	if quotaReached(totalsize) {
		log.Errorln("Gen quota reached!")
		helper.SendResponse(c, http.StatusForbidden, "Gen quota reached!", nil)
		c.Abort()
		return
	}

	log.Successf("Gen quota within limit!")
	c.Next()
}

// ========= Helper methods ==========
// This method recursively reads a dir
// calculates the total size combining
// all dir entries and returns the value
func calcTotalSize(dir string) int64 {
	entries, _ := os.ReadDir(dir)

	var totalsize int64 = 0
	for _, e := range entries {
		// If current entry is a dir
		// read it recursively and fetch
		// the total size
		if e.IsDir() {
			path := path.Join(dir, e.Name())
			totalsize += calcTotalSize(path)
			continue
		}
		i, _ := e.Info()
		totalsize += i.Size()
	}

	return totalsize
}

func quotaReached(size int64) bool {
	quota := util.DecodeSize(settings.MySettings.GEN_QUOTA)
	return size >= quota
}
