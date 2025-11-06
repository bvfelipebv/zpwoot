package constants

import (
	"testing"
)

func TestIsValidEventType(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		want      bool
	}{
		{"Valid message event", "message", true},
		{"Valid connection event", "connected", true},
		{"Valid special event", "all", true},
		{"Invalid event", "invalid_event", false},
		{"Empty string", "", false},
		{"Random string", "random", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEventType(tt.eventType); got != tt.want {
				t.Errorf("IsValidEventType(%q) = %v, want %v", tt.eventType, got, tt.want)
			}
		})
	}
}

func TestIsCriticalEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		want      bool
	}{
		{"Connected is critical", "connected", true},
		{"Disconnected is critical", "disconnected", true},
		{"LoggedOut is critical", "logged_out", true},
		{"Message is not critical", "message", false},
		{"QR is not critical", "qr", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCriticalEvent(tt.eventType); got != tt.want {
				t.Errorf("IsCriticalEvent(%q) = %v, want %v", tt.eventType, got, tt.want)
			}
		})
	}
}

func TestIsMessageEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		want      bool
	}{
		{"Message is message event", "message", true},
		{"Receipt is message event", "receipt", true},
		{"Connected is not message event", "connected", false},
		{"QR is not message event", "qr", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMessageEvent(tt.eventType); got != tt.want {
				t.Errorf("IsMessageEvent(%q) = %v, want %v", tt.eventType, got, tt.want)
			}
		})
	}
}

func TestIsConnectionEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		want      bool
	}{
		{"Connected is connection event", "connected", true},
		{"Disconnected is connection event", "disconnected", true},
		{"QR is connection event", "qr", true},
		{"Message is not connection event", "message", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsConnectionEvent(tt.eventType); got != tt.want {
				t.Errorf("IsConnectionEvent(%q) = %v, want %v", tt.eventType, got, tt.want)
			}
		})
	}
}

func TestGetEventsByCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		wantLen  int
	}{
		{"Messages category", "messages", 5},
		{"Connection category", "connection", 15},
		{"Invalid category", "invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetEventsByCategory(tt.category)
			if len(got) != tt.wantLen {
				t.Errorf("GetEventsByCategory(%q) returned %d events, want %d", tt.category, len(got), tt.wantLen)
			}
		})
	}
}

func TestGetAllCategories(t *testing.T) {
	categories := GetAllCategories()

	if len(categories) == 0 {
		t.Error("GetAllCategories() returned empty slice")
	}

	// Verificar se categorias esperadas existem
	expectedCategories := []string{"messages", "connection", "calls", "presence"}
	for _, expected := range expectedCategories {
		found := false
		for _, cat := range categories {
			if cat == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected category %q not found in categories", expected)
		}
	}
}

func TestGetEventDescription(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		wantEmpty bool
	}{
		{"Message has description", "message", false},
		{"Connected has description", "connected", false},
		{"Invalid event has default description", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetEventDescription(tt.eventType)
			if (got == "") == !tt.wantEmpty {
				t.Errorf("GetEventDescription(%q) = %q, wantEmpty = %v", tt.eventType, got, tt.wantEmpty)
			}
		})
	}
}

func TestValidateEventList(t *testing.T) {
	tests := []struct {
		name        string
		events      []string
		wantValid   int
		wantInvalid int
	}{
		{
			name:        "All valid events",
			events:      []string{"message", "connected", "qr"},
			wantValid:   3,
			wantInvalid: 0,
		},
		{
			name:        "All invalid events",
			events:      []string{"invalid1", "invalid2"},
			wantValid:   0,
			wantInvalid: 2,
		},
		{
			name:        "Mixed valid and invalid",
			events:      []string{"message", "invalid", "connected"},
			wantValid:   2,
			wantInvalid: 1,
		},
		{
			name:        "Empty list",
			events:      []string{},
			wantValid:   0,
			wantInvalid: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, invalid := ValidateEventList(tt.events)
			if len(valid) != tt.wantValid {
				t.Errorf("ValidateEventList() valid = %d, want %d", len(valid), tt.wantValid)
			}
			if len(invalid) != tt.wantInvalid {
				t.Errorf("ValidateEventList() invalid = %d, want %d", len(invalid), tt.wantInvalid)
			}
		})
	}
}

func TestGetEventCategory(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		want      string
	}{
		{"Message is in messages category", "message", "messages"},
		{"Connected is in connection category", "connected", "connection"},
		{"CallOffer is in calls category", "call_offer", "calls"},
		{"Invalid event returns unknown", "invalid", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEventCategory(tt.eventType); got != tt.want {
				t.Errorf("GetEventCategory(%q) = %q, want %q", tt.eventType, got, tt.want)
			}
		})
	}
}

func TestDefaultWebhookEvents(t *testing.T) {
	if len(DefaultWebhookEvents) == 0 {
		t.Error("DefaultWebhookEvents is empty")
	}

	// Verificar se todos os eventos padrão são válidos
	for _, event := range DefaultWebhookEvents {
		if !IsValidEventType(event) {
			t.Errorf("Default event %q is not valid", event)
		}
	}
}

func TestCriticalEvents(t *testing.T) {
	if len(CriticalEvents) == 0 {
		t.Error("CriticalEvents is empty")
	}

	// Verificar se todos os eventos críticos são válidos
	for _, event := range CriticalEvents {
		if !IsValidEventType(event) {
			t.Errorf("Critical event %q is not valid", event)
		}
	}
}

func TestSupportedEventTypes(t *testing.T) {
	if len(SupportedEventTypes) == 0 {
		t.Error("SupportedEventTypes is empty")
	}

	// Verificar se não há duplicatas
	seen := make(map[string]bool)
	for _, event := range SupportedEventTypes {
		if seen[event] {
			t.Errorf("Duplicate event found: %q", event)
		}
		seen[event] = true
	}
}
