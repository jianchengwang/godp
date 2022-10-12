package ssh

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

const (
	sftpWsMsgHeartbeat = "heartbeat"
	sftpWsMsgClose     = "close"
	sftpMsgCmd         = "cmd"
)

type SftpWsMsg struct {
	Type           string      `json:"type"`
	Cmd            string      `json:"cmd"`
	RemoteFilePath string      `json:"remoteFilePath"`
	UpdateString   string      `json:"updateString"`
	Data           interface{} `json:"data"`
}

type LogicSftpWsSession struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
	wsConn     *websocket.Conn
	isAdmin    bool
	IsFlagged  bool `comment:"当前session是否包含禁止命令"`
}

func NewLogicSftpWsSession(isAdmin bool, sshClient *ssh.Client, sftpClient *sftp.Client, wsConn *websocket.Conn) (*LogicSftpWsSession, error) {
	return &LogicSftpWsSession{
		sshClient:  sshClient,
		sftpClient: sftpClient,
		wsConn:     wsConn,
		isAdmin:    isAdmin,
		IsFlagged:  false,
	}, nil
}

//Close 关闭
func (sftpWs *LogicSftpWsSession) Close() {
	if sftpWs.sftpClient != nil {
		sftpWs.sftpClient.Close()
	}
}
func (sftpWs *LogicSftpWsSession) Start(quitChan chan bool) {
	go sftpWs.receiveWsMsg(quitChan)
}

//receiveWsMsg  receive websocket msg do some handling then write into ssh.session.stdin
func (sftpWs *LogicSftpWsSession) receiveWsMsg(exitCh chan bool) {
	wsConn := sftpWs.wsConn
	//tells other go routine quit
	defer setQuit(exitCh)
	for {
		select {
		case <-exitCh:
			return
		default:
			//read websocket msg
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				logrus.WithError(err).Error("reading webSocket message failed")
				return
			}
			//unmashal bytes into struct
			msgObj := SftpWsMsg{}
			if err := json.Unmarshal(wsData, &msgObj); err != nil {
				//logrus.WithError(err).WithField("wsData", string(wsData)).Error("unmarshal websocket message failed")
			}
			switch msgObj.Type {
			case sftpWsMsgHeartbeat:
				continue
			case sftpWsMsgClose:
				sftpWs.Close()
			case sftpMsgCmd:
				//handle xterm.js stdin
				decodeBytes, err := base64.StdEncoding.DecodeString(msgObj.RemoteFilePath)
				if err != nil {
					logrus.WithError(err).Error("websock cmd string base64 decoding failed")
				}
				remoteFilePath := string(decodeBytes)
				cmd := msgObj.Cmd
				switch cmd {
				case "list":
					err, fileList := ListFiles(sftpWs.sftpClient, remoteFilePath)
					if err != nil {
						logrus.WithError(err).Error("websock list " + remoteFilePath + " failed")
						continue
					}
					msgObj.Data = fileList
				case "get":
					err, fetchText := FetchText(sftpWs.sftpClient, remoteFilePath)
					if err != nil {
						logrus.WithError(err).Error("websock get " + remoteFilePath + " failed")
						continue
					}
					msgObj.Data = fetchText
				case "update":
					continue
				}
				if err == nil {
					combo, err := json.Marshal(msgObj)
					if err != nil {
						logrus.WithError(err).Error("sftp json combo output failed")
					}
					err = wsConn.WriteMessage(websocket.TextMessage, combo)
					if err != nil {
						logrus.WithError(err).Error("sftp sending combo output to webSocket failed")
						continue
					}
				}
			}
		}
	}
}

func (sftpWs *LogicSftpWsSession) Wait(quitChan chan bool) {
	if err := sftpWs.sftpClient.Wait(); err != nil {
		logrus.WithError(err).Error("sftp wait failed")
		setQuit1(quitChan)
	}
}

func setQuit1(ch chan bool) {
	ch <- true
}
