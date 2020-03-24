package plan

var PlanSample = &Plan{
	UUID: "4ec48742-0195-11ea-afab-00155d43e1f0",
	Name: "Plan 1",
	Criteria: Criteria{
		{
			UUID:       "a18c8ba0-0195-11ea-938d-00155d43e1f0",
			Name:       "Criteria 1",
			MaxCount:   10,
			Model:      FREEMIUM.String(),
			Action:     CREATE_CHALLENGE.String(),
			TimePeriod: WEEK.String(),
		},
		{
			UUID:       "a18c8ba0-0195-11ea-938d-00155d13e1f0",
			Name:       "Criteria 1",
			MaxCount:   10,
			Model:      FREEMIUM.String(),
			Action:     SHARE_CHALLENGE.String(),
			TimePeriod: WEEK.String(),
		},
	},
}

type MockPlanRepository struct {
	PlanReturned  *Plan
	ErrorReturned error
}

func NewMockPlanRepository() *MockPlanRepository {
	return &MockPlanRepository{}
}

func (m *MockPlanRepository) GetByUUID(uuid string) (*Plan, error) {
	return m.PlanReturned, m.ErrorReturned
}

func (m MockPlanRepository) Update(uuid, updateExpression string, plan *Plan) error {
	return m.ErrorReturned
}

func (m *MockPlanRepository) Store(plan *Plan) error {
	return m.ErrorReturned
}
