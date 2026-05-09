package validation

import (
	"errors"
	"strings"
	"webform-go/internal/domain"
)

func ValidateApplication(application domain.Application) error {

	if strings.TrimSpace(application.Name) == "" {
		return errors.New("Имя не может быть пустым")
	}

	if strings.TrimSpace(application.Email) == "" {
		return errors.New("Email не может быть пустым")
	}

	if !strings.Contains(application.Email, "@") {
		return errors.New("Email должен содержать символ @")
	}

	if len(application.Languages) == 0 {
		return errors.New("выберите хотя бы один язык")
	}

	if !application.Agreement {
		return errors.New("необходимо согласие")
	}
	return nil
}
