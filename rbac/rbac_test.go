package rbac

import "testing"

func TestExecuteRule(t *testing.T) {
	if !executeRule("123", checkingParams{}, "123") {
		t.Error("Ошибкааааа")
	}
}
