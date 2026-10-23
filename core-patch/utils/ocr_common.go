package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/yatori-dev/yatori-go-core/que-core/qentity"
)

var (
	ocrCacheMu sync.Mutex
	ocrCache   = make(map[string]string)
)

// OCRQuestionImages 对题目中的图片做本地 OCR，返回带标签的识别文本
func OCRQuestionImages(question qentity.Question) string {
	type imageItem struct {
		label string
		img   qentity.QuestionImage
	}
	items := make([]imageItem, 0)
	for i, img := range question.ContentImages {
		items = append(items, imageItem{label: "题目图片" + itoa(i+1), img: img})
	}
	keys := make([]string, 0, len(question.OptionImages))
	for letter := range question.OptionImages {
		keys = append(keys, letter)
	}
	sort.Strings(keys)
	for _, letter := range keys {
		for i, img := range question.OptionImages[letter] {
			items = append(items, imageItem{label: "选项" + letter + "图片" + itoa(i+1), img: img})
		}
	}

	var builder strings.Builder
	count := 0
	for _, item := range items {
		if count >= 6 || builder.Len() >= 4000 {
			break
		}
		data := imageBytes(item.img)
		if len(data) == 0 {
			continue
		}
		text := ocrImageCached(data)
		if text == "" {
			continue
		}
		remaining := 4000 - builder.Len()
		if len([]rune(text)) > remaining {
			text = string([]rune(text)[:remaining])
		}
		if builder.Len() > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(item.label)
		builder.WriteString("：")
		builder.WriteString(text)
		count++
	}
	return builder.String()
}

func imageBytes(img qentity.QuestionImage) []byte {
	if img.Base64 != "" {
		data, err := base64.StdEncoding.DecodeString(img.Base64)
		if err == nil {
			return data
		}
	}
	if img.LocalPath != "" {
		data, err := os.ReadFile(img.LocalPath)
		if err == nil {
			return data
		}
	}
	return nil
}

func ocrImageCached(data []byte) string {
	sum := sha1.Sum(data)
	key := hex.EncodeToString(sum[:])
	ocrCacheMu.Lock()
	if text, ok := ocrCache[key]; ok {
		ocrCacheMu.Unlock()
		return text
	}
	ocrCacheMu.Unlock()

	text, err := OCRImageBytes(data)
	if err != nil || strings.TrimSpace(text) == "" {
		return ""
	}
	text = cleanOCRText(text)
	ocrCacheMu.Lock()
	if len(ocrCache) > 300 {
		ocrCache = make(map[string]string)
	}
	ocrCache[key] = text
	ocrCacheMu.Unlock()
	return text
}

func cleanOCRText(s string) string {
	runes := []rune(strings.TrimSpace(s))
	var builder strings.Builder
	for i, r := range runes {
		if r == ' ' && i > 0 && i < len(runes)-1 && isCJK(runes[i-1]) && isCJK(runes[i+1]) {
			continue
		}
		builder.WriteRune(r)
	}
	return strings.TrimSpace(builder.String())
}

func isCJK(r rune) bool {
	return r >= 0x4E00 && r <= 0x9FFF
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
