package calibre

import (
	"github.com/nickalie/go-binwrapper"
	"errors"
)

type MobiType string

const MOBI_TYPE_OLD = MobiType("old")
const MOBI_TYPE_BOTH = MobiType("both")
const MOBI_TYPE_NEW = MobiType("new")

func Convert(source, target string) error {
	return NewConvertWrapper().
		Source(source).
		Target(target).
		Run()
}

type ConvertWrapper struct {
	*binwrapper.BinWrapper
	source    string
	target    string
	mobiType  MobiType
	landscape bool
	grayscale bool
	keepAspectRatio bool
}

func (c *ConvertWrapper) Source(source string) *ConvertWrapper {
	c.source = source
	return c
}

func (c *ConvertWrapper) Target(target string) *ConvertWrapper {
	c.target = target
	return c
}

func (c *ConvertWrapper) Landscape(landscape bool) *ConvertWrapper {
	c.landscape = landscape
	return c
}

func (c *ConvertWrapper) Grayscale(grayscale bool) *ConvertWrapper {
	c.grayscale = grayscale
	return c
}

func (c *ConvertWrapper) KeepAspectRatio(keepAspectRatio bool) *ConvertWrapper {
	c.keepAspectRatio = keepAspectRatio
	return c
}

func (c *ConvertWrapper) MobiType(value MobiType) *ConvertWrapper {
	c.mobiType = value
	return c
}

func (c *ConvertWrapper) Run() error {
	defer c.BinWrapper.Reset()
	c.BinWrapper.Debug()
	c.Arg(c.source, c.target)

	if c.landscape {
		c.Arg("--landscape")
	}

	if !c.grayscale {
		c.Arg("--dont-grayscale")
	}

	if c.mobiType != "" {
		c.Arg("--mobi-file-type", string(c.mobiType))
	}

	if c.keepAspectRatio {
		c.Arg("--keep-aspect-ratio")
	}

	err := c.BinWrapper.Run()

	if err != nil {
		return errors.New(string(c.CombinedOutput()) + "\n" + err.Error())
	} else {
		return nil
	}
}

func NewConvertWrapper() *ConvertWrapper {
	c := &ConvertWrapper{
		BinWrapper: binwrapper.NewBinWrapper().ExecPath("ebook-convert"),
		grayscale:  true}
	return c
}
