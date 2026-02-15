/*
Main engine to generate PDF report from HTML templates
*/

package helper

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"reportgenengine/constant"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type ReportGenEngine struct {
	templatesBaseDir string
	reportBaseDir    string
}

/*
Method to populate HTML template.

Args:

	templpath -> HTML Template path
	injectedtemplpath -> Path where the injected template will be saved
	payload -> Data to be injected

Return:

	error -> any error encountered during process
*/
func (r *ReportGenEngine) InjectData(templpath, injectedtemplpath string, payload map[string]any) error {
	templ, err := template.ParseFiles(templpath)
	if err != nil {
		return err
	}

	// Create a writer for the HTML template inside the directory
	f, err := os.OpenFile(injectedtemplpath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = templ.Execute(f, payload) // Write the HTML template with provided payload
	if err != nil {
		r.RemoveFile(injectedtemplpath) // Remove the file if any error is encountered while writing
		return err
	}

	return nil
}

/*
Generate PDF using HTML template

Args:

	templpath -> Path of the HTML template
	pdfpath -> Path where the PDF will be saved

Return:

	Error -> Any error encountered during the process
*/
func (r *ReportGenEngine) GeneratePdf(templpath, pdfpath string) error {
	f, err := os.Open(templpath)
	if err != nil {
		return fmt.Errorf("Failed to open - %s: %w", templpath, err)
	}
	defer f.Close()

	b, _ := io.ReadAll(f)
	html := string(b) // convert to string

	// Update default alloc opts executed while running chromedp
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium-browser"),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Start chromedp with updated alloc context (headless mode)
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Run chromedp
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),     // Navigate to blank page
		chromedp.EmulateViewport(1240, 1754), // Set a fixed viewport
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx) // Inject html directly into blank DOM
		}),
		// Force PRINT media
		chromedp.ActionFunc(func(ctx context.Context) error {
			return emulation.SetEmulatedMedia().
				WithMedia("print").
				Do(ctx)
		}),
		// Allow layout to settle
		chromedp.Sleep(800*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			b, _, err := page.
				PrintToPDF().
				WithPrintBackground(true).
				WithDisplayHeaderFooter(false). // handle header/footer in HTML
				WithPaperWidth(8.27).           // A4 width (inches)
				WithPaperHeight(11.69).         // A4 height
				WithMarginTop(1.1).
				WithMarginBottom(1.1).
				WithMarginLeft(0.8).
				WithMarginRight(0.8).
				WithPreferCSSPageSize(true).
				Do(ctx) // Convert the current page into pdf data bytes
			if err != nil {
				return err
			}
			return os.WriteFile(pdfpath, b, 0644) // Save the pdf
		}),
	); err != nil {
		return err
	}

	return nil
}

// Merge multiple pdfs
func (r *ReportGenEngine) MergePdf(payloads []string, out string) error {
	if err := api.MergeCreateFile(payloads, out, false, nil); err != nil {
		return err
	}

	return nil
}

// Helper method to trim a given PDF into selected pages
// and save it in dst
func (r *ReportGenEngine) TrimPdf(in, out string, pages []string) error {
	if err := api.TrimFile(in, out, pages, nil); err != nil {
		return err
	}

	return nil
}

/*
Delete a file by its path

Args:

	path -> full path

Return:

	error -> any error encountered while deleting
*/
func (r *ReportGenEngine) RemoveFile(path string) error {
	return os.Remove(path)
}

// Globally handle creation of injected HTML template id
func (r *ReportGenEngine) InjectedtemplateId(id string) string {
	return fmt.Sprintf("injected_%s", id)
}

/*
=============================
Engine initializer
=============================
*/
func ReportGenEngineInit() *ReportGenEngine {
	return &ReportGenEngine{
		templatesBaseDir: constant.TEMPLATE_PATH,
		reportBaseDir:    constant.PDF_PATH,
	}
}
