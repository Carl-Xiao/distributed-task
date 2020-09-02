package common

import "testing"

func TestExtractJobName(t *testing.T) {
	t.Log(ExtractJobName("/cron/jobs/job12321"))
}
