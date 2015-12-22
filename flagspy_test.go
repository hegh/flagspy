package flagspy

import (
	"os"
	"sync"
	"testing"
)

var argc = len(os.Args)

func reset() {
	os.Args = os.Args[:argc] // In case we added any.
	Free()
	once = sync.Once{} // So we re-initialize.
}

func TestNoValue(t *testing.T) {
	reset()
	os.Args = append(os.Args, "-one", "-three")
	if v, ok := Get("one"); v != "" || !ok {
		t.Errorf("Got %q %v, want \"\" true", v, ok)
	}
	if v, ok := Get("three"); v != "" || !ok {
		t.Errorf("Got %q %v, want \"\" true", v, ok)
	}

	reset()
	os.Args = append(os.Args, "--two", "--three")
	if v, ok := Get("two"); v != "" || !ok {
		t.Errorf("Got %q %v, want \"\" true", v, ok)
	}
	if v, ok := Get("three"); v != "" || !ok {
		t.Errorf("Got %q %v, want \"\" true", v, ok)
	}
}

func TestEqualsValue(t *testing.T) {
	reset()
	os.Args = append(os.Args, "-one=val1")
	if v, ok := Get("one"); v != "val1" || !ok {
		t.Errorf("Got %q %v, want \"val1\" true", v, ok)
	}

	reset()
	os.Args = append(os.Args, "--two=val2")
	if v, ok := Get("two"); v != "val2" || !ok {
		t.Errorf("Got %q %v, want \"val2\" true", v, ok)
	}
}

func TestEqualsQuotedValue(t *testing.T) {
	reset()
	os.Args = append(os.Args, "-one='val1'")
	if v, ok := Get("one"); v != "val1" || !ok {
		t.Errorf("Got %q %v, want \"val1\" true", v, ok)
	}

	reset()
	os.Args = append(os.Args, "--two=\"val2\"")
	if v, ok := Get("two"); v != "val2" || !ok {
		t.Errorf("Got %q %v, want \"val2\" true", v, ok)
	}
}

func TestSpaceValue(t *testing.T) {
	reset()
	os.Args = append(os.Args, "-one", "val1")
	if v, ok := Get("one"); v != "val1" || !ok {
		t.Errorf("Got %q %v, want \"val1\" true", v, ok)
	}

	reset()
	os.Args = append(os.Args, "--two", "val2")
	if v, ok := Get("two"); v != "val2" || !ok {
		t.Errorf("Got %q %v, want \"val2\" true", v, ok)
	}
}

func TestSpaceQuotedValue(t *testing.T) {
	reset()
	os.Args = append(os.Args, "-one", "'val1'")
	if v, ok := Get("one"); v != "val1" || !ok {
		t.Errorf("Got %q %v, want \"val1\" true", v, ok)
	}

	reset()
	os.Args = append(os.Args, "--two", "\"val2\"")
	if v, ok := Get("two"); v != "val2" || !ok {
		t.Errorf("Got %q %v, want \"val2\" true", v, ok)
	}
}

func TestBreak(t *testing.T) {
	reset()
	os.Args = append(os.Args, "-one", "--", "val1", "-two")
	if v, ok := Get("one"); v != "" || !ok {
		t.Errorf("Got %q %v, want \"\" true", v, ok)
	}
	if v, ok := Get("two"); v != "" || ok {
		t.Errorf("Got %q %v, want \"\" false", v, ok)
	}

	reset()
	os.Args = append(os.Args, "--one", "--", "val1", "--two=val2")
	if v, ok := Get("one"); v != "" || !ok {
		t.Errorf("Got %q %v, want \"\" true", v, ok)
	}
	if v, ok := Get("two"); v != "" || ok {
		t.Errorf("Got %q %v, want \"\" false", v, ok)
	}
}
