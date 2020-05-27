package pcrawler

import (
	"reflect"
	"testing"
)

func TestCreateRawDocument(t *testing.T) {
	type args struct {
		fromUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    *PDocRaw
		wantErr bool
	}{
		{
			"JustTestDraft",
			args{"https://www.ptt.cc/bbs/Gossiping/M.1559093660.A.946.html"},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseSingleRawDocument(tt.args.fromUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSingleRawDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("ParseSingleRawDocument() = %v, want %v", got, tt.want)
			//}
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
			"Another Test",
			args{"Gossiping", 50, 52},
			nil,
			false,
		},
		{
			"Another Test2",
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

func TestParseRangeDocument(t *testing.T) {
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
			"Another Test3",
			args{"Gossiping", 100, 110},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseRangeDocument(tt.args.board, tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				//t.Errorf("ParseRangeDocument() = %v, want %v", got, tt.want)
				t.Logf("This test is only meant to be complete running.")
			}
		})
	}
}
