package pcrawler

import (
	_ "embed"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"reflect"
	"testing"
)

//go:embed test_resources/M.1559093660.A.946.html.json
var expectedDoc string

func TestParseSingleRawDocument(t *testing.T) {

	defer gock.Off()
	gock.New("https://www.ptt.cc/bbs/Gossiping").Get("/M.1559093660.A.946.html").Reply(200).File("test_resources/M.1559093660.A.946.html")

	type args struct {
		fromUrl string
	}
	tests := []struct {
		name                   string
		args                   args
		expectedResultFilename string
		wantErr                bool
	}{
		{
			"TestParseSingleDocument",
			args{"https://www.ptt.cc/bbs/Gossiping/M.1559093660.A.946.html"},
			"M.1559093660.A.946.html.json",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ParseSingleRawDocument(tt.args.fromUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSingleRawDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			expected := &PDocRaw{}
			err = json.Unmarshal([]byte(expectedDoc), expected)
			if err != nil {
				t.Errorf("error in unit test file, err : %v", err)
				return
			}

			expected.ProcessTime = got.ProcessTime
			equal := reflect.DeepEqual(*got, *expected)
			assert.True(t, equal)

		})
	}
}

func TestFetchArticleList(t *testing.T) {
	type args struct {
		board string
		start int
		end   int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"Fetch article with normal case",
			args{"Gossiping", 50, 52},
			nil,
			false,
		},
		{
			"Fetch article with start > end : should be auto repaired",
			args{"Gossiping", 52, 50},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FetchArticleList(tt.args.board, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchArticleList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("FetchArticleList() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func BenchmarkParseRangeDocument(b *testing.B) {
	type args struct {
		board string
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
		want []*PDocRaw
	}{
		{
			"Fetch article with single thread",
			args{"Gossiping", 100, 102},
			nil,
		},
	}
	for _, bm := range tests {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseRangeDocument(bm.args.board, bm.args.start, bm.args.end)
			}
		})
	}
}

func BenchmarkParseRangeDocumentAsync(b *testing.B) {
	type args struct {
		board string
		start int
		end   int
	}
	tests := []struct {
		name    string
		args    args
		wantRet []*PDocRaw
	}{
		{
			name: "Async go routine test",
			args: args{
				board: "Gossiping",
				start: 100,
				end:   102,
			},
			wantRet: nil,
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseRangeDocumentAsync(tt.args.board, tt.args.start, tt.args.end)
			}
		})
	}
}
