package model

type MetaData struct {
	ID        string `json:"id" bson:"_id"`
	FileName  string `json:"fileName" bson:"fileName"`
	MimeType  string `json:"mimeType" bson:"mimeType"`
	Size      int64  `json:"size" bson:"size"`
	BlurHash  string `json:"blurHash" bson:"blurHash"`
	FilePath  string `json:"filePath" bson:"filePath"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
}


