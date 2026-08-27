package api

import (
	"net/http"
	"time"
	"github.com/wyw14/cry-153/internal/model"
)

type RecoveryPayload struct { IncidentID string `json:"incident_id"`; Ready bool `json:"ready"`; CheckedAt time.Time `json:"checked_at"` }
func EncodeRecovery(w http.ResponseWriter, incident string, decision model.RecoveryDecision) { Encode(w,http.StatusOK,RecoveryPayload{IncidentID:incident,Ready:decision.Ready,CheckedAt:time.Now().UTC()}) }
func RecoveryStatus(decision model.RecoveryDecision) string { if decision.Ready{return "ready"};return "hold" }
func RecoveryHeaders(w http.ResponseWriter) { w.Header().Set("Cache-Control","no-store");w.Header().Set("X-RiverSentinel-Recovery","v1") }
func RecoveryAllowed(decision model.RecoveryDecision, now time.Time) bool { return decision.Ready&&!now.IsZero() }
