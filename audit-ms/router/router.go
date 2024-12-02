package router

import (
	"audit-ms/handlers"
	"audit-ms/storage"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	storage := storage.New()
	handler := &handlers.AuditHandler{Storage: storage}

	r := mux.NewRouter()
	r.HandleFunc("/audit", handler.CreateAuditLog).Methods("POST")
	r.HandleFunc("/audit", handler.GetAuditLogsByTenantID).Methods("GET")
	r.HandleFunc("/audits", handler.GetAllAuditLogs).Methods("GET")
	return r
}
