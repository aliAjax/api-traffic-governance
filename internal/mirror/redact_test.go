package mirror

import "testing"

func TestRedactDoesNotMutateCapturedPayload(t *testing.T) {
	payload := map[string]any{"token": "secret", "path": "/v1/items"}
	out := Redact(payload, []string{"token"})
	if out["token"] != "[REDACTED]" {
		t.Fatalf("redacted value = %#v", out["token"])
	}
	if payload["token"] != "secret" {
		t.Fatalf("original payload was changed: %#v", payload)
	}
}
