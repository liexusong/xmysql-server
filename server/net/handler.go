/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package net

import (
	"errors"
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server"
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/engine"
	"github.com/zhukovaskychina/xmysql-server/server/protocol"
	"sync"
	//"xmysql-serverimpl/serverimpl/net/service"
)

import (
	log "github.com/AlexStocks/log4go"
)

const (
	WritePkgTimeout = 1e8
)

var (
	errTooManySessions = errors.New("Too many MySQL sessions!")
)
var (
	ErrNotEnoughStream = errors.New("packet stream is not enough")
	ErrTooLargePackage = errors.New("package length is exceed the echo package's legal maximum length.")
	ErrIllegalMagic    = errors.New("package magic is not right.")
)

type PackageHandler interface {
	Handle(Session, *MySQLPackage) error
}

type MySQLPackageHandler struct {
}

func (h *MySQLPackageHandler) Handle(session Session, pkg *MySQLPackage) error {
	log.Debug("get echo package{%s}", pkg)
	// write echo message handle logic here.
	return session.WritePkg(nil, WritePkgTimeout)
}

type MySQLMessageHandler struct {
	rwlock       sync.RWMutex
	cfg          *conf.Cfg
	sessionMap   map[Session]server.MySQLServerSession //内存区，用于存储mysql的session
	XMySQLEngine *engine.XMySQLEngine
}

func NewMySQLMessageHandler(cfg *conf.Cfg) *MySQLMessageHandler {
	var mySQLMessageHandler = new(MySQLMessageHandler)
	mySQLMessageHandler.sessionMap = make(map[Session]server.MySQLServerSession)
	mySQLMessageHandler.cfg = cfg
	mySQLMessageHandler.XMySQLEngine = engine.NewXMySQLEngine(cfg)
	return mySQLMessageHandler
}

func (m *MySQLMessageHandler) OnOpen(session Session) error {
	var (
		err error
	)

	m.rwlock.RLock()

	if m.cfg.SessionNumber <= len(m.sessionMap) {
		err = errTooManySessions
	}
	m.rwlock.RUnlock()
	if err != nil {
		return err
	}
	log.Info("got session:%s", session.Stat())
	m.rwlock.Lock()

	m.sessionMap[session] = NewMySQLServerSession(session)
	m.rwlock.Unlock()
	//主动与客户端握手
	m.sessionMap[session].SendHandleOk()
	return nil
}

func (m *MySQLMessageHandler) OnClose(session Session) {
	session.Close()
	delete(m.sessionMap, session)
}

func (m *MySQLMessageHandler) OnError(session Session, err error) {
	fmt.Println("", err)
	session.Close()
	delete(m.sessionMap, session)
}

func (m *MySQLMessageHandler) OnCron(session Session) {
	fmt.Println("session 检查")
}

func (m *MySQLMessageHandler) OnMessage(session Session, pkg interface{}) {
	//	req, res := wrapper.NewServerRequest(m.sessionMap[session])
	//	req, res := wrapper.NewServerRequest(nil)
	currentMysqlSession := m.sessionMap[session]
	recMySQLPkg := pkg.(*MySQLPackage)
	authStatus := session.GetAttribute("auth_status")
	if authStatus == nil {
		a := new(protocol.AuthPacket)
		var authData = make([]byte, 0)
		authData = append(authData, recMySQLPkg.Header.PacketLength...)
		authData = append(authData, recMySQLPkg.Header.PacketId)
		authData = append(authData, recMySQLPkg.Body...)
		a.DecodeAuth(authData)
		fmt.Println(a.Database)
		fmt.Println(a.User)
		fmt.Println(a.Password)
		session.SetAttribute("auth_status", "success")
		currentMysqlSession.SetParamByName("database", a.Database)
		//	req.SetParamByName("currentDatabase", a.Database)
		m.sessionMap[session] = currentMysqlSession
		buff := make([]byte, 0)

		session.WriteBytes(protocol.EncodeOK(buff, 0, 0, nil))
		return
	}
	packetType := recMySQLPkg.Body[0]
	switch packetType {
	case common.COM_SLEEP:
		{
			break
		}

	case common.COM_QUERY:
		{

			sql := string(recMySQLPkg.Body[1:])
			dbName := currentMysqlSession.GetParamByName("database").(string)
			fmt.Println(dbName)
			m.XMySQLEngine.ExecuteQuery(currentMysqlSession, sql, dbName)
		}
	case common.COM_QUIT:
		{
			fmt.Println("")

		}

	}

}
