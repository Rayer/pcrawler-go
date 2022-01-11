package pcrawler

import "time"

type SentenceInfo struct {
	Sentence  string
	timestamp time.Time
}

func NewSentenceInfo(sentence string, timestamp time.Time) *SentenceInfo {
	return &SentenceInfo{Sentence: sentence, timestamp: timestamp}
}

type UserSentenceInfoCollector struct {
	UserSentenceMap map[string][]*SentenceInfo
}

func NewUserSentenceInfoCollector() *UserSentenceInfoCollector {
	return &UserSentenceInfoCollector{UserSentenceMap: make(map[string][]*SentenceInfo)}
}

func (u *UserSentenceInfoCollector) Collect(doc *PDocRaw) error {
	for _, a := range doc.CommentInfoList {
		u.UserSentenceMap[a.Committer] = append(u.UserSentenceMap[a.Committer], NewSentenceInfo(a.Content, a.Timestamp))
	}
	return nil
}

func (u *UserSentenceInfoCollector) GetMap() map[string][]*SentenceInfo {
	return u.UserSentenceMap
}
