package godo

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/thesunnysky/godo/consts"
	"github.com/thesunnysky/godo/server"
	"github.com/thesunnysky/godo/util"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dataFile string
//var r, _ = regexp.Compile("[[:alnum:]]")
var r, _ = regexp.Compile("^\\s*$")
var he *util.HyperEncrypt

func init() {
	dataFile = ClientConfig.DataFile
	initHyperEncrypt()
}

type GitRepo struct {
	repoPath string
}

func initHyperEncrypt() {
	var err error
	he, err = util.NewHyperEncryptF(ClientConfig.RsaConfig.PublicKeyFile,
		ClientConfig.RsaConfig.PrivateKeyFile)
	if err != nil {
		fmt.Printf("new hyper encrypt tool error:%s\n", err)
		os.Exit(-1)
	}
}

func AddCmdImpl(args []string) {
	var buf bytes.Buffer
	for _, str := range args {
		buf.WriteString(str)
		buf.WriteByte(' ')
	}

	f, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, consts.FileMask)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file := util.File{File: f}
	file.AppendNewLine(buf.Bytes())

	fmt.Println("task add successfully")
}

func DelCmdImpl(args []string) {
	num := make([]int, len(args))
	for _, str := range args {
		i, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("invalid parameter value:%s\n", str)
			os.Exit(consts.InvalidParameterValue)
		} else {
			num = append(num, i)
		}
	}

	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDWR, consts.FileMask)
	defer f.Close()
	file := util.File{File: f}
	if err != nil {
		panic(err)
	}

	fileData := file.ReadFile()
	for _, index := range num {
		idx := index - 1
		if (idx < 0) || (idx > len(fileData)-1) {
			continue
		}
		fileData[idx] = string('\n')
	}

	file.RewriteFile(fileData)

	fmt.Println("delete task successfully")
}

func ListLocalTasks(args []string) {
	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDONLY, consts.FileMask)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	printTask(f)
}

func BackupTaskFile() {
	srcFile := ClientConfig.DataFile
	backupFile := ClientConfig.DataFile + "." + consts.BackupTaskFileSuffix
	if err := util.CopyFile(backupFile, srcFile); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf("backup to %s\n", util.ExtractFileName(backupFile))
}

func DelBackupTaskFile(args []string) {
	var delAll bool
	if (len(args) == 1) && (strings.Compare(strings.ToLower(args[0]), consts.DelAllBackupFileArg) == 0) {
		delAll = true
	}

	argSlc := make([]interface{}, len(args))
	for i, v := range args {
		argSlc[i] = v
	}
	delSet := mapset.NewSetFromSlice(argSlc)

	dataFilePath := ClientConfig.DataFile
	dataFileName := util.ExtractFileName(dataFilePath)
	dataFileDir, err := util.ExtractFileDir(dataFilePath)
	if err != nil {
		fmt.Printf("get task file dir error:%s\n", err)
		os.Exit(-1)
	}

	fileInfos, err := ioutil.ReadDir(dataFileDir)
	if err != nil {
		fmt.Printf("read task file dir err:%s\n", err)
		os.Exit(-1)
	}

	//skip "godo.dat" and "godo.dat.0"
	for _, fileInfo := range fileInfos {
		backupFileName := fileInfo.Name()
		if !fileInfo.Mode().IsRegular() || !strings.HasPrefix(backupFileName, dataFileName) ||
			strings.Compare(backupFileName, dataFileName) == 0 ||
			strings.Compare(backupFileName, dataFileName+"."+consts.BackupTaskFileSuffix) == 0 {
			continue
		}

		suffixPos := strings.LastIndex(backupFileName, ".")
		backFileNoStr := backupFileName[suffixPos+1:]
		if _, err := strconv.Atoi(backFileNoStr); err != nil {
			continue
		}

		if delAll || delSet.Contains(backFileNoStr) {
			if err := os.Remove(dataFileDir + "/" + backupFileName); err != nil {
				fmt.Println(err)
			}
			fmt.Printf("delete file:%s\n", backupFileName)
		}
	}
}

func RecoverTaskFile() {
	dstFile := ClientConfig.DataFile
	backUpFile := ClientConfig.DataFile + "." + consts.BackupTaskFileSuffix
	if err := util.CopyFile(dstFile, backUpFile); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf("recover tasks from %s\n", util.ExtractFileName(backUpFile))
}

