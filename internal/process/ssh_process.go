package process

import (
	"godp/internal/db"
	"godp/pkg/errorcode"
	sshHelper "godp/pkg/ssh"
	"net"
	"strings"
)

func GetIntranetIp(assetsHost db.AssetsHost) (error, string) {

	connectionInfo := sshHelper.ConnectionInfo{
		Host: assetsHost.IP,
		Port: assetsHost.Port,
		User: assetsHost.User,
		Pass: assetsHost.Password,
	}
	err, sshClient := sshHelper.GetConnection(connectionInfo)
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return err, ""
	}
	err, combo := sshHelper.ExecuteCmd(session, "ifconfig eth0 | grep \"inet \" | awk '{print $2}' | cut -c 1-")
	if err != nil {
		return err, ""
	}
	ipv4 := strings.Replace(combo, "\n", "", -1)
	address := net.ParseIP(ipv4)
	if address == nil {
		return errorcode.NewError(500, "ip address error"), ""
	} else {
		return nil, address.String()
	}
}
