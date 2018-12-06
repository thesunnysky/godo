package godo

import (
	"bufio"
	"bytes"
	"fmt"
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

func init() {
	dataFile = ClientConfig.DataFile
}

type GitRepo struct {
	repoPath string
}

var r, _ = regexp.Compile("[[:alnum:]]")

func AddCmdImpl(args []string) {
	var buf bytes.Buffer
	for _, str := range args {
		buf.WriteString(str)
		buf.WriteByte(' ')
	}

	f, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, consts.FILE_MAKS)
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
			os.Exit(consts.INVALID_PARAMETER_VALUE)
		} else {
			num = append(num, i)
		}
	}

	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDWR, consts.FILE_MAKS)
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

func ListCmdImpl(args []string) {
	f, err := os.OpenFile(dataFile, os.O_CREATE|os.O_RDONLY, consts.FILE_MAKS)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	br := bufio.NewReader(f)
	var index int
	for {
		str, err := br.ReadString(consts.LINE_SEPARATOR)
		if err == io.EOF {
			break
		}
		index++
		if !isBlankLine(str) {
			fmt.Printf("%d. %s", index, str)
		}
	}
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
		str, err := br.ReadString(consts.LINE_SEPARATOR)
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
	return !r.MatchString(str)
}

func PullCmd(args []string) {
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
	aesUtil := util.Aes{Key: []byte(ClientConfig.AesGCMConfig.AesGCMKey),
		Nonce: []byte(ClientConfig.AesGCMConfig.AesGCMNonce)}
	decryptedData, err := aesUtil.GcmDecrypt(data)
	if err != nil {
		fmt.Printf("decrypt data error:%s\n", err)
	}

	//write download data to tempFile
	tempFile := ClientConfig.DataFile + ".tmp"
	targetFile := ClientConfig.DataFile
	f, err := os.OpenFile(tempFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, consts.FILE_MAKS)
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
		Key:   ClientConfig.AesGCMConfig.AesGCMKey,
		Nonce: ClientConfig.AesGCMConfig.AesGCMNonce}
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
	publicKey, err := ioutil.ReadFile(ClientConfig.RsaConfig.RsaPublicKeyFile)
	if err != nil {
		fmt.Printf("read rsa public key file error:%s\n", err)
		os.Exit(-1)
	}

	privateKey, err := ioutil.ReadFile(ClientConfig.RsaConfig.RsaPrivateKeyFile)
	if err != nil {
		fmt.Printf("read rsa private key file error:%s\n", err)
		os.Exit(-1)
	}

	he := util.NewHyperEncrypt(publicKey, privateKey)
	cipherData, err := he.Encrypt(textData)
	if err != nil {
		fmt.Printf("encryt task file error:%s\n", err)
		os.Exit(-1)
	}

	gitFile := ClientConfig.GithubRepo + "/" + ClientConfig.DataFile
	if err := ioutil.WriteFile(gitFile, cipherData, consts.FILE_MAKS);
		err != nil {
		fmt.Printf("write git file error:%s\n", gitFile)
		os.Exit(-1)
	}

	gitCmArgs := []string{"commit", "-am", "\"\""}
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
