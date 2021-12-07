package database

type DB interface {
	Insert(tableName string, id int64, data interface{}) error
	Find(tableName string, id int64, target interface{}) error
	Update(tableName string, id int64, data interface{}) error
}
