package clickatell

import (
	"context"
	"testing"

	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func TestSendSMS(t *testing.T) {
	config := util.Config{
		ClickatellURL:    "https://platform.clickatell.com/messages/http/send",
		ClickatellAPIKey: "8SGnNNpZQG2jjl6lyzvaDg==",
		ClickatellAPIID:  "77e9286d0ebc4198a467535f29b013bf",
	}

	err := SendSMS(context.TODO(), &config, "0372312058", "Hello, world!")
	if err != nil {
		t.Errorf("Error sending SMS: %v", err)
	}
}
