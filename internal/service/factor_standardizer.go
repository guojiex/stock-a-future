package service

import (
	"math"
	"sort"
	"time"

	"stock-a-future/internal/models"

	"github.com/shopspring/decimal"
)

// FactorStandardizer 因子标准化器
type FactorStandardizer struct{}

// NewFactorStandardizer 创建因子标准化器
func NewFactorStandardizer() *FactorStandardizer {
	return &FactorStandardizer{}
}

// StandardizeFactors 标准化因子数据
func (fs *FactorStandardizer) StandardizeFactors(factors []models.FundamentalFactor) ([]models.FundamentalFactor, error) {
	if len(factors) == 0 {
		return factors, nil
	}

	// 按因子名称分组
	factorGroups := fs.groupFactorsByName(factors)

	// 对每组因子进行标准化
	var standardizedFactors []models.FundamentalFactor
	for _, group := range factorGroups {
		standardized := fs.standardizeFactorGroup(group)
		standardizedFactors = append(standardizedFactors, standardized...)
	}

	return standardizedFactors, nil
}

// StandardizeFundamentalFactors 标准化基本面因子
func (fs *FactorStandardizer) StandardizeFundamentalFactors(factors []models.FundamentalFactor) ([]models.FundamentalFactor, error) {
	if len(factors) == 0 {
		return factors, nil
	}

	// 标准化各个因子字段
	fs.standardizeFundamentalField(factors, "PE")
	fs.standardizeFundamentalField(factors, "PB")
	fs.standardizeFundamentalField(factors, "PS")
	fs.standardizeFundamentalField(factors, "ROE")
	fs.standardizeFundamentalField(factors, "ROA")
	fs.standardizeFundamentalField(factors, "RevenueGrowth")
	fs.standardizeFundamentalField(factors, "NetProfitGrowth")
	fs.standardizeFundamentalField(factors, "NetMargin")
	fs.standardizeFundamentalField(factors, "GrossMargin")
	fs.standardizeFundamentalField(factors, "DebtToAssets")

	// 计算综合得分
	fs.calculateCompositeScores(factors)

	// 计算排名和分位数
	fs.calculateRankingsAndPercentiles(factors)

	return factors, nil
}

// groupFactorsByName 按因子名称分组
func (fs *FactorStandardizer) groupFactorsByName(factors []models.FundamentalFactor) map[string][]models.FundamentalFactor {
	groups := make(map[string][]models.FundamentalFactor)

	for _, factor := range factors {
		// Group by TSCode since FundamentalFactor doesn't have Name field
		groups[factor.TSCode] = append(groups[factor.TSCode], factor)
	}

	return groups
}

// standardizeFactorGroup 标准化因子组
func (fs *FactorStandardizer) standardizeFactorGroup(factors []models.FundamentalFactor) []models.FundamentalFactor {
	if len(factors) <= 1 {
		// 如果只有一个或没有因子，无法标准化，设置默认分位数
		for i := range factors {
			factors[i].MarketPercentile = models.NewJSONDecimal(decimal.NewFromFloat(50))
			factors[i].IndustryPercentile = models.NewJSONDecimal(decimal.NewFromFloat(50))
			factors[i].UpdatedAt = time.Now()
		}
		return factors
	}

	// 标准化各个因子类别的得分
	fs.standardizeScoreField(factors, "ValueScore")
	fs.standardizeScoreField(factors, "GrowthScore")
	fs.standardizeScoreField(factors, "QualityScore")
	fs.standardizeScoreField(factors, "ProfitabilityScore")

	// 计算综合得分
	fs.calculateCompositeScores(factors)

	// 计算排名和分位数
	fs.calculateRankingsAndPercentiles(factors)

	return factors
}

