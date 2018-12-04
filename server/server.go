package server

import (
	"github.com/thesunnysky/godo/consts"
	"github.com/thesunnysky/godo/util"
	"io"
	"log"
	"net/http"
	"os"
)

var serverConfig = ServerConfig

func Run() {
	//upload file
	http.HandleFunc("/upload", uploadHandle)

	//download file server
	//Ex: http://127.0.0.1:9090/download/godo.dat
	http.Handle("/download/", http.StripPrefix("/download/",
		http.FileServer(http.Dir(serverConfig.DataDir))))

	log.Fatal(http.ListenAndServe(":9090", nil))

}

func uploadHandle(w http.ResponseWriter, r *http.Request) {

	// 根据字段名获取表单文件
	formFile, header, err := r.FormFile(consts.GODO_DATA_FILE)
	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		return
	}
	defer formFile.Close()

	if err := util.CreateDirIsNotExist(serverConfig.DataDir); err != nil {
		log.Printf("create data dir error:%s\n", err)
		return
	}
	// 创建保存文件
	destFile, err := os.Create(serverConfig.DataDir + "/" + header.Filename)
	if err != nil {
		log.Printf("Create failed: %s\n", err)
		return
	}
	defer destFile.Close()

	// 读取表单文件，写入保存文件
	_, err = io.Copy(destFile, formFile)
	if err != nil {
		log.Printf("Write file failed: %s\n", err)
		return
	}

	log.Printf("receive file:%s successfully\n", formFile)
}

/*func handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	//解析参数，默认是不会解析的
	if err := r.ParseForm(); err != nil {
		log.Println("Parse from error:", err.Error())
	}
	log.Println("Recv:", r.RemoteAddr)
	pwd, _ := os.Getwd()
	des := pwd + string(os.PathSeparator) + r.URL.Path[1:]
	desStat, err := os.Stat(des)
	if err != nil {
		log.Println("File Not Exit", des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else if desStat.IsDir() {
		log.Println("File Is Dir", des)
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		fileData, err := ioutil.ReadFile(des)
		if err != nil {
			log.Println("Read File Err:", err.Error())
		} else {
			log.Println("Send File:", des)
			w.Write(fileData)
		}
	}
}*/
