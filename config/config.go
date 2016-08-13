package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	SMTPHost      string
	SMTPPassword  string
	SMTPPort      int
	SMTPRecipient string
	SMTPSender    string
	SMTPUsername  string
}

type sendgridService struct {
	Credentials struct {
		Hostname string
		Password string
		Username string
	}
}

func New() (*Config, error) {
	c := &Config{}
	c.SMTPRecipient = os.Getenv("SMTP_RECIPIENT")
	c.SMTPSender = os.Getenv("SMTP_SENDER")

	if os.Getenv("VCAP_SERVICES") != "" {
		services := make(map[string][]sendgridService)
		if err := json.Unmarshal([]byte(os.Getenv("VCAP_SERVICES")), &services); err != nil {
			return nil, fmt.Errorf("could not parse '%s' as json", os.Getenv("VCAP_SERVICES"))
		}

		sendgridServices := services["sendgrid"]
		if len(sendgridServices) > 0 {
			c.SMTPHost = sendgridServices[0].Credentials.Hostname
			c.SMTPPassword = sendgridServices[0].Credentials.Password
			c.SMTPUsername = sendgridServices[0].Credentials.Username
			c.SMTPPort = 25
		}
	}

	if c.SMTPHost == "" {
		c.SMTPHost = os.Getenv("SMTP_HOST")
	}
	if c.SMTPPassword == "" {
		c.SMTPPassword = os.Getenv("SMTP_PASSWORD")
	}
	if c.SMTPPort == 0 {
		if os.Getenv("SMTP_PORT") != "" {
			port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
			if err != nil {
				return nil, fmt.Errorf("could not convert '%s' to an integer", os.Getenv("SMTP_PORT"))
			}
			c.SMTPPort = port
		}
		c.SMTPUsername = os.Getenv("SMTP_USERNAME")
	}

	return c, nil
}
