package xuexitong

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yatori-dev/yatori-go-core/models/ctype"
	"github.com/yatori-dev/yatori-go-core/que-core/aiq"
	"github.com/yatori-dev/yatori-go-core/que-core/qentity"
	"github.com/yatori-dev/yatori-go-core/que-core/qtype"
	"github.com/yatori-dev/yatori-go-core/utils"
)

// PrepareQuestionImages 下载题目图片并填充 Base64 / 本地缓存路径
func PrepareQuestionImages(cache *XueXiTUserCache, question *qentity.Question) {
	if cache == nil || question == nil || !question.HasImages() {
		return
	}
	cookies := cache.GetCookies()
	for i := range question.ContentImages {
		img := &question.ContentImages[i]
		if img.Base64 != "" || img.Url == "" {
			continue
		}
		if downloaded, err := utils.DownloadImageAsQuestionImage(img.Url, cookies); err == nil {
			*img = downloaded
		}
	}
	for key, imgs := range question.OptionImages {
		for i := range imgs {
			img := &imgs[i]
			if img.Base64 != "" || img.Url == "" {
				continue
			}
			if downloaded, err := utils.DownloadImageAsQuestionImage(img.Url, cookies); err == nil {
				*img = downloaded
			}
		}
		question.OptionImages[key] = imgs
	}
}

// AnswerQuestionWithVision 使用识图模型作答含图片的题目
func AnswerQuestionWithVision(cache *XueXiTUserCache, question qentity.Question, setter AnswerSetter,
	visionUrl, visionModel string, visionType ctype.AiType, visionKey string) error {
	if strings.TrimSpace(visionUrl) == "" || strings.TrimSpace(visionKey) == "" {
		return errors.New("未配置识图模型")
	}
	if !question.HasImages() {
		return errors.New("题目不含图片")
	}

	imageDataMap := make(map[string]qentity.QuestionImage)
	cookies := cache.GetCookies()
	cacheImage := func(img qentity.QuestionImage) {
		if img.Url == "" && img.Base64 == "" {
			return
		}
		if existing, ok := imageDataMap[img.Url]; ok && existing.Base64 != "" {
			return
		}
		if img.Base64 == "" && img.Url != "" {
			if downloaded, err := utils.DownloadImageAsQuestionImage(img.Url, cookies); err == nil {
				imageDataMap[img.Url] = downloaded
				return
			}
		}
		imageDataMap[img.Url] = img
	}
	for _, img := range question.ContentImages {
		cacheImage(img)
	}
	for _, imgs := range question.OptionImages {
		for _, img := range imgs {
			cacheImage(img)
		}
	}

	messages := aiq.BuildVisualAiQuestionMessage(question, imageDataMap)
	reply, err := aiq.AggregationVisualAIApi(visionUrl, visionModel, visionType, messages, visionKey)
	if err != nil {
		return err
	}
	answers := ParseAIStringArray(reply)
	if question.Type == qtype.FillInTheBlank.String() {
		answers = aiq.ParseFillAnswers(reply, len(question.Options))
	}
	if len(answers) == 0 {
		return errors.New("识图AI未返回有效答案")
	}
	setter.SetAnswers(answers)
	return nil
}

// ParseAIStringArray 解析 AI 返回，兼容 JSON 数组、对象、代码块和纯文本
func ParseAIStringArray(reply string) []string {
	return aiq.ParseStringArray(reply)
}

// FormatVisionError 便于日志输出
func FormatVisionError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%v", err)
}
