package config

import (
	"os"
	"path/filepath"
	"strings"
)

// WalkDir 获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) ([]string, error) {
	files := make([]string, 0, 30)
	_, err := os.Stat(dirPth)
	if err == nil || os.IsExist(err) {
		suffix = strings.ToUpper(suffix)
		err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
			// 忽略目录
			if fi.IsDir() {
				return nil
			}

			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
				files = append(files, filename)
			}
			return nil
		})
	}

	return files, err
}
