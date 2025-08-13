package cipher

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword() error = %v", err)
		return
	}

	if hashedPassword == "" {
		t.Error("HashPassword() returned empty string")
	}

	if hashedPassword == password {
		t.Error("HashPassword() returned the same string as input")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "mysecretpassword"

	// Hash the password first
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	// Test correct password
	err = CheckPassword(password, hashedPassword)
	if err != nil {
		t.Errorf("CheckPassword() with correct password returned error = %v", err)
	}

	// Test incorrect password
	wrongPassword := "wrongpassword"
	err = CheckPassword(wrongPassword, hashedPassword)
	if err == nil {
		t.Error("CheckPassword() with wrong password should return error")
	}
}

func TestCheckPasswordWithEmptyInputs(t *testing.T) {
	// Test empty password
	err := CheckPassword("", "somehash")
	if err == nil {
		t.Error("CheckPassword() with empty password should return error")
	}

	// Test empty hash
	err = CheckPassword("somepassword", "")
	if err == nil {
		t.Error("CheckPassword() with empty hash should return error")
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "mysecretpassword"

	// Hash the same password multiple times
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("First HashPassword() error = %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Second HashPassword() error = %v", err)
	}

	// Each hash should be different (due to salt)
	if hash1 == hash2 {
		t.Error("HashPassword() should generate different hashes for the same password")
	}

	// Both hashes should work for the same password
	err = CheckPassword(password, hash1)
	if err != nil {
		t.Errorf("First hash validation failed: %v", err)
	}

	err = CheckPassword(password, hash2)
	if err != nil {
		t.Errorf("Second hash validation failed: %v", err)
	}
}
