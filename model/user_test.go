package model

import "testing"

func TestValidateUser(t *testing.T) {
	user := User{
		Name:     "",
		Password: "password",
	}

	user.Name = ""
	if err := user.Validate(); err == nil {
		t.Fatal("Empty username not allowed")
	}

	user.Name = "miso"
	if err := user.Validate(); err != nil {
		t.Fatal("Regular characters allowed")
	}

	user.Name = "m15o"
	if err := user.Validate(); err != nil {
		t.Fatal("Digits allowed")
	}

	user.Name = "has space"
	if err := user.Validate(); err == nil {
		t.Fatal("Space is not allowed")
	}

	user.Name = "M15O"
	if err := user.Validate(); err == nil {
		t.Fatal("Capital letters aren't allowed")
	}

	characters := []string{"#", ":", "/", "@", "?"}
	for _, c := range characters {
		user.Name = c
		if err := user.Validate(); err == nil {
			t.Fatal("Special characters not allowed")
		}
	}
}
