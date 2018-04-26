package lib

import (
	"log"
	"net/http"
	"strings"

	"github.com/obieq/docker-hub-webhook/config"
)

func SendSlackSwarmRedeploySuccessMessage(repoName string, tagName string) error {
	// is slack configured?  if not, then nothing to do here
	cfg := config.Cfg().Slack
	if cfg == nil {
		return nil
	}

	message := `{"channel": "` + cfg.Channel + `", "username": "` + cfg.UserName + `", "icon_emoji": "` + cfg.EmojiIcon + `", "text": "Swarm Redeploy", "attachments": [{"text": "*redeploy succeeded for:*\n<https://github.com|` + repoName + ":" + tagName + `>", "color": "good"}]}`
	if err := sendSlackMessage(message, cfg.Webhook); err != nil {
		log.Println("SendSlackSwarmRedeploySuccessMessage error: " + err.Error())
		return err
	}

	log.Println("successfully sent SendSlackSwarmRedeploySuccessMessage")

	return nil
}

func SendSlackSwarmRedeployFailureMessage(repoName string, tagName string, errorMessage string) error {
	// is slack configured?  if not, then nothing to do here
	cfg := config.Cfg().Slack
	if cfg == nil {
		return nil
	}

	message := `{"channel": "` + cfg.Channel + `",
		         "username": "` + cfg.UserName + `",
				 "icon_emoji": "` + cfg.EmojiIcon + `",
				 "text": "Swarm Redeploy", "attachments": [{"text": "*redeploy failed for:*\n<https://github.com|` + repoName + ":" + tagName + `>\n*error:* ` + errorMessage + `",
			     "color": "danger"}]}`
	if err := sendSlackMessage(message, cfg.Webhook); err != nil {
		log.Println("SendSlackSwarmRedeployFailureMessage error: " + err.Error())
		return err
	}

	log.Println("successfully sent SendSlackSwarmRedeployFailureMessage")

	return nil
}

func sendSlackMessage(message string, webhook string) error {
	body := strings.NewReader(message)
	req, err := http.NewRequest("POST", webhook, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("SendSlackSwarmRedeploySuccessMessage error: " + err.Error())
	}
	defer resp.Body.Close()

	return nil
}
