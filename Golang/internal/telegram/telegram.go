package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendMessage(
	token string,
	chatID string,
	message string,
) error {

	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		token,
	)

	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
