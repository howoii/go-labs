package main

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

type Account struct {
	AccountID uint64  `db:"account_id"`
	Amount    float32 `db:"curr_amount"`
	Storage   float32 `db:"curr_storage"`
}

//table:
//CREATE TABLE IF NOT EXISTS account_jumpsuit_info
//(
//    `account_id`               BIGINT UNSIGNED NOT NULL comment '用户id',
//    `curr_amount`              DECIMAL(6,1) NOT NULL DEFAULT 0 comment '总进度',
//    `curr_storage`             FlOAT NOT NULL DEFAULT 0 comment '熔金炉储量',
//    `update_time`              INT UNSIGNED NOT NULL DEFAULT 0 comment '挂机刷新时间',
//    `rewarded_box_num`         INT UNSIGNED NOT NULL DEFAULT 0 comment '领取过的宝箱数量',
//    `rewarded_stages`          BLOB comment '领取过的阶段奖励id(数组)',
//    `helped_list`              BLOB comment '我助力过的好友id(数组)',
//
//    PRIMARY KEY(`account_id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci comment = 'jumpsuit2B大活动玩家信息';

func TestDecimal(t *testing.T) {
	db, err := sqlx.Open(driver, source)
	if err != nil {
		t.Error(err)
	}

	in := Account{
		AccountID: 1003,
		Amount:    1 / float32(3),
		Storage:   1 / float32(3),
	}
	insert := fmt.Sprintf(`INSERT INTO account_jumpsuit_info (account_id, curr_amount, curr_storage) VALUES(?,?,?)`)
	_, err = db.Exec(insert, in.AccountID, in.Amount, in.Storage)
	if err != nil {
		t.Error(err)
	}

	out := Account{}
	query := fmt.Sprintf(`SELECT account_id, curr_amount, curr_storage FROM account_jumpsuit_info WHERE account_id=?`)
	err = db.Get(&out, query, in.AccountID)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(in, out)
}
