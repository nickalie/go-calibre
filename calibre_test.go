package calibre

import (
	"testing"
	"os"
	"fmt"
	"net/http"
	"io"
	"github.com/stretchr/testify/assert"
)

func init() {
	downloadFile("https://www.gutenberg.org/ebooks/2600.epub.images?session_id=6326d908280f40b489a0b3be7a2653349aa8774d", "source.epub")
	downloadFile("https://royallib.com/get/fb2/tolstoy_lev/voyna_i_mir_kniga_1.zip", "source.zip")
}

func downloadFile(url, target string) {
	_, err := os.Stat(target)

	if err != nil {
		resp, err := http.Get(url)

		if err != nil {
			fmt.Printf("Error while downloading test image: %v\n", err)
			panic(err)
		}

		defer resp.Body.Close()

		f, err := os.Create(target)

		if err != nil {
			panic(err)
		}

		defer f.Close()

		_, err = io.Copy(f, resp.Body)

		if err != nil {
			panic(err)
		}
	}
}

func TestConvertEpub(t *testing.T) {
	err := Convert("source.epub", "target.mobi")
	assert.Nil(t, err)
}

func TestConvertFB2(t *testing.T) {
	err := Convert("source.zip", "target.mobi")
	assert.Nil(t, err)
}