// standardizeScoreField 标准化得分字段
func (fs *FactorStandardizer) standardizeScoreField(factors []models.FundamentalFactor, fieldName string) {
	// 提取有效值
	var values []float64
	var validIndices []int

	for i, factor := range factors {
		var value decimal.Decimal
		switch fieldName {
		case "ValueScore":
			value = factor.ValueScore.Decimal
		case "GrowthScore":
			value = factor.GrowthScore.Decimal
		case "QualityScore":
			value = factor.QualityScore.Decimal
		case "ProfitabilityScore":
			value = factor.ProfitabilityScore.Decimal
		}

		// 过滤零值
		if !value.IsZero() {
			val, _ := value.Float64()
			values = append(values, val)
			validIndices = append(validIndices, i)
		}
	}

	if len(values) <= 1 {
		return // 无法标准化
	}

	// 计算统计量
	mean := fs.calculateMean(values)
	stdDev := fs.calculateStdDev(values, mean)

	// 标准化处理
	for i, idx := range validIndices {
		val := values[i]

		// Z-score标准化
		var zscore decimal.Decimal
		if stdDev > 0 {
			zscore = decimal.NewFromFloat((val - mean) / stdDev)
		}

		// 分位数排名
		percentile := fs.calculatePercentile(val, values)

		// 更新对应字段的标准化得分
		switch fieldName {
		case "ValueScore":
			factors[idx].ValueScore = models.NewJSONDecimal(zscore)
		case "GrowthScore":
			factors[idx].GrowthScore = models.NewJSONDecimal(zscore)
		case "QualityScore":
			factors[idx].QualityScore = models.NewJSONDecimal(zscore)
		case "ProfitabilityScore":
			factors[idx].ProfitabilityScore = models.NewJSONDecimal(zscore)
		}

		// 更新市场分位数
		factors[idx].MarketPercentile = models.NewJSONDecimal(decimal.NewFromFloat(percentile))
		factors[idx].UpdatedAt = time.Now()
	}
}

// standardizeFundamentalField 标准化基本面因子字段
func (fs *FactorStandardizer) standardizeFundamentalField(factors []models.FundamentalFactor, fieldName string) {
	// 提取有效值
	var values []float64
	var validIndices []int

	for i, factor := range factors {
		var value decimal.Decimal
		switch fieldName {
		case "PE":
			value = factor.PE.Decimal
		case "PB":
			value = factor.PB.Decimal
		case "PS":
			value = factor.PS.Decimal
		case "ROE":
			value = factor.ROE.Decimal
		case "ROA":
			value = factor.ROA.Decimal
		case "RevenueGrowth":
			value = factor.RevenueGrowth.Decimal
		case "NetProfitGrowth":
			value = factor.NetProfitGrowth.Decimal
		case "NetMargin":
			value = factor.NetMargin.Decimal
		case "GrossMargin":
			value = factor.GrossMargin.Decimal
		case "DebtToAssets":
			value = factor.DebtToAssets.Decimal
		}

		// 过滤异常值和零值
		if !value.IsZero() && fs.isValidValue(value, fieldName) {
			val, _ := value.Float64()
			values = append(values, val)
			validIndices = append(validIndices, i)
		}
	}

	if len(values) <= 1 {
		return // 无法标准化
	}

	// 计算统计量
	mean := fs.calculateMean(values)
	stdDev := fs.calculateStdDev(values, mean)

	// 标准化处理
	for i, idx := range validIndices {
		val := values[i]

		// Z-score标准化
		var zscore decimal.Decimal
		if stdDev > 0 {
			zscore = decimal.NewFromFloat((val - mean) / stdDev)
		}

		// 分位数排名
		percentile := fs.calculatePercentile(val, values)

		// 根据因子特性调整Z-score方向
		adjustedZScore := fs.adjustZScoreDirection(zscore, fieldName)

		// 更新对应字段的标准化得分（这里简化处理，实际可以扩展更多字段）
		switch fieldName {
		case "PE", "PB", "PS":
			// 价值因子：越低越好
			factors[idx].ValueScore = models.NewJSONDecimal(adjustedZScore)
		case "ROE", "ROA":
			// 质量因子：越高越好
			factors[idx].QualityScore = models.NewJSONDecimal(adjustedZScore)
		case "RevenueGrowth", "NetProfitGrowth":
			// 成长因子：越高越好
			factors[idx].GrowthScore = models.NewJSONDecimal(adjustedZScore)
		case "NetMargin", "GrossMargin":
			// 盈利因子：越高越好
			factors[idx].ProfitabilityScore = models.NewJSONDecimal(adjustedZScore)
		}

		// 更新行业和市场分位数
		factors[idx].MarketPercentile = models.NewJSONDecimal(decimal.NewFromFloat(percentile))
	}
}

