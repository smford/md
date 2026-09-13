package theme

import "testing"

func TestGetTheme(t *testing.T) {
	for _, name := range AvailableThemes() {
		th := GetTheme(name, false)
		if th == nil {
			t.Errorf("expected theme for %s, got nil", name)
		}
		if th.Name != name {
			t.Errorf("expected theme name %s, got %s", name, th.Name)
		}
	}

	// Plain flag override
	plain := GetTheme("dark", true)
	if plain.Name != "plain" {
		t.Errorf("expected plain theme when plain=true, got %s", plain.Name)
	}
}
