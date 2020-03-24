package client

type EmailRequest struct {
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
	Html      string `json:"html,omitempty"`
	Body      string `json:"body,omitempty"`
	Subject   string `json:"subject"`
}
