package main

import (
	"devops/internal/checker"
	"devops/internal/report"
	"devops/internal/retry"
	"devops/internal/telegram"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	TelegramBotToken string   `json:"telegram_bot_token"`
	TelegramChatID   string   `json:"telegram_chat_id"`
	IntervalSeconds  int      `json:"interval_seconds"`
	RetryCount       int      `json:"retry_count"`
	Domains          []string `json:"domains"`
}

var results []report.CheckResult
var mutex sync.Mutex

func addResult(
	target string,
	status string,
	details string,
) {

	mutex.Lock()

	defer mutex.Unlock()

	results = append(
		results,
		report.CheckResult{
			Timestamp: time.Now().Format(time.RFC3339),
			Target:    target,
			Status:    status,
			Details:   details,
		},
	)
}

func loadConfig() (Config, error) {

	var cfg Config

	file, err := os.ReadFile(
		"configs/config.json",
	)

	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(
		file,
		&cfg,
	)

	return cfg, err
}

func checkDomain(
	domain string,
	cfg Config,
) {

	protocols := []string{
		"http",
		"https",
	}

	for _, protocol := range protocols {

		url := fmt.Sprintf(
			"%s://%s",
			protocol,
			domain,
		)

		err := retry.Retry(
			cfg.RetryCount,
			func() error {

				return checker.CheckURL(
					url,
				)
			},
		)

		if err != nil {

			msg := fmt.Sprintf(
				"Url: %s\nError: %v",
				url,
				err,
			)

			log.Println(msg)

			addResult(
				url,
				"FAILED",
				err.Error(),
			)

			_ = telegram.SendMessage(
				cfg.TelegramBotToken,
				cfg.TelegramChatID,
				msg,
			)

			continue
		}

		log.Println(
			url,
			"OK",
		)

		addResult(
			url,
			"OK",
			"200 OK",
		)
	}

	ips, err := checker.ResolveIPs(
		domain,
	)

	if err != nil {

		addResult(
			domain,
			"FAILED",
			err.Error(),
		)

		return
	}

	for _, ip := range ips {

		url := "http://" + ip

		err := retry.Retry(
			cfg.RetryCount,
			func() error {

				return checker.CheckURL(
					url,
				)
			},
		)

		if err != nil {

			addResult(
				ip,
				"FAILED",
				err.Error(),
			)

			continue
		}

		addResult(
			ip,
			"OK",
			"reachable",
		)
	}
}

func runChecks(
	cfg Config,
) {

	results = nil

	var wg sync.WaitGroup

	for _, domain := range cfg.Domains {

		wg.Add(1)

		go func(d string) {

			defer wg.Done()

			checkDomain(
				d,
				cfg,
			)

		}(domain)
	}

	wg.Wait()

	err := report.Save(
		results,
	)

	if err != nil {
		log.Println(err)
	}
}

func main() {

	cfg, err := loadConfig()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(
		"Network Monitor Started",
	)

	ticker := time.NewTicker(
		time.Duration(
			cfg.IntervalSeconds,
		) * time.Second,
	)

	defer ticker.Stop()

	for {

		runChecks(cfg)

		<-ticker.C
	}
}
