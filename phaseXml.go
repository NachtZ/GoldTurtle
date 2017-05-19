package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"time"
	"log"
	"os"

)

type Gold struct {
	date    string
	price   float64
	open    float64
	buy 	float64
	sell 	float64
	high    float64
	low     float64
	percent string
}

type Direct struct {
	amount   float64
	total    float64
	enter    float64
	unitNum int
}

type BaseData struct {
	n      float64
	tr     float64
	pdc    float64
	pdn    float64
	high10 float64
	low10  float64
	high20 float64
	low20  float64
	high55 float64
	low55  float64
}

type Turtle struct {
	buyData  Direct
	sellData Direct
	base     BaseData
	perdeal  float64
	total    float64
	buy      bool
	sell     bool
}

var from,dest,pwd,server string
var port int
var values = make([]float64,30)
var ints = make([]int,4)
var timenow = "Init"

func max3(a, b, c float64) float64 {
	if a > b {
		if a > c {
			return a
		} else {
			return c
		}
	} else {
		if b > c {
			return b
		} else {
			return c
		}
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b float64) float64 {
	if a > b {
		return b
	} else {
		return a
	}
}

func getHighLow(idx, day int, g []Gold) (high, low float64) {
	high, low = -1, -1
	if idx > len(g)+day {
		return
	}
	high = g[idx].high
	low = g[idx].low
	for i := 0; i < day; i++ {
		high = max(g[idx+i].high, high)
		low = min(g[idx+i].low, low)
	}
	return
}

func runTurtle(g []Gold) {
	dataLen := len(g)
	if dataLen < 22 {
		log.Println("Gold data is too short to run trutle.")
		return
	}
	var i, idx int
	var bamount, bmoney, samount, smoney, total, nextbuy, nextsell float64
	var perdeal float64
	var high55, low55, high20, low20, high10, low10, benter, senter float64
	sell, buy := false, false
	perdeal = 2000
	total = 20000
	//here to get first N value and high and low.
	firstN := 0.0
	// pdc :=0.0
	for i = 1; i <= 21; i++ {
		idx := dataLen - i
		//  pdc = g[idx].open
		firstN += max3(g[idx].open-g[idx].low, g[idx].high-g[idx].open, g[idx].high-g[idx].low)
	}
	firstN /= 20
	log.Println(firstN)
	high20, low20 = getHighLow(dataLen-56, 20, g)
	high55, low55 = getHighLow(dataLen-56, 55, g)
	//second get high, low

	for idx = 300 - 56; idx != 0; idx-- {
		high10, low10 = getHighLow(idx, 10, g)
		high20, low20 = getHighLow(idx, 20, g)
		high55, low55 = getHighLow(idx, 55, g)
		tr := max3(g[idx].open-g[idx].low, g[idx].high-g[idx].open, g[idx].high-g[idx].low)
		firstN = (firstN*19 + tr) / 20 //caculate the n value
		if buy || sell {               //whether to leave
			//earn and lose
			//leave to earn
			if buy {
				if g[idx].low < nextbuy && total >= perdeal && bmoney <= 8*perdeal {
					total -= perdeal
					nextbuy = g[idx].low - firstN/2
					bmoney += perdeal
					bamount += perdeal / g[idx].low
				} else {
					if g[idx].high == high10 {
						//leave with earn
						total += bamount * g[idx].high
						bamount = 0
						bmoney = 0
						buy = false
					} else if benter-2*firstN > g[idx].low {
						//leave with lose
						total += bamount * g[idx].low
						bamount = 0
						bmoney = 0
						buy = false
						benter = 0
					}
				}
			}
			if sell {
				if g[idx].high >= nextsell && total >= perdeal && smoney <= 8*perdeal {
					total -= perdeal
					nextsell = g[idx].high + firstN/2
					smoney += perdeal
					bamount += perdeal / g[idx].high
				} else {
					if g[idx].low == low10 {
						total += smoney + smoney - g[idx].low*samount
						smoney = 0
						samount = 0
						sell = false
						//leave with earn
					} else if senter+2*firstN < g[idx].high {
						//leave with lose
						total += smoney + smoney - g[idx].low*samount
						smoney = 0
						samount = 0
						senter = 0
						sell = false
					}
				}
			}
		} else { //whether to enter
			//get a break in 20 days
			if g[idx].high == high20 || g[idx].low == low20 {
				if high20 == g[idx].high { //sell
					if total < perdeal*4 {
						continue
					}
					sell = true
					senter = high20
					nextsell = senter + firstN/2
					total -= perdeal * 4
					smoney += perdeal * 4
					samount += perdeal * 4 / (g[idx].high-2.6)
				}
				if low20 == g[idx].low { //buy
					if total < perdeal*4 {
						continue
					}
					buy = true
					benter = low20
					nextbuy = senter - firstN/2
					total -= perdeal * 4
					bmoney += perdeal * 4
					bamount += perdeal * 4 / (g[idx].low+2.6)
				}
			}
		}
	}
	log.Println(firstN)
	log.Println(low20, high20)
	log.Println(low55, high55)
	log.Println(total, smoney, bmoney)

}

func getGoldData(path string) []Gold {
	var t xml.Token
	var err error
	var g []Gold
	time := 0
	count := 0
	content, err := ioutil.ReadFile(path)
	if err!= nil{
		log.Println(err)
		return g
	}
	//fmt.Println(content)
	decoder := xml.NewDecoder(bytes.NewBuffer(content))
	flag := true
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token1 := t.(type) {
		case xml.StartElement:
			name1 := token1.Name.Local
			if name1 == "tr" {
				var gold Gold
				count = 0
				flag = true
				for t, err = decoder.Token(); err == nil && flag; t, err = decoder.Token() {
					switch token := t.(type) {
					case xml.StartElement:
						count = count + 1

					case xml.CharData:
						content := string([]byte(token))
						switch count {
						case 1:
							gold.date = content
						case 3:
							gold.price, _ = strconv.ParseFloat(content, 32)
						case 5:
							gold.open, _ = strconv.ParseFloat(content, 32)
						case 7:
							gold.high, _ = strconv.ParseFloat(content, 32)
						case 9:
							gold.low, _ = strconv.ParseFloat(content, 32)
						case 13:
							gold.percent = content
						default:
							//
						}
					case xml.EndElement:
						count++
						name := token.Name.Local
						if name == "tr" {
							time++
							flag = false
						}
					default:

					}
				}
				if gold.low > 0 {
					g = append(g, gold)
				}
			}
		}
	}

	return g
}

