package calibre

import (
	"github.com/nickalie/go-binwrapper"
	"strings"
	"github.com/saintfish/chardet"
	"gopkg.in/iconv.v1"
	"encoding/hex"
	"path/filepath"
	"os"
	"crypto/rand"
	"time"
	"github.com/pkg/errors"
)

type MetaData struct {
	Title  string
	Author string
	Cover  string
	Published time.Time
}

func Meta(file string) (*MetaData, error) {
	b := NewMetaWrapper()
	result := MetaData{}
	result.Cover = tempFileName(".png")

	err := b.Run("--get-cover=" + result.Cover, file)

	if err != nil {
		return nil, errors.Wrap(err, string(b.CombinedOutput()))
	}

	outputs := strings.Split(string(b.CombinedOutput()), "\n")
	isCoverFound := false
	for _, v := range outputs {
		if strings.Contains(v, result.Cover) {
			isCoverFound = true
		}

		parts := strings.Split(v, ":")

		if len(parts) < 2 {
			continue
		}

		prefix := strings.ToLower(parts[0])
		suffix := strings.Trim(strings.Join(parts[1:], ":"), " ")

		if strings.Contains(prefix, "title") {
			result.Title = decodeString(suffix)
		} else if strings.Contains(prefix,"author") {
			result.Author = decodeString(suffix)
		} else if strings.Contains(prefix, "published") {
			t, err := time.Parse(time.RFC3339, suffix)
			if err == nil {
				result.Published = t
			}
		}
	}

	if !isCoverFound {
		result.Cover = ""
	}

	return &result, err
}

func tempFileName(ext string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), hex.EncodeToString(randBytes)+ext)
}

func decodeString(value string) string {
	r, err := chardet.NewTextDetector().DetectBest([]byte(value))

	if err != nil || strings.ToLower(r.Charset) == "utf-8" {
		return value
	}

	cd, err := iconv.Open("utf-8", r.Charset)

	if err != nil {
		return value
	}

	defer cd.Close()
	return cd.ConvString(value)
}

type MetaWrapper struct {
	*binwrapper.BinWrapper
}

func NewMetaWrapper() *MetaWrapper {
	c := &MetaWrapper{
		BinWrapper: binwrapper.NewBinWrapper().ExecPath("ebook-meta"),
	}
	return c
}
