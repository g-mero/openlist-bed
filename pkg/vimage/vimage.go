package vimage

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"path"

	"github.com/cshum/vipsgen/vips"
)

// ImageFormat represents supported image formats
type ImageFormat int

const (
	FormatJPEG ImageFormat = iota + 1
	FormatPNG
	FormatGIF
	FormatWEBP
)

// Image wraps image data with metadata
type Image struct {
	originalData   []byte
	originalFormat ImageFormat
	width          int
	height         int
	FileName       string
}

// TransferOption defines compression settings
type TransferOption struct {
	// Quality factor (1-100), default: 85
	Quality int
	// Target format, if 0, keeps original format
	TargetFormat ImageFormat
}

// DefaultCompressOptions returns default compression settings
func DefaultCompressOptions() *TransferOption {
	return &TransferOption{
		Quality: 85,
	}
}

// LoadFromFile loads image from multipart file
func LoadFromFile(file *multipart.FileHeader) (*Image, error) {
	data, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = data.Close() }()

	buf, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	filename := file.Filename[:len(file.Filename)-len(path.Ext(file.Filename))]
	return LoadFromBuffer(buf, filename)
}

// LoadFromBuffer loads image from byte buffer
func LoadFromBuffer(buf []byte, filename string) (*Image, error) {
	if len(buf) == 0 {
		return nil, errors.New("empty image buffer")
	}

	// Validate and detect format using vips
	format, w, h, err := getInfoFromImgBuf(buf)
	if err != nil {
		return nil, err
	}

	return &Image{
		originalData:   buf,
		originalFormat: format,
		width:          w,
		height:         h,
		FileName:       filename,
	}, nil
}

// getInfoFromImgBuf detects image format using vips
func getInfoFromImgBuf(buf []byte) (ImageFormat, int, int, error) {
	// Use NewImageFromBuffer which is more reliable than NewImageFromSource
	img, err := vips.NewImageFromBuffer(buf, nil)
	if err != nil {
		return 0, 0, 0, err
	}
	defer img.Close()

	// Get format from vips
	vipsFormat := img.Format()
	width := img.Width()
	height := img.Height()

	// Map vips format to our ImageFormat
	switch vipsFormat {
	case vips.ImageTypeJpeg:
		return FormatJPEG, width, height, nil
	case vips.ImageTypePng:
		return FormatPNG, width, height, nil
	case vips.ImageTypeGif:
		return FormatGIF, width, height, nil
	case vips.ImageTypeWebp:
		return FormatWEBP, width, height, nil
	default:
		return 0, 0, 0, errors.New("unsupported image format")
	}
}

// Transfer compresses the image and returns the compressed buffer
func (img *Image) Transfer(opts *TransferOption) ([]byte, error) {
	if opts == nil {
		opts = DefaultCompressOptions()
	}

	targetFormat := img.originalFormat
	if opts.TargetFormat != 0 {
		targetFormat = opts.TargetFormat
	}

	if opts.Quality < 1 || opts.Quality > 100 {
		opts.Quality = 100
	}

	// For animated formats, load all frames
	loadOpts := &vips.LoadOptions{}
	if img.originalFormat == FormatGIF || img.originalFormat == FormatWEBP {
		loadOpts.N = -1 // Load all animation frames
	}

	// Load image from original data
	vipsImg, err := vips.NewImageFromBuffer(img.originalData, loadOpts)
	if err != nil {
		return nil, err
	}
	defer vipsImg.Close()

	// Get page height for animation support
	pageHeight := vipsImg.PageHeight()
	if pageHeight == 0 {
		pageHeight = vipsImg.Height()
	}

	// Save in target format
	switch targetFormat {
	case FormatJPEG:
		return vipsImg.JpegsaveBuffer(&vips.JpegsaveBufferOptions{
			Q:                  opts.Quality,
			OptimizeCoding:     true,
			Interlace:          true,
			TrellisQuant:       true,
			OvershootDeringing: true,
			OptimizeScans:      true,
			QuantTable:         3,
			Keep:               vips.KeepNone,
		})
	case FormatPNG:
		return vipsImg.PngsaveBuffer(&vips.PngsaveBufferOptions{
			Compression: 9,            // 0-9
			Palette:     true,         // use palette can reduce size effectively
			Q:           opts.Quality, // only available when palette is true
			Bitdepth:    8,
			Keep:        vips.KeepNone,
		})
	case FormatWEBP:
		return vipsImg.WebpsaveBuffer(&vips.WebpsaveBufferOptions{
			Q:          opts.Quality,
			PageHeight: pageHeight, // Preserve animation
			Keep:       vips.KeepNone,
		})
	case FormatGIF:
		return vipsImg.GifsaveBuffer(&vips.GifsaveBufferOptions{
			Dither:               1.0,  // Dithering for better quality
			Effort:               7,    // CPU effort (1-10, higher = better compression)
			Bitdepth:             8,    // Color depth
			InterframeMaxerror:   0,    // Frame optimization
			Reuse:                true, // Reuse palette between frames
			InterpaletteMaxerror: 3.0,  // Palette optimization
			Keep:                 vips.KeepNone,
			PageHeight:           pageHeight, // Preserve animation
		})
	default:
		return nil, errors.New("unsupported target format")
	}
}

