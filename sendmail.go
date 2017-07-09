package main

import "net/smtp"

func main() {
	// authenticate, if needed
	auth := smtp.PlainAuth("", "AKI***********", "Ao*******************", "email-smtp.us-east-1.amazonaws.com")
	err := smtp.SendMail("email-smtp.us-east-1.amazonaws.com:587", auth, "cloudsec@dev.mozaws.net", []string{"cloudsec@mozilla.com"}, []byte(`From: cloudsec@dev.mozaws.net
To: cloudsec@mozilla.com
Subject: TLS Observatory runner results

TestCaribou
`))
	if err != nil {
		panic(err)
	}
	return
}
