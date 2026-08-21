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

func TestRedactCopiesNestedLists(t *testing.T) {
	payload := map[string]any{"items": []any{"a", "b"}}
	out := Redact(payload, nil)
	out["items"].([]any)[0] = "changed"
	if payload["items"].([]any)[0] != "a" {
		t.Fatal("nested payload was shared")
	}
}
