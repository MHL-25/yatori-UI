package aiq

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/yatori-dev/yatori-go-core/que-core/qentity"
	"github.com/yatori-dev/yatori-go-core/que-core/qtype"
)

// 单选题处理策略
func handleSingleChoice(topic qentity.Question) AIChatMessages {
	problem := buildProblemHeader(topic.Type, topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: `接下来无论出现任何题目，你都必须只回答题目中某个选项对应的内容，并严格按照以下要求作答：

【回答规则】
1. 最终输出必须严格遵循 JSON 数组格式，例如：["选项内容"]
2. 数组中只能有一个字符串元素。
3. 字符串中不能包含选项前缀，如 A. B. C. D. 等，只能输出选项的纯内容。
4. 不能输出解析、解释步骤、理由、提示语或任何多余文本。
5. 不能输出题目本身、不能输出其他格式，只能输出 JSON 数组。
6. 如果你无法判断正确答案，也必须随机选择一个选项的内容进行输出，不允许回答“我不知道”“无法判断”之类内容。

【格式要求】
- 只能输出 JSON
- 不允许换行，若内容中需要换行必须使用\n
- 不能出现额外的空格、标点或第二层数组

【示例】
题目如下：
试卷名称：考试
题目类型：单选题
题目内容：新中国是什么时候成立的？
A. 1949年10月5日
B. 1949年10月1日
C. 1949年09月1日
D. 2002年10月1日

你必须回答：
["1949年10月1日"]
（注意：不能输出 B 或 B. 等前缀，只能输出选项的内容本身）

请严格按照以上规则回答后续所有题目，只输出 JSON 数组格式内容，不得违反任何一项要求。`},
		//{Role: "system", Content: `接下来你只需要回答选项对应内容即可...回答格式严格遵循json范式，比如：["选项1内容"]`},
		//{Role: "system", Content: "就算你不知道选什么也随机选，答题过程中你无需回答任何解释，只需按照预定json范式输出！！！"},
		//{Role: "system", Content: exampleSingleChoice()},
		{Role: "user", Content: problem},
	}}
}

// 构建AI题目消息
func BuildAiQuestionMessage(topic qentity.Question) AIChatMessages {
	switch qtype.Index(topic.Type) {
	case qtype.SingleChoice: //单选题
		return handleSingleChoice(topic)
	case qtype.MultipleChoice: //多选题
		return handleMultipleChoice(topic)
	case qtype.TrueOrFalse: //判断题
		return handleTrueFalse(topic)
	case qtype.FillInTheBlank: //填空题
		return handleFillInTheBlank(topic)
	case qtype.ShortAnswer:
		return handleShortAnswer(topic)
	case qtype.TermExplanation:
		return handleTermExplanationAnswer(topic)
	case qtype.Essay:
		return handleEssayAnswer(topic)
	case qtype.Matching:
		return handleMatchingAnswer(topic)

	}
	return AIChatMessages{}
}

// 将AI的回复转换到Question里面
func ResponseTurnQuestion(question qentity.Question, response string) qentity.Question {
	var answers []string
	json.Unmarshal([]byte(response), &answers)
	question.Answers = answers
	return question
}

