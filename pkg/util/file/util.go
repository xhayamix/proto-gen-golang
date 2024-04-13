package file

import "os"

// RemoveFile ファイルを削除する
func RemoveFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	if _, statErr := os.Stat(filePath); statErr == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return nil
}

// Exists ファイルの存在判定をする. (存在している場合はtrue)
func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
