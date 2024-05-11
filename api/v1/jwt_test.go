package v1

import (
	"fmt"
	"openfly/conf"
	"openfly/logger"
	"testing"
)

func TestGenToken(t *testing.T) {
	token, err := GenToken("admin", "admin")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(token)
}
func TestParseToken(t *testing.T) {
	conf.ParseConfig("../../config.toml")
	logger.InitLog()

	token, err := GenToken("admin", "admin")
	if err != nil {
		t.Error(err)
		return
	}
	_, gerr := ParseToken(token)
	if gerr != nil {
		t.Error(gerr)
		return
	}
}
