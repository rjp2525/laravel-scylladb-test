package handlers

import (
	"audit-ms/domain"
	"audit-ms/storage"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AuditHandler struct {
	Storage *storage.AuditStorage
}

func (h *AuditHandler) CreateAuditLog(w http.ResponseWriter, r *http.Request) {
	var auditLog domain.AuditLog

	if err := json.NewDecoder(r.Body).Decode(&auditLog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if auditLog.ID == "" {
		auditLog.ID = uuid.New().String()
	}

	if auditLog.CreatedAt.IsZero() {
		auditLog.CreatedAt = time.Now()
	}
	if auditLog.UpdatedAt.IsZero() {
		auditLog.UpdatedAt = time.Now()
	}

	log.Printf("Populated audit log: %+v\n", auditLog)

	if err := h.Storage.SaveAuditLog(&auditLog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuditHandler) GetAuditLogsByTenantID(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "Missing tenant_id parameter", http.StatusBadRequest)
		return
	}

	logs, err := h.Storage.FindAuditLogsByTenantID(tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *AuditHandler) GetAllAuditLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.Storage.FindAllAuditLogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
