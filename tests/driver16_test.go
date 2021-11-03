package main

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	driver = "mysql"
	source = "root:123456@tcp(localhost:3306)/db_redeem?parseTime=1"
)

type AccountInfo struct {
	AccountId uint64 `db:"account_id"`
	Name      string `db:"name"`
	CreateAt  string `db:"create_at"`
}

func (this *AccountInfo) GetAccountId() uint64 {
	return this.AccountId
}

func (this *AccountInfo) GetName() string {
	return this.Name
}

func TestDriver16(t *testing.T) {
	db, err := sqlx.Open(driver, source)
	if err != nil {
		t.Error(err)
	}

	account := AccountInfo{}
	err = db.Get(&account, "SELECT account_id, name, create_at FROM account_info WHERE account_id=?", 10001)
	if err != nil {
		t.Error(errors.WithStack(err))
	}
	fmt.Printf("%#v\n", account.CreateAt)
	fmt.Println(time.ParseInLocation("2006-01-02 15:04:05", account.CreateAt, time.Local))
}
