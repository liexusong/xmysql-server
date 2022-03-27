package handler

//
//import (
//	getty "github.com/AlexStocks/getty/transport"
//	"github.com/zhukovaskychina/xmysql-server/server/net"
//)
//
////
////func ProcessMySQLPacketFromClient(buff []byte, conn xnet.WriteCloser, singleEngine *engine.XEngine) {
////	packetType := buff[4]
////	switch packetType {
////	case common.COM_SLEEP:
////		{
////			break
////		}
////
////	case common.COM_QUERY:
////		{
////			bytes := buff[5:]
////			sql := string(bytes)
////			fmt.Println(sql)
////			stmt, err := sqlparser.Parse(sql)
////
////			if err != nil {
////				fmt.Println(err)
////			}
////			switch stmt := stmt.(type) {
////			case *sqlparser.Select:
////				{
////					singleEngine.SelectResultsWithLock(conn, conn.GetConnectDBName(), sql, stmt)
////					break
////				}
////			case *sqlparser.Show:
////				{
////
////				}
////			case *sqlparser.Insert:
////				{
////					singleEngine.InsertTableExecutorWithLock(conn, conn.GetConnectDBName(), stmt)
////					break
////				}
////
////			case *sqlparser.Set:
////				{
////					setExprs := stmt.Exprs
////					fmt.Println(setExprs)
////					conn.WriteOk()
////					break
////				}
////			case *sqlparser.Update:
////				{
////
////				}
////			case *sqlparser.DBDDL:
////				{
////					//	stmt.Action
////					action := stmt.Action
////					switch action {
////					case sqlparser.CreateStr:
////						singleEngine.CreateDBExecutorWithoutLock(conn, stmt)
////					case sqlparser.DropStr:
////
////					}
////				}
////			case *sqlparser.DDL:
////				{
////					action := stmt.Action
////					switch action {
////					case sqlparser.CreateStr:
////						singleEngine.CreateFrmExecutorWithNoLock(conn, conn.GetConnectDBName(), stmt)
////					case sqlparser.DropStr:
////					case sqlparser.RenameStr:
////					case sqlparser.AlterStr:
////					case sqlparser.CreateVindexStr:
////					case sqlparser.AddColVindexStr:
////					case sqlparser.DropColVindexStr:
////					default:
////					}
////				}
////
////			case *sqlparser.Use:
////				{
////					dbName := stmt.DBName.String()
////					if dbName != "" {
////						conn.SetDataBaseName(dbName)
////						conn.WriteOk()
////					}
////					conn.WriteError("Database Name 为空")
////				}
////			}
////			break
////		}
////	case common.COM_QUIT:
////		{
////			conn.Close()
////		}
////	}
////
////}
//
////定义事件处理handler
//type DispatchHandler struct {
//}
//
//func (h *DispatchHandler) Handle(session getty.Session, pkg *net.MySQLPackage) error {
//	//log.Debug("get echo heartbeat package{%s}", pkg)
//	//
//	//var rspPkg EchoPackage
//	//rspPkg.H = pkg.H
//	//rspPkg.B = echoHeartbeatResponseString
//	//rspPkg.H.Len = uint16(len(rspPkg.B) + 1)
//
//	return session.WritePkg(nil, 0)
//}