// 多选题处理策略
func handleMultipleChoice(topic qentity.Question) AIChatMessages {
	problem := buildProblemHeader("多选题", topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: `接下来无论出现任何题目，你都必须只回答选项对应的内容，且必须严格按照以下格式输出：

【最终输出格式】
["选项内容1","选项内容2", ...]
- JSON 数组只能包含字符串元素。
- 每个元素对应一个被选中的选项内容。
- 严禁携带 A. B. C. D. 等前缀，只能输出纯内容。
- 不得输出解析、解释、思考过程、题目内容或任何无关文本。
- 如果你无法判断正确选项，也必须随机选择多个选项内容填入数组。

【格式注意事项】
1. 只能输出 JSON 数组，不允许输出其他结构。
2. 数组内字符串之间用英文逗号分隔。
3. 字符串内部如需换行，必须使用 \n，不允许真实换行。
4. 不允许出现空元素、null、其他非字符串类型。

【示例】
题目：
试卷名称：考试
题目类型：多选题
题目内容：马克思关于资本积累的学说是剩余价值理论的重要组成部分...
A. 资本主义扩大再生产的源泉
B. 资本有机构成呈现不断降低趋势的根本原因
C. 社会财富占有两极分化的重要原因
D. 资本主义社会失业现象产生的根源

正确回答应该是：
["资本主义扩大再生产的源泉","社会财富占有两极分化的重要原因","资本主义社会失业现象产生的根源"]
（注意：不能携带 A/B/C/D 前缀，只能输出选项的纯内容）

请严格按照以上规则回答后续所有题目。`},
		//{Role: "system", Content: `接下来你只需要回答选项对应内容即可...格式：["选项1","选项2"]`},
		//{Role: "system", Content: `就算你不知道选什么也随机选...无需回答任何解释！！！`},
		//{Role: "system", Content: exampleMultipleChoice()},
		{Role: "user", Content: problem},
	}}
}

// 判断题处理策略
func handleTrueFalse(topic qentity.Question) AIChatMessages {
	problem := buildProblemHeader("判断题", topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: `接下来你只需要回答“正确”或者“错误”即可...格式：["正确"]`},
		{Role: "system", Content: "就算你不知道选什么也随机选...无需回答任何解释！！！"},
		{Role: "system", Content: exampleTrueFalse()},
		{Role: "user", Content: problem},
	}}
}

// 填空题处理策略
func handleFillInTheBlank(topic qentity.Question) AIChatMessages {
	problem := buildProblemHeader("填空题", topic)
	blankCount := len(topic.Options)
	if blankCount <= 0 {
		blankCount = 1
	}
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: fmt.Sprintf(`本题一共有 %d 个空，题目下方会按顺序列出“第1空”到“第%d空”。
你必须只输出一个 JSON 字符串数组，并且数组长度必须恰好为 %d：
- 第 1 个元素填第 1 空，第 2 个元素填第 2 空，依次对应，顺序不能错。
- 只输出 JSON 数组，禁止输出编号、开场白、解释、题目原文、答案总结或任何多余文字。
- 每个元素只写该空对应的简短答案，不要带“第N空”“1.”“（1）”等前缀。
- 不要写“此处应绘制”“此处应提供图形”“略”“如图”之类的说明，要给出最可能的简短文字答案。
- 如果某个空位要求“画出……图”“拆画……图”而无法真正画图，请填写该图的名称（例如“A-A断面图”“零件2夹套的零件图”），不要填材料牌号、尺寸等无关信息。
- 逐个空位分别作答，不要机械重复同一个答案，除非题目原文确实要求重复。
- 即使不确定也必须填写，不能留空。`, blankCount, blankCount, blankCount)},
		{Role: "system", Content: exampleFillInTheBlank()},
		{Role: "user", Content: problem},
	}}
}

// 简答题处理策略
func handleShortAnswer(topic qentity.Question) AIChatMessages {
	problem := buildProblemHeader("简答题", topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: `这是一个简答题，回答时请严格遵循json格式，包括换行等特殊符号也要遵循json语法：["答案"]，注意不要拆分答案！！！`},
		//{Role: "system", Content: "就算你不知道选什么也随机选...无需回答任何解释！！！"},
		{Role: "system", Content: exampleShortAnswer()},
		{Role: "user", Content: problem},
	}}
}

