package s3

import "strings"

func parseExt(a string) string {
	switch {
	case strings.Contains(a, "/json"):
		return ".json"
	case strings.Contains(a, "/html"):
		return ".html"
	case strings.Contains(a, "/jpeg"):
		return ".jpg"
	case strings.Contains(a, "/png"):
		return ".png"
	case strings.Contains(a, "/gif"):
		return ".gif"
	case strings.Contains(a, "/bmp"):
		return ".bmp"
	case strings.Contains(a, "/tiff"):
		return ".tiff"
	case strings.Contains(a, "/plain"):
		return ".txt"
	case strings.Contains(a, "/rtf"):
		return ".rtf"
	case strings.Contains(a, "/msword"):
		return ".doc"
	case strings.Contains(a, "/zip"):
		return ".zip"
	case strings.Contains(a, "/mpeg"):
		return ".mp4"
	case strings.Contains(a, "/pdf"):
		return ".pdf"
	case strings.Contains(a, "/stylesheet"):
		return ".css"
	case strings.Contains(a, "/javascript"):
		return ".js"
	default:
		return ""
	}
}
