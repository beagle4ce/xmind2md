package markdown

import (
	"strconv"
	"strings"
	"xmind2md/models"
)

const (
	ROOT_TOPIC_LEVEL = 0
	FIRST_HEADER     = 2
	SECOND_HEADER    = 3
	THIRD_LEVEL      = 4
)

// 主方法: 解析Sheet
func ToMarkDown(sheets []models.Sheet) map[string][]string {
	mdMap := make(map[string][]string)
	for idx, sheet := range sheets {
		// 获取单个sheet内容,并生成markdown字符串切片
		mdStrSlice := sheetContentAnalisis(&sheet)
		// mdMap-Key设置为sheet的title+sheet的index作为文件名
		// mdMap-Value设置为markdown字符串切片
		mdMap[sheet.Title+strconv.Itoa(idx)] = mdStrSlice
	}
	return mdMap
}

// 子方法:生成markdown字符串切片
func sheetContentAnalisis(sheet *models.Sheet) []string {
	// 建立Relationships的映射表
	// relMap := make(map[string]models.Relationship)
	mdStrSlice := make([]string, 0, 10*1024)
	// 渲染一级标题和对应内容
	mdStrSlice = append(mdStrSlice, rootTopicRender(&sheet.RootTopic)...)
	// 只有存在子节点时才渲染
	mdStrSlice = append(mdStrSlice, headerChildrenTopicRender(&sheet.RootTopic.Children, FIRST_HEADER)...)
	// 渲染自由主题
	mdStrSlice = append(mdStrSlice, detachedTopicRender(&sheet.RootTopic.Children, FIRST_HEADER)...)
	return mdStrSlice
}

// 渲染一级标题和对应内容
func rootTopicRender(rootTopic *models.RootTopic) []string {
	rootTopicSlice := []string{"# ", rootTopic.Title, "\n"}
	// 渲染其余内容,笔记,标签,链接等
	rootTopicSlice = append(rootTopicSlice, labelsRender(rootTopic.Labels, ROOT_TOPIC_LEVEL)...)
	rootTopicSlice = append(rootTopicSlice, notesRender(&rootTopic.Notes)...)
	return rootTopicSlice
}

// 渲染二级,三级和后续标题的对应内容
func headerChildrenTopicRender(children *models.Children, level int) []string {
	// 若不存在子节点,则返回空字符串切片
	if children == nil {
		return []string{}
	}
	childrenTopicSlice := make([]string, 0, 10)
	for _, topic := range children.Attached {
		// 只限制渲染二三级Header Topic,后续级别标题均使用 * 号开头
		switch level {
		case FIRST_HEADER:
			// 渲染二级标题, 二级标题一定是##开头
			childrenTopicSlice = append(childrenTopicSlice, "\n## ", topic.Title, "\n")
			childrenTopicSlice = append(childrenTopicSlice, labelsRender(topic.Labels, level)...)
			childrenTopicSlice = append(childrenTopicSlice, notesRender(&topic.Notes)...)
			// 递归一层, 渲染三级标题
			childrenTopicSlice = append(childrenTopicSlice, headerChildrenTopicRender(&topic.Children, level+1)...)
		case SECOND_HEADER:
			// 渲染三级标题, 三级标题一定是###开头
			childrenTopicSlice = append(childrenTopicSlice, "\n### ", topic.Title, "\n")
			childrenTopicSlice = append(childrenTopicSlice, labelsRender(topic.Labels, level)...)
			childrenTopicSlice = append(childrenTopicSlice, notesRender(&topic.Notes)...)
			// 递归渲染后续标题
			childrenTopicSlice = append(childrenTopicSlice, nodeChildrenTopicRender(&topic.Children, level+1)...)
		}
	}

	return childrenTopicSlice
}

// 自由主题渲染
func detachedTopicRender(children *models.Children, level int) []string {
	// 若不存在子节点,则返回空字符串切片
	if children == nil || children.Detached == nil {
		return []string{}
	}
	childrenTopicSlice := make([]string, 0, 10)
	for _, topic := range children.Detached {
		// 只限制渲染二三级Header Topic,后续级别标题均使用 * 号开头
		switch level {
		case FIRST_HEADER:
			// 渲染二级标题, 二级标题一定是##开头
			childrenTopicSlice = append(childrenTopicSlice, "\n## ", topic.Title, "\n")
			childrenTopicSlice = append(childrenTopicSlice, labelsRender(topic.Labels, level)...)
			childrenTopicSlice = append(childrenTopicSlice, notesRender(&topic.Notes)...)
			// 递归一层, 渲染三级标题
			childrenTopicSlice = append(childrenTopicSlice, detachedTopicRender(&topic.Children, level+1)...)
		case SECOND_HEADER:
			// 渲染三级标题, 三级标题一定是###开头
			childrenTopicSlice = append(childrenTopicSlice, "\n### ", topic.Title, "\n")
			childrenTopicSlice = append(childrenTopicSlice, labelsRender(topic.Labels, level)...)
			childrenTopicSlice = append(childrenTopicSlice, notesRender(&topic.Notes)...)
			// 递归渲染后续标题
			childrenTopicSlice = append(childrenTopicSlice, nodeChildrenTopicRender(&topic.Children, level+1)...)
		}
	}

	return childrenTopicSlice
}

// 渲染后续级数标题和对应内容
func nodeChildrenTopicRender(children *models.Children, level int) []string {
	// 若不存在子节点,则返回空字符串切片
	if children == nil {
		return []string{}
	}
	childrenTopicSlice := make([]string, 0, 10)
	for _, topic := range children.Attached {
		// 每个Topic都不同,可能有多个子节点, 渲染风格也会不同, 需要及时计算
		childrenTopicSlice = append(childrenTopicSlice, culculatePrefix(level, "*"), topic.Title, "\n")
		// 递归解析子节点, 递归深度标记+1
		childrenTopicSlice = append(childrenTopicSlice, nodeChildrenTopicRender(&topic.Children, level+1)...)
	}
	// TODO: 其余内容,笔记,标签,链接等
	return childrenTopicSlice
}

/*
计算非Header Topic的渲染前缀, 默认为 * 号开头.
*/
func culculatePrefix(level int, symbol string) string {
	prefix := []string{"\n"}
	for i := 0; i < level-THIRD_LEVEL; i++ {
		prefix = append(prefix, "  ")
	}
	return strings.Join(append(prefix, symbol, " "), "")
}

// 渲染标签
func labelsRender(lables []string, level int) []string {
	if lables == nil {
		return []string{}
	}
	newLabels := make([]string, 0, 3)
	var strbuilder strings.Builder
	for idx, label := range lables {
		// 使用markdown ==语法== 渲染
		strbuilder.WriteString("==")
		strbuilder.WriteString(label)
		strbuilder.WriteString("==")
		lables[idx] = strbuilder.String()
		strbuilder.Reset()
	}
	// 渲染前缀, 计算缩进量
	newLabels = append(newLabels, culculatePrefix(level, " "))
	newLabels = append(newLabels, strings.Join(lables, ", "), "\n")
	return newLabels
}

// 渲染笔记, 将笔记内容渲染为markdown的引用语法
func notesRender(notes *models.Notes) []string {
	if notes == nil {
		return []string{}
	}
	return []string{}
}

// 优先解析Relationships块, 双重映射, endId1和endId2都映射为map的key
// 解析过程中, 遇到Relationship块就用占位符记录, "rel:<ID>"
// 进入RootTopic块后,每个Topic都要检查Summaries

// RootTopic不存在Summary和Summaries,
// Summaries和Summary一定与同层级的children配对