// calculateCompositeScores 计算综合得分
func (fs *FactorStandardizer) calculateCompositeScores(factors []models.FundamentalFactor) {
	// 默认权重配置
	weights := map[string]decimal.Decimal{
		"value":         decimal.NewFromFloat(0.25), // 价值因子权重
		"growth":        decimal.NewFromFloat(0.25), // 成长因子权重
		"quality":       decimal.NewFromFloat(0.25), // 质量因子权重
		"profitability": decimal.NewFromFloat(0.25), // 盈利因子权重
	}

	for i := range factors {
		// 计算加权综合得分
		compositeScore := factors[i].ValueScore.Decimal.Mul(weights["value"]).
			Add(factors[i].GrowthScore.Decimal.Mul(weights["growth"])).
			Add(factors[i].QualityScore.Decimal.Mul(weights["quality"])).
			Add(factors[i].ProfitabilityScore.Decimal.Mul(weights["profitability"]))

		factors[i].CompositeScore = models.NewJSONDecimal(compositeScore)
	}
}

// calculateRankingsAndPercentiles 计算排名和分位数
func (fs *FactorStandardizer) calculateRankingsAndPercentiles(factors []models.FundamentalFactor) {
	// 按综合得分排序
	sort.Slice(factors, func(i, j int) bool {
		return factors[i].CompositeScore.Decimal.GreaterThan(factors[j].CompositeScore.Decimal)
	})

	// 设置市场排名和分位数
	for i := range factors {
		factors[i].MarketRank = i + 1
		percentile := decimal.NewFromInt(int64(len(factors) - i)).
			Div(decimal.NewFromInt(int64(len(factors)))).
			Mul(decimal.NewFromInt(100))
		factors[i].MarketPercentile = models.NewJSONDecimal(percentile)
	}
}

// CalculateIndustryRankings 计算行业排名
func (fs *FactorStandardizer) CalculateIndustryRankings(factors []models.FundamentalFactor, industryField string) ([]models.FundamentalFactor, error) {
	// 按行业分组（这里简化处理，实际需要行业信息）
	// 由于FundamentalFactor结构中没有行业字段，这里模拟处理

	// 按综合得分在行业内排序
	sort.Slice(factors, func(i, j int) bool {
		return factors[i].CompositeScore.Decimal.GreaterThan(factors[j].CompositeScore.Decimal)
	})

	// 设置行业排名（这里简化为全市场排名）
	for i := range factors {
		factors[i].IndustryRank = i + 1
		percentile := decimal.NewFromInt(int64(len(factors) - i)).
			Div(decimal.NewFromInt(int64(len(factors)))).
			Mul(decimal.NewFromInt(100))
		factors[i].IndustryPercentile = models.NewJSONDecimal(percentile)
	}

	return factors, nil
}

// NormalizeFactorsToRange 将因子标准化到指定范围
func (fs *FactorStandardizer) NormalizeFactorsToRange(factors []models.FundamentalFactor, minVal, maxVal float64) []models.FundamentalFactor {
	if len(factors) <= 1 || minVal >= maxVal {
		return factors
	}

	// 找到当前值的最大值和最小值
	var currentMin, currentMax float64
	first := true

	for _, factor := range factors {
		val, _ := factor.PE.Decimal.Float64()
		if first {
			currentMin = val
			currentMax = val
			first = false
		} else {
			if val < currentMin {
				currentMin = val
			}
			if val > currentMax {
				currentMax = val
			}
		}
	}

	if currentMax == currentMin {
		// 所有值相同，设置为范围中点
		midPoint := (minVal + maxVal) / 2
		for i := range factors {
			factors[i].PE = models.NewJSONDecimal(decimal.NewFromFloat(midPoint))
		}
		return factors
	}

	// 线性缩放到目标范围
	scale := (maxVal - minVal) / (currentMax - currentMin)

	for i := range factors {
		val, _ := factors[i].PE.Decimal.Float64()
		normalizedVal := minVal + (val-currentMin)*scale
		factors[i].PE = models.NewJSONDecimal(decimal.NewFromFloat(normalizedVal))
		factors[i].UpdatedAt = time.Now()
	}

	return factors
}

