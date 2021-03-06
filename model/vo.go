package model

import "time"

// TopicVO 接收前端发送的创建投票数据
type TopicVO struct {
	// topic
	Title       string `binding:"required"`
	Description string
	Deadline    time.Time `binding:"required"`

	// topic_set
	SelectType int
	Anonymous  int
	ShowResult int
	Password   string

	// topic_option
	Option []struct {
		OptionContent string `binding:"required"`
	}
}

// TopicFriendly 返回友好数据方便前端使用，主要用在投票总览页
type TopicFriendly struct {
	Id           int
	Title        string
	StuName      string
	SelectType   string // 直接返回 string 而不是 int 标识
	Anonymous    string
	ShowResult   string
	NeedPassword string
	Deadline     string
}

type VoteSingleVO struct {
	Record  *VoteRecord
	Votes   int32 // 票数
	TopicId int
}

type VoteMultipleVO struct {
	Record  []*VoteRecord
	Votes   int32 // 票数
	TopicId int
}

// VoteResultVO 投票结果
type VoteResultVO struct {
	OptionId      string  `json:"option_id"`
	OptionContent string  `json:"option_content"` // 选项
	Votes         int     `json:"votes"`          // 票数
	Percentage    float32 `json:"percentage"`     // 占比
}
