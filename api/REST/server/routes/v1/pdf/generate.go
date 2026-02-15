package pdf

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"reportgenengine/constant"
	"reportgenengine/helper"
	"reportgenengine/log"
	"reportgenengine/util"

	"github.com/gin-gonic/gin"
)

// Handle PDF generation from HTML template
func HandlePDFGen(c *gin.Context) {
	templid := c.Param("id") // template ID
	ip := c.ClientIP()

	reportgenEngine := helper.ReportGenEngineInit()

	templname := fmt.Sprintf("%s.html", templid)
	pdfname := fmt.Sprintf("%s.pdf", util.UUIDGen())

	templpath := path.Join(constant.TEMPLATE_PATH, ip, templname)
	pdfpath := path.Join(constant.PDF_PATH, ip, pdfname)

	// Make sure necessary dirs are created
	if err := os.MkdirAll(path.Join(constant.PDF_PATH, ip), 0777); err != nil {
		log.Errorf("%v", err)
		helper.SendResponse(c, http.StatusInternalServerError, "", nil)
		c.Abort()
		return
	}

	err := reportgenEngine.GeneratePdf(templpath, pdfpath)
	if err != nil {
		log.Errorf("Failed to generate PDF: %v", err)
		helper.SendResponse(c, http.StatusInternalServerError, "Failed to generate PDF!", nil)
		c.Abort()

		// Remove any incomplete written file
		reportgenEngine.RemoveFile(pdfpath)

		return
	}

	// Return the id of the generated PDF
	helper.SendResponse(c, http.StatusCreated, "PDF generated!", helper.RmFileExt(pdfname))
}

// Handle standlone PDF generation using real-time
// uploaded HTML template
func HandleStandalonePDFGen(c *gin.Context) {
	f, _ := c.FormFile("template")
	ip := c.ClientIP()

	// Save the template in template dir under a tmp name
	tmpname := fmt.Sprintf("%s.html", util.UUIDGen())
	tmptemplpath := path.Join(constant.TEMPLATE_PATH, ip, tmpname)
	err := c.SaveUploadedFile(f, tmptemplpath)
	if err != nil {
		log.Errorf("Failed to save uploaded file at %s: %v", tmptemplpath, err)
		helper.SendResponse(c, http.StatusInternalServerError, "", nil)
		return
	}

	reportgenEngine := helper.ReportGenEngineInit()

	pdfname := fmt.Sprintf("%s.pdf", util.UUIDGen())
	pdfpath := path.Join(constant.PDF_PATH, ip, pdfname)

	// Make sure necessary dirs are created
	if err := os.MkdirAll(path.Join(constant.PDF_PATH, ip), 0777); err != nil {
		log.Errorf("%v", err)
		helper.SendResponse(c, http.StatusInternalServerError, "", nil)
		c.Abort()
		return
	}

	err = reportgenEngine.GeneratePdf(tmptemplpath, pdfpath)
	if err != nil {
		log.Errorf("Failed to generate PDF: %v", err)
		helper.SendResponse(c, http.StatusInternalServerError, "Failed to generate PDF!", nil)
		c.Abort()

		// Remove any incomplete written file
		reportgenEngine.RemoveFile(pdfpath)
		// Remove the temporary html template
		reportgenEngine.RemoveFile(tmptemplpath)

		return
	}

	// Return the id of the generated PDF
	helper.SendResponse(c, http.StatusCreated, "PDF generated!", helper.RmFileExt(pdfname))

	// Remove the tmp html template
	reportgenEngine.RemoveFile(tmptemplpath)
}
