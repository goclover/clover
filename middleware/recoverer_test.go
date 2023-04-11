package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goclover/clover"
)

func panicingHandler(http.ResponseWriter, *http.Request) { panic("foo") }

func TestRecoverer(t *testing.T) {
	r := clover.New()

	oldRecovererErrorWriter := RecovererErrorWriter
	defer func() { RecovererErrorWriter = oldRecovererErrorWriter }()
	buf := &bytes.Buffer{}
	RecovererErrorWriter = buf

	r.Use(Recoverer)
	r.MethodFunc("GET", "/", panicingHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, "GET", "/", nil)
	assertEqual(t, res.StatusCode, http.StatusInternalServerError)

	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "->") {
			if !strings.Contains(line, "panicingHandler") {
				t.Fatalf("First func call line should refer to panicingHandler, but actual line:\n%v\n", line)
			}
			return
		}
	}
	t.Fatal("First func call line should start with ->.")
}

func TestRecovererAbortHandler(t *testing.T) {
	defer func() {
		rcv := recover()
		if rcv != http.ErrAbortHandler {
			t.Fatalf("http.ErrAbortHandler should not be recovered")
		}
	}()

	w := httptest.NewRecorder()

	r := clover.New()
	r.Use(Recoverer)

	r.MethodFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		panic(http.ErrAbortHandler)
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(w, req)
}
