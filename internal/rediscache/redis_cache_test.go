package rediscache

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRedisCache_SetGetDelAccount(t *testing.T) {
	rc, err := NewRedisCache(&Option{
		Address:         "test-redis:6379",
		Timeout:         20 * time.Millisecond,
		Retries:         3,
		PoolSize:        20,
		TokenExpiration: time.Minute,
	})
	Convey("test set account", t, func() {
		So(err, ShouldBeNil)
		So(rc, ShouldNotBeNil)

		err = rc.DelAccount("f80f8d59a1694602b89efa24a9028282")
		So(err, ShouldBeNil)

		Convey("set account and get account", func() {
			err = rc.SetAccount("f80f8d59a1694602b89efa24a9028282", &Account{
				ID:       666,
				Email:    "hpifu@foxmail.com",
				Phone:    "+8612345678901",
				Password: "e010597fcf126d58fdfa36e636f8fc9e",
			})
			So(err, ShouldBeNil)

			account, err := rc.GetAccount("f80f8d59a1694602b89efa24a9028282")
			So(err, ShouldBeNil)
			So(account.ID, ShouldEqual, 666)
			So(account.Email, ShouldEqual, "hpifu@foxmail.com")
			So(account.Phone, ShouldEqual, "+8612345678901")
			So(account.Password, ShouldEqual, "e010597fcf126d58fdfa36e636f8fc9e")
		})
	})
}
