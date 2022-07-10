package env

import (
	"CP_Discussion/log"
	"os"
	"strings"
)

const DBUrl = "localhost"
const DBPort = "9487"
const DBName = "cp-discussion-db"
const DBUsername = "CPDiscussion"
const DBPassword = "94879487"

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
