package service

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go4s/iam/internal/db"
	"github.com/go4s/iam/internal/model"
)

// FormatField 格式字段定义
type FormatField struct {
	Name    string  `json:"name"`
	Label   string  `json:"label"`
	Type    string  `json:"type"`
	Ref     string  `json:"ref,omitempty"`
	Visible bool    `json:"visible"`
	Fold    *bool   `json:"fold,omitempty"`
}

// FormatRegistry 全局格式注册表
var (
	formatRegistry = make(map[string][]FormatField)
	formatMutex    sync.RWMutex
)

// LoadFormats 从数据库加载格式配置
func LoadFormats() error {
	var formats []model.EntityFormat
	if err := db.Engine.Find(&formats); err != nil {
		return err
	}

	formatMutex.Lock()
	defer formatMutex.Unlock()

	formatRegistry = make(map[string][]FormatField)
	for _, f := range formats {
		var fields []FormatField
		if err := json.Unmarshal([]byte(f.Fields), &fields); err != nil {
			continue
		}
		key := fmt.Sprintf("%s:%s", f.Template, f.Mode)
		formatRegistry[key] = fields
	}
	return nil
}

// GetFormat 获取指定模板和模式的格式配置
func GetFormat(template, mode string) ([]FormatField, bool) {
	formatMutex.RLock()
	defer formatMutex.RUnlock()
	key := fmt.Sprintf("%s:%s", template, mode)
	fields, ok := formatRegistry[key]
	return fields, ok
}

// ReloadFormats 重新加载格式配置
func ReloadFormats() error {
	return LoadFormats()
}

// FormatEntity 根据格式配置格式化实体
func FormatEntity(template string, fields []FormatField, data map[string]any) map[string]any {
	result := make(map[string]any)
	for _, field := range fields {
		if !field.Visible {
			continue
		}
		if val, ok := data[field.Name]; ok {
			result[field.Name] = val
			// 引用字段自动添加 fold
			if field.Type == "entity" || field.Type == "entities" {
				foldKey := field.Name + "-fold"
				if field.Fold != nil {
					result[foldKey] = *field.Fold
				} else {
					result[foldKey] = true
				}
			}
		}
	}
	return result
}
