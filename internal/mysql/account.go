package mysql

import (
	"fmt"
	"github.com/hpifu/go-account/internal/c"
	"github.com/hpifu/go-account/internal/rule"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type Account struct {
	ID        int       `gorm:"type:bigint(20) auto_increment;primary_key" json:"id"`
	Email     string    `gorm:"type:varchar(64);index:email_idx" json:"email"`
	Phone     string    `gorm:"type:varchar(64);index:phone_idx" json:"phone"`
	FirstName string    `gorm:"type:varchar(32);" json:"firstName"`
	LastName  string    `gorm:"type:varchar(32);" json:"lastName"`
	Password  string    `gorm:"type:varchar(32);" json:"password"`
	Birthday  time.Time `gorm:"type:timestamp default '1970-01-02 00:00:00'" json:"birthday"`
	Gender    c.Gender  `gorm:"type:int(1);" json:"gender"`
	Avatar    string    `gorm:"type:varchar(512);" json:"avatar"`
	Role      int       `gorm:"type:bigint(20) default 0;not null" json:"role"`
}

func (m *Mysql) SelectAccountByPhoneOrEmail(key string) (*Account, error) {
	if err := rule.ValidPhone(key); err == nil {
		return m.SelectAccountByPhone(key)
	}
	if err := rule.ValidEmail(key); err == nil {
		return m.SelectAccountByEmail(key)
	}

	return nil, fmt.Errorf("key [%v] is not a valid phone or email", key)
}

func (m *Mysql) SelectAccountByPhone(phone string) (*Account, error) {
	account := &Account{}
	if err := m.db.Where("phone=?", phone).First(account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return account, nil
}

func (m *Mysql) SelectAccountByEmail(email string) (*Account, error) {
	account := &Account{}
	if err := m.db.Where("email=?", email).First(account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return account, nil
}

func (m *Mysql) UpdateAccountName(id int, firstName string, lastName string) (bool, error) {
	account := &Account{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
	}
	if err := m.db.Model(account).Where("id=?", account.ID).Update(account).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (m *Mysql) UpdateAccountEmail(id int, email string) (bool, error) {
	account := &Account{
		ID:    id,
		Email: email,
	}
	accountDB := &Account{}
	err := m.db.Where("email=?", account.Email).First(accountDB).Error
	if err == nil {
		return false, fmt.Errorf("email [%v] is already exists", accountDB.Email)
	}
	if err != gorm.ErrRecordNotFound {
		return false, err
	}
	if err := m.db.Model(account).Where("id=?", account.ID).Update(account).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (m *Mysql) UpdateAccountPhone(id int, phone string) (bool, error) {
	account := &Account{
		ID:    id,
		Phone: phone,
	}
	accountDB := &Account{}
	err := m.db.Where("phone=?", account.Phone).First(accountDB).Error
	if err == nil {
		return false, fmt.Errorf("phone [%v] is already exists", accountDB.Phone)
	}
	if err != gorm.ErrRecordNotFound {
		return false, err
	}
	if err := m.db.Model(account).Where("id=?", account.ID).Update(account).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (m *Mysql) UpdateAccountBirthday(id int, birthday time.Time) (bool, error) {
	account := &Account{
		ID:       id,
		Birthday: birthday,
	}
	if err := m.db.Model(account).Where("id=?", account.ID).Update(account).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (m *Mysql) UpdateAccountPassword(id int, password string) (bool, error) {
	account := &Account{
		ID:       id,
		Password: password,
	}
	if err := m.db.Model(account).Where("id=?", account.ID).Update(account).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (m *Mysql) UpdateAccountGender(id int, gender c.Gender) (bool, error) {
	if err := m.db.Model(&Account{}).Where("id=?", id).Update("gender", gender).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (m *Mysql) UpdateAccountAvatar(id int, avatar string) (bool, error) {
	account := &Account{
		ID:     id,
		Avatar: avatar,
	}
	if err := m.db.Model(account).Where("id=?", account.ID).Update(account).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (m *Mysql) InsertAccount(account *Account) (bool, error) {
	if account.Email == "" && account.Phone == "" {
		return false, fmt.Errorf("email or phone are is null, account [%#v]", account)
	}

	accountDB := &Account{}
	var conditions []string
	if account.ID != 0 {
		conditions = append(conditions, fmt.Sprintf("id=%v", account.ID))
	}
	if account.Phone != "" {
		conditions = append(conditions, fmt.Sprintf("phone='%v'", account.Phone))
	}
	if account.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email='%v'", account.Email))
	}
	err := m.db.Where(strings.Join(conditions, " OR ")).First(accountDB).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if err != gorm.ErrRecordNotFound {
		if accountDB.ID == account.ID {
			return false, fmt.Errorf("accountID [%v] is already exists", accountDB.ID)
		}
		if accountDB.Phone == account.Phone {
			return false, fmt.Errorf("phone [%v] is already exists", accountDB.Phone)
		}
		if accountDB.Email == account.Email {
			return false, fmt.Errorf("email [%v] is already exists", accountDB.Email)
		}
	}

	if err := m.db.Create(account).Error; err != nil {
		return false, err
	}

	return true, nil
}
