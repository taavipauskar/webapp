package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestForm_Has(t *testing.T) {
	form := NewForm(nil)
	has := form.Has("whatever")
	if has {
		t.Error("form should not have any values")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = NewForm(postedData)

	has = form.Has("a")
	if !has {
		t.Error("form should have a value")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/path", nil)
	form := NewForm(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form should not be valid if fields are missing values")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r = httptest.NewRequest(http.MethodPost, "/path", strings.NewReader(postedData.Encode()))
	r.PostForm = postedData
	form = NewForm(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form should be valid if fields are present")
	}
}

func TestForm_Check(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password is required")
	if form.Valid() {
		t.Error("Valid should not be valid if fields are not present")
	}
}

func TestForm_ErrorGet(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password is required")
	s := form.Errors.Get("password")
	if len(s) == 0 {
		t.Error("form should have an error returned from get")
	}

	s = form.Errors.Get("whatever")
	if len(s) != 0 {
		t.Error("form should not an error for nonexisting field")
	}
}
