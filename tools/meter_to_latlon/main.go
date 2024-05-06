package main

import (
	"fmt"
	"math"
)
var (
	//Lattitude
	POLE_RADIUS = 6356752.314
	JAPAN_LATITUDE_AVG = 35.0
	JAPAN_LATITUDE_MIN = 20.2531 // 沖ノ鳥島
	//JAPAN_LATITUDE_MIN = 20.0 // 沖ノ鳥島
	JAPAN_LATITUDE_MAX = 45.3326 // 択捉島
	//JAPAN_LATITUDE_MAX = 46.0 // 択捉島
	
	//Longtude
	EQATOR_RADIUS = 6378137.0
	JAPAN_LONGTUDE_MAX = 153.5912
	JAPAN_LONGTUDE_MIN = 122.5557
)

//https://easyramble.com/latitude-and-longitude-per-kilometer.html
// LONG -90.0000:21bit
//      int 7bit dp 14bit
// LAT  360.0000:22bit
//	uint 9bit dp 14bit

func main() {
	//地球 極方向 円周	
	lat_cir:= 2*math.Pi* POLE_RADIUS
	// 1kmあたり緯度
	lat_per_1km := 360 * 1000 / lat_cir	

	
	//地球 断面方向 円周 (断面計算をどの緯度で実施するかで調整が必要)
	jap_lon_cir_avg := 2*math.Pi * EQATOR_RADIUS * math.Cos(JAPAN_LATITUDE_AVG * math.Pi / 180.0)
	lon_per_1km_avg := 360 * 1000 / jap_lon_cir_avg
	jap_lon_cir_min := 2*math.Pi * EQATOR_RADIUS * math.Cos(JAPAN_LATITUDE_MIN * math.Pi / 180.0)
	lon_per_1km_min := 360 * 1000 / jap_lon_cir_min
	jap_lon_cir_max := 2*math.Pi * EQATOR_RADIUS * math.Cos(JAPAN_LATITUDE_MAX * math.Pi / 180.0)
	lon_per_1km_max := 360 * 1000 / jap_lon_cir_max
	fmt.Printf("lat/km:%v\n",lat_per_1km)
	fmt.Printf("lon/km(avg):%v lon/km(min):%v lon/km(max):%v\n",
		lon_per_1km_avg,lon_per_1km_min,lon_per_1km_max)

	
}
