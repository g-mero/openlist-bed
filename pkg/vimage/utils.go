package vimage

func FormatToContentType(format ImageFormat) string {
	switch format {
	case FormatJPEG:
		return "image/jpeg"
	case FormatPNG:
		return "image/png"
	case FormatGIF:
		return "image/gif"
	case FormatWEBP:
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
