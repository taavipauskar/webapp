package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/asd", http.StatusNotFound},
	}

	routes := app.routes()

	// create a test server

	ts := httptest.NewTLSServer(routes)
	// close server when done
	defer ts.Close()

	for _, test := range theTests {
		resp, err := ts.Client().Get(ts.URL + test.url)
		if err != nil {
			t.Log(err)
			t.Fatalf("GET %s: %s", test.url, err)
		}
		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("for %s, GET %s: expected %d, got %d", test.name, test.url, test.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestAppHome(t *testing.T) {
	var tests = []struct {
		name         string
		putInSession string
		expectedHtml string
	}{
		{"First visit", "", "From session"},
		{"Second visit", "hello world!", "From session hello world!"},
	}

	for _, e := range tests {
		req, _ := http.NewRequest("GET", "/", nil)

		req = addContextAndSessionToRequest(req, app)
		_ = app.Session.Destroy(req.Context())

		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Home)

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
		}

		body, _ := io.ReadAll(rr.Body)
		if !strings.Contains(string(body), e.expectedHtml) {
			t.Errorf("%s: handler returned unexpected body: got %v want %v", e.name, string(body), e.expectedHtml)
		}
	}
}

func TestAppp_renderWithBadTemplate(t *testing.T) {
	pathToTemplates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()
	err := app.render(rr, req, "bad.page.gohtml", &TemplateData{})
	if err == nil {
		t.Error("Expected an error")
	}
}

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))

	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)
}
