package route

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"token/web/controller"
)

func Init() {
	router := mux.NewRouter()

	router.HandleFunc("/generate", controller.GenerateToken).Methods("POST")
	router.HandleFunc("/assign", controller.AssignToken).Methods("GET")
	router.HandleFunc("/unblock", controller.UnblockToken).Methods("GET")
	router.HandleFunc("/delete", controller.DeleteToken).Methods("DELETE")
	router.HandleFunc("/keepalive", controller.KeepTokenAlive).Methods("GET")
	router.HandleFunc("/health", controller.Health).Methods("GET")

	router.Handle("/", router)

	// Listen on port 8082
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Println("Error starting server:", err)
	}
}