// 名词解释处理策略
func handleTermExplanationAnswer(topic qentity.Question) AIChatMessages {
	problem := buildProblemHeader(topic.Type, topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: `最终输出必须是一个合法 JSON 数组格式：["答案内容"]
数组中只能包含一个字符串元素，答案必须完整写在同一个字符串里，不能拆分成多个元素。
字符串内如需换行必须写为 \n，不能出现真正的换行符。
字符串内如出现双引号必须转义为 \"。
答案内容必须是连贯的完整论述，不得包含解析、题目、注释或生成说明。
除 JSON 数组外严禁输出任何其他内容。
示例（仅结构示例，非真实内容）：
["这里填写你的完整名词解释题答案……"]
请开始回答名词解释题。`},
		//{Role: "system", Content: `这是一个名词解释题，回答时请严格遵循json格式，包括换行等特殊符号也要遵循json语法：["答案"]，注意不要拆分答案！！！`},
		//{Role: "system", Content: exampleTermExplanationAnswer()},
		{Role: "user", Content: problem},
	}}
}

// 论述题处理策略
func handleEssayAnswer(topic qentity.Question) AIChatMessages {
	problem := buildProblemHeader(topic.Type, topic)
	return AIChatMessages{Messages: []Message{
		//{Role: "system", Content: `这是一个论述题，回答时请严格遵循json格式，包括换行等特殊符号也要遵循json语法：["答案"]，注意不要拆分答案！！！`},
		{Role: "system", Content: `最终输出必须是一个合法 JSON 数组格式：["答案内容"]
数组中只能包含一个字符串元素，答案必须完整写在同一个字符串里，不能拆分成多个元素。
字符串内如需换行必须写为 \n，不能出现真正的换行符。
字符串内如出现双引号必须转义为 \"。
答案内容必须是连贯的完整论述，不得包含解析、题目、注释或生成说明。
答案字数不少于 500 字。
除 JSON 数组外严禁输出任何其他内容。
示例（仅结构示例，非真实内容）：
["这里填写你的完整论述题答案，必须超过500字……"]
请开始回答论述题。`},
		//{Role: "system", Content: exampleEssayAnswer()},
		{Role: "user", Content: problem},
	}}
}

// 连线题处理策略
func handleMatchingAnswer(topic qentity.Question) AIChatMessages {
	problem := buildProblemHeader(topic.Type, topic)
	return AIChatMessages{Messages: []Message{
		{Role: "system", Content: `接下来你只需要以json格式回答选项对应内容即可，比如：["xxx->xxx","xxx->xxx"]`},
		{Role: "system", Content: "就算你不知道选什么也随机按指定要求格式回答...无需回答任何解释！！！"},
		{Role: "system", Content: exampleMatchingAnswer()},
		{Role: "user", Content: problem},
	}}
}

// 构建题目头部信息
func buildProblemHeader(topicType string, topic qentity.Question) string {
	selectStr := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	if qtype.Index(topic.Type) == qtype.FillInTheBlank {
		return buildFillProblemHeader(topicType, topic)
	}
	sprintf := fmt.Sprintf(`题目类型：%s
题目内容：
%s\n`, topicType, topic.Content)
	switch qtype.Index(topic.Type) {
	case qtype.SingleChoice: //单选题
		for i, option := range topic.Options {
			sprintf += selectStr[i] + "." + optionTextWithImageHint(topic, selectStr[i], option) + "\n"
		}
	case qtype.MultipleChoice: //多选题
		for i, option := range topic.Options {
			sprintf += selectStr[i] + "." + optionTextWithImageHint(topic, selectStr[i], option) + "\n"
		}
	case qtype.TrueOrFalse: //判断题
		for i, option := range topic.Options {
			sprintf += selectStr[i] + "." + optionTextWithImageHint(topic, selectStr[i], option) + "\n"
		}
	case qtype.FillInTheBlank: //填空题
		for i := range topic.Options {
			sprintf += fmt.Sprintf("第%d空：____\n", i+1)
		}
	case qtype.ShortAnswer: //简答题
	case qtype.TermExplanation: //名词解释
	case qtype.Essay: //论述题
	case qtype.Matching:
		sprintf += "组别一：\n"
		for _, option := range topic.Options {
			if strings.HasPrefix(option, "[1]") {
				text := strings.Replace(option, "[1]", "", 1)
				if strings.TrimSpace(text) == "" {
					text = "[图片项，请查看下方匹配项图片]"
				}
				sprintf += text + "\n"
			}
		}

		sprintf += "组别二：\n"
		for _, option := range topic.Options {
			if strings.HasPrefix(option, "[2]") {
				text := strings.Replace(option, "[2]", "", 1)
				if strings.TrimSpace(text) == "" {
					text = "[图片项，请查看下方匹配项图片]"
				}
				sprintf += text + "\n"
			}
		}

	}

	return sprintf
}

