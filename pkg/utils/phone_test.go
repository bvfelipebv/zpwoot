package utils

import "testing"

func TestValidatePhone(t *testing.T) {
    cases := []struct{
        in string
        ok bool
    }{
        {"+5511999999999", true},
        {"5511999999999", true},
        {"12345", false},
        {"+abcd", false},
        {"", false},
    }

    for _, c := range cases {
        got := ValidatePhone(c.in)
        if got != c.ok {
            t.Fatalf("ValidatePhone(%q) = %v, want %v", c.in, got, c.ok)
        }
    }
}
