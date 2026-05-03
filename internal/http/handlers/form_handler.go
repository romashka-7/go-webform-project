package handlers

import (
	"net/http"
	"webform-go/internal/domain"
	"webform-go/internal/validation.go"
)

func FormHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	application := domain.Application {
		Name : r.FormValue("name"),
		Email : r.FormValue("email"),
	}

	err := validation.ValidateApplication(application)

	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	response := "Заявка успешно принята!\n" + "Имя: " + application.Name + " Email: " + application.Email

	w.Write([]byte(response))
}