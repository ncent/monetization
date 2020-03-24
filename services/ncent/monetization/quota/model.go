package quota

type QuotaQuery struct {
	UserUUID   string `json:"uuid"`
	Action     string `json:"action"`
	Amount     int    `json:"amount"`
	TimePeriod string `json:"time_period"`
}

type QuotaQueries []*QuotaQuery

type QuotaUsage struct {
	UserUUID  string `json:"uuid"`
	Action    string `json:"action"`
	TimeStamp string `json:"timestamp"`
}

type QuotaUsages []*QuotaUsage
