package store

import (
	"errors"
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
)

type DataBaseImpl struct {
	schemas.Database
	dictionarySys *DictionarySys
	infos         schemas.InfoSchemas
	dataBaseName  string
	conf          *conf.Cfg
	tableCache    map[string]schemas.Table
}

func NewDataBaseImpl(infos schemas.InfoSchemas, conf *conf.Cfg, databaseName string) (schemas.Database, error) {
	var database = new(DataBaseImpl)
	database.infos = infos
	database.conf = conf
	isExist := database.infos.GetSchemaExist(databaseName)
	if isExist {
		return nil, errors.New("数据库已经存在")
	}
	database.tableCache = make(map[string]schemas.Table)
	database.dataBaseName = databaseName
	database.infos.PutDatabaseCache(database)
	return database, nil
}

func (d *DataBaseImpl) Name() string {
	return d.dataBaseName
}

func (d *DataBaseImpl) GetTable(name string) (schemas.Table, error) {
	if d.tableCache[name] == nil {
		return nil, errors.New("没有该表")
	}
	return d.tableCache[name], nil
}

func (d *DataBaseImpl) ListTables() []schemas.Table {
	var tableArrays = make([]schemas.Table, 0)

	for _, v := range d.tableCache {
		tableArrays = append(tableArrays, v)
	}
	return tableArrays
}

func (d *DataBaseImpl) DropTable(name string) error {
	delete(d.tableCache, name)
	return nil
}

func (d *DataBaseImpl) ListTableName() []string {
	panic("implement me")
}

func (d *DataBaseImpl) CreateTable(conf *conf.Cfg, stmt *sqlparser.DDL) (schemas.Table, error) {
	isExist := stmt.IfExists
	tableName := ""
	if isExist {
		tableName = stmt.NewName.Name.String()
	} else {
		tableName = stmt.Table.Name.String()
	}
	if stmt.Table.Name.String() == "" {
		tableName = stmt.NewName.Name.String()
	}

	colDefs := stmt.TableSpec.Columns

	indexDefs := stmt.TableSpec.Indexes

	tableTuple := d.createtable(conf, tableName, indexDefs, colDefs)

	tablespace := NewTableSpaceFile(conf, d.dataBaseName, tableName, d.dictionarySys.currentSpaceId, false)

	fmt.Println(tablespace.GetSpaceId())

	err := d.dictionarySys.CreateTable(d.dataBaseName, tableTuple)

	return nil, err
}

func (d *DataBaseImpl) createtable(conf *conf.Cfg, tableName string, indexDefs []*sqlparser.IndexDefinition, colDefs []*sqlparser.ColumnDefinition) *TableTupleMeta {
	tableTuple := NewTupleMeta(d.dataBaseName, tableName, conf).(*TableTupleMeta)

	tableTuple.WriteIndexDefinitions(indexDefs)
	tableTuple.WriteTupleColumns(colDefs)
	tableTuple.FlushToDisk()
	return tableTuple
}
