package ssh

import (
	"bufio"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

type SftpFileList struct {
	Name    string `json:"name"`
	ModTime string `json:"modTime"`
	Size    string `json:"size"`
}

type ConnectionInfo struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func Test(host string, port uint16, user string, pass string) {
	err, conn := GetConnection(ConnectionInfo{host, port, user, pass})
	if err != nil {

	}
	defer conn.Close()
	err, sc := GetSftpClient(conn)
	defer sc.Close()
	ListFiles(sc, "docker-compose-sxl-test")
	fmt.Fprintf(os.Stdout, "\n")

	session, err := conn.NewSession()
	defer session.Close()
	cmd := "whoami; cd /; ls -al;echo hello"
	ExecuteCmd(session, cmd)

	//UploadFile(*sc, "./local.txt", "./remote.txt")
	//fmt.Fprintf(os.Stdout, "\n")
	//
	//DownloadFile(*sc, "./remote.txt", "./download.txt")
	//fmt.Fprintf(os.Stdout, "\n")
}

func GetConnection(ci ConnectionInfo) (error, *ssh.Client) {
	//hostKey := getHostKey(host)
	host := ci.Host
	port := ci.Port
	user := ci.User
	pass := ci.Pass

	fmt.Fprintf(os.Stdout, "Connecting to %s ...\n", host)

	var auths []ssh.AuthMethod

	// Try to use $SSH_AUTH_SOCK which contains the path of the unix file socket that the sshd agent uses
	// for communication with other processes.
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}

	// Use password authentication if provided
	if pass != "" {
		auths = append(auths, ssh.Password(pass))
	}

	// Initialize client configuration
	config := ssh.ClientConfig{
		User: user,
		Auth: auths,
		// Uncomment to ignore host key check
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	// Connect to server
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connecto to [%s]: %v\n", addr, err)
	}
	return err, conn
}

func GetSftpClient(conn *ssh.Client) (error, *sftp.Client) {
	// Create new SFTP client
	sc, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start SFTP subsystem: %v\n", err)
	}
	return err, sc
}

// ExecuteCmd 执行远程命令
func ExecuteCmd(session *ssh.Session, cmd string) (error, string) {
	//执行远程命令
	combo, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Println("远程执行cmd 失败", err)
	}
	log.Println("命令输出:", string(combo))
	return err, string(combo)
}

// ListFiles 列出目录文件列表
func ListFiles(sc *sftp.Client, remoteDir string) (error, []SftpFileList) {

	var fileList = make([]SftpFileList, 0)

	fmt.Fprintf(os.Stdout, "Listing [%s] ...\n\n", remoteDir)

	files, err := sc.ReadDir(remoteDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to list remote dir: %v\n", err)
		return err, nil
	}

	for _, f := range files {
		var name, modTime, size string

		name = f.Name()
		modTime = f.ModTime().Format("2006-01-02 15:04:05")
		size = fmt.Sprintf("%12d", f.Size())

		if f.IsDir() {
			name = name + "/"
			modTime = ""
			size = "PRE"
		}
		// Output each file name and size in bytes
		//fmt.Fprintf(os.Stdout, "%19s %12s %s\n", modTime, size, name)
		fileList = append(fileList, SftpFileList{name, modTime, size})
	}

	return err, fileList
}

func UploadDirectory(sshClient *ssh.Client, sftpClient *sftp.Client, localPath string, remotePath string) error {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal("read dir list fail ", err)
		return err
	}

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			sftpClient.Mkdir(remoteFilePath)
			err = UploadDirectory(sshClient, sftpClient, localFilePath, remoteFilePath)
			if err != nil {
				return err
			}
		} else {
			err = UploadFile(sshClient, sftpClient, path.Join(localPath, backupDir.Name()), remotePath)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println(localPath + "  copy directory to remote server finished!")
	return nil
}

