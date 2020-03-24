package challenge

type MockValidateChallengeService struct {
	ErrResult error
}

func NewMockValidateChallengeService() *MockValidateChallengeService {
	return &MockValidateChallengeService{}
}

func (m *MockValidateChallengeService) Validate(userUUID, action string) error {
	return m.ErrResult
}
