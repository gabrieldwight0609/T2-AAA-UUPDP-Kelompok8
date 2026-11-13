package handler

import (
	"aaa_bekend/internal/db"
	"aaa_bekend/internal/lmodel"
	"encoding/json"
	"net/http"
	"time"
)

func GetMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO: Ambil data profil user dari database
	profile, err := db.GetUserByID(r.Context(), user.ID)
	// TODO: Handle error
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func UpdateMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO: Decode request body
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Update data user di database
	if err := db.UpdateMyData(r.Context(), user.ID, req.Username); err != nil {
		http.Error(w, "Failed to update data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Data updated successfully",
		"username": req.Username,
	})
}

func GetMyAuditLog(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO: Ambil audit logs user dari database
	logs, err := db.GetUserAuditLogs(r.Context(), user.ID)
	// TODO: Handle error
	if err != nil {
		http.Error(w, "Failed to fetch audit logs", http.StatusInternalServerError)
		return
	}
	if logs == nil {
		logs = []lmodel.AuditLog{}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user.Username,
		"total": len(logs),
		"logs":  logs,
	})
}

func ExportMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO: Ambil data profil user
	profile, err := db.GetUserByID(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	// TODO: Ambil semua proyek dan audit logs milik user
	projects, err := db.GetUserProjects(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}
	logs, err := db.GetUserAuditLogs(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch audit logs", http.StatusInternalServerError)
		return
	}

	export := map[string]interface{}{
		"exported_at": time.Now(),
		"profile":     profile,
		"projects":    projects,
		"audit_logs":  logs,
	}

	w.Header().Set("Content-Disposition", "attachment; filename=my_data_export.json")
	json.NewEncoder(w).Encode(export)
}
func DeleteMyAccount(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO: Hapus akun user dan semua data terkait dari database
	if err := db.DeleteUserAccount(r.Context(), user.ID); err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Your account and all data have been permanently deleted",
	})
}
