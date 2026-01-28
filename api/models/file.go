package models

import (
	"time"
)

// FileType 文件类型
const (
	FileTypeAvatar   = "avatar"   // 头像
	FileTypeDocument = "document" // 文档
	FileTypeImage    = "image"    // 图片
	FileTypeAudio    = "audio"    // 音频
	FileTypeVideo    = "video"    // 视频
	FileTypeOther    = "other"    // 其他
)

// File 文件表
type File struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FileName     string    `gorm:"type:varchar(255);not null;comment:文件名" json:"file_name"`
	OriginalName string    `gorm:"type:varchar(255);not null;comment:原始文件名" json:"original_name"`
	FilePath     string    `gorm:"type:varchar(500);not null;comment:文件路径" json:"file_path"`
	FileURL      string    `gorm:"type:varchar(500);not null;comment:文件URL" json:"file_url"`
	FileSize     int64     `gorm:"not null;comment:文件大小(字节)" json:"file_size"`
	FileType     string    `gorm:"type:varchar(20);not null;index;comment:文件类型" json:"file_type"`
	MimeType     string    `gorm:"type:varchar(100);comment:MIME类型" json:"mime_type"`
	UploaderID   uint      `gorm:"not null;index;comment:上传者ID" json:"uploader_id"`
	RelationID   uint      `gorm:"index;comment:关联ID" json:"relation_id"`
	RelationType string    `gorm:"type:varchar(50);index;comment:关联类型:order/session/message" json:"relation_type"`
	Status       int       `gorm:"not null;default:1;comment:状态:1-正常,0-删除" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	Uploader User `gorm:"foreignKey:UploaderID" json:"uploader,omitempty"`
}
