package pcrawler

import (
	"encoding/json"
	"github.com/magiconair/properties/assert"
	diff "github.com/yudai/gojsondiff"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"reflect"
	"testing"
)

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

			//Compare result and expected
			aString, err := json.Marshal(got)
			if err != nil {
				t.Error(err.Error())
			}
			bString, err := ioutil.ReadFile("test_resources/M.1559093660.A.946.html.json")
			if err != nil {
				t.Error(err.Error())
			}

			differ := diff.New()
			d, err := differ.Compare(aString, bString)
			if err != nil {
				t.Error(err.Error())
			}
			//expect only "ProcessTime" different
			if d == nil {
				t.Error("Compare is nil, kinda strange....")
				return
			}
			assert.Equal(t, len(d.Deltas()), 1, "Parsed item is not match!")

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
			args{"Gossiping", 100, 105},
			nil,
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			if got := ParseRangeDocument(tt.args.board, tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				//t.Errorf("ParseRangeDocument() = %v, want %v", got, tt.want)
				b.Logf("This test is only meant to be complete running.")
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
				end:   105,
			},
			wantRet: nil,
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			if gotRet := ParseRangeDocumentAsync(tt.args.board, tt.args.start, tt.args.end); !reflect.DeepEqual(gotRet, tt.wantRet) {
				//t.Errorf("ParseRangeDocumentAsync() = %v, want %v", gotRet, tt.wantRet)
				b.Logf("This test is only meant to be complete running.")
			}
		})
	}
}
