package models

import (
	"time"
)

// File 文件表
type File struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FileName     string    `gorm:"type:varchar(255);not null;comment:文件名" json:"file_name"`
	OriginalName string    `gorm:"type:varchar(255);not null;comment:原始文件名" json:"original_name"`
	FilePath     string    `gorm:"type:varchar(500);not null;comment:文件路径" json:"file_path"`
	FileSize     int64     `gorm:"not null;comment:文件大小(字节)" json:"file_size"`
	FileType     string    `gorm:"type:varchar(100);not null;comment:文件类型" json:"file_type"`
	MimeType     string    `gorm:"type:varchar(100);not null;comment:MIME类型" json:"mime_type"`
	MD5          string    `gorm:"type:varchar(32);uniqueIndex;comment:文件MD5" json:"md5"`
	StorageType  string    `gorm:"type:varchar(20);default:local;comment:存储类型:local/oss/qiniu" json:"storage_type"`
	BucketName   string    `gorm:"type:varchar(100);comment:OSS桶名" json:"bucket_name"`
	UploadedBy   uint      `gorm:"not null;index;comment:上传人ID" json:"uploaded_by"`
	IsDeleted    bool      `gorm:"default:false;index;comment:是否删除" json:"is_deleted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	Uploader User `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
}
