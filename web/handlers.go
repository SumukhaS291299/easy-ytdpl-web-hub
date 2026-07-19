package web

// Web package for UI

import (
	"encoding/json"
	"html/template"
	"net/http"

	"ytdpldownloader/utils"
)

var templates = template.Must(
	template.ParseFS(
		Files,
		"templates/*.html",
	),
)

func Home(w http.ResponseWriter, r *http.Request) {
	success := r.URL.Query().Get("success") == "true"
	err := templates.ExecuteTemplate(w, "index.html", struct{ Success bool }{Success: success})
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func ProbeytdplHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	url := r.FormValue("url")

	data := utils.Probeytdpl(url)
	if data == nil {
		http.Error(w, "Failed to probe URL", http.StatusInternalServerError)
		return
	}

	var formats []utils.FormatInfo
	if err := json.Unmarshal(data, &formats); err != nil {
		http.Error(w, "Failed to parse formats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := struct {
		URL     string
		Formats []utils.FormatInfo
	}{
		URL:     url,
		Formats: formats,
	}

	err := templates.ExecuteTemplate(w, "formatSelection.html", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DownloadytdplHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	url := r.FormValue("url")
	format := r.FormValue("format")
	// quality := r.FormValue("quality")

	utils.Downloadytdpl(url, format)

	http.Redirect(w, r, "/?success=true", http.StatusSeeOther)
}
