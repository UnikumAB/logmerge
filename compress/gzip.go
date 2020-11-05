package compress

import (
	"bytes"
	"encoding/binary"
	"net/http"
	"os"

	"github.com/UnikumAB/logmerge/utils"
	"github.com/pkg/errors"
)

func DetectGzip(filename string) (bool, int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, 0, errors.WithMessagef(err, "Cannot open file %s", filename)
	}
	defer utils.CheckedClose(file)
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return false, 0, errors.WithMessagef(err, "Failed reading from %s", filename)
	}
	contentType := http.DetectContentType(buff)
	stat, err := file.Stat()
	if err != nil {
		return false, 0, errors.WithMessagef(err, "failed to get stats for %v", filename)
	}

	switch contentType {
	case "application/x-gzip", "application/zip":
		lastBytes := make([]byte, 4)
		_, err = file.ReadAt(lastBytes, stat.Size()-4)
		if err != nil {
			return false, 0, errors.WithMessagef(err, "Failed to read last 4 bytes from %v", filename)
		}
		buf := bytes.NewBuffer(lastBytes)
		var decompressedSize int32
		err = binary.Read(buf, binary.LittleEndian, &decompressedSize)
		if err != nil {
			return false, 0, errors.WithMessagef(err, "Failed to decode filesize for %v", filename)
		}
		return true, int64(decompressedSize), nil
	default:
		return false, stat.Size(), nil
	}
}
