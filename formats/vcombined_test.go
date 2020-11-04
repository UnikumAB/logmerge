package formats

import (
	"testing"
	"time"
)

func Test_vCommonLine_ParseLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Simple Logline",
			args:    args{line: "test.example.com:443 10.245.49.59 - - [01/Nov/2020:06:59:19 +0100] \"GET /app/page.html?__id=123456&currentUrl=start.html%3F__id%3D123456 HTTP/2.0\" 303 0 \"https://example.com/app/page?__id=123456\" \"Mozilla/5.0 (Linux; Android 10; SM-G965F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.114 Mobile Safari/537.36\""},
			want:    time.Unix(1604210359, 0),
			wantErr: false,
		},
		{
			name:    "October Line",
			args:    args{line: "test.example.com:443 10.174.22.71 - - [31/Oct/2020:04:42:49 +0100] \"GET /app/subscribe/calendar/af0628e7-9b04-4fb5-a8bb-96584d9004ee HTTP/2.0\" 200 281 \"-\" \"iOS/14.0.1 (18A393) dataaccessd/1.0\""},
			want:    time.Unix(1604115769, 0),
			wantErr: false,
		},
	}
	V, err := NewVCombinedParser()
	if err != nil {
		t.Fatalf("Cannot instanciate VCommonParser")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := V.ParseLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Source != tt.args.line {
				t.Errorf("ParseLine() Source got = %v, want %v", got.Source, tt.args.line)
			}
			if !got.Timestamp.Equal(tt.want) {
				t.Errorf("ParseLine() got = %v, want %v", got.Timestamp, tt.want)
			}
		})
	}
}
