package new_subscription

import (
	"bytes"
	"text/template"
)

type NewSubscriptionEmailBodyVars struct {
	Email string
}

func GenerateNewSubscriptionEmailBody(vars NewSubscriptionEmailBodyVars) (string, error) {
	t, err := template.New("new_subscription_email").Parse(`
		<html>
			<head>
				<meta http-equiv="refresh" content="0; URL='%s'" />
				<style>
					body {
						margin: 0;
					}
					.background {
						width: 100%%;
						height: 100vh;
						display: flex;
						background-size: contain;
						background-color: #18191B;
						background-repeat: no-repeat;
						background-position: center;
						flex-direction: column;
						justify-content: center;
						align-items: center;
					}
					.forwardImage {
						fill: #FFFFFF;
					}
					.forwardButton svg {
						width: 100px;
						height: 100px;
						margin: 0px;
					}
					.forwardButton {
						margin: auto;
						width: 200px;
						height: 200px;
						background-color: #b71c1b;
						background-repeat:no-repeat;
						cursor:pointer;
						overflow: hidden;
						outline:none;
						padding: 0px;
						border: none;
					}
					.forwardButton:hover {
						background-color: #9a1312;   
					}
					.text {
						font-size: 50px;
						color: #FFFFFF;
						margin-bottom: 100px;
					}
				</style>
			</head>
			<body >
				<div class="background">
					<div>
						<h1 class="text"New Subscription Created For {{.Email}}</h1>
					</div>
				</div>
			</body>
		</html>
	`)

	var tpl bytes.Buffer
	err = t.Execute(&tpl, vars)

	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}
