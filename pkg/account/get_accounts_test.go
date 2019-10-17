package account

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGETAccounts(t *testing.T) {
	Convey("test get accounts", t, func() {
		cli := NewClient("127.0.0.1:16060", 4, 20*time.Millisecond, 20*time.Millisecond)
		accounts, err := cli.GETAccounts("123", "1c15b6b0b18aa0d3a5d2de37484f992c", []int{666})
		t.Log(err)
		t.Log(accounts)
	})
}
