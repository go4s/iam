package cmd

import (
    "testing"
    
    "golang.org/x/crypto/bcrypt"
)

func TestGenAdminPasswd(t *testing.T) {
    pwd, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
    if err != nil {
        t.Fatalf("failed to generate password: %v", err)
    }
    t.Logf("admin password: %s", string(pwd))
    
    if ret := bcrypt.CompareHashAndPassword(pwd, []byte("admin123")); ret != nil {
        t.Fatalf("invalid password")
    }
}
