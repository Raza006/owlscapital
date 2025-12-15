package assets

import _ "embed"

//go:embed supportBanner.png
var SupportBanner []byte

//go:embed owlsfooter.png
var OwlsFooter []byte

const (
	SupportBannerFilename = "supportBanner.png"
	OwlsFooterFilename    = "owlsfooter.png"
)
