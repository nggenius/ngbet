package toolkit

import (
	"math"
	"math/rand"
	"time"
)

// 判断两个float64是否相等
func IsEqual64(f1, f2 float64) bool {
	return math.Dim(f1, f2) < 0.000001
}

// 判断两个float32是否相等
func IsEqual32(f1, f2 float32) bool {
	return math.Dim(float64(f1), float64(f2)) < 0.000001
}

// 随机数据种子
func RandSeed() {
	rand.Seed(time.Now().Unix())
}

// 获取[min, max)之间的int32随机数
func RandRangeI32(min int32, max int32) int32 {
	if max < min || max-min == 0 {
		return 0
	}
	return min + rand.Int31n(max-min)
}

// 获取[min, max)之间的int64随机数
func RandRangeI64(min int64, max int64) int64 {
	if max < min || max-min == 0 {
		return 0
	}
	return min + rand.Int63n(max-min)
}

// 获取[min, max)之间的int随机数
func RandRange(min int, max int) int {
	if max < min || max-min == 0 {
		return 0
	}
	return min + rand.Intn(max-min)
}

// 获取[min, max)之间的float32随机数
func RandRangef(min float32, max float32) float32 {
	return min + (max-min)*rand.Float32()
}

// 获取字符串hash值
func DJBHash(str string) int32 {
	hash := 5381
	for _, c := range str {
		hash += (hash << 5) + int(c)
	}

	return int32(hash & 0x7FFFFFFF)
}