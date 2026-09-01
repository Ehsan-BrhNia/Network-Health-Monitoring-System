package checker

import (
	"fmt"
	"net/http"
	"time"
)

func CheckURL(url string) error {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf(
			"status code %d",
			resp.StatusCode,
		)
	}

	return nil
}
