package server

import (
	"fmt"
	"net/http"
	"strconv"
)

// https://www.youtube.com/watch?v=5BIylxkudaE
func StartServer() {
	println("\nstarting http server...")
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/getpoint", getPointHandler())
	http.ListenAndServe(":8080", nil)
}

func getPointHandler() http.HandlerFunc {
	pointHandler := func(writer http.ResponseWriter, request *http.Request) {
		inputLon, _ := strconv.ParseFloat(request.URL.Query()["lon"][0], 64)
		inputLat, _ := strconv.ParseFloat(request.URL.Query()["lat"][0], 64)

		//TODO GET CLOSEST POINT HERE
		outputLon := inputLon
		outputLat := inputLat

		outputString := fmt.Sprintf("{lon: %f, lat: %f}", outputLon, outputLat)
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(outputString))
	}

	return pointHandler
}

func getRouteHandler() {

}
