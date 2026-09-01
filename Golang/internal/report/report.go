package report

import (
	"encoding/json"
	"os"
	"time"
)

type CheckResult struct {
	Timestamp string `json:"timestamp"`
	Target    string `json:"target"`
	Status    string `json:"status"`
	Details   string `json:"details"`
}

func Save(results []CheckResult) error {

	fileName := "health_report_" +
		time.Now().Format("20060102_150405") +
		".json"

	file, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(results)
}
