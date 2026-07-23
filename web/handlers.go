package web

// Web package for UI

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"

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
	getPlaylist := r.FormValue("playlist") == "true"

	if !getPlaylist {
		singleVideo(url, w, r)
	} else {
		MultiVideo(url, w, r)
	}
}

func PlaylistDownloadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	videos := r.Form["videos"]
	format := r.FormValue("format")

	log.Debug("Selected videos:", len(videos))

	log.Debug(videos)

	for _, url := range videos {
		log.Debug("Downloading URL:\t" + url)
		// TODO Worker pool or async calls
		utils.Downloadytdpl(url, format)
	}
	http.Redirect(w, r, "/?success=true", http.StatusSeeOther)
}

func MultiVideo(rawURL string, w http.ResponseWriter, r *http.Request) {
	log.Info("Using URL:\t:" + rawURL)
	playlist := utils.FlatPlaylist(rawURL)

	err := templates.ExecuteTemplate(w, "playlistSelection.html", playlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func singleVideo(rawURL string, w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(rawURL)
	if err != nil {
		http.Error(w, "Failed to parse formats: "+err.Error(), http.StatusInternalServerError)
	}
	q := u.Query()
	q.Del("list")
	q.Del("index") // optional
	u.RawQuery = q.Encode()
	singleURL := u.String()
	log.Info("Using URL:\t:" + singleURL)
	data := utils.Probeytdpl(singleURL)
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
		URL:     singleURL,
		Formats: formats,
	}

	err = templates.ExecuteTemplate(w, "formatSelection.html", templateData)
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
