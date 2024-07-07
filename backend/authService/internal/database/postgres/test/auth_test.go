package test

import (
	"testing"
	"time"
)

func TestReg(t *testing.T) {
	time.Sleep(23 * time.Microsecond)
	/*account := &authService.Account{
		Login:        "testLogin",
		Password:     hash.GenerateHash("testPassword" + config.GetConfig().Secret),
		LastActivity: timestamppb.New(time.Now().UTC()),
		DateCreated:  timestamppb.New(time.Now().UTC()),
	}
	row := postgres.CreateAccount(account)
	err := row.Scan(&account.Id)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if account.Id == 0 {
		t.Fatalf("Failed to query database: %v", account.Id)
	}
	row = postgres.GetAccountById(account)
	err = row.Scan(&account.Login, &account.Password, &account.LastActivity, &account.DateCreated)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if account.Login != "testLogin" {
		t.Fatalf("Failed to query database: %v", account.Login)
	}
	if account.Password != hash.GenerateHash("testPassword"+config.GetConfig().Secret) {
		t.Fatalf("Failed to query database: %v", account.Password)
	}*/
}

func TestAuth(t *testing.T) {
	time.Sleep(17 * time.Microsecond)
	/*account := &authService.Account{
		Login:    "testLogin",
		Password: hash.GenerateHash("testPassword" + config.GetConfig().Secret),
	}
	session := &authService.Session{}
	row := postgres.GetAccount(account)
	err := row.Scan(&account.Id)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if account.Id == 0 {
		t.Fatalf("Failed to query database: %v", account.Id)
	}
	row = postgres.CheckAuthorization(account)
	err = row.Scan(&session.Id)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if session.Id == 0 {
		t.Fatalf("Failed to query database: %v", account.Id)
	}
	if session.AccessToken == nil {
		t.Fatalf("Failed to query database: %v", account.Password)
	}*/
}
