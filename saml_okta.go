package main

import (
	"fmt"

	"github.com/mattbaird/gosaml"
)

func main() {
	// Configure the app and account settings
	appSettings := saml.NewAppSettings("https://some.random.service.mozilla.com", "issuer")
	accountSettings := saml.NewAccountSettings("cert", "https://mozilla.okta.com/app/mozilla_some.random.service_1/k2urmue5SFSOISCLADEQ/sso/saml")

	// Construct an AuthnRequest
	authRequest := saml.NewAuthorizationRequest(*appSettings, *accountSettings)

	// Return a SAML AuthnRequest as a string
	saml, err := authRequest.GetRequest(false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(saml)

	samlUrl, err := authRequest.GetRequestUrl()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(samlUrl)
}
