package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driver = "mysql"
	source = "root:haiwei@100@tcp(localhost:3306)/db_redeem"
)

func Query(db *sql.DB) {
	var (
		code   string
		usedBy uint64
	)
	// 带有Query的方法表示查询结果可能会有多行
	rows, err := db.Query(`SELECT code, used_by FROM redeem_record WHERE used=?`, 1)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("error no rows")
		} else {
			log.Fatal(err)
		}
	}
	// rows没被关闭之前连接不会被释放（放回连接池中）
	defer rows.Close()
	for rows.Next() {
		// scan函数自动完成类型转换
		err := rows.Scan(&code, &usedBy)
		if err != nil {
			log.Fatalln("scan error")
		}
		log.Println(code, usedBy)
	}
	// 检查Next是否出错
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func QueryRow(db *sql.DB) {
	var code string
	// QueryRow只返回一行，即使有多个结果
	// QueryRow不返回error，只有scan的时候才返回error
	// sql.ErrNoRows只有使用QueryRow的情况下才出现，其他情况（比如Query）不用判断sql.ErrNoRows
	err := db.QueryRow(`SELECT code from redeem_record WHERE used=?`, 1).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no redeem code used")
			return
		} else {
			log.Fatal(err)
		}
	}
	log.Println(code)
}

func Prepare(db *sql.DB) {
	var (
		code   string
		usedBy uint64
	)
	stmt, err := db.Prepare(`SELECT code, used_by FROM redeem_record WHERE used=?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("error no rows")
		} else {
			log.Fatal(err)
		}
	}
	defer rows.Close()
	for rows.Next() {
		// scan函数自动完成类型转换
		err := rows.Scan(&code, &usedBy)
		// error时rows会自动关闭
		if err != nil {
			log.Fatalln("scan error")
		}
		log.Println(code, usedBy)
	}
	// 检查Next是否出错
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func Exec(db *sql.DB) {
	stmt, err := db.Prepare(`INSERT INTO redeem_record(code, expire_time) VALUES (?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	// Exec用于执行 INSERT,UPDATE,DELETE 等不需要返回行的语句
	res, err := stmt.Exec("888", time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		log.Fatal(err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)
}

func Transaction(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	updateStmt, err := tx.Prepare(`UPDATE redeem_record SET used=1 WHERE code=?`)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := tx.Query(`SELECT code FROM redeem_record WHERE used=?`, 0)
	if err != nil {
		log.Fatal(err)
	}
	var code string
	var codes []string
	for rows.Next() {
		if err := rows.Scan(&code); err != nil {
			log.Fatal(err)
		}
		codes = append(codes, code)
	}
	log.Println(codes)
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	for _, v := range codes {
		_, err := updateStmt.Exec(v)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatal("transaction commit error", err)
	}
	updateStmt.Close()
}

func main() {
	// Open不会建立连接，也不会检查source合法性
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal(err)
	}
	// 只有当进程结束了才建议关闭连接，通常情况下在一个进程中复用db
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	// 此时与数据库建立tcp连接
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	//Query(db)
	//Prepare(db)
	//Exec(db)
	//Transaction(db)
	QueryRow(db)
}
