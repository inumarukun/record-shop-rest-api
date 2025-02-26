package validator

import (
	"fmt"
	"record-shop-rest-api/model"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IRecordValidator interface {
	RecordValidate(record model.Record) error
}

type recordValidator struct{}

func NewRecordValidator() IRecordValidator {
	return &recordValidator{}
}

func ValidateReleaseYear(value interface{}) error {
	year, ok := value.(int)
	if !ok {
		return fmt.Errorf("release year must be an integer")
	}
	if year < 1000 || year > 9999 {
		return fmt.Errorf("release year must be a 4-digit number")
	}
	if year > time.Now().Year() {
		return fmt.Errorf("release year must not be in the future")
	}
	return nil
}

func (rv *recordValidator) RecordValidate(record model.Record) error {
	// is.Digit、validation.Lengthは数値の評価が出来ない
	// releaseYearStr := fmt.Sprintf("%d", record.ReleaseYear)
	return validation.ValidateStruct(&record,
		validation.Field(
			&record.Artist,
			validation.Required.Error("artist is required."),
		),
		validation.Field(
			&record.Title,
			validation.Required.Error("title is required."),
		),
		validation.Field(
			&record.Genre,
			validation.Required.Error("genre is required."),
		),
		validation.Field(
			&record.Style,
			validation.Required.Error("style is required."),
		),
		validation.Field(
			&record.ReleaseYear, // Fieldという縛りがあるので変数にできない
			validation.Required.Error("release year is required."),
			validation.By(ValidateReleaseYear), // カスタムバリデーションを適用
			// is.Digit.Error("release year must be a numeric value."),
			// validation.Length(4, 4).Error("release year must be a 4-digit number."),
			// validation.Max(time.Now().Year()).Error("release year must not be in the future.\n"),
		),
	)
}
