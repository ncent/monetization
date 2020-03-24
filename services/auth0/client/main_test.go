package client

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create Auth0 Suite", func() {
	var (
		auth0Client IAuth0Client
	)

	JustBeforeEach(func() {
		auth0Client = NewMockAuth0Client()
		_ = auth0Client
	})

	Context("Given auth0 client", func() {
		It("It request token properly", func() {
			err := auth0Client.RequestToken()
			Expect(err).To(BeNil())
		})

		It("It creates user properly", func() {
			userDetails := &UserDetails{
				Email: "foobar@example.org",
			}

			err := auth0Client.CreateNewUser(userDetails)
			Expect(err).To(BeNil())
		})

		It("Generate wallet keys", func() {
			GeneratePrivateAndPubKeyHex()
		})
	})
})

func TestAuth0Client(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Auth0 Suite")
}
