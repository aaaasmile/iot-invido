package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/aaaasmile/iot-invido/util"
	_ "github.com/mattn/go-sqlite3"
)

type LiteDB struct {
	connDb       *sql.DB
	DebugSQL     bool
	SqliteDBPath string
}

type UserInfo struct {
	ID        int
	Username  string
	Password  string
	Salt      string
	Active    bool
	Timestamp time.Time
}

func (ld *LiteDB) GetConnDB() *sql.DB {
	return ld.connDb
}

func (ld *LiteDB) OpenSqliteDatabase() error {
	var err error
	// Source control should be only an empty navrepo.db.
	dbname := util.GetFullPath(ld.SqliteDBPath)
	log.Println("Using the sqlite file: ", dbname)
	ld.connDb, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	return nil
}

func (ld *LiteDB) GetNewTransaction() (*sql.Tx, error) {
	tx, err := ld.connDb.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (ld *LiteDB) DeleteUser(tx *sql.Tx, recID int) error {
	q := fmt.Sprintf(`DELETE FROM User WHERE id=%d;`, recID)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec()
	return err
}

func (ld *LiteDB) UpdateUser(tx *sql.Tx, recID int, username, password string, active bool) error {
	var q string
	if recID == 0 {
		panic("RecID is null")
	}
	if username == "" {
		return fmt.Errorf("Username is empty")
	}
	if password == "" {
		return fmt.Errorf("Password is empty")
	}
	q = fmt.Sprintf(`UPDATE User SET Username=?,Password=?,Active=? WHERE id=%d;`, recID)

	if ld.DebugSQL {
		log.Println("Query is", q)
	}
	updateMore, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(updateMore).Exec(username, password, active)

	if err != nil {
		log.Println("Error in UpdateUser")
		return err
	}
	if ld.DebugSQL {
		log.Println("User updated OK: ", username)
	}
	return nil
}

func (ld *LiteDB) FetchUser(username string) ([]*UserInfo, error) {
	q := `SELECT id,Username,Password,Salt,Active,Timestamp FROM User WHERE Username = ?;`
	if ld.DebugSQL {
		log.Println("Query is", q)
	}
	rows, err := ld.connDb.Query(q, username)
	if err != nil {
		return nil, err
	}
	result := make([]*UserInfo, 0)
	defer rows.Close()
	for rows.Next() {
		item := UserInfo{}
		var ts int64

		if err := rows.Scan(&item.ID, &item.Username, &item.Password, &item.Salt, &item.Active, &ts); err != nil {
			log.Println("Error in scan lite nav ", err)
			return nil, err
		}
		item.Timestamp = time.Unix(ts, 0)

		result = append(result, &item)
	}
	return result, nil
}

func (ld *LiteDB) InsertUser(tx *sql.Tx, ui *UserInfo) error {
	if ui.Username == "" || ui.Password == "" || ui.Salt == "" {
		return fmt.Errorf("at least one of username, password or salt is empty")
	}
	q := `INSERT INTO User(Username,Password,Salt,Active,Timestamp) VALUES (?,?,?,?,?);`
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	insertMore, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(insertMore).Exec(ui.Username, ui.Password, ui.Salt, ui.Active, ui.Timestamp)
	if err != nil {
		return err
	}
	if ld.DebugSQL {
		log.Println("User added OK: ", ui)
	}
	return nil
}
