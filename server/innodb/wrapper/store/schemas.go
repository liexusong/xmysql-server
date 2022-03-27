package store

import (
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
	"github.com/zhukovaskychina/xmysql-server/util"
	"io/ioutil"
	"path"
)

type InfoSchemaManager struct {
	conf          *conf.Cfg
	sysTableSpace TableSpace
	dictionarySys *DictionarySys
	schemaMap     map[string]schemas.Database
}

func NewInfoSchemaManager(conf *conf.Cfg) schemas.InfoSchemas {
	var infoSchemaManager = new(InfoSchemaManager)
	infoSchemaManager.conf = conf
	infoSchemaManager.schemaMap = make(map[string]schemas.Database)
	infoSchemaManager.sysTableSpace = NewSysTableSpace(conf)
	infoSchemaManager.initSysSchemas()
	return infoSchemaManager
}

func (i *InfoSchemaManager) initSysSchemas() {

	dataDir := i.conf.DataDir
	fs, errors := ioutil.ReadDir(dataDir)
	if errors != nil {
		panic("出现异常")
	}
	for _, v := range fs {
		if v.IsDir() {
			dirName := v.Name()
			dataBaseDir := path.Join(dataDir, "/", dirName)
			//如果db.opt存在，则是mysql的数据库
			isExist, err := util.PathExists(path.Join(dataBaseDir, "db.opt"))
			if err != nil {
				panic("出现异常")
			}
			if isExist {
				currentDataBase, databaseError := NewDataBaseImpl(i, i.conf, dirName)
				if databaseError != nil {
					panic("出现异常")
				}
				i.schemaMap[dirName] = currentDataBase
			}
		}
	}

	//加载数据字典
	content, _ := i.sysTableSpace.LoadPageByPageNumber(7)
	//数据字典页面
	sysDict := NewDataDictWrapperByBytes(content)

	//加载或者初始化系统表，预热
	i.dictionarySys = NewDictionarySysByWrapper(sysDict)
	i.dictionarySys.initDictionary(i.sysTableSpace.(*SysTableSpace))

	sysDb := NewInfoSchemasDB().(*InfoSchemasDB)
	memorySystemTable := NewMemoryInnodbSysTable(i.dictionarySys)
	sysDb.addSystemTable(common.INNODB_SYS_TABLES, memorySystemTable)
	memoryColumnTable := NewMemoryInnodbSysColumns(i.dictionarySys)
	sysDb.addSystemTable(common.INNODB_SYS_COLUMNS, memoryColumnTable)
	i.schemaMap[common.INFORMATION_SCHEMAS] = sysDb

	//遍历各个表
	//tbIterator, _ := i.dictionarySys.SysTable.BTree.Iterate()
	//
	//if tbIterator != nil {
	//
	//	for _, currentRow, _, tbIterator := tbIterator(); tbIterator != nil; _, currentRow, _, tbIterator = tbIterator() {
	//
	//		var keyColumn innodb.Value
	//		//tableId
	//		keyColumn = currentRow.ReadValueByIndex(3)
	//		fmt.Println(keyColumn.Raw())
	//		tableName := string(currentRow.ReadValueByIndex(4).ToByte())
	//		fmt.Println("tableName=" + tableName)
	//		tColumnIterator, _ := i.dictionarySys.SysColumns.BTree.Iterate()
	//		fmt.Println(string(currentRow.ReadValueByIndex(4).ToByte()))
	//		i.processColumns(keyColumn, tColumnIterator)
	//
	//	}
	//}

}

func (i *InfoSchemaManager) processColumns(tableId innodb.Value, tColumnIterator Iterator) {

	if tColumnIterator != nil {
		for _, currentRow, err, tColumnIterator := tColumnIterator(); tColumnIterator != nil; _, currentRow, err, tColumnIterator = tColumnIterator() {
			if err != nil {
				fmt.Println(err)
			}

			tableIdValue := currentRow.ReadValueByIndex(3)
			isEq, equalErrors := tableIdValue.Equal(tableId)
			fmt.Println(equalErrors)
			if isEq.Raw().(bool) {

				fmt.Println(string(currentRow.ReadValueByIndex(4).ToByte()))

			} else {
				fmt.Println("-----迭代不成功---")
			}

		}
		fmt.Println("------------AAAAAAAAAAAAAAAAA----------------------")
	}
}

//初始化加载已经有了的表
func (i *InfoSchemaManager) initLoadTable() {

}

func (i *InfoSchemaManager) GetSchemaByName(schemaName string) (schemas.Database, error) {

	return i.schemaMap[schemaName], nil
}

func (i *InfoSchemaManager) GetSchemaExist(schemaName string) bool {
	if i.schemaMap[schemaName] != nil {
		return true
	}
	return false
}

func (i *InfoSchemaManager) GetTableByName(schema string, tableName string) (schemas.Table, error) {
	database, databaseNotExistError := i.GetSchemaByName(schema)
	if databaseNotExistError != nil {
		return nil, databaseNotExistError
	}
	table, tableNotExistError := database.GetTable(tableName)
	if tableNotExistError != nil {
		return nil, tableNotExistError
	}
	return table, nil
}

func (i *InfoSchemaManager) GetTableExist(schemaName string, tableName string) bool {
	panic("implement me")
}

func (i *InfoSchemaManager) GetAllSchemaNames() []string {
	panic("implement me")
}

func (i *InfoSchemaManager) GetAllSchemas() []schemas.Database {
	panic("implement me")
}

func (i *InfoSchemaManager) GetAllSchemaTablesByName(schemaName string) []schemas.Table {
	panic("implement me")
}

func (i *InfoSchemaManager) PutDatabaseCache(databaseCache schemas.Database) {
	i.schemaMap[databaseCache.Name()] = databaseCache
}
