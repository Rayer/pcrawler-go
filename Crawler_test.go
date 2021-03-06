package pcrawler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func getGockBaseUrl(name string) string {
	return "https://www.ptt.cc/bbs/" + name
}

func TestCrawler_IndexRawTest(t *testing.T) {
	c := NewCrawler("Gossiping")
	u, err := url.Parse("https://www.ptt.cc/bbs/Gossiping/index100.html")
	if err != nil {
		t.FailNow()
	}
	content, err := c.getContent(u)
	if err != nil {
		t.FailNow()
	}
	print(content)
}

func TestCrawler_RegexCapture(t *testing.T) {
	link := "/bbs/Gossiping/index12164.html"
	r := regexp.MustCompile("/bbs/.+/index(\\d+)\\.html")
	var index int
	if find := r.FindStringSubmatch(link); len(find) > 1 {
		index, _ = strconv.Atoi(find[1])
		index += 1
	}
	assert.Equal(t, 12165, index)
}

func TestCrawler_ParseDocumentRaw(t *testing.T) {
	c := NewCrawler("NothingToDo")
	targetUrl, _ := url.Parse("https://www.ptt.cc/bbs/Gossiping/M.1569751115.A.5A7.html")

	doc, _ := c.ParseDocument(targetUrl)
	fmt.Println(doc)
}

func TestCrawler_ParseDocument(t *testing.T) {
	defer gock.Off()

	gock.New(getGockBaseUrl("case1")).Get("/case1.html").Reply(200).File("test_resources/M.1569901516.A.5F2.html")

	t1, _ := url.Parse(getGockBaseUrl("case1") + "/case1.html")

	type args struct {
		url *url.URL
	}
	tests := []struct {
		name    string
		fields  *Crawler
		args    args
		want    *PDocRaw
		wantErr bool
	}{
		{
			name:   "Common documents(M.1569901516.A.5F2.html)",
			fields: NewCrawler("case1"),
			args:   args{t1},
			want: &PDocRaw{
				UniqueID:          "",
				Board:             "",
				Title:             "[問卦] 有無永齡基金會的八卦？",
				Author:            "Godfrey0216 (Godfrey)",
				Identifier:        "",
				PublicUrl:         "https://www.ptt.cc/bbs/case1/case1.html",
				PublishTime:       time.Time{},
				ProcessTime:       time.Time{},
				CommitterInfoList: make([]CommitInfo, 7),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				Board: tt.fields.Board,
				//IndexInfo:    tt.fields.IndexInfo,
				sharedClient: tt.fields.sharedClient,
			}
			got, err := c.ParseDocument(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("Strange, 'got' is nil but err is nil too!?")
				return
			}

			fmt.Printf("%+v\n", got.CommitterInfoList)
			assert.Equal(t, tt.want.PublicUrl, got.PublicUrl)
			assert.Equal(t, tt.want.Author, got.Author)
			assert.Equal(t, len(tt.want.CommitterInfoList), len(got.CommitterInfoList))

		})
	}
}

func TestCrawler_createIndexUrl(t *testing.T) {
	type fields struct {
		Board        string
		sharedClient http.Client
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *url.URL
	}{
		{
			name: "ShouldGenerateHeadIndex",
			fields: fields{
				Board:        "Gossiping",
				sharedClient: http.Client{},
			},
			args: args{
				index: 0,
			},
			want: ErrorStripper(url.Parse("https://www.ptt.cc/bbs/Gossiping/index.html")).(*url.URL),
		},
		{
			name: "ShouldGenerateCorrectIndex",
			fields: fields{
				Board:        "Gossiping",
				sharedClient: http.Client{},
			},
			args: args{
				index: 125,
			},
			want: ErrorStripper(url.Parse("https://www.ptt.cc/bbs/Gossiping/index125.html")).(*url.URL),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				Board:        tt.fields.Board,
				sharedClient: tt.fields.sharedClient,
			}
			if got := c.createIndexUrl(tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createIndexUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
