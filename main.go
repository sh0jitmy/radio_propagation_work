package main


import (
	"fmt"
	"math"
	"net/http"


	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

var (
	Freq float64
	Power float64
	Ht float64
	Hr float64
	
	DistXAxis []float64
	FsPlot []opts.LineData 
	TwPlot []opts.LineData 
	DsfsPlot []opts.LineData
)

func main() {
	freq := float64(225*1000*1000) //Hz
	//dist := 1000.0 //m
	power := 40.0//dBm
	ht := 30.0
	hr := 30.0
	diststep :=50.0
	distance :=100.0
	distmax := 30000.0

	Freq = freq
	Power = power
	Ht = ht
	Hr = hr

	InitialPlotData()

	for distance < distmax + diststep {
		PredictDistance(freq,distance,power,ht,hr)
		distance += diststep
	}
	fmt.Printf("http server start\n")
	http.HandleFunc("/", GraphHandle)
	http.Handle("/assets/",http.StripPrefix("/assets",http.FileServer(http.Dir("assets"))))
	http.ListenAndServe(":8081", nil)
}

func PredictDistance(freq,dist,power,ht,hr float64) {
	lambda := (3.0 * math.Pow(10,8)) / freq
	/***** free space prop *****/
	fsL := (math.Pow(4*math.Pi*dist/(lambda),2)) 
	//other calc
	//fsL := 20*math.Log10(dist)+20*math.Log10(freq) - 147.6 
	
	fsLdB := 10*math.Log10(fsL) 
	//fmt.Printf("lambda:%v\n",lambda)
	//fmt.Printf("freespace loss:%v dB,recvlevel:%v dBm\n",fsLdB,power-fsLdB)
	
	/***** 2waves prop *****/
	twL := 0.0
	delta_l := 2*ht*hr/dist
	//fmt.Printf("delta_l:%v\n",delta_l)

	//https://www.apmc-mwe.org/mwe2005/src/TL/TL05-01.pdf 4.4 (12)
	twL_m := math.Pow(lambda / (2 * math.Pi * dist) * math.Sin(math.Pi/lambda*delta_l),2)
	twL = 1 / twL_m //gain to loss change 
	twLdB := 10*math.Log10(twL)
	//fmt.Printf("2Wave loss:%v dB,recvlevel:%v dBm\n",twLdB,power+twLdB)
	
	/***** Diff spherical earth prop *****/
	/* https://www.ieice.org/cs/ap/misc/denpan-db/prop_model_db/model_list/spherical_earth_diffraction/  */
	dsfsLdB:=Dsfs(freq,dist,ht,hr,lambda)

	//output scale 
	//fmt.Printf("distance:%v m, freespace:%v dBm, 2wave:%v dBm, diff spherical earth:%v dBm \n",
	//	dist,power-fsLdB,power-twLdB,power-dsfsLdB)

	InsertData(dist,power-fsLdB,power-twLdB,power-dsfsLdB)
}

func InitialPlotData(){
	DistXAxis = make([]float64,0)
	FsPlot = make([]opts.LineData,0)
	TwPlot = make([]opts.LineData,0)
	DsfsPlot = make([]opts.LineData,0)
	
	InsertData(0.0,0.0,0.0,0.0)
}


func InsertData(dist,fsRecv,twRecv,dsfsRecv float64) {
	DistXAxis = append(DistXAxis,dist)
	PlotLowerLimit(&fsRecv,&twRecv,&dsfsRecv)
	FsPlot = append(FsPlot,opts.LineData{Value:fsRecv})
	TwPlot = append(TwPlot,opts.LineData{Value:twRecv})
	DsfsPlot = append(DsfsPlot,opts.LineData{Value:dsfsRecv})
}

func PlotLowerLimit(fsRecv,twRecv,dsfsRecv *float64) {
	if *fsRecv < -150.0 { *fsRecv = -150.0 }
	if *twRecv < -150.0 { *twRecv = -150.0 }
	if *dsfsRecv < -150.0 { *dsfsRecv = -150.0 }
}

func GraphHandle(w http.ResponseWriter, _ *http.Request) {
	PropCond := fmt.Sprintf("周波数:%vMHz 送信出力:%vW 送信アンテナ高:%vm 受信アンテナ高:%vm",Freq/1000/1000,Power,Ht,Hr)
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(
			opts.Initialization{PageTitle:"RadioPropagation",Theme:types.ThemeWesteros,AssetsHost:"http://localhost:8081/assets/"},
		),
		charts.WithTitleOpts(opts.Title{
			Title: "Radio Propatation",
			Subtitle: PropCond,
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
	)
	line.SetXAxis(DistXAxis).
		AddSeries("FreeSpace",FsPlot).
		AddSeries("2Wave",TwPlot).
		AddSeries("Dsfs",DsfsPlot).
		SetSeriesOptions(
			charts.WithLineChartOpts(
				opts.LineChart{
					ShowSymbol: false,
					Smooth: true,
				}),
			charts.WithLabelOpts(opts.Label{
					//Show: true,
			}),
		)
	line.Render(w)
}


func Dsfs(freq,dist,ht,hr,lambda float64)(float64) {
	sigma:=10*math.Pow(10,-3)//%
	k:=float64(4/3)
	erad:=6378137.0//m
	ae:=float64(k*erad)
	K:=math.Sqrt(6.89*(sigma/(math.Pow(k,2/3)*math.Pow(freq,5/3))))

	beta:=(1+1.6*math.Pow(K,2)+0.67*math.Pow(K,4))/(1+4.5*math.Pow(K,2)+1.53*math.Pow(K,4))
	
	//fmt.Printf("dsfs K:%v beta:%v\n",K,beta)

	X:=dist*beta*math.Cbrt(math.Pi/(lambda*math.Pow(ae,2)))
	Y1:=2*ht*beta*math.Cbrt(math.Pow(math.Pi,2)/(math.Pow(lambda,2)*math.Pow(ae,2)))
	Y2:=2*hr*beta*math.Cbrt(math.Pow(math.Pi,2)/(math.Pow(lambda,2)*math.Pow(ae,2)))
	//fmt.Printf("dsfs x:%v y1:%v y2:%v\n",X,Y1,Y2)
	
	B1:=beta*Y1
	B2:=beta*Y2

	FX:=0.0
	GY1:=0.0
	GY2:=0.0
	
	if X>=1.6 {
		FX=11+10*math.Log10(X)-17.6*X
	} else {
		FX=-20*math.Log10(X)-5.6488*math.Pow(X,1.425)
	}
	if B1 > 2 {
		GY1=math.Sqrt(17.6*(B1-1.1))-5*math.Log10(B1-1.1)-8
	} else {
		GY1=20*math.Log10(B1+0.1*math.Pow(B1,3))
	}
	if B2 > 2 {
		GY2=math.Sqrt(17.6*(B2-1.1))-5*math.Log10(B2-1.1)-8
	} else {
		GY2=20*math.Log10(B2+0.1*math.Pow(B2,3))
	}
	loss := -(FX+GY1+GY2)
	//fmt.Printf("dsfs fx:%v gy1:%v gy2:%v\n",FX,GY1,GY2)
	//fmt.Printf("dsfs loss:%v\n",loss)
	return loss
}

	