func crawlGoldNow() (Gold,error) {
	var g Gold
	url := "http://www.icbc.com.cn/ICBCDynamicSite/Charts/GoldTendencyPicture.aspx"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(1, err)
		return g,err
	}
	body,_ := ioutil.ReadAll(resp.Body)
	html := string(body)
	html = strings.Replace(html, " ", "", -1)
	html = strings.Replace(html, "	", "", -1)
	html = strings.Replace(html, "\n", "", -1)
	html = strings.Replace(html, "\r", "", -1)
	re := regexp.MustCompile("<tdstyle=\"width:.*?%;height:23px\"align=\"middle\">(.*?)</td>")
	res := re.FindAllStringSubmatch(html, -1)
	t := time.Now()
	hour,min,sec := t.Clock()
	year,mon,day := t.Date()
	g.date = fmt.Sprintf("%d-%d-%d %d:%d:%d",year,mon,day,hour,min,sec)
	g.buy,_ = strconv.ParseFloat(res[38][1],32)
	g.sell,_ = strconv.ParseFloat(res[39][1],32)
	g.price,_ = strconv.ParseFloat(res[40][1],32)
	g.high,_ = strconv.ParseFloat(res[41][1],32)
	g.low,_ = strconv.ParseFloat(res[42][1],32)
	return g,nil
}

