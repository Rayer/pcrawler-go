package pcrawler

import (
	"time"
)

type CommentInfo struct {
	Type      int
	Committer string
	Timestamp time.Time
	Content   string
}

type AnalyzedInfo struct {
	DetectedFWNames    []string
	DetectedFWUPs      int
	DetectedFWDowns    int
	DetectedFWComments int
	TransformedHTML    string
}

//PTT Document raw contents
type PDocRaw struct {
	UniqueID            string
	Board               string
	Title               string
	Author              string
	RawArticleHtml      string
	Identifier          string
	PublicUrl           string
	PublishTime         time.Time
	ProcessTime         time.Time
	CommentInfoList     []CommentInfo
	CommentUpCount      int
	CommentDownCount    int
	CommentNeutralCount int
}
