package http_method

import (
	"net/http"
	"strings"
	"testing"
)

func TestUnknownMethod(t *testing.T) {
	converted := StringToHttpMethod("unknown")
	if converted != "UNKNOWN" {
		t.Fatalf("expected 'UNKNOWN', got '%s'", converted)
	}
}

func TestConvertGetMethod(t *testing.T) {
	method := "Get"
	converted := StringToHttpMethod(method)
	if converted != http.MethodGet {
		t.Fatalf("expected '%s', got '%s'", http.MethodGet, converted)
	}
	converted = StringToHttpMethod(strings.ToLower(method))
	if converted != http.MethodGet {
		t.Fatalf("expected '%s', got '%s'", http.MethodGet, converted)
	}
	converted = StringToHttpMethod(strings.ToUpper(method))
	if converted != http.MethodGet {
		t.Fatalf("expected '%s', got '%s'", http.MethodGet, converted)
	}
}

func TestConvertPostMethod(t *testing.T) {
	method := "Post"
	converted := StringToHttpMethod(method)
	if converted != http.MethodPost {
		t.Fatalf("expected '%s', got '%s'", http.MethodPost, converted)
	}
	converted = StringToHttpMethod(strings.ToLower(method))
	if converted != http.MethodPost {
		t.Fatalf("expected '%s', got '%s'", http.MethodPost, converted)
	}
	converted = StringToHttpMethod(strings.ToUpper(method))
	if converted != http.MethodPost {
		t.Fatalf("expected '%s', got '%s'", http.MethodPost, converted)
	}
}

func TestConvertPutMethod(t *testing.T) {
	method := "Put"
	converted := StringToHttpMethod(method)
	if converted != http.MethodPut {
		t.Fatalf("expected '%s', got '%s'", http.MethodPut, converted)
	}
	converted = StringToHttpMethod(strings.ToLower(method))
	if converted != http.MethodPut {
		t.Fatalf("expected '%s', got '%s'", http.MethodPut, converted)
	}
	converted = StringToHttpMethod(strings.ToUpper(method))
	if converted != http.MethodPut {
		t.Fatalf("expected '%s', got '%s'", http.MethodPut, converted)
	}
}

func TestConvertPatchMethod(t *testing.T) {
	method := "Patch"
	converted := StringToHttpMethod(method)
	if converted != http.MethodPatch {
		t.Fatalf("expected '%s', got '%s'", http.MethodPatch, converted)
	}
	converted = StringToHttpMethod(strings.ToLower(method))
	if converted != http.MethodPatch {
		t.Fatalf("expected '%s', got '%s'", http.MethodPatch, converted)
	}
	converted = StringToHttpMethod(strings.ToUpper(method))
	if converted != http.MethodPatch {
		t.Fatalf("expected '%s', got '%s'", http.MethodPatch, converted)
	}
}

func TestConvertDeleteMethod(t *testing.T) {
	method := "Delete"
	converted := StringToHttpMethod(method)
	if converted != http.MethodDelete {
		t.Fatalf("expected '%s', got '%s'", http.MethodDelete, converted)
	}
	converted = StringToHttpMethod(strings.ToLower(method))
	if converted != http.MethodDelete {
		t.Fatalf("expected '%s', got '%s'", http.MethodDelete, converted)
	}
	converted = StringToHttpMethod(strings.ToUpper(method))
	if converted != http.MethodDelete {
		t.Fatalf("expected '%s', got '%s'", http.MethodDelete, converted)
	}
}
