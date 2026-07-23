package web

import (
	"io/fs"
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(Files, "static")
	if err != nil {
		panic(err)
	}

	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.FS(staticFS)),
		),
	)

	mux.HandleFunc("/", Home)
	mux.HandleFunc("/ytdplprobe", ProbeytdplHandler)
	mux.HandleFunc("/playlist/download", PlaylistDownloadHandler)
	mux.HandleFunc("/ytdpldownload", DownloadytdplHandler)

	return mux
}