func (t *Turtle) run(g Gold) {
	var act string
	var action,actionType int//action do noting 0 buy 1,sell 2, actionType : 加仓0，建仓1，盈利清仓2，止损清仓3
	var price,amount,earn,watermark float64
	var action1,actionType1 int//action do noting 0 buy 1,sell 2, actionType : 加仓0，建仓1，盈利清仓2，止损清仓3
	var price1,amount1,earn1,watermark1 float64

	times := time.Now()
	hour,min,sec := times.Clock()
	year,mon,day := times.Date()
	tr := max3(t.base.pdc-g.low, g.high-t.base.pdc, g.high-g.low)
	t.base.n = (tr +19*t.base.pdn)/20
	act = fmt.Sprintf("%d-%d-%d %d:%d:%d",year,mon,day,hour,min,sec)
	if t.sell || t.buy {
		if t.buy {
			if g.sell < t.buyData.enter-t.base.n/2 && t.total >= t.perdeal && t.buyData.unitNum < 8 {
				t.total -= t.perdeal
				t.buyData.enter = g.sell
				t.buyData.total += t.perdeal
				amount = t.perdeal / g.sell
				t.buyData.amount += amount
				t.buyData.unitNum++
				actionType = 0
				action = 1
				price = g.sell
			} else {
				if g.buy > t.base.high10 {
					t.total += g.sell * float64(t.buyData.amount)
					t.buy = false
					amount = t.buyData.amount
					t.buyData.amount = 0
					t.buyData.enter = 0
					t.buyData.total = 0
					t.buyData.unitNum = 0
					actionType = 2
					action = 1
					price = g.buy
				} else if t.buyData.enter-2*t.base.n > g.buy { //here need to check the n's meaning is firstN
					t.total += g.sell * float64(t.buyData.amount)
					t.buy = false
					amount = t.buyData.amount
					t.buyData.amount = 0
					t.buyData.enter = 0
					t.buyData.total = 0
					t.buyData.unitNum = 0
					actionType = 3
					action = 1
					price = g.sell
				}
			}
		}
		if t.sell {
			if g.buy > t.sellData.enter+t.base.n/2 && t.total >= t.perdeal && t.sellData.unitNum < 8 {
				t.total -= t.perdeal
				t.sellData.enter = g.buy
				t.sellData.total += g.buy
				amount1 = t.perdeal / g.buy
				t.sellData.amount += amount
				t.buyData.unitNum++
				actionType1 = 0
				action1 = 2
				price1 = g.buy
			} else {
				if g.sell < t.base.low10 {
					t.total += 2*t.sellData.total - g.sell*float64(t.sellData.amount)
					amount1 = t.sellData.amount
					t.sellData.total = 0
					t.sellData.amount = 0
					t.sellData.enter = 0
					t.sellData.unitNum = 0
					t.sell = false
					actionType1 = 2
					action1 = 2
					price1 = g.sell
				} else if t.sellData.enter+2*t.base.n < g.sell {
					t.total += 2*t.sellData.total - g.sell*float64(t.sellData.amount)
					amount1 = t.sellData.amount
					t.sellData.total = 0
					t.sellData.amount = 0
					t.sellData.enter = 0
					t.sellData.unitNum = 0
					t.sell = false
					actionType1 = 3
					action1 = 2
					price1 = g.sell
				}
			}
		}
	}
	 if t.sell == false || t.buy == false{
		if t.total < t.perdeal*4 {
			return //herr may need to give an alart.
		}
		if g.buy > t.base.high20 || g.sell < t.base.low20 {
			if g.buy > t.base.high20 && t.sell == false {
				t.sell = true
				amount1 = 4 * t.perdeal /g.buy
				t.sellData.amount = amount
				t.sellData.total = 4 * t.perdeal
				t.sellData.enter = g.buy
				t.sellData.unitNum = 4
				t.total -= 4*t.perdeal
				price1 = g.buy
				action1 = 2
				actionType1 = 1
			} 
			if g.sell < t.base.low20 && t.buy == false{
				t.buy = true
				amount = 4 * t.perdeal /g.sell
				t.buyData.amount = amount
				t.total -= 4*t.perdeal
				t.buyData.total = 4 * t.perdeal
				t.buyData.enter = g.sell
				t.buyData.unitNum = 4
				price = g.sell
				action = 1
				actionType = 1
			}
		}
	}
	log.Println(t.base.n)
	log.Println("Buy",act,action,actionType,price,amount,earn,watermark)
	log.Println("Sell",act,action1,actionType1,price1,amount1,earn1,watermark1)
	values[0] = t.total+t.buyData.total+t.sellData.total
	values[1] = t.total
	values[2] = t.buyData.total
	values[3] = t.buyData.enter
	values[4] = t.buyData.amount
	values[5] = t.sellData.total
	values[6] = t.sellData.enter
	values[7] = t.sellData.amount
	values[8] = g.buy
	values[9] = g.sell
	values[10] = (g.buy+g.sell)/2
	values[11] = g.high
	values[12] = g.low
	values[13] = t.base.pdn
	values[14] = t.base.n
	values[15] = tr
	values[24] = t.base.high55
	values[25] = t.base.low55
	values[26] = t.base.high20
	values[27] = t.base.low20
	values[28] = t.base.high10
	values[29] = t.base.low10
	timenow = act
	if action == 0 && action1 == 0{
		return//nothing to do, no need to log.
	}

	

	if action != 0{
		watermark = t.buyData.total
		err := writeLog(act,action,actionType,price,amount,earn,watermark)
		if err!= nil{
			log.Println(err)
		}
		go sendMail(phaseAction(act,action,actionType,price,amount,earn,watermark))
		ints[0] = 1
		ints[1] = actionType
		values[16] = price
		values[17] = amount
		values[18] = earn
		values[19] = watermark
	}
	if action1 != 0{
		watermark1 = t.sellData.total
		act = fmt.Sprintf("%d-%d-%d %d:%d:%d",year,mon,day,hour,min,sec+1)
		err := writeLog(act,action1,actionType1,price1,amount1,earn1,watermark1)
		if err!= nil{
			log.Println(err)
		}
		go sendMail(phaseAction(act,action1,actionType1,price1,amount1,earn1,watermark1))
		ints[2] = 2
		ints[3] = actionType1
		values[20] = price1
		values[21] = amount1
		values[22] = earn1
		values[23] = watermark1
	}

}

