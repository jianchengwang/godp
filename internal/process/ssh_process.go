package process

import (
	"fmt"
	"godp/internal/pojo"
	"godp/pkg/errorcode"
	"godp/pkg/ssh"
	"net"
	"strings"
)

func GetIntranetIp(projectConfig pojo.ProjectConfig, host string) (error, string) {
	err, ipAddress, sshClient, sftpClient, _ := PreGetInfo(projectConfig, host)
	if err != nil {
		return err, ""
	}
	fmt.Println(ipAddress)
	defer sshClient.Close()
	defer sftpClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return err, ""
	}
	err, combo := ssh.ExecuteCmd(session, "ifconfig eth0 | grep \"inet \" | awk '{print $2}' | cut -c 1-")
	if err != nil {
		return err, ""
	}
	ipv4 := strings.Replace(combo, "\n", "", -1)
	address := net.ParseIP(ipv4)
	if address == nil {
		return errorcode.NewError(500, "ip地址不合法"), ""
	} else {
		return nil, address.String()
	}
}
