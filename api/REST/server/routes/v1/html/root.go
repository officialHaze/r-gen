package html

import (
	"reportgenengine/api/REST/server/routes/v1/hooks"

	"github.com/gin-gonic/gin"
)

func MapRoutes(r *gin.RouterGroup) {
	template := r.Group("/template")

	// Upload HTML template
	uploadgrp := template.Group("/upload")
	hooks.FileUploadVerify(uploadgrp)
	uploadgrp.POST("/", HandleTemplateUpload)

	// Replace/Edit a HTML template by uploading
	// a new template and providing a previously
	// generated ID
	updategrp := template.Group("/update")
	hooks.FileUploadVerify(updategrp)
	updategrp.PUT("/:id", HandleTemplateUpdate)

	// List all templates
	template.GET("/list", HandleListingTempl)

	// Inject data into a HTML template
	template.POST("/inject/:id", HandleInjection)

	// Preview a HTML template
	template.GET("/preview/:id", PreviewTemplate)

	// Delete a template
	template.DELETE("/delete/:id", HandleTemplDelete)
}
