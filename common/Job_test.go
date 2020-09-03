package common

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"testing"
	"time"
)

func TestExtractJobName(t *testing.T) {
	t.Log(ExtractJobName("/cron/jobs/job12321"))
}

func TestTimeEncoder(t *testing.T) {
	var recentTime *time.Time

	now := time.Now()
	recentTime = &now

	timeAfter := (*recentTime).Sub(now)

	fmt.Println(timeAfter)
}

func TestCron(t *testing.T) {
	cron := "* * * * *"
	var (
		cronExpr *cronexpr.Expression
		err      error
	)
	if cronExpr, err = cronexpr.Parse(cron); err != nil {
		Error(err.Error())
		return
	}
	fmt.Println(cronExpr.Next(time.Now()))
}