func ListBackupFiles(args []string) {
	filepath := ClientConfig.DataFile
	filename := util.ExtractFileName(filepath)
	fileDir, err := util.ExtractFileDir(filepath)
	if err != nil {
		fmt.Printf("get task file dir error:%s\n", err)
		os.Exit(-1)
	}

	fileInfos, err := ioutil.ReadDir(fileDir)
	if err != nil {
		fmt.Printf("read task file dir err:%s\n", err)
		os.Exit(-1)
	}

	for _, fileInfo := range fileInfos {
		backupFileName := fileInfo.Name()
		if !fileInfo.Mode().IsRegular() || !strings.HasPrefix(backupFileName, filename) {
			continue
		}

		suffixPos := strings.LastIndex(backupFileName, ".")
		backFileNo, err := strconv.Atoi(backupFileName[suffixPos+1:])
		if err != nil {
			continue
		}
		fmt.Printf("%d. %s\n", backFileNo, backupFileName)
	}
}

func ListBackupTasks(backupFileNo int) {
	backupFile := ClientConfig.DataFile + "." + strconv.Itoa(backupFileNo)
	f, err := os.OpenFile(backupFile, os.O_RDONLY, consts.FileMask)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer f.Close()

	printTask(f)
}

func ListRemoteTasks(args []string) {
	remoteTasks, err := decryptRemoteTasks()
	if err != nil {
		fmt.Printf("decrypt remote tasks error:%s\n", err)
		os.Exit(-1)
	}

	printTask(bytes.NewReader(remoteTasks))
}

func decryptRemoteTasks() ([]byte, error) {
	remoteTaskFileName := util.ExtractFileName(ClientConfig.DataFile)
	remoteTaskFilePath := ClientConfig.GithubRepo + "/" + remoteTaskFileName

	remoteTaskCipherData, err := ioutil.ReadFile(remoteTaskFilePath)
	if err != nil {
		fmt.Printf("read file error:%s\n", remoteTaskFilePath)
		os.Exit(-1)
	}

	remoteTaskData, err := he.Decrypt(remoteTaskCipherData)
	if err != nil {
		fmt.Printf("decrypte file error:%s\n", remoteTaskFilePath)
		return nil, err
	}
	return remoteTaskData, nil
}

func TidyCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	file := util.File{File: f}
	defer file.Close()

	//read and filter task file
	br := bufio.NewReader(file.File)
	var fileData []string
	for {
		str, err := br.ReadString(consts.LineSeparator)
		if err == io.EOF {
			break
		}
		if !isBlankLine(str) {
			//remove empty line
			fileData = append(fileData, str)
		}
	}

	//rewrite task myfile
	file.RewriteFile(fileData)

	fmt.Println("tidy task file successfully")
}

func isBlankLine(str string) bool {
	return r.MatchString(str)
}

func PullServerCmd(args []string) {
	filename := ClientConfig.DataFile
	index := strings.LastIndex(filename, "/")
	if index == -1 {
		index = 0
	}
	fileName := filename[index:]

	apiClient := server.ApiClient{Url: ClientConfig.GodoServerUrl}
	reader, err := apiClient.PullFile(fileName)
	if err != nil {
		fmt.Printf("download file:%s error:%s\n", fileName, err)
		os.Exit(-1)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("read data from response error:%s\n", err)
	}
	aesUtil := util.Aes{Key: []byte(ClientConfig.AesGCMConfig.Key),
		Nonce: []byte(ClientConfig.AesGCMConfig.Nonce)}
	decryptedData, err := aesUtil.GcmDecrypt(data)
	if err != nil {
		fmt.Printf("decrypt data error:%s\n", err)
	}

	//write download data to tempFile
	tempFile := ClientConfig.DataFile + ".tmp"
	targetFile := ClientConfig.DataFile
	f, err := os.OpenFile(tempFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, consts.FileMask)
	if err != nil {
		fmt.Printf("open file error:%s\n", err)
	}
	defer f.Close()
	if _, err := f.Write(decryptedData); err != nil {
		fmt.Printf("open file error:%s\n", err)
		os.Exit(-1)
	}

	// rename tempFile to targetFile
	if util.PathExist(targetFile) {
		backupFile := targetFile + ".bak"
		if err := os.Rename(targetFile, backupFile); err != nil {
			fmt.Printf("remove old task file error%s\n", err)
			os.Exit(-1)
		}
	}

	// rename tempFile to targetFile
	if err := os.Rename(tempFile, targetFile); err != nil {
		fmt.Printf("remove old task file error%s\n", err)
		os.Exit(-1)
	}
	fmt.Println("pull task file successfully")
}

