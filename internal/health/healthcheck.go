package health

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func StartHealthEndpoint() {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))
	r.HandleFunc("/health", send200).Methods(http.MethodGet)
	err := http.ListenAndServe("0.0.0.0:8080", r)
	if err != nil {
		log.WithFields(log.Fields{
			"Caller": "StartHealthEndpoint",
		}).Error(fmt.Sprintf("Error creating health endpoint: %s", err.Error()))
	}
}

func send200(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte{})
	if err != nil {
		log.WithFields(log.Fields{
			"Caller": "send200",
		}).Error(fmt.Sprintf("Error answering health endpoint: %s", err.Error()))
	}
}
