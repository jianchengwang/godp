package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"godp/internal/db"
	"godp/internal/pojo"
	"godp/pkg/helper"
	sshHelper "godp/pkg/ssh"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func UseWsSshRouter(r *gin.RouterGroup) {
	r.GET("/ws/ssh/:sessionId", WsSsh)
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handle webSocket connection.
// first,we establish a ssh connection to ssh server when a webSocket comes;
// then we deliver ssh data via ssh connection between browser and ssh server.
// That is, read webSocket data from browser (e.g. 'ls' command) and send data to ssh server via ssh connection;
// the other hand, read returned ssh data from ssh server and write back to browser via webSocket API.
func WsSsh(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if helper.HandleError(c, err) {
		return
	}
	defer wsConn.Close()

	sessionId := c.Param("sessionId")
	fmt.Println(sessionId)
	if helper.WsHandleError(wsConn, err) {
		return
	}
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if helper.WsHandleError(wsConn, err) {
		return
	}
	host := c.DefaultQuery("host", "127.0.0.1")

	err, projectInfo := db.ProjectInfoDb.GetById(uint(projectId))
	if helper.WsHandleError(wsConn, err) {
		return
	}
	var ipAddress pojo.IpAddressStruct
	if host == projectInfo.ProjectConfig.CI.IP {
		ipAddress = projectInfo.ProjectConfig.CI
	} else {
		err, ipAddress = pojo.FindIpAddressByIp(projectInfo.ProjectConfig.IPAddressArr, host)
		if helper.WsHandleError(wsConn, err) {
			return
		}
	}
	port := ipAddress.Port
	user := ipAddress.User
	pass := ipAddress.Password

	err, ssConn := sshHelper.GetConnection(sshHelper.ConnectionInfo{Host: host, Port: port, User: user, Pass: pass})
	if helper.WsHandleError(wsConn, err) {
		log.Println("创建ssh connection 失败", err)
		return
	}
	defer ssConn.Close()

	if strings.HasPrefix(sessionId, "sftp-") {
		err, sftpClient := sshHelper.GetSftpClient(ssConn)
		if helper.WsHandleError(wsConn, err) {
			log.Println("创建sftp connection 失败", err)
			return
		}
		defer sftpClient.Close()

		sftpws, err := sshHelper.NewLogicSftpWsSession(true, ssConn, sftpClient, wsConn)
		if helper.WsHandleError(wsConn, err) {
			return
		}
		defer sftpws.Close()

		quitChan := make(chan bool, 3)
		sftpws.Start(quitChan)
		go sftpws.Wait(quitChan)
		<-quitChan
	} else {
		sws, err := sshHelper.NewLogicSshWsSession(80, 40, true, ssConn, wsConn)
		if helper.WsHandleError(wsConn, err) {
			return
		}
		defer sws.Close()

		quitChan := make(chan bool, 3)
		sws.Start(quitChan)
		go sws.Wait(quitChan)
		<-quitChan
	}
}
