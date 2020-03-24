package client

import (
	"encoding/json"
)

type MockAuth0Client struct {
	Token         string
	ErrorReturned error
}

func NewMockAuth0Client() *MockAuth0Client {
	return &MockAuth0Client{}
}

func (mc *MockAuth0Client) RequestToken() error {
	if mc.ErrorReturned != nil {
		return mc.ErrorReturned
	}

	log.Info("Preparing to request access token")

	var auth0Return struct {
		AccessToken string `json:"access_token"`
	}

	mockBody := `{
		"access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlFVWTBNRVF4UkVJNFJrTTJRMFEyTXpRd016WXpORE14UmtaRlFqYzNOa1U0T0RoQ1FUUTJOZyJ9.eyJpc3MiOiJodHRwczovL2Rldi1iNzRtbDIwZS5hdXRoMC5jb20vIiwic3ViIjoiM1V3T0V4cjlXN2JaWnd0U3I2RGFUZmR5YnFsTnJZSk1AY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vZGV2LWI3NG1sMjBlLmF1dGgwLmNvbS9hcGkvdjIvIiwiaWF0IjoxNTc5MjA1OTkzLCJleHAiOjE1NzkyOTIzOTMsImF6cCI6IjNVd09FeHI5VzdiWlp3dFNyNkRhVGZkeWJxbE5yWUpNIiwic2NvcGUiOiJyZWFkOmNsaWVudF9ncmFudHMgY3JlYXRlOmNsaWVudF9ncmFudHMgZGVsZXRlOmNsaWVudF9ncmFudHMgdXBkYXRlOmNsaWVudF9ncmFudHMgcmVhZDp1c2VycyB1cGRhdGU6dXNlcnMgZGVsZXRlOnVzZXJzIGNyZWF0ZTp1c2VycyByZWFkOnVzZXJzX2FwcF9tZXRhZGF0YSB1cGRhdGU6dXNlcnNfYXBwX21ldGFkYXRhIGRlbGV0ZTp1c2Vyc19hcHBfbWV0YWRhdGEgY3JlYXRlOnVzZXJzX2FwcF9tZXRhZGF0YSBjcmVhdGU6dXNlcl90aWNrZXRzIHJlYWQ6Y2xpZW50cyB1cGRhdGU6Y2xpZW50cyBkZWxldGU6Y2xpZW50cyBjcmVhdGU6Y2xpZW50cyByZWFkOmNsaWVudF9rZXlzIHVwZGF0ZTpjbGllbnRfa2V5cyBkZWxldGU6Y2xpZW50X2tleXMgY3JlYXRlOmNsaWVudF9rZXlzIHJlYWQ6Y29ubmVjdGlvbnMgdXBkYXRlOmNvbm5lY3Rpb25zIGRlbGV0ZTpjb25uZWN0aW9ucyBjcmVhdGU6Y29ubmVjdGlvbnMgcmVhZDpyZXNvdXJjZV9zZXJ2ZXJzIHVwZGF0ZTpyZXNvdXJjZV9zZXJ2ZXJzIGRlbGV0ZTpyZXNvdXJjZV9zZXJ2ZXJzIGNyZWF0ZTpyZXNvdXJjZV9zZXJ2ZXJzIHJlYWQ6ZGV2aWNlX2NyZWRlbnRpYWxzIHVwZGF0ZTpkZXZpY2VfY3JlZGVudGlhbHMgZGVsZXRlOmRldmljZV9jcmVkZW50aWFscyBjcmVhdGU6ZGV2aWNlX2NyZWRlbnRpYWxzIHJlYWQ6cnVsZXMgdXBkYXRlOnJ1bGVzIGRlbGV0ZTpydWxlcyBjcmVhdGU6cnVsZXMgcmVhZDpydWxlc19jb25maWdzIHVwZGF0ZTpydWxlc19jb25maWdzIGRlbGV0ZTpydWxlc19jb25maWdzIHJlYWQ6ZW1haWxfcHJvdmlkZXIgdXBkYXRlOmVtYWlsX3Byb3ZpZGVyIGRlbGV0ZTplbWFpbF9wcm92aWRlciBjcmVhdGU6ZW1haWxfcHJvdmlkZXIgYmxhY2tsaXN0OnRva2VucyByZWFkOnN0YXRzIHJlYWQ6dGVuYW50X3NldHRpbmdzIHVwZGF0ZTp0ZW5hbnRfc2V0dGluZ3MgcmVhZDpsb2dzIHJlYWQ6c2hpZWxkcyBjcmVhdGU6c2hpZWxkcyBkZWxldGU6c2hpZWxkcyByZWFkOmFub21hbHlfYmxvY2tzIGRlbGV0ZTphbm9tYWx5X2Jsb2NrcyB1cGRhdGU6dHJpZ2dlcnMgcmVhZDp0cmlnZ2VycyByZWFkOmdyYW50cyBkZWxldGU6Z3JhbnRzIHJlYWQ6Z3VhcmRpYW5fZmFjdG9ycyB1cGRhdGU6Z3VhcmRpYW5fZmFjdG9ycyByZWFkOmd1YXJkaWFuX2Vucm9sbG1lbnRzIGRlbGV0ZTpndWFyZGlhbl9lbnJvbGxtZW50cyBjcmVhdGU6Z3VhcmRpYW5fZW5yb2xsbWVudF90aWNrZXRzIHJlYWQ6dXNlcl9pZHBfdG9rZW5zIGNyZWF0ZTpwYXNzd29yZHNfY2hlY2tpbmdfam9iIGRlbGV0ZTpwYXNzd29yZHNfY2hlY2tpbmdfam9iIHJlYWQ6Y3VzdG9tX2RvbWFpbnMgZGVsZXRlOmN1c3RvbV9kb21haW5zIGNyZWF0ZTpjdXN0b21fZG9tYWlucyByZWFkOmVtYWlsX3RlbXBsYXRlcyBjcmVhdGU6ZW1haWxfdGVtcGxhdGVzIHVwZGF0ZTplbWFpbF90ZW1wbGF0ZXMgcmVhZDptZmFfcG9saWNpZXMgdXBkYXRlOm1mYV9wb2xpY2llcyByZWFkOnJvbGVzIGNyZWF0ZTpyb2xlcyBkZWxldGU6cm9sZXMgdXBkYXRlOnJvbGVzIHJlYWQ6cHJvbXB0cyB1cGRhdGU6cHJvbXB0cyByZWFkOmJyYW5kaW5nIHVwZGF0ZTpicmFuZGluZyIsImd0eSI6ImNsaWVudC1jcmVkZW50aWFscyJ9.kUdejpDINJy5jZYSEFX1V6R2L36PtCQhZ7FCmPFeYKTRsshc9rzfHLkUIRX7fBY7fuIAZogGXzEOh7qNQrEZTJ9wC8MnNtJ3qWDMV1hkO_qNa3s3cE22Paw2AoKnx1MZCeNjUvdOarpUQR9ebXS6xDlNixhf1yobgorXLhDCx6qnFpDHvswYrTp_V8SLY8OVUugGHeOOoQnarmsZg-MnRP1LG3-_IkjhjG1z6DjLk4i08YX1_BJe2P7gFP3sE7K58i7Ij8jegzeLI1E4W6Ik4HW5BBIRIRTNKqoY_gPmKxbRmwUZaUnuzjdghtY0v7J24-9Mc6AQ6QnpvmJA9k-zMQ",
		"scope": "read:client_grants create:client_grants delete:client_grants update:client_grants read:users update:users delete:users create:users read:users_app_metadata update:users_app_metadata delete:users_app_metadata create:users_app_metadata create:user_tickets read:clients update:clients delete:clients create:clients read:client_keys update:client_keys delete:client_keys create:client_keys read:connections update:connections delete:connections create:connections read:resource_servers update:resource_servers delete:resource_servers create:resource_servers read:device_credentials update:device_credentials delete:device_credentials create:device_credentials read:rules update:rules delete:rules create:rules read:rules_configs update:rules_configs delete:rules_configs read:email_provider update:email_provider delete:email_provider create:email_provider blacklist:tokens read:stats read:tenant_settings update:tenant_settings read:logs read:shields create:shields delete:shields read:anomaly_blocks delete:anomaly_blocks update:triggers read:triggers read:grants delete:grants read:guardian_factors update:guardian_factors read:guardian_enrollments delete:guardian_enrollments create:guardian_enrollment_tickets read:user_idp_tokens create:passwords_checking_job delete:passwords_checking_job read:custom_domains delete:custom_domains create:custom_domains read:email_templates create:email_templates update:email_templates read:mfa_policies update:mfa_policies read:roles create:roles delete:roles update:roles read:prompts update:prompts read:branding update:branding",
		"expires_in": 86400,
		"token_type": "Bearer"
	}`

	err := json.Unmarshal([]byte(mockBody), &auth0Return)

	if err != nil {
		log.Errorf("Request token failed %v", err)
		return err
	}

	mc.Token = auth0Return.AccessToken

	log.Infoln("Access token requested with success", mc.Token)

	return nil
}

func (mc MockAuth0Client) formatCreateUserPayload(userDetails *UserDetails) (string, error) {
	jsonResult, err := json.Marshal(userDetails)

	if err != nil {
		log.Errorf("Failed to marshal user details %v", err)
		return "", err
	}

	return string(jsonResult), nil
}

func (mc MockAuth0Client) CreateNewUser(userDetails *UserDetails) error {
	log.Info("User created with success on Auth0")

	pk, pa, err := GeneratePrivateAndPubKeyHex()

	if err != nil {
		log.Errorf("Failed to create private and pub key %v", err)
		return err
	}

	wallet := Wallet{
		PrivateKey: pk,
		PublicKey:  pa,
	}

	metadata := Metadata{
		Wallets{
			wallet,
		},
	}

	userDetails.Metadata = metadata

	payload, _ := mc.formatCreateUserPayload(userDetails)

	log.Infof("User details %s", payload)

	return mc.ErrorReturned
}
