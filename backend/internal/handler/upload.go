package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ionutale/restaurant-orders-manager/internal/auth"
)

type UploadHandler struct {
	db       *pgxpool.Pool
	uploadDir string
}

func NewUploadHandler(db *pgxpool.Pool, uploadDir string) *UploadHandler {
	os.MkdirAll(uploadDir, 0755)
	return &UploadHandler{db: db, uploadDir: uploadDir}
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	if claims == nil {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, "file too large or invalid", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		respondError(w, "allowed: jpg, jpeg, png, webp, gif", http.StatusBadRequest)
		return
	}

	name := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.TrimSuffix(header.Filename, ext), ext)
	path := filepath.Join(h.uploadDir, name)

	dst, err := os.Create(path)
	if err != nil {
		respondError(w, "could not save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondError(w, "could not write file", http.StatusInternalServerError)
		return
	}

	url := "/uploads/" + name

	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "file.uploaded", "file", nil, map[string]string{"name": name, "url": url})
	}

	respondJSON(w, map[string]string{"url": url, "name": name})
}
