package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/yatori-dev/yatori-go-core/que-core/qentity"
)

const maxImageBytes = 5 * 1024 * 1024

var (
	imageCacheMu  sync.RWMutex
	imageCacheDir = "./data/images"
)

// SetImageCacheDir 设置题目图片缓存目录（建议设置为软件目录下的 data/images）
func SetImageCacheDir(dir string) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return
	}
	imageCacheMu.Lock()
	imageCacheDir = dir
	imageCacheMu.Unlock()
	_ = os.MkdirAll(dir, 0755)
}

// GetImageCacheDir 获取当前图片缓存目录
func GetImageCacheDir() string {
	imageCacheMu.RLock()
	defer imageCacheMu.RUnlock()
	return imageCacheDir
}

func saveImageToCache(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	dir := GetImageCacheDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return ""
	}
	sum := sha1.Sum(data)
	name := hex.EncodeToString(sum[:]) + imageExt(data)
	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); err == nil {
		return path
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return ""
	}
	return path
}

func imageExt(data []byte) string {
	contentType := http.DetectContentType(data)
	switch {
	case strings.Contains(contentType, "png"):
		return ".png"
	case strings.Contains(contentType, "jpeg"), strings.Contains(contentType, "jpg"):
		return ".jpg"
	case strings.Contains(contentType, "gif"):
		return ".gif"
	case strings.Contains(contentType, "webp"):
		return ".webp"
	case strings.Contains(contentType, "bmp"):
		return ".bmp"
	default:
		return ".img"
	}
}

// CleanImageCache 清理超过 maxAge 的图片缓存，返回删除数量
func CleanImageCache(maxAge time.Duration) (int, error) {
	dir := GetImageCacheDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	cutoff := time.Now().Add(-maxAge)
	removed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(filepath.Join(dir, entry.Name())); err == nil {
				removed++
			}
		}
	}
	return removed, nil
}

// StartImageCacheCleaner 启动图片缓存定时清理：启动后立即清理一次，之后每隔 interval 清理一次
func StartImageCacheCleaner(interval, maxAge time.Duration) {
	if interval <= 0 {
		interval = 6 * time.Hour
	}
	if maxAge <= 0 {
		maxAge = 3 * 24 * time.Hour
	}
	go func() {
		if removed, err := CleanImageCache(maxAge); err == nil && removed > 0 {
			log.Printf("图片缓存清理完成，删除 %d 个过期文件", removed)
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if removed, err := CleanImageCache(maxAge); err == nil && removed > 0 {
				log.Printf("图片缓存清理完成，删除 %d 个过期文件", removed)
			}
		}
	}()
}

func DownloadImageAsQuestionImage(url string, cookies []*http.Cookie) (qentity.QuestionImage, error) {
	result := qentity.QuestionImage{Url: url}

	if strings.HasPrefix(url, "data:image/") {
		parts := strings.SplitN(url, ",", 2)
		if len(parts) == 2 {
			data, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				return result, err
			}
			if len(data) > maxImageBytes {
				return result, errors.New("图片超过5MB限制")
			}
			result.Base64 = parts[1]
			result.FileSize = len(data)
			result.LocalPath = saveImageToCache(data)
			img, _, err := image.Decode(strings.NewReader(string(data)))
			if err == nil {
				bounds := img.Bounds()
				result.Width = bounds.Dx()
				result.Height = bounds.Dy()
			}
			return result, nil
		}
	}

	tr := &http.Transport{
		IdleConnTimeout:   10 * time.Second,
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: tr, Timeout: 15 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(io.LimitReader(resp.Body, maxImageBytes+1))
	if err != nil {
		return result, err
	}

	if len(data) > maxImageBytes {
		return result, errors.New("图片超过5MB限制")
	}
	if !strings.HasPrefix(http.DetectContentType(data), "image/") {
		return result, errors.New("下载内容不是图片")
	}

	result.Base64 = base64.StdEncoding.EncodeToString(data)
	result.FileSize = len(data)
	result.LocalPath = saveImageToCache(data)

	img, _, err := image.Decode(strings.NewReader(string(data)))
	if err == nil {
		bounds := img.Bounds()
		result.Width = bounds.Dx()
		result.Height = bounds.Dy()
	}

	return result, nil
}
