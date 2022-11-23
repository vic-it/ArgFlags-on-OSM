package server

import "net/http"

//https://www.youtube.com/watch?v=5BIylxkudaE
func StartServer() {
	http.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi there"))
	})
	http.ListenAndServe(":8080", nil)
}
