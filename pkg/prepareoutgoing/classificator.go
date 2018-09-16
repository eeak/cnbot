package prepareoutgoing

import (
	"bytes"
	"errors"
	"regexp"
	"unicode/utf8"
)

var (
	RX_JSON = regexp.MustCompile(
		`(?s:^` + // go fmt padding :-/
			`[[:space:]]*` +
			`(` +
			`sendMessage` +
			`|` +
			`deleteMessage` +
			`|` +
			`editMessageText` +
			`)` +
			`[[:space:]]*` +
			`({[[:space:]]*".*})` +
			`[[:space:]]*$)`,
	)
	RX_DOT       = regexp.MustCompile(`^[[:space:]]*\.[[:space:]]*$`)
	RX_NOT_SPACE = regexp.MustCompile(`[^[:space:]]`)
	// Details: https://en.wikipedia.org/wiki/List_of_file_signatures
	FP_GIF   = []byte{0x47, 0x49, 0x46, 0x38}
	FP_PNG   = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	FP_JPG_A = []byte{0xFF, 0xD8, 0xFF, 0xDB}
	FP_JPG_B = []byte{0xFF, 0xD8, 0xFF, 0xE0}
	FP_JPG_C = []byte{0xFF, 0xD8, 0xFF, 0xE1}
	// Errors
	errorMessageTooLong = errors.New("Message too long")
	errorInvalidUTF8    = errors.New("Invalid UTF8 string")
)

func classifyData(data []byte) (
	isEmpty bool,
	leftIt bool,
	isRaw bool,
	rawMethod string,
	rawPayload []byte,
	isImage bool,
	imageType string,
	err error,
) {
	if len(data) == 0 {
		isEmpty = true
	} else {
		if bytes.HasPrefix(data, FP_GIF) {
			isImage = true
			imageType = "gif"
		} else if bytes.HasPrefix(data, FP_PNG) {
			isImage = true
			imageType = "png"
		} else if bytes.HasPrefix(data, FP_JPG_A) ||
			bytes.HasPrefix(data, FP_JPG_B) ||
			bytes.HasPrefix(data, FP_JPG_C) {
			isImage = true
			imageType = "jpeg"
		} else {
			if utf8.Valid(data) {
				if !RX_NOT_SPACE.Match(data) {
					isEmpty = true
				} else if RX_DOT.Match(data) {
					leftIt = true
				} else if r := RX_JSON.FindSubmatch(data); r != nil {
					isRaw = true
					rawMethod = string(r[1])
					rawPayload = r[2]
				} else if utf8.RuneCount(data) > 4000 {
					err = errorMessageTooLong
				}
			} else {
				err = errorInvalidUTF8
			}
		}
	}
	return
}
