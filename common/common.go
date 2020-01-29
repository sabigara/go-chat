package common

import "bytes"

const DELIMITER = '\n'

func EndsWithDelimiter(b []byte) bool {
	return b[len(b)-1:][0] == DELIMITER
}

func TrimBytes(b []byte) []byte {
	b = bytes.TrimRight(b, "\x00")
	b = bytes.TrimRight(b, "\x0a")
	b = bytes.TrimRight(b, "\x0d")
	return b
}