func PushServerCmd(args []string) {

	apiClient := server.ApiClient{Url: ClientConfig.GodoServerUrl,
		Key:   ClientConfig.AesGCMConfig.Key,
		Nonce: ClientConfig.AesGCMConfig.Nonce}
	if err := apiClient.PushFile(consts.GODO_DATA_FILE, ClientConfig.DataFile);
		err != nil {
		fmt.Printf("push task file to server error:%s\n", err)
		os.Exit(-1)
	}

	fmt.Println("pull task file successfully")
}

func PushGitCmd(args []string) {
	textData, err := ioutil.ReadFile(dataFile)
	if err != nil {
		fmt.Printf("read godo data file error:%s\n", err)
		os.Exit(-1)
	}

	cipherData, err := he.Encrypt(textData)
	if err != nil {
		fmt.Printf("encryt task file error:%s\n", err)
		os.Exit(-1)
	}

	gitFileName := util.ExtractFileName(ClientConfig.DataFile)
	gitFile := ClientConfig.GithubRepo + "/" + gitFileName
	if err := ioutil.WriteFile(gitFile, cipherData, consts.FileMask);
		err != nil {
		fmt.Printf("write git file error:%s\n", gitFile)
		os.Exit(-1)
	}

	msg := time.Now().Format("2006-01-02 15:04:05")
	gitCmArgs := []string{"commit", "-am", "\"" + msg + "\""}
	if err := g.GitCmd(gitCmArgs); err != nil {
		fmt.Printf("git commit -am error:%s\n", err)
		os.Exit(-1)
	}

	gitPushArgs := []string{"push"}
	if err := g.GitCmd(gitPushArgs); err != nil {
		fmt.Printf("git push error:%s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("goto push sucessfully")
}

func (r *GitRepo) PullGitCmd(args []string) {
	gitCmArgs := []string{"pull"}
	if err := g.GitCmd(gitCmArgs); err != nil {
		fmt.Printf("git pull error:%s\n", err)
		os.Exit(-1)
	}
}

func (r *GitRepo) GitCmd(args []string) (err error) {
	cmd := exec.Command("git", args...)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Dir = r.repoPath

	if err = cmd.Start(); err != nil {
		fmt.Printf("git command start error:%s\n", err)
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	timeout := 30 * time.Second
	select {
	case <-time.After(timeout):
		if cmd.Process != nil && cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
			if err = cmd.Process.Kill(); err != nil {
				fmt.Printf("command process kill error:%s\n", err)
				os.Exit(-1)
			}
		}
		<-done
		fmt.Println("command execute timeout")

	case err = <-done:
	}

	if err != nil {
		fmt.Printf("[command execute error]: %s\n", err)
	}

	if stderr.Len() > 0 {
		fmt.Print(string(stderr.Bytes()))
	}

	if stdout.Len() > 0 {
		fmt.Print(string(stdout.Bytes()))
	}

	return err
}

func UpdateCmd(args []string) {
	remoteTasks, err := decryptRemoteTasks()
	if err != nil {
		fmt.Printf("decryt remote tasks error:%s\n", remoteTasks)
		os.Exit(-1)
	}

	if err := util.BackupFile(ClientConfig.DataFile); err != nil {
		fmt.Printf("backup file %s error\n", ClientConfig.DataFile)
	}

	if err := util.RewriteFile(ClientConfig.DataFile, remoteTasks); err != nil {
		fmt.Printf("rewrite local task file error:%s\n", err)
		os.Exit(-1)
	}

	fmt.Println("update local task file successfully")
}

func printTask(r io.Reader) {
	br := bufio.NewReader(r)
	var index int
	for {
		str, err := br.ReadString(consts.LineSeparator)
		if err == io.EOF {
			break
		}
		index++
		if !isBlankLine(str) {
			fmt.Printf("%d. %s", index, str)
		}
	}

}
