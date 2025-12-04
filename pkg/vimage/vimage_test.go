package vimage_test

import (
	"openlist-bed/pkg/vimage"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var testImages = []string{
	"test_data/human.jpg",
	"test_data/sunset.jpg",
	"test_data/words.png",
	"test_data/animate.gif",
	"test_data/animated-flower.webp",
	"test_data/monalisa.webp",
	"test_data/with-alpha.png",
	"test_data/ios.heic",
}

func saveOutputFile(t *testing.T, filename string, data []byte) {
	t.Helper()
	dir := path.Join("test_output", path.Dir(filename))
	err := os.MkdirAll(dir, os.ModePerm)
	require.NoError(t, err)
	err = os.WriteFile(path.Join("test_output", filename), data, 0644)
	require.NoError(t, err)
}

func getFilename(p string) string {
	return strings.Split(path.Base(p), ".")[0]
}

func TestTransferOnly(t *testing.T) {
	for _, imgPath := range testImages {
		data, err := os.ReadFile(imgPath)
		require.NoError(t, err)

		img, err := vimage.LoadFromBuffer(data, getFilename(imgPath))
		require.NoError(t, err)

		compressed, err := img.Transfer(&vimage.TransferOption{
			Quality: 80,
		})
		require.NoError(t, err)
		require.NotEmpty(t, compressed)
		t.Logf("%s, before: %d bytes, after: %d bytes, less: %g %%", img.FullName(), len(data), len(compressed),
			100.0*(1.0-float64(len(compressed))/float64(len(data))))

		saveOutputFile(t, "transfer/"+img.FullName(), compressed)
	}
}

func TestTransferToWebP(t *testing.T) {
	for _, imgPath := range testImages {
		data, err := os.ReadFile(imgPath)
		require.NoError(t, err)

		img, err := vimage.LoadFromBuffer(data, getFilename(imgPath))
		require.NoError(t, err)

		webp, err := img.Transfer(&vimage.TransferOption{
			TargetFormat: vimage.FormatWEBP,
			Quality:      80,
		})
		require.NoError(t, err)
		require.NotEmpty(t, webp)
		t.Logf("%s, before: %d bytes, after: %d bytes, less: %g %%", img.FullName(), len(data), len(webp),
			100.0*(1.0-float64(len(webp))/float64(len(data))))

		saveOutputFile(t, "webp/"+img.FullNameWithFormat(vimage.FormatWEBP), webp)
	}
}

func TestSmartCompress(t *testing.T) {
	for _, imgPath := range testImages {
		data, err := os.ReadFile(imgPath)
		require.NoError(t, err)

		img, err := vimage.LoadFromBuffer(data, getFilename(imgPath))
		require.NoError(t, err)

		compressed, format, err := img.SmartCompress(false)
		require.NoError(t, err)
		require.NotEmpty(t, compressed)
		t.Logf("%s, before: %d bytes, after: %d bytes, less: %g %%", img.FullName(), len(data), len(compressed),
			100.0*(1.0-float64(len(compressed))/float64(len(data))))
		saveOutputFile(t, "smart/"+img.FullNameWithFormat(format), compressed)
	}
}

func TestInvalidImage(t *testing.T) {
	_, err := vimage.LoadFromBuffer([]byte{0x00, 0x01}, "invalid")
	require.Error(t, err)
}
