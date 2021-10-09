package database_Connect

import "testing"

func TestConnectDB(t *testing.T) {
	var user, post = ConnectDB()
	if user == nil {
		t.Errorf("User data not available")
	}
	if post == nil {
		t.Errorf("Post data not available")
	}
}
