package pdf

import (
	"reportgenengine/api/REST/server/routes/v1/hooks"

	"github.com/gin-gonic/gin"
)

func MapRoutes(r *gin.RouterGroup) {
	// Merge multiple pdfs by template ID
	r.POST("/merge", HandlePDFMerge)

	// Download a PDF by ID
	r.GET("/download/:id", HandlePDFDownload)

	// Preview a PDF by ID
	// In case of preview the file is not removed after downloading
	r.GET("/preview/:id", HandlePDFPreview)

	// List all PDFs
	r.GET("/list", HandleListingPDF)

	// Delete a pdf
	r.DELETE("/delete/:id", HandlePDFDelete)

	// Gen route group
	gengroup := r.Group("/generate")
	hooks.PDFGenVerify(gengroup)

	// Generate pdf from a HTML template
	gengroup.GET("/:id", HandlePDFGen)

	// Generate standalone pdf by uploading
	// real-time HTML template
	standalonegrp := gengroup.Group("/standlone")
	hooks.FileUploadVerify(standalonegrp)
	standalonegrp.POST("/", HandleStandalonePDFGen)
}
