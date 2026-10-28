package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// 用于生成IMEI
func TokenHex(len int) string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// 极端情况下退化为时间种子，避免因随机源异常导致进程崩溃
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	token := hex.EncodeToString(b)
	return token
}
