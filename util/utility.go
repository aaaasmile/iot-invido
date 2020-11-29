package util

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/kardianos/osext"
)

var rootPath string

const (
	useRelativeRoot = true
)

func GetEncryptedPasswordFromPrompt() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter the secret to encode")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	buf := []byte(response)
	return base64.StdEncoding.EncodeToString(buf)
}

func GetUserLogFile(serviceName string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fn := filepath.Join(usr.HomeDir, "Appdata", "Local", "Temp", fmt.Sprintf("%s.log", serviceName))
	log.Println("Logfile in use is ", fn)
	return fn
}

func GetFullPath(relPath string) string {
	if useRelativeRoot {
		return relPath
	}

	if rootPath == "" {
		var err error
		rootPath, err = osext.ExecutableFolder()
		if err != nil {
			log.Fatalf("ExecutableFolder failed: %v", err)
		}
		log.Println("Executable folder (rootdir) is ", rootPath)
		//rootPath, _ = filepath.Split(os.Args[0]) // os.Args[0] can be "faked". (https://github.com/kardianos/osext)
	}
	r := filepath.Join(rootPath, relPath)
	return r
}

func WriteJsonResp(w http.ResponseWriter, resp interface{}) error {
	blobresp, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(blobresp)

	return nil
}
