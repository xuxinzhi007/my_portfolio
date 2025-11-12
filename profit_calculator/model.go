package profit_calculator

import (
	"time"

	"github.com/google/uuid"
)

// Investor 投资者结构
type Investor struct {
	ID               string    `json:"id"`                // 唯一标识符 (UUID)
	Name             string    `json:"name"`              // 投资者姓名
	InvestmentAmount float64   `json:"investment_amount"` // 投资金额
	CreatedAt        time.Time `json:"created_at"`        // 创建时间
}

// MonthlyProfit 月度收益记录结构
type MonthlyProfit struct {
	ID            string             `json:"id"`            // 唯一标识符 (UUID)
	Date          time.Time          `json:"date"`          // 收益日期
	TotalProfit   float64            `json:"total_profit"`  // 总收益金额
	Distributions map[string]float64 `json:"distributions"` // 投资者ID -> 分配金额
	CreatedAt     time.Time          `json:"created_at"`    // 创建时间
}

// ProfitCalculatorData 整体数据容器
type ProfitCalculatorData struct {
	Investors      []Investor      `json:"investors"`
	MonthlyProfits []MonthlyProfit `json:"monthly_profits"`
}

// InvestorStats 投资者统计信息
type InvestorStats struct {
	InvestorID       string
	InvestorName     string
	InvestmentAmount float64
	InvestmentRatio  float64 // 投资比例 (0-1)
	TotalProfit      float64 // 累计收益
	FinalAmount      float64 // 最终金额 (投资 + 收益)
	ProfitCount      int     // 收益记录数
}

// OverallStats 整体统计信息
type OverallStats struct {
	TotalInvestment   float64
	TotalProfit       float64
	InvestorCount     int
	ProfitRecordCount int
}

// NewInvestor 创建新的投资者
func NewInvestor(name string, amount float64) *Investor {
	return &Investor{
		ID:               uuid.New().String(),
		Name:             name,
		InvestmentAmount: amount,
		CreatedAt:        time.Now(),
	}
}

// NewMonthlyProfit 创建新的月度收益记录
func NewMonthlyProfit(date time.Time, totalProfit float64, distributions map[string]float64) *MonthlyProfit {
	return &MonthlyProfit{
		ID:            uuid.New().String(),
		Date:          date,
		TotalProfit:   totalProfit,
		Distributions: distributions,
		CreatedAt:     time.Now(),
	}
}

// CalculateTotalInvestment 计算总投资
func CalculateTotalInvestment(investors []Investor) float64 {
	total := 0.0
	for _, investor := range investors {
		total += investor.InvestmentAmount
	}
	return total
}

// CalculateInvestmentRatio 计算投资比例
func CalculateInvestmentRatio(investor Investor, totalInvestment float64) float64 {
	if totalInvestment == 0 {
		return 0
	}
	return investor.InvestmentAmount / totalInvestment
}

// DistributeProfit 分配收益给所有投资者
func DistributeProfit(totalProfit float64, investors []Investor) map[string]float64 {
	distributions := make(map[string]float64)
	
	if len(investors) == 0 {
		return distributions
	}
	
	totalInvestment := CalculateTotalInvestment(investors)
	if totalInvestment == 0 {
		return distributions
	}
	
	for _, investor := range investors {
		ratio := CalculateInvestmentRatio(investor, totalInvestment)
		distributions[investor.ID] = totalProfit * ratio
	}
	
	return distributions
}

// CalculateInvestorStats 计算单个投资者的统计信息
func CalculateInvestorStats(investorID string, investors []Investor, profits []MonthlyProfit) InvestorStats {
	stats := InvestorStats{
		InvestorID: investorID,
	}
	
	// 查找投资者信息
	var investor *Investor
	for i := range investors {
		if investors[i].ID == investorID {
			investor = &investors[i]
			break
		}
	}
	
	if investor == nil {
		return stats
	}
	
	stats.InvestorName = investor.Name
	stats.InvestmentAmount = investor.InvestmentAmount
	
	// 计算总投资和投资比例
	totalInvestment := CalculateTotalInvestment(investors)
	stats.InvestmentRatio = CalculateInvestmentRatio(*investor, totalInvestment)
	
	// 计算累计收益
	stats.TotalProfit = 0
	stats.ProfitCount = 0
	for _, profit := range profits {
		if amount, exists := profit.Distributions[investorID]; exists {
			stats.TotalProfit += amount
			stats.ProfitCount++
		}
	}
	
	// 计算最终金额
	stats.FinalAmount = stats.InvestmentAmount + stats.TotalProfit
	
	return stats
}

// CalculateOverallStats 计算整体统计信息
func CalculateOverallStats(data *ProfitCalculatorData) OverallStats {
	stats := OverallStats{
		InvestorCount:     len(data.Investors),
		ProfitRecordCount: len(data.MonthlyProfits),
	}
	
	// 计算总投资
	stats.TotalInvestment = CalculateTotalInvestment(data.Investors)
	
	// 计算累计总收益
	stats.TotalProfit = 0
	for _, profit := range data.MonthlyProfits {
		stats.TotalProfit += profit.TotalProfit
	}
	
	return stats
}
