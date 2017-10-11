package main

const (
	txt  = "txt"
	html = "html"
	css  = "css"
	js   = "js"
	jpg  = "jpg"
	jpeg = "jpeg"
	gif  = "gif"
	png  = "png"
	swf  = "swf"
)

func GetContentType(ext string) string {
	var contentType string
	switch ext {
	case txt:
		contentType = "text/plain"
	case html:
		contentType = "text/html"
	case css:
		contentType = "text/css"
	case js:
		contentType = "text/javascript"
	case jpg:
		contentType = "image/jpeg"
	case jpeg:
		contentType = "image/jpeg"
	case gif:
		contentType = "image/gif"
	case png:
		contentType = "image/png"
	case swf:
		contentType = "application/x-shockwave-flash"
	}
	return contentType
}
