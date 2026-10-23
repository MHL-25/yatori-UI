package studyplan

import (
	"time"

	"github.com/yatori-dev/yatori-go-core/aggregation/xuexitong"
)

// LockType 章节锁类型
type LockType string

const (
	LockOpen    LockType = "open"    // 可学习
	LockTime    LockType = "time"    // 未到开放时间
	LockPrereq  LockType = "prereq"  // 前置任务点未完成
	LockExpired LockType = "expired" // 已过截止时间
	LockClosed  LockType = "closed"  // 已关闭且无明确时间
)

// NodePlan 单个章节节点的学习计划
type NodePlan struct {
	Index     int
	Node      xuexitong.KnowledgeItem
	Lock      LockType
	UnlockAt  time.Time
	Reason    string
	Completed bool
}

// HasRestriction 是否存在会中断并发学习的限制
func (p NodePlan) HasRestriction() bool {
	if p.Completed {
		return false
	}
	return p.Lock != LockOpen
}

// Analyze 按学习通章节字段分析锁类型。
// 规则：
//  1. status=time 或 beginTime 未到 → 定时锁
//  2. status=close 且 endTime 已过 → 已截止
//  3. status=close 且前面存在未完成任务点 → 前置锁
//  4. status=close 其他情况 → 已关闭
//  5. 其余 → 开放
func Analyze(chapters []xuexitong.KnowledgeItem, now time.Time) []NodePlan {
	nowMs := now.UnixMilli()
	plans := make([]NodePlan, 0, len(chapters))
	for i := range chapters {
		ch := chapters[i]
		if ch.Layer != 2 || ch.JobCount <= 0 {
			continue
		}
		plan := NodePlan{Index: i, Node: ch, Lock: LockOpen}
		// 已完成节点不参与锁判定，避免把“已完成但已过截止时间”的章节显示为受限
		if ch.PointTotal > 0 && ch.PointFinished >= ch.PointTotal {
			plan.Completed = true
			plan.Reason = "已完成"
			plans = append(plans, plan)
			continue
		}
		switch ch.Status {
		case "time":
			plan.Lock = LockTime
			if ch.BeginTime > 0 {
				plan.UnlockAt = time.UnixMilli(ch.BeginTime)
				plan.Reason = "未到开放时间，将于 " + plan.UnlockAt.Format("01-02 15:04") + " 开放"
			} else {
				plan.Reason = "未到开放时间"
			}
		case "close":
			switch {
			case ch.EndTime > 0 && ch.EndTime < nowMs:
				plan.Lock = LockExpired
				plan.Reason = "已过截止时间 " + time.UnixMilli(ch.EndTime).Format("01-02 15:04")
			case hasUnfinishedBefore(chapters, i):
				plan.Lock = LockPrereq
				plan.Reason = "存在前置任务点未完成"
			default:
				plan.Lock = LockClosed
				plan.Reason = "章节已关闭"
			}
		default:
			if ch.BeginTime > 0 && ch.BeginTime > nowMs {
				plan.Lock = LockTime
				plan.UnlockAt = time.UnixMilli(ch.BeginTime)
				plan.Reason = "未到开放时间，将于 " + plan.UnlockAt.Format("01-02 15:04") + " 开放"
			}
		}
		plans = append(plans, plan)
	}
	return plans
}

func hasUnfinishedBefore(chapters []xuexitong.KnowledgeItem, index int) bool {
	for i := 0; i < index; i++ {
		ch := chapters[i]
		if ch.Layer != 2 || ch.JobCount <= 0 {
			continue
		}
		if ch.PointTotal > 0 && ch.PointFinished < ch.PointTotal {
			return true
		}
	}
	return false
}

// Restricted 判断计划中是否含限制节点
func Restricted(plans []NodePlan) bool {
	for _, p := range plans {
		if p.HasRestriction() {
			return true
		}
	}
	return false
}

// NextUnlock 返回最近的解锁时间（仅定时锁）
func NextUnlock(plans []NodePlan) (time.Time, bool) {
	var next time.Time
	found := false
	for _, p := range plans {
		if p.Lock == LockTime && !p.UnlockAt.IsZero() {
			if !found || p.UnlockAt.Before(next) {
				next = p.UnlockAt
				found = true
			}
		}
	}
	return next, found
}
