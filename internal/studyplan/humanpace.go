package studyplan

import (
	"math/rand"
	"runtime"
	"time"
)

// AutoConcurrency 按 CPU 逻辑核数给出建议并发上限。
// 借鉴 Python 版 parallel.py 的“按机器性能调参”思路：
// 并发数不超过 CPU 核数的 2 倍，且至少为 1。
func AutoConcurrency(configured int) int {
	if configured <= 0 {
		return 1
	}
	maxByCPU := runtime.NumCPU() * 2
	if maxByCPU < 1 {
		maxByCPU = 1
	}
	if configured > maxByCPU {
		return maxByCPU
	}
	return configured
}

// HumanPause 生成拟人化的操作间隔（毫秒级抖动）。
// 借鉴 Python 版 human_pace.py 的对数正态分布 + 疲劳曲线思路：
// 大多数间隔集中在 base 附近，偶发较长停顿。
func HumanPause(base time.Duration) time.Duration {
	if base <= 0 {
		return 0
	}
	// 0.6 ~ 2.4 倍，偏向下限，长尾
	factor := 0.6 + rand.Float64()*1.8
	if rand.Intn(100) < 8 {
		factor *= 2.2 // 8% 概率出现较长停顿
	}
	return time.Duration(float64(base) * factor)
}

// AnswerPause 答题前的思考停顿，单题 1.2~4.5 秒
func AnswerPause() time.Duration {
	return HumanPause(1600 * time.Millisecond)
}
