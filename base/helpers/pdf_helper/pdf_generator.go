package pdf_helper

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PdfGenerator interface {
	GenerateHTML(htmlContent string, footerContent string, writer io.Writer) error
	RenderTemplateAsset(assetName string, data interface{}) (string, error)
}

type PdfGeneratorImpl struct{}

/*
NewPdfGenerator returns a new instance of PdfGenerator, which is implemented by PdfGeneratorImpl.
*/
func NewPdfGenerator() PdfGenerator {
	return &PdfGeneratorImpl{}
}

/*
GenerateHTML generates a PDF file from the provided HTML content, with a custom footer, and saves it to the specified output file path.

Parameters:
- htmlContent: The HTML content to be converted to PDF.
- footerContent: The content to be displayed in the footer of each page.
- outputFilePath: The file path where the generated PDF will be saved.

Returns:
- An error if any issue occurs during PDF generation or file writing, nil otherwise.
*/
func (generator *PdfGeneratorImpl) GenerateHTML(htmlContent string, footerContent string, writer io.Writer) error {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
	}

	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeLetter)
	pdfg.MarginTop.Set(12)
	pdfg.MarginLeft.Set(12)
	pdfg.MarginBottom.Set(12)
	pdfg.MarginRight.Set(12)

	page := wkhtmltopdf.NewPageReader(strings.NewReader(htmlContent))

	page.Encoding.Set("UTF-8") // Ensure UTF-8 encoding

	page.DisableSmartShrinking.Set(true)
	page.FooterRight.Set(footerContent)
	page.FooterFontName.Set("Arial")
	page.FooterLine.Set(true)
	page.FooterFontSize.Set(6)

	pdfg.AddPage(page)

	// Set the output destination to the file
	pdfg.SetOutput(writer)

	return pdfg.Create()
}

/*
RenderTemplateAsset renders a PDF template using the specified asset name and data.

Parameters:
- assetName: the name of the asset containing the HTML template to render.
- data: the data to be injected into the template.

Returns:
- string: the rendered PDF content as a string.
- error: an error if any occurred during the rendering process.
*/
func (generator *PdfGeneratorImpl) RenderTemplateAsset(assetName string, data interface{}) (string, error) {
	templateData, err := os.ReadFile(fmt.Sprintf("assets/pdf_template/%s.html", assetName))
	if err != nil {
		return "", err
	}

	t, err := template.New("pdfTemplate").Funcs(template.FuncMap{
		"amount": amount,
	}).Parse(string(templateData))

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
