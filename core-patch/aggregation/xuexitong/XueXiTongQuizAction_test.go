package xuexitong

import "testing"

func TestParseQuizTimeline(t *testing.T) {
	body := `[{"datas":[{"resourceId":123,"description":"1+1=?","questionType":"single","memberinfo":"m1","answerContent":"","options":[{"content":"1"},{"content":"2"}]}]}]`
	items := ParseQuizTimeline(body)
	if len(items) != 1 {
		t.Fatalf("expect 1 item, got %d", len(items))
	}
	if items[0].ResourceID != "123" || items[0].Description != "1+1=?" || len(items[0].RawOptions) != 2 {
		t.Fatalf("unexpected item: %+v", items[0])
	}
}

func TestParseVideoAttachment(t *testing.T) {
	html := `<script>mArg = ""; try{ mArg = {"attachments":[{"type":"document","mid":"d1"},{"type":"video","mid":"v1","objectId":"o1","property":{"mid":"v1"}}]} }catch(e){}</script>`
	mid, oid, ok := ParseVideoAttachment(html)
	if !ok || mid != "v1" || oid != "o1" {
		t.Fatalf("got mid=%s oid=%s ok=%v", mid, oid, ok)
	}
}

func TestNormalizeQuizAnswer(t *testing.T) {
	cases := map[string]string{
		"A":       "A",
		"答案是B":    "B",
		"A,C":     "A,C",
		"选 A 和 C": "A,C",
	}
	for in, want := range cases {
		got, err := NormalizeQuizAnswer(in)
		if err != nil || got != want {
			t.Fatalf("%s => %s, err=%v, want %s", in, got, err, want)
		}
	}
}
