//go:build windows

package utils

import (
	"bytes"
	"context"
	"errors"
	"image"
	stddraw "image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	xdraw "golang.org/x/image/draw"
)

const winOCRPowerShell = `param([string]$Path, [string]$Out)
$ErrorActionPreference = "Stop"
Add-Type -AssemblyName System.Runtime.WindowsRuntime
$asTaskGeneric = ([System.WindowsRuntimeSystemExtensions].GetMethods() | Where-Object { $_.Name -eq 'AsTask' -and $_.GetParameters().Count -eq 1 -and $_.GetParameters()[0].ParameterType.Name -eq 'IAsyncOperation` + "`" + `1' })[0]
function Await($WinRtTask, $ResultType) {
  $asTask = $asTaskGeneric.MakeGenericMethod($ResultType)
  $netTask = $asTask.Invoke($null, @($WinRtTask))
  $netTask.Wait(-1) | Out-Null
  $netTask.Result
}
$null = [Windows.Storage.StorageFile, Windows.Storage, ContentType=WindowsRuntime]
$null = [Windows.Graphics.Imaging.BitmapDecoder, Windows.Graphics.Imaging, ContentType=WindowsRuntime]
$null = [Windows.Media.Ocr.OcrEngine, Windows.Foundation, ContentType=WindowsRuntime]
$file = Await ([Windows.Storage.StorageFile]::GetFileFromPathAsync($Path)) ([Windows.Storage.StorageFile])
$stream = Await ($file.OpenAsync([Windows.Storage.FileAccessMode]::Read)) ([Windows.Storage.Streams.IRandomAccessStream])
$decoder = Await ([Windows.Graphics.Imaging.BitmapDecoder]::CreateAsync($stream)) ([Windows.Graphics.Imaging.BitmapDecoder])
$bitmap = Await ($decoder.GetSoftwareBitmapAsync()) ([Windows.Graphics.Imaging.SoftwareBitmap])
$engine = [Windows.Media.Ocr.OcrEngine]::TryCreateFromUserProfileLanguages()
if ($engine -eq $null) {
  [System.IO.File]::WriteAllText($Out, "NO_ENGINE", (New-Object System.Text.UTF8Encoding($false)))
  exit 2
}
$result = Await ($engine.RecognizeAsync($bitmap)) ([Windows.Media.Ocr.OcrResult])
[System.IO.File]::WriteAllText($Out, $result.Text, (New-Object System.Text.UTF8Encoding($false)))
`

// prepareOCRImage 透明底转白底，并对小图放大，提升 Windows OCR 识别率
func prepareOCRImage(data []byte) []byte {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return data
	}
	base := image.NewRGBA(image.Rect(0, 0, w, h))
	stddraw.Draw(base, base.Bounds(), image.White, image.Point{}, stddraw.Src)
	stddraw.Draw(base, base.Bounds(), img, bounds.Min, stddraw.Over)

	scale := 1
	if w < 400 {
		scale = 3
	} else if w < 800 {
		scale = 2
	}
	if h < 40 && scale < 3 {
		scale = 3
	}

	var target image.Image = base
	if scale > 1 {
		scaled := image.NewRGBA(image.Rect(0, 0, w*scale, h*scale))
		xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), base, base.Bounds(), stddraw.Over, nil)
		target = scaled
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, target); err != nil {
		return data
	}
	return buf.Bytes()
}

// OCRImageBytes 使用 Windows 自带 OCR 识别图片文字（需要系统安装对应语言包）
func OCRImageBytes(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("空图片")
	}
	dir := GetImageCacheDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	scriptPath := filepath.Join(dir, "win_ocr.ps1")
	if info, err := os.Stat(scriptPath); err != nil || time.Since(info.ModTime()) > 24*time.Hour {
		if err := os.WriteFile(scriptPath, []byte(winOCRPowerShell), 0644); err != nil {
			return "", err
		}
	}

	imgFile, err := os.CreateTemp(dir, "ocr-img-*.png")
	if err != nil {
		return "", err
	}
	imgPath := imgFile.Name()
	defer os.Remove(imgPath)
	if _, err := imgFile.Write(prepareOCRImage(data)); err != nil {
		imgFile.Close()
		return "", err
	}
	imgFile.Close()

	outFile, err := os.CreateTemp(dir, "ocr-out-*.txt")
	if err != nil {
		return "", err
	}
	outPath := outFile.Name()
	outFile.Close()
	defer os.Remove(outPath)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath, "-Path", imgPath, "-Out", outPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	var stderr strings.Builder
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	text, readErr := os.ReadFile(outPath)
	if readErr != nil {
		if runErr != nil {
			return "", runErr
		}
		return "", readErr
	}
	result := strings.TrimSpace(string(text))
	if result == "NO_ENGINE" {
		return "", errors.New("系统未安装 OCR 语言包")
	}
	if result == "" {
		if ctx.Err() != nil {
			return "", errors.New("OCR 识别超时")
		}
		if stderr.Len() > 0 {
			return "", errors.New(strings.TrimSpace(stderr.String()))
		}
	}
	return result, nil
}
