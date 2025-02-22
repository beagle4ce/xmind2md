package models

// 一个Sheet
type Sheet struct {
	Id            string         `json:"id"`
	Class         string         `json:"class"`
	Title         string         `json:"title"`
	RootTopic     RootTopic      `json:"rootTopic"`
	Relationships []Relationship `json:"relationships"`
}

// 记录根节点
type RootTopic struct {
	Topic
	Class string `json:"class"`
}

// 记录Topic
type Topic struct {
	Id        string     `json:"id"`
	Title     string     `json:"title"`
	Notes     Notes      `json:"notes"`
	Labels    []string   `json:"labels"`
	Children  Children   `json:"children"`
	Summaries []Summarie `json:"summaries"`
}

// 记录Topic的子节点
type Children struct {
	Attached []Topic `json:"attached"`
	Detached []Topic `json:"detached"`
	Summary  []Topic `json:"summary"`
}

type Summarie struct {
	Id      string `json:"id"`
	rangee  string `json:"range"`
	TopicId string `json:"topicId"`
}

// 记录笔记
type Notes struct {
	Plain    Plain    `json:"plain"`
	RealHTML RealHTML `json:"realHTML"`
}

// 笔记内容(无渲染)
type Plain struct {
	Content string `json:"content"`
}

// 笔记内容(有渲染)
type RealHTML Plain

// 记录Topic关系
type Relationship struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	End1Id string `json:"end1Id"`
	End2Id string `json:"end2Id"`
}
