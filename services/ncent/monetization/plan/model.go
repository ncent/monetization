package plan

type Model int

const (
	FREEMIUM Model = iota
	SUBSCRIPTION
)

func (m Model) String() string {
	return [...]string{"FREEMIUM", "SUBSCRIPTION"}[m]
}

type Action int

const (
	CREATE_CHALLENGE Action = iota
	SHARE_CHALLENGE
)

func (a Action) String() string {
	return [...]string{"CREATE_CHALLENGE", "SHARE_CHALLENGE"}[a]
}

type TimePeriod int

const (
	SECOND TimePeriod = iota
	MINUTE
	HOUR
	DAY
	WEEK
	MONTH
	YEAR
)

func (t TimePeriod) String() string {
	return [...]string{"SECOND", "MINUTE", "HOUR", "DAY", "WEEK", "MONTH", "YEAR"}[t]
}

type Criterion struct {
	UUID       string `json:"uuid,omitempty"`
	Name       string `json:"name"`
	MaxCount   int    `json:"max_count"`
	Model      string `json:"model"`
	Action     string `json:"action"`
	TimePeriod string `json:"time_period"`
}

func NewCriterion(name string, maxCount int, model, action, timePeriod string) *Criterion {
	return &Criterion{
		Name: name, MaxCount: maxCount,
		Model: model, Action: action, TimePeriod: timePeriod,
	}
}

type Criteria []*Criterion

type Plan struct {
	UUID     string `json:"uuid,omitempty"`
	Name     string `json:"name"`
	Criteria `json:"criteria,omitempty"`
}

type Plans []*Plan

func NewPlan(name string, criteria Criteria) *Plan {
	return &Plan{Name: name, Criteria: criteria}
}
