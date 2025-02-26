package common

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// 任意の型のスライスを変換する汎用関数
func MapSlice[T, U any](input []T, transformFunc func(T) U) []U {
	var result []U
	for _, item := range input {
		result = append(result, transformFunc(item))
	}
	return result
}

// validateエラーをフォーマット（複数発生すると;で連結表示されるため、改行で整える）
func HandleValidationError(err error) string {
	if errs, ok := err.(validation.Errors); ok {
		var messages []string
		// errsはvalidation.Errors型のMap、よってkeyとvalueが返される
		for _, validationErr := range errs {
			// messages = append(messages, fmt.Sprintf("%s: %s", field, validationErr.Error()))
			messages = append(messages, validationErr.Error())
		}
		return strings.Join(messages, "\n")
	}
	return ""
}
