package mail

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	Convey("test", t, func() {
		mc := &MailClient{
			Password: "iwvxfcvxljrxbfed",
			From:     "hatlonely@foxmail.com",
			Server:   "smtp.qq.com",
			Port:     25,
		}

		err := mc.Send("hatlonely@foxmail.com", "hatlonely 账号验证", NewAuthCodeTpl("悟空", "孙", "654321"))
		fmt.Println(err)
		So(err, ShouldBeNil)
	})
}
