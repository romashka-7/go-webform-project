package handlers

import (
	"html/template"
	"net/http"
)

func AdminPanelHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/admin.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки admin.html", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
