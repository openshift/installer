package extender

import "testing"

func TestExtensionsString(t *testing.T) {
	e := Extensions{}

	s := e.String()
	if s != "" {
		t.Errorf("Expected %s and got %s", "", s)
	}

	e["foo"] = "bar"
	e["baz"] = "qux"

	s = e.String()
	if s != "baz:qux, foo:bar" {
		t.Errorf("Expected %s and got %s", "baz:qux, foo:bar", s)
	}

}

func TestExtensionsSet(t *testing.T) {
	e := Extensions{}

	_, ok := e["foo"]
	if ok {
		t.Errorf("Expected %t and got %t", false, ok)
	}

	err := e.Set("foo:bar")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)

	}
	v, ok := e["foo"]
	if !ok {
		t.Errorf("Expected %t and got %t", true, ok)
	}
	if v != "bar" {
		t.Errorf("Expected %s and got %s", "bar", v)

	}

	err = e.Set("foo:bar:baz")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)

	}
	v, ok = e["foo"]
	if !ok {
		t.Errorf("Expected %t and got %t", true, ok)
	}
	if v != "bar:baz" {
		t.Errorf("Expected %s and got %s", "bar:baz", v)

	}

	err = e.Set("foo")
	if err == nil {
		t.Errorf("Expected error with extensions")

	}
}
