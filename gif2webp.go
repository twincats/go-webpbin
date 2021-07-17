package webpbin

import (
	"errors"
	"fmt"
	"image/gif"
	"io"

	"github.com/nickalie/go-binwrapper"
)

// Gif2WebP converts a GIF image to a WebP image.
// https://developers.google.com/speed/webp/docs/gif2webp
type Gif2WebP struct {
	*binwrapper.BinWrapper
	inputFile  string
	inputImage *gif.GIF
	input      io.Reader
	outputFile string
	output     io.Writer
	quality    int
}

// NewGif2WebP creates new Gif2WebP instance.
func NewGif2WebP(optionFuncs ...OptionFunc) *Gif2WebP {
	bin := &Gif2WebP{
		BinWrapper: createBinWrapper(optionFuncs...),
		quality:    -1,
	}
	bin.ExecPath("gif2webp")

	return bin
}

// Version returns gif2webp version.
func (c *Gif2WebP) Version() (string, error) {
	return version(c.BinWrapper)
}

// InputFile sets image file to convert.
// Input or InputImage called before will be ignored.
func (c *Gif2WebP) InputFile(file string) *Gif2WebP {
	c.input = nil
	c.inputImage = nil
	c.inputFile = file
	return c
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *Gif2WebP) Input(reader io.Reader) *Gif2WebP {
	c.inputFile = ""
	c.inputImage = nil
	c.input = reader
	return c
}

// InputImage sets image to convert.
// InputFile or Input called before will be ignored.
func (c *Gif2WebP) InputImage(img *gif.GIF) *Gif2WebP {
	c.inputFile = ""
	c.input = nil
	c.inputImage = img
	return c
}

// OutputFile specify the name of the output WebP file.
// Output called before will be ignored.
func (c *Gif2WebP) OutputFile(file string) *Gif2WebP {
	c.output = nil
	c.outputFile = file
	return c
}

// Output specify writer to write webp file content.
// OutputFile called before will be ignored.
func (c *Gif2WebP) Output(writer io.Writer) *Gif2WebP {
	c.outputFile = ""
	c.output = writer
	return c
}

// Quality specify the compression factor for RGB channels between 0 and 100. The default is 75.
//
// A small factor produces a smaller file with lower quality. Best quality is achieved by using a value of 100.
func (c *Gif2WebP) Quality(quality uint) *Gif2WebP {
	if quality > 100 {
		quality = 100
	}

	c.quality = int(quality)
	return c
}

// Run starts gif2webp with specified parameters.
func (c *Gif2WebP) Run() error {
	defer c.BinWrapper.Reset()

	if c.quality > -1 {
		c.Arg("-q", fmt.Sprintf("%d", c.quality))
	}

	output, err := c.getOutput()

	if err != nil {
		return err
	}

	c.Arg("-o", output)

	err = c.setInput()

	if err != nil {
		return err
	}

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	return nil
}

// Reset all parameters to default values
func (c *Gif2WebP) Reset() *Gif2WebP {
	c.quality = -1
	return c
}

func (c *Gif2WebP) setInput() error {
	if c.input != nil {
		c.Arg("--").Arg("-")
		c.StdIn(c.input)
	} else if c.inputImage != nil {
		r, err := createReaderFromGIF(c.inputImage)

		if err != nil {
			return err
		}

		c.Arg("--").Arg("-")
		c.StdIn(r)
	} else if c.inputFile != "" {
		c.Arg(c.inputFile)
	} else {
		return errors.New("Undefined input")
	}

	return nil
}

func (c *Gif2WebP) getOutput() (string, error) {
	if c.output != nil {
		return "-", nil
	} else if c.outputFile != "" {
		return c.outputFile, nil
	} else {
		return "", errors.New("Undefined output")
	}
}
