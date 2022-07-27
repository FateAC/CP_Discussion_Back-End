package fileManager

import (
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const (
	PostUrl    = "http://localhost:8080/post"
	PostPath   = "data/post"
	AvatarUrl  = "http://localhost:8080/avatar"
	AvatarPath = "data/avatar"
)

func SaveFile(filePath string, src io.Reader) error {
	fileDir := filepath.Dir(filePath)
	if _, err := os.Stat(fileDir); err != nil {
		err = os.MkdirAll(fileDir, os.ModePerm)
	}
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	err = dst.Close()
	if err != nil {
		return err
	}
	return nil
}

func BuildAvatarUrl(filename string) string {
	avatarUrl, _ := url.Parse(AvatarUrl)
	avatarUrl.Path = path.Join(avatarUrl.Path, filename)
	return avatarUrl.String()
}

func BuildAvatarPath(filename string) string {
	return filepath.Join(AvatarPath, filename)
}

func BuildPostUrl(year int, semester int, filename string) string {
	postUrl, _ := url.Parse(PostUrl)
	postUrl.Path = path.Join(
		postUrl.Path,
		strconv.FormatInt(int64(year), 10),
		strconv.FormatInt(int64(semester), 10),
		filename,
	)
	return postUrl.String()
}

func BuildPostPath(year int, semester int, filename string) string {
	return filepath.Join(
		PostPath,
		strconv.FormatInt(int64(year), 10),
		strconv.FormatInt(int64(semester), 10),
		filename,
	)
}