var fillBlankTextPattern = regexp.MustCompile(`[\x{00A0} ]{2,}|_{2,}`)

// buildFillProblemHeader 把题干中的空位替换成 [第N空]，让模型看清每个空的上下文
func buildFillProblemHeader(topicType string, topic qentity.Question) string {
	marked := markFillBlanks(topic.Content, len(topic.Options))
	var builder strings.Builder
	fmt.Fprintf(&builder, "题目类型：%s\n题目内容（[第N空]为需要填写的位置）：\n%s\n", topicType, marked)
	for i := range topic.Options {
		fmt.Fprintf(&builder, "第%d空：____\n", i+1)
	}
	return builder.String()
}

func markFillBlanks(content string, blankCount int) string {
	index := 0
	return fillBlankTextPattern.ReplaceAllStringFunc(content, func(match string) string {
		if index >= blankCount {
			return match
		}
		index++
		return fmt.Sprintf("[第%d空]", index)
	})
}

func optionTextWithImageHint(topic qentity.Question, letter, option string) string {
	if strings.TrimSpace(option) != "" {
		return option
	}
	if len(topic.OptionImages[letter]) > 0 {
		return "[图片选项，请查看下方选项" + letter + "的图片]"
	}
	return "[空]"
}

// 单选题示例
func exampleSingleChoice() string {
	return `比如：
试卷名称：考试
题目类型：单选
题目内容：新中国是什么时候成立的
A. 1949年10月5日
B. 1949年10月1日
C. 1949年09月1日
D. 2002年10月1日

那么你应该回答选项B的内容：["1949年10月1日"]。注意不要携带A，B，C，D等选项前缀。`
}

// 多选题示例
func exampleMultipleChoice() string {
	return `比如：
试卷名称：考试
题目类型：多选题
题目内容：马克思关于资本积累的学说是剩余价值理论的重要组成部分...
A. 资本主义扩大再生产的源泉
B. 资本有机构成呈现不断降低趋势的根本原因
C. 社会财富占有两极分化的重要原因
D. 资本主义社会失业现象产生的根源

那么你应该回答选项A、B、D的内容：["资本主义扩大再生产的源泉","社会财富占有两极分化的重要原因","资本主义社会失业现象产生的根源"]，注意不要携带A，B，C，D等选项前缀。`
}

// 判断题示例
func exampleTrueFalse() string {
	return `比如：
试卷名称：考试
题目类型：判断
题目内容：新中国是什么时候成立是1949年10月1日吗？
A. 正确
B. 错误

那么你应该回答选项A的内容：["正确"]。注意不要携带A，B，C，D等选项前缀。`
}

// 填空题示例
func exampleFillInTheBlank() string {
	return ` 比如：
试卷名称：考试
题目类型：填空题
题目内容：1.主视图表达了部件的____和____。2.画出A-A断面图____。
第1空：____
第2空：____
第3空：____

那么你应该回答：["形状","尺寸","A-A断面图"]
注意：数组长度必须和空的数量完全一致，不能把整句话、编号或解释放进数组。`
}

func exampleShortAnswer() string {
	return `比如：
试卷名称：考试
题目类型：简答
题目内容：请简述中国和外国的国别 differences
答案：中国和外国的国别 differences

那么你应该回答： ["中国和外国的国别 differences"]`
}

