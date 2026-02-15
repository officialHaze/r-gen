package interfaces

type PDFMergeSchema struct {
	Path  string
	Pages []string
}

type PDFMergePayload []PDFMergeSchema
