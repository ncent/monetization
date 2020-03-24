package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const apiPath = "dev-b74ml20e.auth0.com"

var log *logrus.Entry

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"service": "service.auth0.client",
	})
}

type (
	IAuth0Client interface {
		RequestToken() error
		CreateNewUser(userDetails *UserDetails) error
	}

	Auth0Client struct {
		ClientID     string
		ClientSecret string
		Token        string
	}
)

func NewAuth0Client() *Auth0Client {
	auth0ClientID, ok := os.LookupEnv("AUTH0_CLIENT_ID")

	if !ok {
		log.Fatal("Auth0 ClientID not found")
	}

	auth0ClientSecret, ok := os.LookupEnv("AUTH0_SECRET")

	if !ok {
		log.Fatal("Auth0 Secret not found")
	}

	return &Auth0Client{
		ClientID:     auth0ClientID,
		ClientSecret: auth0ClientSecret,
	}
}

func (ac *Auth0Client) formatRequestTokenPayload() *strings.Reader {
	return strings.NewReader(fmt.Sprintf(`{
		"client_id":"%s",
		"client_secret":"%s",
		"audience":"https://%s/api/v2/",
		"grant_type":"client_credentials"
	}`, ac.ClientID, ac.ClientSecret, apiPath))
}

func (ac *Auth0Client) RequestToken() error {
	log.Info("Preparing to request access token")

	url := fmt.Sprintf("https://%s/oauth/token", apiPath)

	req, err := http.NewRequest("POST", url, ac.formatRequestTokenPayload())

	if err != nil {
		log.Errorf("Request token failed %v", err)
		return err
	}

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Errorf("Request token failed %v", err)
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Errorf("Request token failed %v", err)
		return err
	}

	var auth0Return struct {
		AccessToken string `json:"access_token"`
	}

	err = json.Unmarshal(body, &auth0Return)

	if err != nil {
		log.Errorf("Request token failed %v", err)
		return err
	}

	ac.Token = auth0Return.AccessToken

	log.Info("Access token requested with success")

	return nil
}

func (ac *Auth0Client) formatCreateUserPayload(userDetails *UserDetails) (*strings.Reader, error) {
	jsonResult, err := json.Marshal(userDetails)

	if err != nil {
		log.Errorf("Failed to marshal user details %v", err)
		return nil, err
	}

	return strings.NewReader(string(jsonResult)), nil
}

func (ac *Auth0Client) CreateNewUser(userDetails *UserDetails) error {
	if ac.Token == "" {
		return fmt.Errorf("API Token not found, first request the API Token")
	}

	url := fmt.Sprintf("https://%s/api/v2/users", apiPath)

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
	userDetails.Password = uuid.NewV4().String()
	userDetails.Connection = "Username-Password-Authentication"

	payload, err := ac.formatCreateUserPayload(userDetails)

	if err != nil {
		log.Errorf("Failed to create new user %v", err)
		return err
	}

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		log.Errorf("Failed to create new user %v", err)
		return err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", ac.Token))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Errorf("Failed to create new user %v", err)
		return err
	}

	if res.StatusCode != 201 {
		return fmt.Errorf("Failed to created user on Auth0 Status Code %d (StatusCode != 201)", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	log.Infoln("Body request", string(body))

	log.Info("User created with success on Auth0")

	return nil
}
