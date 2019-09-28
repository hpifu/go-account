package mysql

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMysqlDB_InsertAccount(t *testing.T) {
	m, err := NewMysql("hatlonely:keaiduo1@tcp(test-mysql:3306)/hads?charset=utf8&parseTime=True&loc=Local")
	Convey("test mysqldb insert account", t, func() {
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)

		m.db.Where("email='hatlonely1@foxmail.com'").
			Or("phone='13112345678'").Delete(&Account{})

		Convey("insert account", func() {
			birthday, _ := time.Parse("2006-01-02", "1992-01-01")
			ok, err := m.InsertAccount(&Account{
				Email:     "hatlonely1@foxmail.com",
				Phone:     "13112345678",
				Password:  "123456",
				FirstName: "孙",
				LastName:  "悟空",
				Birthday:  birthday,
				Gender:    1,
			})
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			account, err := m.SelectAccountByPhoneOrEmail("hatlonely1@foxmail.com")
			So(err, ShouldBeNil)
			So(account.Email, ShouldEqual, "hatlonely1@foxmail.com")
			So(account.Phone, ShouldEqual, "13112345678")
			So(account.Password, ShouldEqual, "123456")
			So(account.FirstName, ShouldEqual, "孙")
			So(account.LastName, ShouldEqual, "悟空")
			So(account.Birthday, ShouldEqual, birthday)
			So(account.Gender, ShouldEqual, 1)

			Convey("insert dup email", func() {
				ok, err := m.InsertAccount(&Account{
					Email:     "hatlonely1@foxmail.com",
					Phone:     "13812345678",
					Password:  "123456",
					FirstName: "孙",
					LastName:  "悟空",
					Birthday:  birthday,
					Gender:    1,
				})
				So(err, ShouldNotBeNil)
				So(ok, ShouldBeFalse)
			})

			Convey("insert dup phone", func() {
				ok, err := m.InsertAccount(&Account{
					Email:     "hatlonely2@foxmail.com",
					Phone:     "13112345678",
					Password:  "123456",
					FirstName: "孙",
					LastName:  "悟空",
					Birthday:  birthday,
					Gender:    1,
				})
				So(err, ShouldNotBeNil)
				So(ok, ShouldBeFalse)
			})
		})
	})
}

