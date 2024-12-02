package storage

import (
	"audit-ms/domain"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

type AuditStorage struct {
	cluster *gocql.ClusterConfig
}

func New() *AuditStorage {
	cluster := NewCluster()
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatalf("Failed to create ScyllaDB session: %v", err)
	}
	defer session.Close()

	if err := TryToCreateKeyspace(&session); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	if err := TryToCreateTable(&session); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return &AuditStorage{cluster: cluster}
}

func (s *AuditStorage) SaveAuditLog(auditLog *domain.AuditLog) error {
	session, err := gocqlx.WrapSession(s.cluster.CreateSession())
	if err != nil {
		return err
	}
	defer session.Close()

	if auditLog.CreatedAt.IsZero() {
		auditLog.CreatedAt = time.Now()
	}
	if auditLog.UpdatedAt.IsZero() {
		auditLog.UpdatedAt = time.Now()
	}

	// Map for holding non-empty fields
	fields := make(map[string]interface{})
	if auditLog.ID != "" {
		fields["id"] = auditLog.ID
	}
	if auditLog.UserType != "" {
		fields["user_type"] = auditLog.UserType
	}
	if auditLog.UserID != "" {
		fields["user_id"] = auditLog.UserID
	}
	if auditLog.Event != "" {
		fields["event"] = auditLog.Event
	}
	if auditLog.AuditableType != "" {
		fields["auditable_type"] = auditLog.AuditableType
	}
	if auditLog.AuditableID != "" {
		fields["auditable_id"] = auditLog.AuditableID
	}
	if auditLog.OldValues != "" {
		fields["old_values"] = auditLog.OldValues
	}
	if auditLog.NewValues != "" {
		fields["new_values"] = auditLog.NewValues
	}
	if auditLog.URL != "" {
		fields["url"] = auditLog.URL
	}
	if auditLog.IPAddress != "" {
		fields["ip_address"] = auditLog.IPAddress
	}
	if auditLog.UserAgent != "" {
		fields["user_agent"] = auditLog.UserAgent
	}
	if auditLog.Tags != "" {
		fields["tags"] = auditLog.Tags
	}
	fields["created_at"] = auditLog.CreatedAt.Format(time.RFC3339)
	fields["updated_at"] = auditLog.UpdatedAt.Format(time.RFC3339)

	// Build query dynamically
	columns := []string{}
	values := []string{}
	for key, value := range fields {
		columns = append(columns, key)
		switch v := value.(type) {
		case string:
			if key == "user_id" || key == "id" || key == "auditable_id" {
				// UUIDs should not be quoted
				values = append(values, v)
			} else {
				values = append(values, fmt.Sprintf("'%s'", escapeString(v)))
			}
		default:
			values = append(values, fmt.Sprintf("'%v'", value))
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO audit.logs (%s) VALUES (%s)",
		join(columns, ", "), // Joins column names with a comma
		join(values, ", "),  // Joins values with a comma
	)

	log.Printf("Saving audit log: %+v\n", auditLog)
	log.Printf("Generated query: %s\n", query)

	// Execute the query
	err = session.Query(query, nil).Exec()
	if err != nil {
		log.Printf("Query execution failed: %v\n", err)
		return err
	}

	log.Println("Audit log saved successfully")
	return nil
}

// Helper function to join strings with a delimiter
func join(items []string, delimiter string) string {
	return strings.Join(items, delimiter)
}

func (s *AuditStorage) FindAuditLogsByTenantID(tenantID string) ([]domain.AuditLog, error) {
	session, err := gocqlx.WrapSession(s.cluster.CreateSession())
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var logs []domain.AuditLog

	query := fmt.Sprintf("SELECT * FROM audit.logs WHERE tenant_id = '%s' ALLOW FILTERING", tenantID)

	err = session.Query(query, nil).Select(&logs)
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return nil, fmt.Errorf("tenant not found")
	}

	return logs, nil
}

func (s *AuditStorage) FindAllAuditLogs() ([]domain.AuditLog, error) {
	session, err := gocqlx.WrapSession(s.cluster.CreateSession())
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var logs []domain.AuditLog

	query := "SELECT * FROM audit.logs"

	err = session.Query(query, nil).Select(&logs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func escapeString(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(value, "'", "''"), "\\", "\\\\")
}
