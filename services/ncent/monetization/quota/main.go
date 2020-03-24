package quota

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gitlab.com/ncent/monetization/services"
)

const dbDriver = "postgres"

var (
	databaseName        string
	databaseUser        string
	databasePasswd      string
	databaseHost        string
	databasePort        int
	quotaUsageTableName string
	log                 *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"service": "service.ncent.monetization.quota",
	})

	databaseName = "ncentdb"
	databaseUser = "ncentuser"
	databasePort = 5439
	quotaUsageTableName = "challenge_quota_usage"

	dbPasswd, ok := os.LookupEnv("DB_PASSWD")

	if !ok {
		log.Warnf("Database Password for Redshift needed")
	} else {
		databasePasswd = dbPasswd
	}

	hostPasswd, ok := os.LookupEnv("DB_HOST")

	if !ok {
		log.Warnf("Database Host for Redshift needed")
	} else {
		databaseHost = hostPasswd
	}
}

type IQuotaUsageRepository interface {
	GetQuotaUsage(uuid string, action string, amount int, timePeriod string) (int, error)
}

type QuotaUsageRepository struct {
	client services.ISQLDB
}

func NewQuotaUsageRepository() *QuotaUsageRepository {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		databaseUser,
		databasePasswd,
		databaseHost,
		databasePort,
		databaseName,
	)

	client, err := sql.Open(dbDriver, connStr)

	if err != nil {
		log.Fatal("Failed to open connection with database")
		return nil
	}

	return &QuotaUsageRepository{client: client}
}

func (qs QuotaUsageRepository) GetQuotaUsage(uuid string, action string, amount int, timePeriod string) (int, error) {
	log.Infof("Get quota usage from UUID [%s] with action [%s] and amount [%d] with time [%s]",
		uuid, action, amount, timePeriod,
	)

	query := `
	SELECT
			COUNT(*) AS count
	FROM
			challenge_usage_quota
	WHERE
			challenge_usage_quota.uuid = '%s'
			AND action = '%s'
			AND timestamp between (GETDATE() - '%d %s' ::INTERVAL)
			AND GETDATE();`

	queryStr := fmt.Sprintf(query, uuid, action, amount, timePeriod)
	log.Infof("Query %s", queryStr)

	rows, err := qs.client.Query(queryStr)

	if err != nil {
		log.Errorf("Error on query %v", err)
		return 0, err
	}

	var qtds []int
	for rows.Next() {
		var p int
		if err := rows.Scan(&p); err != nil {
			log.Errorf("Error on scan row %v", err)
			return 0, err
		}
		qtds = append(qtds, p)
	}

	if err == sql.ErrNoRows || len(qtds) == 0 {
		return 0, errors.New("No result returned on query")
	}

	log.Infof("UUID %s with Action %s and amount %d for period of %s has %d",
		uuid, action, amount, timePeriod, qtds[0],
	)

	return qtds[0], nil
}
