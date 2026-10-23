package qentity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type ResultQuestion struct {
	Question
	Replier string `json:"replier"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
}

type QuestionImage struct {
	Url       string `json:"url"`
	Base64    string `json:"base64"`
	LocalPath string `json:"localPath,omitempty"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	FileSize  int    `json:"fileSize"`
	AltText   string `json:"altText"`
}

type Question struct {
	Md5           string                       `gorm:"column:md5" json:"md5"`
	Type          string                       `gorm:"column:type" json:"type"`
	Content       string                       `gorm:"column:content" json:"content"`
	ContentImages []QuestionImage              `gorm:"-" json:"contentImages,omitempty"`
	Options       StringArray                  `gorm:"column:options" json:"options"`
	OptionImages  map[string][]QuestionImage   `gorm:"-" json:"optionImages,omitempty"`
	Answers       StringArray                  `gorm:"column:answers;type:TEXT" json:"answers"`
}

func (q *Question) HasImages() bool {
	if len(q.ContentImages) > 0 {
		return true
	}
	for _, imgs := range q.OptionImages {
		if len(imgs) > 0 {
			return true
		}
	}
	return false
}

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte: %T", value)
	}
	return json.Unmarshal(bytes, s)
}
