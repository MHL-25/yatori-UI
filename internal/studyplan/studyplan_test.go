package studyplan

import (
	"testing"
	"time"

	"github.com/yatori-dev/yatori-go-core/aggregation/xuexitong"
)

// 使用 2026-10-23 抓取到的真实章节数据做分类校验
func TestAnalyzeRealCourse(t *testing.T) {
	now := time.Date(2026, 10, 23, 21, 0, 0, 0, time.Local)
	chapters := []xuexitong.KnowledgeItem{
		{ID: 1, Label: "4.1", Name: "正手发高远球", Layer: 2, JobCount: 2, Status: "open",
			BeginTime: time.Date(2026, 8, 30, 20, 11, 0, 0, time.Local).UnixMilli(),
			EndTime:   time.Date(2026, 9, 4, 20, 11, 0, 0, time.Local).UnixMilli(),
			PointTotal: 2, PointFinished: 2},
		{ID: 2, Label: "4.2", Name: "击打高远球", Layer: 2, JobCount: 1, Status: "close",
			BeginTime: time.Date(2026, 9, 6, 20, 12, 0, 0, time.Local).UnixMilli(),
			EndTime:   time.Date(2026, 9, 11, 20, 12, 0, 0, time.Local).UnixMilli(),
			PointTotal: 1, PointFinished: 0},
		{ID: 3, Label: "4.6", Name: "正反手挑后场", Layer: 2, JobCount: 1, Status: "open",
			BeginTime: time.Date(2026, 10, 18, 20, 15, 0, 0, time.Local).UnixMilli(),
			EndTime:   time.Date(2026, 10, 23, 20, 15, 0, 0, time.Local).UnixMilli(),
			PointTotal: 1, PointFinished: 1},
		{ID: 4, Label: "4.8", Name: "正反手推球", Layer: 2, JobCount: 1, Status: "time",
			BeginTime: time.Date(2026, 10, 25, 20, 16, 0, 0, time.Local).UnixMilli(),
			EndTime:   time.Date(2026, 10, 30, 20, 16, 0, 0, time.Local).UnixMilli(),
			PointTotal: 1, PointFinished: 0},
		{ID: 5, Label: "3.2", Name: "场地介绍", Layer: 2, JobCount: 3, Status: "open",
			PointTotal: 3, PointFinished: 1},
		{ID: 6, Label: "3", Name: "理论知识", Layer: 1, JobCount: 0, Status: "open"},
	}
	plans := Analyze(chapters, now)
	if len(plans) != 5 {
		t.Fatalf("expect 5 node plans, got %d", len(plans))
	}
	want := map[string]LockType{
		"4.1": LockOpen,
		"4.2": LockExpired,
		"4.6": LockOpen,
		"4.8": LockTime,
		"3.2": LockOpen,
	}
	for _, p := range plans {
		if want[p.Node.Label] != p.Lock {
			t.Errorf("%s expect %s, got %s", p.Node.Label, want[p.Node.Label], p.Lock)
		}
	}
	if !Restricted(plans) {
		t.Error("course should be restricted")
	}
	next, ok := NextUnlock(plans)
	if !ok || next.Day() != 25 {
		t.Errorf("expect next unlock on 10-25, got %v %v", next, ok)
	}
}

func TestAnalyzePrerequisite(t *testing.T) {
	now := time.Date(2026, 10, 23, 21, 0, 0, 0, time.Local)
	chapters := []xuexitong.KnowledgeItem{
		{ID: 1, Label: "1.1", Layer: 2, JobCount: 1, Status: "open", PointTotal: 1, PointFinished: 0},
		{ID: 2, Label: "1.2", Layer: 2, JobCount: 1, Status: "close",
			EndTime: time.Date(2026, 12, 1, 0, 0, 0, 0, time.Local).UnixMilli()},
	}
	plans := Analyze(chapters, now)
	if plans[1].Lock != LockPrereq {
		t.Fatalf("expect prereq lock, got %s", plans[1].Lock)
	}
}

func TestAutoConcurrency(t *testing.T) {
	if got := AutoConcurrency(0); got != 1 {
		t.Errorf("expect 1, got %d", got)
	}
	if got := AutoConcurrency(3); got != 3 {
		t.Errorf("expect 3, got %d", got)
	}
	if got := AutoConcurrency(1000); got <= 0 {
		t.Errorf("expect positive, got %d", got)
	}
}
