package routes

import "net/http"

func YoutubeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from YT"))
}
