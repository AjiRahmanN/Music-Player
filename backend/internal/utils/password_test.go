package utils

import "testing"

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name    string
		pw      string
		wantErr bool
	}{
		{"valid", "Abcdef12", false},
		{"too short", "Ab1", true},
		{"no upper", "abcdefg1", true},
		{"no lower", "ABCDEFG1", true},
		{"no digit", "Abcdefgh", true},
		{"empty", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.pw)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidatePassword(%q) err=%v, wantErr=%v", tc.pw, err, tc.wantErr)
			}
		})
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("Abcdef12")
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckPassword(hash, "Abcdef12"); err != nil {
		t.Fatalf("password benar seharusnya cocok: %v", err)
	}
	if err := CheckPassword(hash, "Salah123"); err == nil {
		t.Fatal("password salah seharusnya ditolak")
	}
}

func TestBurnPasswordCheckDoesNotPanic(t *testing.T) {
	BurnPasswordCheck("apa-saja")
	BurnPasswordCheck("")
}
