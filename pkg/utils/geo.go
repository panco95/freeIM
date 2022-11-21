package utils

import (
	"math"
	"strconv"
	"strings"
)

// 获取经纬度中心点[纬度,经度]
// 参数 points = ['30.86660,104.390740', '30.861961,104.386963', '30.842287,104.388079', '点的纬度,点的经度'……];
func GetPointsCenter(points []string) []float64 {
	point_num := len(points)
	var (
		X float64 = 0
		Y float64 = 0
		Z float64 = 0
	)
	for i := 0; i < point_num; i++ {
		if points[i] == "" {
			continue
		}
		point := strings.Split(points[i], ",")
		var (
			lat float64
			lng float64
			x   float64
			y   float64
			z   float64
		)
		point1, _ := strconv.ParseFloat(point[0], 64)
		point2, _ := strconv.ParseFloat(point[1], 64)
		lat = point1 * math.Pi / 180
		lng = point2 * math.Pi / 180
		x = math.Cos(lat) * math.Cos(lng)
		y = math.Cos(lat) * math.Sin(lng)
		z = math.Sin(lat)
		X += x
		Y += y
		Z += z
	}
	X = X / float64(point_num)
	Y = Y / float64(point_num)
	Z = Z / float64(point_num)

	tmp_lng := math.Atan2(Y, X)
	tmp_lat := math.Atan2(Z, math.Sqrt(X*X+Y*Y))
	return []float64{tmp_lat * 180 / math.Pi, tmp_lng * 180 / math.Pi}
}
