package base

type FileStruct struct {
	FileId   string `json:"file_id" bson:"file_id"`
	FileName string `json:"file_name" bson:"file_name"`
	FileSize string `json:"file_size" bson:"file_size"`
	FilePath string `json:"file_path" bson:"file_path"`
	UserId   string `json:"user_id" bson:"user_id"`
	FileTime string `json:"file_time" bson:"file_time"`
	FileHash string `json:"file_hash" bson:"file_hash"`
	OssPath  string `json:"oss_path" bson:"oss_path"`
}

type FilesResponse struct {
	Code  int           `json:"code"`
	Files []*FileStruct `json:"files"`
}

type UserCapacity struct {
	TotalCapacity  int `json:"total_capacity"`
	UsedCapacity   int `json:"used_capacity"`
	RemainCapacity int `json:"remain_capacity"`
}
