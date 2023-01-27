//go:build blackfriday
// +build blackfriday

package syntax

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"strings"
)

func Convert(gmi string, wrap bool) string {
	clearedString := strings.ReplaceAll(gmi, "\r\n", "\n")
	unsafeString := blackfriday.Run([]byte(clearedString), blackfriday.WithExtensions(blackfriday.CommonExtensions))
	saferString := bluemonday.UGCPolicy().SanitizeBytes(unsafeString)
	return string(saferString)
}
