package auth

import "testing"

func TestGenerateURL(t *testing.T) {
	actual := GenerateURL("github")
	expected := "expected"
	if actual != expected {
		t.Errorf("Generate github oauth2 url %s, expect: %s", actual, expected)
	}

	actual = GenerateURL("gitee")
	if actual != expected {
		t.Errorf("Generate github oauth2 url %s, expect: %s", actual, expected)
	}
}