// 名词解释
func exampleTermExplanationAnswer() string {
	return `比如：
试卷名称：考试
题目类型：名词解释
题目内容：绿色设计
答案：绿色设计是指在产品、建筑、工程或系统设计的全过程中，将环境保护和可持续发展理念融入其中的一种设计方法。

那么你应该回答： ["绿色设计是指在产品、建筑、工程或系统设计的全过程中，将环境保护和可持续发展理念融入其中的一种设计方法。"]`
}

// 论述题
func exampleEssayAnswer() string {
	return `比如：
试卷名称：考试
题目类型：论述题
题目内容：试述设计艺术的构成元素
答案：设计艺术的构成元素包括点、线、面、形体、色彩、质感与空间等。它们相互依存、互为补充，通过合理的组织和运用，形成和谐、统一而富有美感的设计作品。

那么你应该回答（回答字数不能少于500字）： ["设计艺术的构成元素包括点、线、面、形体、色彩、质感与空间等。它们相互依存、互为补充，通过合理的组织和运用，形成和谐、统一而富有美感的设计作品。"]`
}

// 连线题
func exampleMatchingAnswer() string {
	return `比如：
试卷名称：考试
题目类型：连线题
题目内容：
5.[连线题] 下列认知心理学家与其所做的经典研究之间的关系：

1、桑代克 ()
2、威特金 ()
3、凯利 ()
4、卡特尔 ()

A、迷箱实验
B、角色建构测验
C、16PF
D、棒框实验

答案：第一空：桑代克->迷箱实验、威特金->棒框实验、凯利->角色建构测验、卡特尔->16PF

那么你应该回答： ["桑代克->迷箱实验","威特金->棒框实验","凯利->角色建构测验","卡特尔->16PF"]`

}

func BuildVisualAiQuestionMessage(topic qentity.Question, imageDataMap map[string]qentity.QuestionImage) AIChatMessages {
	textMessages := BuildAiQuestionMessage(topic)
	if len(textMessages.Messages) == 0 {
		return textMessages
	}

	optionKeys := make([]string, 0, len(topic.OptionImages))
	for letter := range topic.OptionImages {
		optionKeys = append(optionKeys, letter)
	}
	sort.Strings(optionKeys)

	hasImage := len(topic.ContentImages) > 0
	for _, letter := range optionKeys {
		if len(topic.OptionImages[letter]) > 0 {
			hasImage = true
			break
		}
	}
	if !hasImage {
		return textMessages
	}

	var visualMessages []Message
	for _, msg := range textMessages.Messages {
		textStr := ContentToString(msg.Content)
		if msg.Role != "user" {
			visualMessages = append(visualMessages, NewTextMessage(msg.Role, textStr))
			continue
		}
		items := []MessageContent{{Type: "text", Text: textStr}}
		for index, img := range topic.ContentImages {
			items = append(items, MessageContent{
				Type: "text",
				Text: fmt.Sprintf("\n题目图片%d（请结合题干理解）：", index+1),
			})
			items = append(items, imageContent(img, imageDataMap))
		}
		for _, letter := range optionKeys {
			for index, img := range topic.OptionImages[letter] {
				items = append(items, MessageContent{
					Type: "text",
					Text: fmt.Sprintf("\n选项%s的图片%d：", letter, index+1),
				})
				items = append(items, imageContent(img, imageDataMap))
			}
		}
		visualMessages = append(visualMessages, NewVisualMessage(msg.Role, items))
	}

	return AIChatMessages{Messages: visualMessages}
}

func imageContent(img qentity.QuestionImage, imageDataMap map[string]qentity.QuestionImage) MessageContent {
	url := img.Url
	if data, ok := imageDataMap[img.Url]; ok {
		if data.Base64 != "" {
			url = "data:image/png;base64," + data.Base64
		} else if data.Url != "" {
			url = data.Url
		}
	}
	if url == "" && img.Base64 != "" {
		url = "data:image/png;base64," + img.Base64
	}
	return MessageContent{
		Type:     "image_url",
		ImageUrl: &ImageURL{URL: url},
	}
}
