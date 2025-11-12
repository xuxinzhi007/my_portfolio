package profit_calculator

import (
	"encoding/json"
	"os"
)

// Storage 存储接口
type Storage interface {
	Load() (*ProfitCalculatorData, error)
	Save(data *ProfitCalculatorData) error
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

// Load 从JSON文件加载数据
func (s *JSONStorage) Load() (*ProfitCalculatorData, error) {
	// 检查文件是否存在
	if _, err := os.Stat(s.filepath); os.IsNotExist(err) {
		// 文件不存在，返回空数据
		return &ProfitCalculatorData{
			Investors:      []Investor{},
			MonthlyProfits: []MonthlyProfit{},
		}, nil
	}

	// 读取文件
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return nil, err
	}

	// 如果文件为空，返回空数据
	if len(data) == 0 {
		return &ProfitCalculatorData{
			Investors:      []Investor{},
			MonthlyProfits: []MonthlyProfit{},
		}, nil
	}

	// 解析JSON
	var profitData ProfitCalculatorData
	err = json.Unmarshal(data, &profitData)
	if err != nil {
		return nil, err
	}

	// 确保切片不为nil
	if profitData.Investors == nil {
		profitData.Investors = []Investor{}
	}
	if profitData.MonthlyProfits == nil {
		profitData.MonthlyProfits = []MonthlyProfit{}
	}

	return &profitData, nil
}

// Save 保存数据到JSON文件
func (s *JSONStorage) Save(data *ProfitCalculatorData) error {
	// 序列化为JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	err = os.WriteFile(s.filepath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
