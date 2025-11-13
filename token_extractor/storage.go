package token_extractor

import (
	"encoding/json"
	"os"
)

// Storage 存储接口
type Storage interface {
	// SaveHistory 保存历史记录
	SaveHistory(record HistoryRecord) error

	// GetHistory 获取历史记录
	GetHistory(limit int) ([]HistoryRecord, error)

	// ClearHistory 清空历史
	ClearHistory() error
}

// JSONStorage JSON文件存储实现
type JSONStorage struct {
	filePath string
}

// NewJSONStorage 创建JSON存储
func NewJSONStorage(filePath string) *JSONStorage {
	return &JSONStorage{
		filePath: filePath,
	}
}

// SaveHistory 保存历史记录
func (s *JSONStorage) SaveHistory(record HistoryRecord) error {
	// 读取现有历史
	history, err := s.GetHistory(0)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// 添加新记录
	history = append([]HistoryRecord{record}, history...)

	// 限制历史记录数量（最多保存100条）
	if len(history) > 100 {
		history = history[:100]
	}

	// 保存到文件
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0600)
}

// GetHistory 获取历史记录
func (s *JSONStorage) GetHistory(limit int) ([]HistoryRecord, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []HistoryRecord{}, nil
		}
		return nil, err
	}

	var history []HistoryRecord
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	if limit > 0 && len(history) > limit {
		return history[:limit], nil
	}

	return history, nil
}

// ClearHistory 清空历史
func (s *JSONStorage) ClearHistory() error {
	return os.Remove(s.filePath)
}
