package file

// FileInfo 文件基本信息
type FileInfo struct {
	Filename string `json:"filename"`
	Metadata struct {
		UserID   string   `json:"userId" bson:"userId"`
		Username string   `json:"username"`
		Tags     []string `json:"tags" bson:"tags,omitempty"`
		Public   bool     `json:"public"`
	}
}
