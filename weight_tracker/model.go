package weight_tracker

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// WeightRecord 体重记录结构
type WeightRecord struct {
	ID         string    `json:"id"`
	Weight     float64   `json:"weight"`
	Date       time.Time `json:"date"`
	Change     float64   `json:"change"`
	ChangeType string    `json:"change_type"` // "increase", "decrease", "stable", "first"
}

// CalculateChange 计算体重变化
func CalculateChange(current, previous float64) (change float64, changeType string) {
	change = current - previous

	if change > 0 {
		changeType = "increase"
	} else if change < 0 {
		changeType = "decrease"
	} else {
		changeType = "stable"
	}

	return change, changeType
}

// NewWeightRecord 创建新的体重记录
func NewWeightRecord(weight float64, previousRecord *WeightRecord) *WeightRecord {
	record := &WeightRecord{
		ID:     uuid.New().String(),
		Weight: weight,
		Date:   time.Now(),
	}

	if previousRecord == nil {
		record.Change = 0
		record.ChangeType = "first"
	} else {
		record.Change, record.ChangeType = CalculateChange(weight, previousRecord.Weight)
	}

	return record
}

// FormatChange 格式化变化显示文本
func (r *WeightRecord) FormatChange() string {
	switch r.ChangeType {
	case "first":
		return "● 首次记录"
	case "increase":
		return fmt.Sprintf("↑ +%.1f kg", r.Change)
	case "decrease":
		return fmt.Sprintf("↓ %.1f kg", r.Change)
	case "stable":
		return "● 持平"
	default:
		return ""
	}
}

// FormatDate 格式化日期显示
func (r *WeightRecord) FormatDate() string {
	return r.Date.Format("2006-01-02 15:04")
}

// WeightStats 体重统计信息
type WeightStats struct {
	TotalRecords  int
	CurrentWeight float64
	StartWeight   float64
	TotalChange   float64
	HighestWeight float64
	LowestWeight  float64
}

// CalculateStats 计算统计信息
func CalculateStats(records []WeightRecord) *WeightStats {
	if len(records) == 0 {
		return &WeightStats{}
	}

	stats := &WeightStats{
		TotalRecords:  len(records),
		CurrentWeight: records[0].Weight,
		StartWeight:   records[len(records)-1].Weight,
		HighestWeight: records[0].Weight,
		LowestWeight:  records[0].Weight,
	}

	// 计算最高和最低体重
	for _, record := range records {
		if record.Weight > stats.HighestWeight {
			stats.HighestWeight = record.Weight
		}
		if record.Weight < stats.LowestWeight {
			stats.LowestWeight = record.Weight
		}
	}

	// 计算总变化
	stats.TotalChange = stats.CurrentWeight - stats.StartWeight

	return stats
}
