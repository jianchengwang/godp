package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"godp/internal/config"
	"godp/internal/db"
	"godp/pkg/api/helper"
	"godp/pkg/file"
	sshHelper "godp/pkg/ssh"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var SessionSftpWsMap = make(map[string]*sshHelper.LogicSftpWsSession)

func UseWsSshRouter(r *gin.RouterGroup) {
	r.GET("/ws/ssh/:sessionId", WsSsh)
	r.POST("/ws/ssh/:sessionId/uploadFile", WsSshUploadFile)
	r.GET("/ws/ssh/:sessionId/downloadFile", WsSshDownloadFile)
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
	host := c.Query("host")
	err, assetsHost := db.AssetsHostDb.GetByIP(host, 0)
	if helper.WsHandleError(wsConn, err) {
		return
	}
	port := assetsHost.Port
	user := assetsHost.User
	pass := assetsHost.Password

	err, ssConn := sshHelper.GetConnection(sshHelper.ConnectionInfo{Host: host, Port: port, User: user, Pass: pass})
	if helper.WsHandleError(wsConn, err) {
		log.Println("create ssh connection failed", err)
		return
	}
	defer ssConn.Close()

	if strings.HasPrefix(sessionId, "sftp-") {
		err, sftpClient := sshHelper.GetSftpClient(ssConn)
		if helper.WsHandleError(wsConn, err) {
			log.Println("create sftp connection failed", err)
			return
		}
		defer sftpClient.Close()

		sftpws, err := sshHelper.NewLogicSftpWsSession(true, ssConn, sftpClient, wsConn)
		if helper.WsHandleError(wsConn, err) {
			return
		}
		defer sftpws.Close()

		SessionSftpWsMap[sessionId] = sftpws

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

func WsSshUploadFile(c *gin.Context) {
	sessionId := c.Param("sessionId")
	remoteFilePathEncoding, _ := c.GetPostForm("remoteFilePath")
	if sftpws, ok := SessionSftpWsMap[sessionId]; ok {
		decodeBytes, err := base64.StdEncoding.DecodeString(remoteFilePathEncoding)
		if err != nil {
			helper.HandleError(c, err)
		}
		remoteFilePath := string(decodeBytes)
		uploadFile, err := c.FormFile("file")
		localFilePath := filepath.Join(config.Config.App.TempPath, sessionId, remoteFilePath, uploadFile.Filename)
		os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm)
		err = c.SaveUploadedFile(uploadFile, localFilePath)
		if helper.HandleError(c, err) {
			return
		}
		// 上传到目标服务器
		err = sshHelper.UploadFile(sftpws.SshClient, sftpws.SftpClient, localFilePath, remoteFilePath)
		if helper.HandleError(c, err) {
			return
		}
	} else {
		helper.HandleError(c, errors.New("sessionId not found"))
	}

}

func WsSshDownloadFile(c *gin.Context) {
	sessionId := c.Param("sessionId")
	remoteFilePathEncoding := c.Query("remoteFilePath")
	if sftpws, ok := SessionSftpWsMap[sessionId]; ok {
		decodeBytes, err := base64.StdEncoding.DecodeString(remoteFilePathEncoding)
		if err != nil {
			helper.HandleError(c, err)
		}
		remoteFilePath := string(decodeBytes)
		localFilePath := filepath.Join(config.Config.App.TempPath, sessionId, remoteFilePath)
		sshHelper.DownloadFile(sftpws.SftpClient, remoteFilePath, localFilePath)
		if _, ok := file.IsDir(localFilePath); ok {
			file.Zip(localFilePath, localFilePath+".zip")
			localFilePath = localFilePath + ".zip"
		}
		_, filename := filepath.Split(localFilePath)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(localFilePath)
	} else {
		helper.HandleError(c, errors.New("sessionId not found"))
	}
}
