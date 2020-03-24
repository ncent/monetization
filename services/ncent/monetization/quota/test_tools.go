package quota

type MockQuotaUsageRepository struct {
	QuotaUsageReturned int
	ErrorReturned      error
}

func NewMockQuotaUsageRepository() *MockQuotaUsageRepository {
	return &MockQuotaUsageRepository{}
}

func (m *MockQuotaUsageRepository) GetQuotaUsage(uuid string, action string, amount int, timePeriod string) (int, error) {
	return m.QuotaUsageReturned, m.ErrorReturned
}
