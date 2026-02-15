package pdf

import (
	"fmt"
	"net/http"
	"path"
	"reportgenengine/constant"
	"reportgenengine/helper"
	"reportgenengine/interfaces"
	"reportgenengine/log"
	"reportgenengine/settings"
	"reportgenengine/util"
	"strings"

	"github.com/gin-gonic/gin"
)

type payload struct {
	TemplateID string   `json:"templateID"`
	Pages      []string `json:"pages"`
}

// Handle PDF merge
func HandlePDFMerge(c *gin.Context) {
	ip := c.ClientIP()
	_, keepSingles := c.GetQuery("keepSingles") // This query decides whether to keep
	// the original single PDFs after merging is complete. By default they are deleted

	maxallowed := settings.MySettings.PDF_MAX_MERGE_ALLOWED
	payloads := make([]payload, maxallowed) // fixed length slice of template IDs
	if err := c.BindJSON(&payloads); err != nil {
		log.Errorf("%v", err)
		c.Abort()
		return
	}

	if len(payloads) > maxallowed {
		log.Errorf("Payload exceeds maximum PDF allowed to merge - %d", maxallowed)
		helper.SendResponse(c, http.StatusBadRequest, "Payload exceeds maximum PDF allowed to merge!", maxallowed)
		c.Abort()
		return
	}

	reportengine := helper.ReportGenEngineInit()

	outfilename := fmt.Sprintf("merged_%s.pdf", util.UUIDGen())
	outpath := path.Join(constant.PDF_PATH, ip, outfilename)

	// Generate pdf merge payload
	mergepayload := make(interfaces.PDFMergePayload, 0)
	for _, p := range payloads {
		path := pdfpathGen(ip, p.TemplateID)
		mergepayload = append(mergepayload, interfaces.PDFMergeSchema{
			Path:  path,
			Pages: p.Pages,
		})
	}

	// For each payload trim the pdf by selected pages
	// and generate a tmp pdf
	tmpPdfPaths := make([]string, 0)
	for _, p := range mergepayload {
		originalpdfpath := p.Path
		pagescombined := strings.Join(p.Pages, ",")
		pagescombined = strings.ReplaceAll(pagescombined, "*", "") // eg:- replace 1-* to 1- (valid pdfcpu syntax)
		// Re-split
		pages := strings.Split(pagescombined, ",")

		outdir := path.Join(constant.PDF_PATH, ip)
		tmpfile := fmt.Sprintf("tmp_%s.pdf", util.UUIDGen())
		tmpout := path.Join(outdir, tmpfile)

		if err := reportengine.TrimPdf(originalpdfpath, tmpout, pages); err != nil {
			log.Errorf("%v", err)
			continue
		}
		tmpPdfPaths = append(tmpPdfPaths, path.Join(outdir, tmpfile))
	}

	// No need to proceed if trimmed tmp pdf list is empty
	if len(tmpPdfPaths) <= 0 {
		helper.SendResponse(c, http.StatusBadRequest, "Nothing to merge!", nil)
		c.Abort()
		return
	}

	if err := reportengine.MergePdf(tmpPdfPaths, outpath); err != nil {
		log.Errorf("%v", err)
		helper.SendResponse(c, http.StatusInternalServerError, "Failed to merge PDFs!", nil)
		c.Abort()

		// Remove all un-necessary files
		reportengine.RemoveFile(outpath)            // any incomplete merge attempt
		removeAllTmpPdfs(reportengine, tmpPdfPaths) // all tmp pdf paths created

		return
	}

	// Return the ID of the newly created PDF
	helper.SendResponse(c, http.StatusCreated, "PDFs merged successfully!", helper.RmFileExt(outfilename))

	// Remove trimmed tmp pdfs
	removeAllTmpPdfs(reportengine, tmpPdfPaths)
	// If explicitly it is said to keep singles
	// then end without deleting the original single pdfs
	if keepSingles {
		return
	}

	// By default, original single pdfs will be deleted
	for _, p := range mergepayload {
		err := reportengine.RemoveFile(p.Path)
		if err != nil {
			log.Errorf("Failed to remove %s: %v", p.Path, err)
		}
	}
}

// ========= Helper methods ==========
func pdfpathGen(ip, templID string) string {
	filename := fmt.Sprintf("%s.pdf", templID)

	return path.Join(constant.PDF_PATH, ip, filename)
}

func removeAllTmpPdfs(re *helper.ReportGenEngine, tmpPaths []string) {
	// Remove all tmp pdfs
	for _, t := range tmpPaths {
		if err := re.RemoveFile(t); err != nil {
			log.Errorf("Failed to remove %s", t)
		}
	}
}
