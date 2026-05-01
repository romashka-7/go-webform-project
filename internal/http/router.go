package http

import (
	"html/template"
	"net/http"
)
func NewRouter() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/form", formHandler)

	return mux

}

func homeHandler (w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("web/templates/form.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func healthHandler (w http.ResponseWriter, r *http.Request){
	w.Write([]byte("OK"))
}

func formHandler (w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}
	
	name := r.FormValue("name")
	email := r.FormValue("email")

	response := "Имя: " + name + " Email: " + email

	w.Write([]byte(response))
}
// рома любит алиса и еще кроче алису и вообще красоту алисы а еще красоту короче все