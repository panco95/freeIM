package models

import (
	"database/sql/driver"
	"encoding/json"
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt *time.Time     `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"-"`
}

type ID struct {
	ID uint `uri:"id" binding:"required"`
}

type IDBulk struct {
	IDs []uint `form:"ids" bindling:"required"`
}

type FileReq struct {
	File *multipart.FileHeader `json:"file" form:"file"`
}

type JsonDownload struct {
	Content     []byte `json:"content"`
	ContentType string `json:"contentType"`
	Filename    string `json:"filename"`
}

type ArrayFieldString []string

func (p ArrayFieldString) Value() (driver.Value, error) {
	return json.Marshal(p)
}
func (p *ArrayFieldString) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type ArrayFieldUint []uint32

func (p ArrayFieldUint) Value() (driver.Value, error) {
	return json.Marshal(p)
}
func (p *ArrayFieldUint) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type ArrayFieldFloat64 []float64

func (p ArrayFieldFloat64) Value() (driver.Value, error) {
	return json.Marshal(p)
}
func (p *ArrayFieldFloat64) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}
