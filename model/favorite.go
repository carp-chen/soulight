package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
)

/*
   This code is generated by gendry
*/

// Favorite is a mapping object for favorite table in mysql
type Favorite struct {
	UserID     int       `json:"user_id"`
	AdviserID  int       `json:"adviser_id"`
	CreateTime time.Time `json:"create_time"`
}

//GetOne gets one record from table favorite by condition "where"
func GetOneFavorite(db *sql.DB, where map[string]interface{}) (*Favorite, error) {
	if nil == db {
		return nil, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildSelect("favorite", where, nil)
	if nil != err {
		return nil, err
	}
	row, err := db.Query(cond, vals...)
	if nil != err || nil == row {
		return nil, err
	}
	defer row.Close()
	var res *Favorite
	err = scanner.Scan(row, &res)
	return res, err
}

//GetMulti gets multiple records from table favorite by condition "where"
func GetMultiFavorite(db *sql.DB, where map[string]interface{}) ([]*Favorite, error) {
	if nil == db {
		return nil, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildSelect("favorite", where, nil)
	if nil != err {
		return nil, err
	}
	row, err := db.Query(cond, vals...)
	if nil != err || nil == row {
		return nil, err
	}
	defer row.Close()
	var res []*Favorite
	err = scanner.Scan(row, &res)
	return res, err
}

//Insert inserts an array of data into table favorite
func InsertFavorite(db *sql.DB, data []map[string]interface{}) (int64, error) {
	if nil == db {
		return 0, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildInsert("favorite", data)
	if nil != err {
		return 0, err
	}
	result, err := db.Exec(cond, vals...)
	if nil != err || nil == result {
		return 0, err
	}
	return result.LastInsertId()
}

//Update updates the table favorite
func UpdateFavorite(db *sql.DB, where, data map[string]interface{}) (int64, error) {
	if nil == db {
		return 0, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildUpdate("favorite", where, data)
	if nil != err {
		return 0, err
	}
	result, err := db.Exec(cond, vals...)
	if nil != err {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete deletes matched records in favorite
func DeleteFavorite(db *sql.DB, where, data map[string]interface{}) (int64, error) {
	if nil == db {
		return 0, errors.New("sql.DB object couldn't be nil")
	}
	cond, vals, err := builder.BuildDelete("favorite", where)
	if nil != err {
		return 0, err
	}
	result, err := db.Exec(cond, vals...)
	if nil != err {
		return 0, err
	}
	return result.RowsAffected()
}