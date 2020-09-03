package common

import (
	"fmt"
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
