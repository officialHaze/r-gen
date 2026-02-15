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

func HandleInjection(c *gin.Context) {
	templateId := c.Param("id")
	ip := c.ClientIP()

	var payload map[string]any
	if err := c.BindJSON(&payload); err != nil {
		log.Errorf("%v", err)
		c.Abort()
		return
	}

	reportgenengine := helper.ReportGenEngineInit()
	injectedtemplid := reportgenengine.InjectedtemplateId(templateId)

	templpath := path.Join(constant.TEMPLATE_PATH, ip, fmt.Sprintf("%s.html", templateId))
	injectedtemplpath := path.Join(constant.TEMPLATE_PATH, ip, fmt.Sprintf("%s.html", injectedtemplid))

	// Inject data
	err := reportgenengine.InjectData(templpath, injectedtemplpath, payload)
	if err != nil {
		log.Errorf("%v", err)
		helper.SendResponse(c, http.StatusInternalServerError, "Server error! Please try after sometime.", nil)
		c.Abort()
		return
	}

	helper.SendResponse(c, http.StatusOK, "Data Injected!", injectedtemplid)
}
