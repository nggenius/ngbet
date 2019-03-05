package main

import (
	"config"
	"testing"

	"github.com/levigross/grequests"
)

func Test_main(t *testing.T) {
	config.LoadConfig()

	r, err := grequests.Get(config.Setting.Bet365.Host, nil)
	if err != nil {
		t.Fatalf("get session failed: %s ", err.Error())
	}

	res := r.RawResponse.Cookies()
	if len(res) == 0 {
		t.Fatalf("get session failed")
	}

	session := res[1].Value
	t.Log("Sessionid=" + session)
}
