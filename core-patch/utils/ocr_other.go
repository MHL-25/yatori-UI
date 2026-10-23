//go:build !windows

package utils

import "errors"

// OCRImageBytes 非 Windows 平台暂不支持本地 OCR
func OCRImageBytes(data []byte) (string, error) {
	return "", errors.New("当前平台不支持本地 OCR")
}
