package routes

import (
	"github.com/SDI-Boston/filemanager_go_node/controller"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/files/{userId}/{fileId}", controller.DownloadFileHandler).Methods("GET")
	router.HandleFunc("/files_backup/{userId}/{fileId}", controller.DownloadFileHandler).Methods("GET")
	return router
}
