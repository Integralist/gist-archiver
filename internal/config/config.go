package config

const (
	APIBaseURL     = "https://api.github.com"
	GistsPerPage   = 100
	GithubUsername = "integralist" // Replace if needed

	WebDeployDir    = "."
	AssetsSubDir    = "assets"
	CssSubDir       = "css"
	GistsHtmlSubDir = "gists"

	MarkdownOutputDir = "markdown_files"
	GistsMdSubDir     = "gists"

	CssFileName          = "styles.css"
	MaxConcurrentFetches = 10

	JsSubDir            = "js"
	JsFilterFileName    = "filter.js"
	JsTagsFileName      = "tags.js"
	JsHighlightFileName = "highlight.min.js"
	JsNotesFileName     = "notes.js"
)
