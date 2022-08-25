package parser

var (
	//Formats contains all supported file formats
	Formats = make(map[string]struct{})

	archive = []string{
		"tar", "zip", "zipx", "rar", "7z",
	}

	audio = []string{
		"aif", "flac", "m3u", "m4a",
		"mid", "mp3", "ogg", "wav", "wma",
	}

	document = []string{
		"bin", "csv", "obb",
		"doc", "docx", "log", "odt",
		"rtf", "txt", "pdf", "ppt",
		"pptx", "xml", "xls", "xlsx",
	}

	executable = []string{
		"apk", "app", "bat", "bin",
		"cmd", "exe",
	}

	image = []string{
		"jpg", "jpeg", "jpe", "jif", "jfif", "jfi",
		"png", "gif", "webp", "tiff", "tif",
		".ind", "indd", "indt", "psd", "ps",
		"raw", "arw", "cr", "rw2", "nrw", "k25",
		"svg", "svgz",
	}

	video = []string{
		"3gp", "asf", "avi", "flw",
		"m4v", "mov", "mp4", "mpeg", "wmv",
	}

	web = []string{
		"csr", "css", "js", "json", "jsp", "php",
	}
)

func init() {
	allFormats := make([]string, 0, 100)
	allFormats = append(allFormats, archive...)
	allFormats = append(allFormats, audio...)
	allFormats = append(allFormats, document...)
	allFormats = append(allFormats, executable...)
	allFormats = append(allFormats, image...)
	allFormats = append(allFormats, video...)
	allFormats = append(allFormats, web...)

	for _, v := range allFormats {
		Formats[v] = struct{}{}
	}
}

//IsResourceFormat ...
func IsResourceFormat(s string) bool {

	_, has := Formats[s]
	return has
}
