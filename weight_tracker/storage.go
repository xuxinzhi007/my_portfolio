package weight_tracker

import (
	"encoding/json"
	"os"
)

// Storage 存储接口
type Storage interface {
	Load() ([]WeightRecord, error)
	Save(records []WeightRecord) error
}

// JSONStorage JSON文件存储实现
type JSONStorage struct {
	filepath string
}

// NewJSONStorage 创建新的JSON存储
func NewJSONStorage(filepath string) *JSONStorage {
	return &JSONStorage{
		filepath: filepath,
	}
}

// Load 从JSON文件加载记录
func (s *JSONStorage) Load() ([]WeightRecord, error) {
	// 检查文件是否存在
	if _, err := os.Stat(s.filepath); os.IsNotExist(err) {
		// 文件不存在，返回空列表
		return []WeightRecord{}, nil
	}

	// 读取文件
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return nil, err
	}

	// 如果文件为空，返回空列表
	if len(data) == 0 {
		return []WeightRecord{}, nil
	}

	// 解析JSON
	var records []WeightRecord
	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

// Save 保存记录到JSON文件
func (s *JSONStorage) Save(records []WeightRecord) error {
	// 序列化为JSON
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	err = os.WriteFile(s.filepath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