// 辅助方法

// calculateMean 计算均值
func (fs *FactorStandardizer) calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// calculateStdDev 计算标准差
func (fs *FactorStandardizer) calculateStdDev(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += math.Pow(v-mean, 2)
	}
	return math.Sqrt(sum / float64(len(values)-1)) // 使用样本标准差
}

// calculatePercentile 计算分位数排名
func (fs *FactorStandardizer) calculatePercentile(value float64, values []float64) float64 {
	if len(values) == 0 {
		return 50.0
	}

	count := 0
	for _, v := range values {
		if v <= value {
			count++
		}
	}

	return float64(count) / float64(len(values)) * 100.0
}

// isValidValue 检查数值是否有效
func (fs *FactorStandardizer) isValidValue(value decimal.Decimal, fieldName string) bool {
	val, _ := value.Float64()

	// 检查是否为无穷大或NaN
	if math.IsInf(val, 0) || math.IsNaN(val) {
		return false
	}

	// 根据因子类型设置合理范围
	switch fieldName {
	case "PE":
		return val > 0 && val < 1000 // PE应该为正数且不会太大
	case "PB":
		return val > 0 && val < 100 // PB应该为正数
	case "PS":
		return val > 0 && val < 100 // PS应该为正数
	case "ROE", "ROA":
		return val > -100 && val < 100 // ROE/ROA在-100%到100%之间
	case "RevenueGrowth", "NetProfitGrowth":
		return val > -100 && val < 1000 // 增长率在合理范围内
	case "NetMargin", "GrossMargin":
		return val > -100 && val < 100 // 利润率在-100%到100%之间
	case "DebtToAssets":
		return val >= 0 && val <= 100 // 资产负债率在0-100%之间
	default:
		return true
	}
}

// adjustZScoreDirection 根据因子特性调整Z-score方向
func (fs *FactorStandardizer) adjustZScoreDirection(zscore decimal.Decimal, fieldName string) decimal.Decimal {
	switch fieldName {
	case "PE", "PB", "PS", "DebtToAssets":
		// 这些因子越低越好，所以取负值
		return zscore.Neg()
	case "ROE", "ROA", "RevenueGrowth", "NetProfitGrowth", "NetMargin", "GrossMargin":
		// 这些因子越高越好，保持原值
		return zscore
	default:
		return zscore
	}
}

// GetFactorStatistics 获取因子统计信息
func (fs *FactorStandardizer) GetFactorStatistics(factors []models.FundamentalFactor) map[string]FactorStats {
	stats := make(map[string]FactorStats)

	// 按因子名称分组
	groups := fs.groupFactorsByName(factors)

	for name, group := range groups {
		var values []float64
		for _, factor := range group {
			val, _ := factor.PE.Decimal.Float64()
			values = append(values, val)
		}

		if len(values) > 0 {
			stats[name] = FactorStats{
				Count:  len(values),
				Mean:   fs.calculateMean(values),
				StdDev: fs.calculateStdDev(values, fs.calculateMean(values)),
				Min:    fs.findMin(values),
				Max:    fs.findMax(values),
			}
		}
	}

	return stats
}

// FactorStats 因子统计信息
type FactorStats struct {
	Count  int     `json:"count"`   // 样本数量
	Mean   float64 `json:"mean"`    // 均值
	StdDev float64 `json:"std_dev"` // 标准差
	Min    float64 `json:"min"`     // 最小值
	Max    float64 `json:"max"`     // 最大值
}

// findMin 找最小值
func (fs *FactorStandardizer) findMin(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// findMax 找最大值
func (fs *FactorStandardizer) findMax(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}