func UploadFile(sshClient *ssh.Client, sftpClient *sftp.Client, localFilePath string, remotePath string) error {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		return err
	}
	defer srcFile.Close()
	sftpClient.MkdirAll(remotePath)
	remoteFileName := filepath.Base(localFilePath)
	remotePath = filepath.Join(remotePath, remoteFileName)
	dstFile, err := sftpClient.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		fmt.Println("sftpClient.Create error : ", remotePath)
		return err
	}
	defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		return err
	}
	dstFile.Write(ff)
	fmt.Println(localFilePath + " copy file to remote server finished!")
	// 如果是脚本文件，进行chmod+x dos2unix
	if strings.HasSuffix(remoteFileName, ".sh") {
		session, err := sshClient.NewSession()
		if err != nil {
			return err
		}
		ExecuteCmd(session, "yum -y install dos2unix;chmod +x "+remotePath+";dos2unix "+remotePath+";")
	}
	return nil
}

// DownloadFile 下载文件
func DownloadFile(sshClient *sftp.Client, remoteFile, localFile string) (err error) {

	fmt.Fprintf(os.Stdout, "Downloading [%s] to [%s] ...\n", remoteFile, localFile)
	// Note: SFTP To Go doesn't support O_RDWR mode
	srcFile, err := sshClient.OpenFile(remoteFile, os.O_RDONLY)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open remote file: %v\n", err)
		return
	}
	defer srcFile.Close()
	err = os.MkdirAll(filepath.Dir(localFile), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to mkdir local file: %v\n", err)
		return
	}
	dstFile, err := os.Create(localFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open local file: %v\n", err)
		return
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to download remote file: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%d bytes copied\n", bytes)

	return
}

// FetchText 获取文本信息
func FetchText(sftpClient *sftp.Client, remoteFile string) (error, string) {

	fmt.Fprintf(os.Stdout, "FetchText [%s] ...\n", remoteFile)
	// Note: SFTP To Go doesn't support O_RDWR mode
	srcFile, err := sftpClient.OpenFile(remoteFile, os.O_RDONLY)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open remote file: %v\n", err)
		return err, ""
	}
	defer srcFile.Close()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to stat remote file: %v\n", err)
		return err, ""
	}
	buffer := make([]byte, fileInfo.Size())
	for {
		_, err := srcFile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("已读取到文件末尾")
				break
			} else {
				fmt.Println("读取文件出错", err)
				return err, ""
			}
		}
	}

	return nil, string(buffer)
}

func UpdateText(sshClient *ssh.Client, sftpClient *sftp.Client, remoteFilePath string, updateString string) error {
	paths, remoteFileName := path.Split(remoteFilePath)
	sftpClient.MkdirAll(paths)
	dstFile, err := sftpClient.OpenFile(remoteFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		fmt.Println("sftpClient.Create error : ", remoteFilePath)
		return err
	}
	defer dstFile.Close()
	ff := []byte(updateString)
	dstFile.Write(ff)
	fmt.Println(" update text to remote server finished!")
	// 如果是脚本文件，进行chmod+x dos2unix
	if strings.HasSuffix(remoteFileName, ".sh") {
		session, err := sshClient.NewSession()
		if err != nil {
			return err
		}
		ExecuteCmd(session, "yum -y install dos2unix;chmod +x "+remoteFilePath+";dos2unix "+remoteFilePath+";")
	}
	return nil
}

// Get host key from local known hosts
func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	fmt.Println("HOME:" + os.Getenv("HOME"))
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Home dir:", u.HomeDir)
	file, err := os.Open(filepath.Join(u.HomeDir, ".ssh", "known_hosts"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read known_hosts file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing %q: %v\n", fields[2], err)
				os.Exit(1)
			}
			break
		}
	}

	if hostKey == nil {
		fmt.Fprintf(os.Stderr, "No hostkey found for %s", host)
		os.Exit(1)
	}

	return hostKey
}
