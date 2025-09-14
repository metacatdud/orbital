package components

import (
	"fmt"
	"strings"
)

type RegKey string

func (key RegKey) String() string {
	return string(key)
}

func (key RegKey) WithExtra(sep string, parts ...any) RegKey {
	joined := fmt.Sprint(parts...)
	result := fmt.Sprintf("%s%s%s", key, sep, joined)
	result = strings.ReplaceAll(result, " ", "")
	return RegKey(result)
}
