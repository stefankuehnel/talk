package message

import "testing"

func TestRender(t *testing.T) {
	t.Parallel()

	t.Run("renders template with data", func(t *testing.T) {
		t.Parallel()

		rendered, err := Render("Hello {{.Name}}", map[string]string{"Name": "Stefan"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rendered != "Hello Stefan" {
			t.Fatalf("expected rendered template %q, got %q", "Hello Stefan", rendered)
		}
	})

	t.Run("returns an error for malformed template syntax", func(t *testing.T) {
		t.Parallel()

		_, err := Render("Hello {{.Name", map[string]string{"Name": "Stefan"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns an error for missing template key", func(t *testing.T) {
		t.Parallel()

		_, err := Render("Hello {{.Name}}", map[string]string{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
