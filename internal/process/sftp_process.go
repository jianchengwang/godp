package process

import (
	"fmt"
	"github.com/pkg/sftp"
	"godp/internal/config"
	"godp/internal/pojo"
	"godp/pkg/errorcode"
	sshHelper "godp/pkg/ssh"
	"golang.org/x/crypto/ssh"
)

func SftpProcessUploadDevelopFile(projectConfig pojo.ProjectConfig, host string, filename string) error {
	err, ipAddress, sshClient, sftpClient, remoteWorkPath := PreGetInfo(projectConfig, host)
	if err != nil {
		return err
	}
	fmt.Println(ipAddress)
	defer sshClient.Close()
	defer sftpClient.Close()

	localFilePath := config.Config.App.WorkPath + "/" + projectConfig.ProjectApp + "/"
	remoteFilePath := remoteWorkPath
	localFilePath += filename
	remoteFilePath += filename

	err = sshHelper.UploadDirectory(sshClient, sftpClient, localFilePath, remoteFilePath)
	return err
}

func SftpFetchText(projectConfig pojo.ProjectConfig, host string, filename string, remoteFilePath string) (error, string) {
	err, ipAddress, sshClient, sftpClient, remoteWorkPath := PreGetInfo(projectConfig, host)
	if err != nil {
		return err, ""
	}
	fmt.Println(ipAddress)
	defer sshClient.Close()
	defer sftpClient.Close()

	if remoteFilePath == "" {
		remoteFilePath = getRemoteFilePath(remoteWorkPath, filename)
	}

	if remoteFilePath == "" {
		return errorcode.NewError(400, "参数错误"), ""
	}

	return sshHelper.FetchText(sftpClient, remoteFilePath)
}

func SftpUpdateText(projectConfig pojo.ProjectConfig, host string, filename string, remoteFilePath string, updateString string) error {
	err, ipAddress, sshClient, sftpClient, remoteWorkPath := PreGetInfo(projectConfig, host)
	if err != nil {
		return err
	}
	fmt.Println(ipAddress)
	defer sshClient.Close()
	defer sftpClient.Close()

	if remoteFilePath == "" {
		remoteFilePath = getRemoteFilePath(remoteWorkPath, filename)
	}

	if remoteFilePath == "" {
		return errorcode.NewError(400, "参数错误")
	}

	return sshHelper.UpdateText(sshClient, sftpClient, remoteFilePath, updateString)
}

func getRemoteFilePath(remoteWorkPath, filename string) string {
	switch filename {
	case "install.sh":
		{
			return remoteWorkPath + "docker-install/install.sh"
		}
	case "gitInit.sh":
		{
			return remoteWorkPath + "build/gitInit.sh"
		}
	case "build.sh":
		{
			return remoteWorkPath + "build/build.sh"
		}
	case "nginx.admin.conf":
		{
			return remoteWorkPath + "docker-compose-server/config/nginx/conf.d/admin.conf"
		}
	case "nginx.client.conf":
		{
			return remoteWorkPath + "docker-compose-server/config/nginx/conf.d/client.conf"
		}
	case "nginx.h5.conf":
		{
			return remoteWorkPath + "docker-compose-server/config/nginx/conf.d/h5.conf"
		}
	case "docker-compose-server.yml":
		{
			return remoteWorkPath + "docker-compose-server/docker-compose-server.yml"
		}
	case "docker-compose-admin.yml":
		{
			return remoteWorkPath + "docker-compose-app/docker-compose-admin.yml"
		}
	case "docker-compose-client.yml":
		{
			return remoteWorkPath + "docker-compose-app/docker-compose-client.yml"
		}
	case "admin.application-prod.yml":
		{
			return remoteWorkPath + "docker-compose-app/app/admin/config/application-prod.yml"
		}
	case "client.application-prod.yml":
		{
			return remoteWorkPath + "docker-compose-app/app/client/config/application-prod.yml"
		}
	}
	return ""
}

func PreGetInfo(projectConfig pojo.ProjectConfig, host string) (error, pojo.IpAddressStruct, *ssh.Client, *sftp.Client, string) {
	var err error
	var ipAddress pojo.IpAddressStruct
	if host == projectConfig.CI.IP {
		ipAddress = projectConfig.CI
	} else {
		err, ipAddress = pojo.FindIpAddressByIp(projectConfig.IPAddressArr, host)
		if err != nil {
			return err, pojo.IpAddressStruct{}, &ssh.Client{}, &sftp.Client{}, ""
		}
	}
	connectionInfo := sshHelper.ConnectionInfo{
		Host: host,
		Port: ipAddress.Port,
		User: ipAddress.User,
		Pass: ipAddress.Password,
	}
	err, sshClient := sshHelper.GetConnection(connectionInfo)
	if err != nil {
		return err, pojo.IpAddressStruct{}, &ssh.Client{}, &sftp.Client{}, ""
	}

	err, sftpClient := sshHelper.GetSftpClient(sshClient)
	if err != nil {
		sshClient.Close()
		return err, pojo.IpAddressStruct{}, &ssh.Client{}, &sftp.Client{}, ""
	}

	deployDir := "/root"
	if ipAddress.DeployDir != "" {
		deployDir = ipAddress.DeployDir
	}
	remoteWorkPath := deployDir + "/" + projectConfig.ProjectApp + "/"

	return nil, ipAddress, sshClient, sftpClient, remoteWorkPath
}
