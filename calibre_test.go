package calibre

import (
	"testing"
	"os"
	"fmt"
	"net/http"
	"io"
	"github.com/stretchr/testify/assert"
	"archive/zip"
	"path/filepath"
)

func init() {
	downloadFile("https://royallib.com/get/epub/tolstoy_lev/voyna_i_mir_tom_1.zip", "source1.zip", true)
	downloadFile("https://royallib.com/get/fb2/tolstoy_lev/voyna_i_mir_kniga_1.zip", "source2.zip", false)
}

func downloadFile(url, target string, extract bool) {
	_, err := os.Stat(target)

	if err == nil {
		if extract {
			extractEpub(target)
		}
		return
	}

	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Error while downloading test file: %v\n", err)
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

	if extract {
		extractEpub(target)
	}
}

func extractEpub(value string) {
	zipReader, err := zip.OpenReader(value)

	if err != nil {
		panic(err)
	}

	for _, v := range zipReader.File {
		if filepath.Ext(v.Name) != ".epub" {
			continue
		}

		zipToFile(v, "source.epub")
	}
}

func zipToFile(zipFile *zip.File, target string) {
	f, err := os.Create(target)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	zipReader, err := zipFile.Open()

	if err != nil {
		panic(err)
	}

	defer zipReader.Close()

	_, err = io.Copy(f, zipReader)

	if err != nil {
		panic(err)
	}
}

func TestConvertEpub(t *testing.T) {
	err := Convert("source.epub", "target.mobi")
	assert.Nil(t, err)
}

func TestConvertFB2(t *testing.T) {
	err := Convert("source2.zip", "target.mobi")
	assert.Nil(t, err)
}

func TestConvertError(t *testing.T) {
	err := Convert("source3.zip", "target2.mobi")
	assert.NotNil(t, err)
}

func TestMeta(t *testing.T) {
	r, err := Meta("source1.zip")
	assert.Equal(t, r.Title, "Война и мир. Том 1")
	assert.Equal(t, r.Author, "Лев Николаевич Толстой")
	assert.Equal(t, r.Published.Year(), 1867)
	assert.Equal(t, int(r.Published.Month()), 3)
	assert.NotEqual(t, r.Cover, "")
	assert.Nil(t, err)
}
