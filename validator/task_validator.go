package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/tasuke/udemy/model"
)

type ITaskValidator interface {
	TaskValidate(task model.Task) error
}

// 依存するものは特になし
type taskValidator struct{}

func NewTaskValidator() ITaskValidator {
	return &taskValidator{}
}

func (tv *taskValidator) TaskValidate(task model.Task) error {
	return validation.ValidateStruct(&task,
		validation.Field(
			&task.Title,
			validation.Required.Error("title is required"),
			validation.RuneLength(1, 10).Error("limited max 10 char"),
		),
	)
}
