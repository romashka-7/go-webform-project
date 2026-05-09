package handlers

import (
	"html/template"
	"net/http"
)

func FormHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"web/templates/form.html",
	)

	if err != nil {
		http.Error(
			w,
			"Ошибка загрузки шаблона",
			http.StatusInternalServerError,
		)

		return
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(
			w,
			"Ошибка рендера шаблона",
			http.StatusInternalServerError,
		)

		return
	}
}
