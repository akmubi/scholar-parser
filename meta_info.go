package scholarScraper

type HTMLMetaInfo struct {
	Authors		[]string
	WebSite		string
	Name		string
	URL			string
}

// func ExtractMetaHTML(article Article) (pdfMeta PDFMetaInfo, err error) {

// }

type PDFMetaInfo struct {
	Authors		[]string
	Conclusion	string
	Abstract	string
	Name		string
	URL			string
	PageCount	int
	Number		int
}

// func ExtractMetaPDF(article Article) (pdfMeta PDFMetaInfo, err error) {

// }