func phaseAction(act string,action,actionType int,price,amount,earn,watermark float64)string{
	mailContent := act + ":\n"
	tmp := []string{"加仓","建仓","盈利清仓","止损清仓"}
	if action == 1{
		mailContent += "buy "
	}else{
		mailContent += "sell "
	}
	mailContent += tmp[actionType] + "\n"
	mailContent += fmt.Sprintf("操作价格:%f ",price)
	mailContent += fmt.Sprintf(" 数目：%f", amount)
	mailContent += fmt.Sprintf("目前仓位：%f\n" , watermark)
	return mailContent
}

func NewTurtle()*Turtle{
	t := &Turtle{
		perdeal:1000.0,
		total:20000.0,
	}
	return t
}

func(t * Turtle)updateBase(pdc float64){
	g,err := readGoldDay(60)
	if err!=nil{
		log.Println(err)
		return 
	}
	t.base.high10,t.base.low10 = getHighLow(0,10,g)
	t.base.high20,t.base.low20 = getHighLow(0,20,g)
	t.base.high55,t.base.low55 = getHighLow(0,55,g)
	firstN :=0.0
	if t.base.pdn == 0{//means the new turtle, need to caculate new base data.
		for i := 1; i <= 20; i++ {
		//  pdc = g[idx].open
			firstN += max3(g[i-1].open-g[i].low, g[i].high-g[i-1].open, g[i].high-g[i].low)
		}
		firstN += max3(pdc - g[0].low,g[0].high-pdc,g[0].high-g[0].low)
		t.base.pdn = firstN/20
	}else{
		t.base.pdn = t.base.n
	}
	t.base.pdc = pdc
}

func(t * Turtle)saveTurtle(){
	file,err := os.Create("tutle.dat")
	if err!= nil{
		log.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(file,"%f %f %f %d\n",t.buyData.amount,t.buyData.total,t.buyData.enter,t.buyData.unitNum)
	fmt.Fprintf(file,"%f %f %f %d\n",t.sellData.amount,t.sellData.total,t.sellData.enter,t.sellData.unitNum)
	fmt.Fprintf(file,"%f %f %f %f %f %f %f %f %f %f\n",t.base.n,t.base.tr,t.base.pdc,t.base.pdn,t.base.high10,t.base.low10,t.base.high20,t.base.low20,t.base.high55,t.base.low55)
	fmt.Fprintf(file,"%f\n",t.perdeal)
	fmt.Fprintf(file,"%f\n",t.total)
	fmt.Fprintln(file,t.buy)
	fmt.Fprintln(file,t.sell)
}

func(t * Turtle)readTurtle(){
	file,err := os.Open("tutle.dat")
	if err!= nil{
		log.Println(err)
		return
	}
	defer file.Close()
	n,err := fmt.Fscanf(file,"%f %f %f %d\n",&t.buyData.amount,&t.buyData.total,&t.buyData.enter,&t.buyData.unitNum)
	fmt.Println(n,err)
	fmt.Fscanf(file,"%f %f %f %d\n",&t.sellData.amount,&t.sellData.total,&t.sellData.enter,&t.sellData.unitNum)
	fmt.Fscanf(file,"%f %f %f %f %f %f %f %f %f %f\n",&t.base.n,&t.base.tr,&t.base.pdc,&t.base.pdn,&t.base.high10,&t.base.low10,&t.base.high20,&t.base.low20,&t.base.high55,&t.base.low55)
	fmt.Fscanf(file,"%f\n",&t.perdeal)
	fmt.Fscanf(file,"%f\n",&t.total)
	fmt.Fscanln(file,&t.buy)
	fmt.Fscanln(file,&t.sell)
}

func sendMail(content string) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", from, pwd, server)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{dest}
	msg := []byte("To: nachtz<"+dest+">\r\n" +
        "From: Gold<"+from+">\r\n"+
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		content + "\r\n")
	err := smtp.SendMail(server +":"+ strconv.Itoa(port), auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func initMail(path string){
    file,err := os.Open(path)
    if err!=nil{
        log.Println(err)
    }
    fmt.Fscanf(file,"%s\n",&from)
    fmt.Fscanf(file,"%s\n",&dest)
    fmt.Fscanf(file,"%s\n",&pwd)
    fmt.Fscanf(file,"%s\n",&server)
    fmt.Fscanf(file,"%d",&port)
}

func mainp() {
	g := getGoldData("gold.txt")
	log.Println(len(g))
	runTurtle(g)
	//test()
}
