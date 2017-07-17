/*
Storage проверяет квоты в веб-службе, которые обеспечивают сетевые
хранилища пользовательских файлов.
Когда пользователи превышают 90% квоты, система отправляет им предупреждение
по электронной почте
*/

package storage

import (
	"fmt"
	"log"
	"net/smtp"
)

var usage = make(map[string]int64)

func bytesInUse(username string) int64 {
	return usage[username]
}

// Настройка отправителя e-mail.
// Примечание: никогда не помещайте пароль в исходники программы!
const sender = "notifications@onflow.studio"
const passwd = "qwerty123"
const hostname = "smtp.onflow.studio"
const template = `Внимание, вы использовали %d байт хранилища, %d%% вашей квоты.`

var notifyUser = func(username, msg string) {
	auth := smtp.PlainAuth("", sender, passwd, hostname)
	err := smtp.SendMail(hostname+":587", auth, sender, []string{username}, []byte(msg))
	if err != nil {
		log.Printf("Сбой smtp.SendMail(%s): %s", username, err)
	}
}

// CheckQuota проверяет процент.
func CheckQuota(username string) {
	used := bytesInUse(username)
	const quota = 1000000000 // 1GB
	percent := 100 * used / quota
	if percent < 90 {
		return // OK
	}

	msg := fmt.Sprintf(template, used, percent)
	notifyUser(username, msg)
}