func TestMysqlDB_SelectAccountByUsernameOrTelephoneOrEmail(t *testing.T) {
	m, err := NewMysql("hatlonely:keaiduo1@tcp(test-mysql:3306)/hads?charset=utf8&parseTime=True&loc=Local")
	Convey("test mysqldb select account by username or phone or email", t, func() {
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)

		m.db.Where("email='hatlonely1@foxmail.com'").
			Or("phone='13112345678'").Delete(&Account{})

		Convey("select account use empty key", func() {
			account, err := m.SelectAccountByPhoneOrEmail("")
			So(account, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("select account not phone or email", func() {
			account, err := m.SelectAccountByPhoneOrEmail("hatlonely1")
			So(account, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		birthday, _ := time.Parse("2006-01-02", "1992-01-01")
		ok, err := m.InsertAccount(&Account{
			Email:     "hatlonely1@foxmail.com",
			Phone:     "13112345678",
			Password:  "123456",
			FirstName: "孙",
			LastName:  "悟空",
			Birthday:  birthday,
			Gender:    1,
		})
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)

		Convey("select account by phone", func() {
			account, err := m.SelectAccountByPhoneOrEmail("hatlonely1@foxmail.com")
			So(err, ShouldBeNil)
			So(account.Email, ShouldEqual, "hatlonely1@foxmail.com")
			So(account.Phone, ShouldEqual, "13112345678")
			So(account.Password, ShouldEqual, "123456")
			So(account.FirstName, ShouldEqual, "孙")
			So(account.LastName, ShouldEqual, "悟空")
			So(account.Birthday, ShouldEqual, birthday)
			So(account.Gender, ShouldEqual, 1)
		})

		Convey("select account by email", func() {
			account, err := m.SelectAccountByPhoneOrEmail("13112345678")
			So(err, ShouldBeNil)
			So(account.Email, ShouldEqual, "hatlonely1@foxmail.com")
			So(account.Phone, ShouldEqual, "13112345678")
			So(account.Password, ShouldEqual, "123456")
			So(account.FirstName, ShouldEqual, "孙")
			So(account.LastName, ShouldEqual, "悟空")
			So(account.Birthday, ShouldEqual, birthday)
			So(account.Gender, ShouldEqual, 1)
		})
	})
}

func TestMysqlDB_UpdateAccount(t *testing.T) {
	m, err := NewMysql("hatlonely:keaiduo1@tcp(test-mysql:3306)/hads?charset=utf8&parseTime=True&loc=Local")
	Convey("test mysqldb update account", t, func() {
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)

		m.db.Where("email='hatlonely1@foxmail.com'").
			Or("phone='13112345678'").Delete(&Account{})

		birthday, _ := time.Parse("2006-01-02", "1992-01-01")
		ok, err := m.InsertAccount(&Account{
			Email:     "hatlonely1@foxmail.com",
			Phone:     "13112345678",
			Password:  "123456",
			FirstName: "孙",
			LastName:  "悟空",
			Birthday:  birthday,
			Gender:    1,
		})
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)

		account, err := m.SelectAccountByPhoneOrEmail("13112345678")
		So(err, ShouldBeNil)
		So(account, ShouldNotBeNil)

		Convey("update name", func() {
			ok, err = m.UpdateAccountName(account.ID, "猪", "八戒")
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
			acc, err := m.SelectAccountByPhoneOrEmail("13112345678")
			So(err, ShouldBeNil)
			So(acc, ShouldNotBeNil)
			So(acc.Email, ShouldEqual, "hatlonely1@foxmail.com")
			So(acc.Phone, ShouldEqual, "13112345678")
			So(acc.Password, ShouldEqual, "123456")
			So(acc.FirstName, ShouldEqual, "猪")
			So(acc.LastName, ShouldEqual, "八戒")
			So(acc.Birthday, ShouldEqual, birthday)
			So(acc.Gender, ShouldEqual, 1)
		})

		Convey("update email", func() {
			ok, err = m.UpdateAccountEmail(account.ID, "hatlonely2@foxmail.com")
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
			acc, err := m.SelectAccountByPhoneOrEmail("13112345678")
			So(err, ShouldBeNil)
			So(acc, ShouldNotBeNil)
			So(acc.Email, ShouldEqual, "hatlonely2@foxmail.com")
			So(acc.Phone, ShouldEqual, "13112345678")
			So(acc.Password, ShouldEqual, "123456")
			So(acc.FirstName, ShouldEqual, "孙")
			So(acc.LastName, ShouldEqual, "悟空")
			So(acc.Birthday, ShouldEqual, birthday)
			So(acc.Gender, ShouldEqual, 1)
		})

		Convey("update phone", func() {
			ok, err = m.UpdateAccountPhone(account.ID, "13111112222")
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
			acc, err := m.SelectAccountByPhoneOrEmail("hatlonely1@foxmail.com")
			So(err, ShouldBeNil)
			So(acc, ShouldNotBeNil)
			So(acc.Email, ShouldEqual, "hatlonely1@foxmail.com")
			So(acc.Phone, ShouldEqual, "13111112222")
			So(acc.Password, ShouldEqual, "123456")
			So(acc.FirstName, ShouldEqual, "孙")
			So(acc.LastName, ShouldEqual, "悟空")
			So(acc.Birthday, ShouldEqual, birthday)
			So(acc.Gender, ShouldEqual, 1)
		})

		Convey("update birthday", func() {
			b, _ := time.Parse("2006-01-02", "1990-12-12")
			ok, err = m.UpdateAccountBirthday(account.ID, b)
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
			acc, err := m.SelectAccountByPhoneOrEmail("13112345678")
			So(err, ShouldBeNil)
			So(acc, ShouldNotBeNil)
			So(acc.Email, ShouldEqual, "hatlonely1@foxmail.com")
			So(acc.Phone, ShouldEqual, "13112345678")
			So(acc.Password, ShouldEqual, "123456")
			So(acc.FirstName, ShouldEqual, "孙")
			So(acc.LastName, ShouldEqual, "悟空")
			So(acc.Birthday, ShouldEqual, b)
			So(acc.Gender, ShouldEqual, 1)
		})

		Convey("update password", func() {
			ok, err = m.UpdateAccountPassword(account.ID, "11112222")
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
			acc, err := m.SelectAccountByPhoneOrEmail("13112345678")
			So(err, ShouldBeNil)
			So(acc, ShouldNotBeNil)
			So(acc.Email, ShouldEqual, "hatlonely1@foxmail.com")
			So(acc.Phone, ShouldEqual, "13112345678")
			So(acc.Password, ShouldEqual, "11112222")
			So(acc.FirstName, ShouldEqual, "孙")
			So(acc.LastName, ShouldEqual, "悟空")
			So(acc.Birthday, ShouldEqual, birthday)
			So(acc.Gender, ShouldEqual, 1)
		})

		Convey("update gender", func() {
			ok, err = m.UpdateAccountGender(account.ID, 0)
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
			acc, err := m.SelectAccountByPhoneOrEmail("13112345678")
			So(err, ShouldBeNil)
			So(acc, ShouldNotBeNil)
			So(acc.Email, ShouldEqual, "hatlonely1@foxmail.com")
			So(acc.Phone, ShouldEqual, "13112345678")
			So(acc.Password, ShouldEqual, "123456")
			So(acc.FirstName, ShouldEqual, "孙")
			So(acc.LastName, ShouldEqual, "悟空")
			So(acc.Birthday, ShouldEqual, birthday)
			So(acc.Gender, ShouldEqual, 0)
		})
	})
}
