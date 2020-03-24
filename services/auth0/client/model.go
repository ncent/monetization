package client

type Wallet struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type Wallets []Wallet

type Metadata struct {
	Wallets `json:"wallets"`
}

type UserDetails struct {
	Email      string `json:"email"`
	Connection string `json:"connection"`
	Password   string `json:"password"`
	Metadata   `json:"user_metadata"`
}
