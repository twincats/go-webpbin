package webpbin

import (
	"image/gif"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	DetectUnsupportedPlatforms()
	downloadFile("https://upload.wikimedia.org/wikipedia/commons/thumb/b/b7/CanadaLakeLouise_Bird.gif/240px-CanadaLakeLouise_Bird.gif", "source.gif")
}

func TestGifEncodeImage(t *testing.T) {
	c := NewGif2WebP()
	f, err := os.Open("source.gif")
	assert.Nil(t, err)
	img, err := gif.DecodeAll(f)
	assert.Nil(t, err)
	c.InputImage(img)
	c.OutputFile("target.webp")
	err = c.Run()
	assert.Nil(t, err)
	validateGifWebp(t)
}

func TestGifEncodeReader(t *testing.T) {
	c := NewGif2WebP()
	f, err := os.Open("source.gif")
	assert.Nil(t, err)
	c.Input(f)
	c.OutputFile("target.webp")
	err = c.Run()
	assert.Nil(t, err)
	validateGifWebp(t)
}

func TestGifEncodeFile(t *testing.T) {
	c := NewGif2WebP()
	c.InputFile("source.gif")
	c.OutputFile("target.webp")
	err := c.Run()
	assert.Nil(t, err)
	validateGifWebp(t)
}

func TestGifEncodeWriter(t *testing.T) {
	f, err := os.Create("target.webp")
	assert.Nil(t, err)
	defer f.Close()

	c := NewGif2WebP()
	c.InputFile("source.gif")
	c.Output(f)
	err = c.Run()
	assert.Nil(t, err)
	f.Close()
	validateGifWebp(t)
}

func TestGifVersionCWebP(t *testing.T) {
	c := NewGif2WebP()
	_, err := c.Version()
	assert.Nil(t, err)
}

func validateGifWebp(t *testing.T) {
	defer os.Remove("target.webp")
	// We don't decode animated webp for check since it's implemented for now.
	// fSource, err := os.Open("source.gif")
	// assert.Nil(t, err)
	// imgSource, err := gif.Decode(fSource)
	// assert.Nil(t, err)
	// fTarget, err := os.Open("target.webp")
	// assert.Nil(t, err)
	// defer fTarget.Close()
	// imgTarget, err := webp.Decode(fTarget)
	// assert.Nil(t, err)
	// assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}
