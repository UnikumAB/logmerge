package compress

import "testing"

func TestDetectGzip(t *testing.T) {
	type args struct {
		filename string
	}

	tests := []struct {
		name       string
		args       args
		wantedSize int64
		wantisGz   bool
		wantErr    bool
	}{
		{
			name:       "uncompressed file",
			args:       args{filename: "testdata/lorem2.txt"},
			wantisGz:   false,
			wantedSize: int64(336),
			wantErr:    false,
		},
		{
			name:       "compressed file",
			args:       args{filename: "testdata/lorem.txt.gz"},
			wantisGz:   true,
			wantedSize: int64(336),
			wantErr:    false,
		},
		{
			name:       "bigger compressed file",
			args:       args{filename: "testdata/lorem3.txt.gz"},
			wantisGz:   true,
			wantedSize: int64(75600),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			isGzip, fileSize, err := DetectGzip(tt.args.filename)

			if (err != nil) != tt.wantErr {
				t.Errorf("DetectGzip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if isGzip != tt.wantisGz {
				t.Errorf("DetectGzip() isGzip = %v, wantisGz %v", isGzip, tt.wantisGz)
			}

			if fileSize != tt.wantedSize {
				t.Errorf("DetectGzip() fileSize = %v, wantisGz %v", fileSize, tt.wantedSize)
			}
		})
	}
}
