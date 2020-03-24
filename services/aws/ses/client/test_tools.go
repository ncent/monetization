package client

type MockSESService struct {
	SendMailError error
}

func NewMockSESService() *MockSESService {
	return &MockSESService{}
}

func (ms MockSESService) SendEmail(er EmailRequest) error {
	return ms.SendMailError
}
