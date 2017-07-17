package storage

import (
	"strings"
	"testing"
)

func TestCheckQuotaNotifiesUser(t *testing.T) {
	var notifiedUser, notifiedMsg string
	notifyUser = func(user, msg string) {
		notifiedUser, notifiedMsg = user, msg
	}
	// имитация условий 980-ти Мбайтной занятости.
	const user = "matz@github.com"
	CheckQuota(user)
	if notifiedUser == "" && notifiedMsg == "" {
		t.Fatalf("notifyUser not called")
	}
	if notifiedUser != user {
		t.Errorf("Уведомлен (%s) вместо %s", notifiedUser, user)
	}
	const wantSubstring = "98% вашей квоты"
	if !strings.Contains(notifiedMsg, wantSubstring) {
		t.Errorf("неожиданное уведомление <<%s>>, "+"want substring %q",
			notifiedMsg, wantSubstring)
	}
}
