package quota

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/DATA-DOG/go-sqlmock"
)

var _ = Describe("Create Quotas Usage Query Service Suite", func() {
	var (
		quotaUsageRepository *QuotaUsageRepository
	)

	JustBeforeEach(func() {
		quotaUsageRepository = NewQuotaUsageRepository()

		db, mock, err := sqlmock.New()

		if err != nil {
			log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		quotaUsageRepository.client = db

		var query = `^SELECT
				COUNT\(\*\) AS count
		FROM
				challenge_usage_quota
		WHERE
				challenge_usage_quota.uuid = '.+'
				AND action = '.+'
				AND timestamp between \(GETDATE\(\) - '.+' ::INTERVAL\)
				AND GETDATE\(\);$`

		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
		mock.ExpectQuery(query).WillReturnRows(rows)
	})
})

func TestQuotaUsagesDataService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Quotas Usage Query Service Suite")
}
