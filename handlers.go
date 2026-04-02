package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) PostHeartbeat(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("device_id")

	var req HeartbeatRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to decode heartbeat request")
		return
	}

	sentAt, err := parseSentAt(req.SentAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.store.AddHeartbeat(deviceID, sentAt); err != nil {
		if errors.Is(err, ErrDeviceNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PostStats(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("device_id")

	var req UploadStatsRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to decode upload stats request")
		return
	}

	if _, err := parseSentAt(req.SentAt); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.store.AddUploadStat(deviceID, req.UploadTime); err != nil {
		if errors.Is(err, ErrDeviceNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("device_id")

	resp, err := h.store.GetStats(deviceID)
	if err != nil {
		switch {
		case errors.Is(err, ErrDeviceNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrNoStatsAvailable):
			w.WriteHeader(http.StatusNoContent)
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func decodeJSON(r *http.Request, v any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func parseSentAt(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid sent_at value: %q", value)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Msg: msg})
}