func (img *Image) SmartCompress(allowWebp bool) ([]byte, ImageFormat, error) {
	opt := &TransferOption{
		Quality:      80,
		TargetFormat: 0,
	}

	// compress gif is not effective, return original data
	if img.originalFormat == FormatGIF && !allowWebp {
		return img.originalData, FormatGIF, nil
	}

	if allowWebp {
		opt.TargetFormat = FormatWEBP
	}

	// if webp is not allowed, convert webp to gif
	if !allowWebp && img.originalFormat == FormatWEBP {
		opt.TargetFormat = FormatGIF
	}

	compressed, err := img.Transfer(opt)
	if err != nil {
		return nil, 0, err
	}

	targetFormat := img.originalFormat
	if opt.TargetFormat != 0 {
		targetFormat = opt.TargetFormat
	}

	return compressed, targetFormat, nil
}

// OriginalData returns the original image data
func (img *Image) OriginalData() []byte {
	return img.originalData
}

// OriginalReader returns a reader for the original image data
func (img *Image) OriginalReader() io.Reader {
	return bytes.NewReader(img.originalData)
}

// OriginalFormat returns the detected original format
func (img *Image) OriginalFormat() ImageFormat {
	return img.originalFormat
}

// FullName returns filename with original extension
func (img *Image) FullName() string {
	return img.FileName + img.extensionForFormat(img.originalFormat)
}

// FullNameWithFormat returns filename with specified format extension
func (img *Image) FullNameWithFormat(format ImageFormat) string {
	return img.FileName + img.extensionForFormat(format)
}

// extensionForFormat returns file extension for given format
func (img *Image) extensionForFormat(format ImageFormat) string {
	switch format {
	case FormatJPEG:
		return ".jpg"
	case FormatPNG:
		return ".png"
	case FormatGIF:
		return ".gif"
	case FormatWEBP:
		return ".webp"
	default:
		return ""
	}
}

// ContentType returns MIME type for original format
func (img *Image) ContentType() string {
	return img.contentTypeForFormat(img.originalFormat)
}

// ContentTypeForFormat returns MIME type for given format
func (img *Image) ContentTypeForFormat(format ImageFormat) string {
	return img.contentTypeForFormat(format)
}

// contentTypeForFormat returns MIME type for given format
func (img *Image) contentTypeForFormat(format ImageFormat) string {
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

// Dimensions returns image width and height
func (img *Image) Dimensions() (width, height int) {

	return img.width, img.height
}

// Width returns image width
func (img *Image) Width() int {
	return img.width
}

// Height returns image height
func (img *Image) Height() int {
	return img.height
}
