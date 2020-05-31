package pcrawler

import (
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	diff "github.com/yudai/gojsondiff"
	"io"
	"io/ioutil"
	"testing"
)

func TestParseIndexContent(t *testing.T) {

	type args struct {
		contentIo io.ReadCloser
	}
	tests := []struct {
		name                   string
		args                   args
		expectedResultFilename string
		wantErr                bool
	}{
		{
			name: "Pinned object",
			args: args{
				contentIo: GetResourceFromFile("test_resources/index_common.html"),
			},
			expectedResultFilename: "test_resources/index_common.html.json",
			wantErr:                false,
		},
		{
			name: "Non-pinned object",
			args: args{
				contentIo: GetResourceFromFile("test_resources/index_without_pinned.html"),
			},
			expectedResultFilename: "test_resources/index_without_pinned.html.json",
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIndexContent(tt.args.contentIo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIndexContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//bytes, _ := json.Marshal(got)
			//ioutil.WriteFile("/tmp/ti" + tt.name, bytes, 0644)

			differ := diff.New()
			aBytes, _ := json.Marshal(got)
			bBytes, _ := ioutil.ReadFile(tt.expectedResultFilename)
			result, _ := differ.Compare(aBytes, bBytes)

			assert.Equal(t, len(result.Deltas()), 0, fmt.Sprintf("Delta : %+v", result.Deltas()))

		})
	}
}
