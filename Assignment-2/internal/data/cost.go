package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidCostFormat = errors.New("invalid cost format")

type Cost int32

func (r Cost) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d dollars", r)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

func (r *Cost) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidCostFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "dollars" {
		return ErrInvalidCostFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidCostFormat
	}

	*r = Cost(i)
	return nil
}
