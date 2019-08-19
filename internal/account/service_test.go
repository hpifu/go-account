package account

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	Convey("test", t, func() {
		So(len(NewToken()), ShouldEqual, 32)
		So(len(NewCode()), ShouldEqual, 6)
	})
}
