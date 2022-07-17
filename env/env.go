package env

import (
	"CP_Discussion/log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

var DBInfo = getDBInfo()

func getDBInfo() map[string]string {
	const projectDirName = "CP_Discussion_Back-End"
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Error.Fatal("Error loading .env file")
	}
	ret := make(map[string]string)
	ret["DBUrl"] = os.Getenv("DBUrl")
	ret["DBPort"] = os.Getenv("DBPort")
	ret["DBName"] = os.Getenv("DBName")
	ret["DBUsername"] = os.Getenv("DBUsername")
	ret["DBPassword"] = os.Getenv("DBPassword")
	return ret
}

var JWTKey = getJWTKey()

func getJWTKey() string {
	dat, err := os.ReadFile("./env/jwtKey")
	if err != nil {
		log.Error.Fatal("Cannot get jwtKey!!")
	}
	res := strings.ReplaceAll(string(dat), "\r\n", "")
	res = strings.ReplaceAll(res, "\n", "")
	return res
}
