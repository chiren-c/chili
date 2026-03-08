package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Status  ArticleStatus
	Content string
	// 作者
	Author Author
	Ctime  time.Time
	Utime  time.Time
}

// Abstract 取部分作为摘要
func (a Article) Abstract() string {
	cs := []rune(a.Content)
	if len(cs) < 100 {
		return a.Content
	}
	return string(cs[:100])
}

func (a Article) Published() bool {
	return a.Status == ArticleStatusPublished
}

type ArticleAuthor struct {
	Id      int64
	Title   string
	Status  ArticleStatus
	Content string
	// 作者
	Author Author
	Ctime  time.Time
	Utime  time.Time
}

type ArticleReader struct {
	Id      int64
	Title   string
	Status  ArticleStatus
	Content string
	// 作者
	Author Author
	Ctime  time.Time
	Utime  time.Time
}

type ArticleStatus uint8

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

func (s ArticleStatus) ToString() string {
	switch s {
	case ArticleStatusUnpublished:
		return "未发表"
	case ArticleStatusPublished:
		return "已发表"
	case ArticleStatusPrivate:
		return "仅自己可见"
	default:
		return "未知状态"
	}
}

const (
	// ArticleStatusUnknown 未知状态
	ArticleStatusUnknown ArticleStatus = iota
	// ArticleStatusUnpublished 未发表
	ArticleStatusUnpublished
	// ArticleStatusPublished 已发表
	ArticleStatusPublished
	// ArticleStatusPrivate 仅自己可见
	ArticleStatusPrivate
)

type Author struct {
	Id   int64
	Name string
}
