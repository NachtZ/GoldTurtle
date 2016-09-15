package main

import (
    "fmt"
    "log"
    "net/http"
)

const (
    pageTop    = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>Gold Crawler && Turtle</title>
<body><h3>Result</h3>
<p>Show the result of Gold Crawler and Turtle</p>`
    form       = `<form action="/" method="POST">
<input type="submit" value="Fresh">
</form>`
    pageBottom = `</body></html>`
    anError    = `<p class="error">%s</p>`
)

type statistics struct {
    numbers []float64
    mean    float64
    median  float64
}

func startHttp() {
    http.HandleFunc("/", homePage)
    if err := http.ListenAndServe(":9001", nil); err != nil {
        log.Fatal("failed to start server", err)
    }
}

func homePage(writer http.ResponseWriter, request *http.Request) {
    
    err := request.ParseForm() // Must be called before writing response
    fmt.Fprint(writer, pageTop, form)
    if err != nil {
        fmt.Fprintf(writer, anError, err)
    } else {
        if  message, _ := processRequest(request); true {
            values := make([]float64,22)
            ints := make([]int,4)
            for i:=0;i<22;i++{
                values[i] = 1.0*float64(i)
            }
            for i:=0;i<4;i++{
                ints[i] = i
            }
            fmt.Fprint(writer, formatStats())
        } else if message != "" {
            fmt.Fprintf(writer, anError, message)
        }
    }
    fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ( string, bool) {
    return  "", true
}

func formatStats() string {
    table1 := fmt.Sprintf(`<br>Time Now: %s</br><table border="1">
<tr><th colspan="4">Turtle States</th></tr>
<tr><td>总资金</td><td>%v</td>
<td>余额</td><td>%f</td></tr>
<tr><td>买入总资金</td><td>%f</td>
<td>买入均价</td><td>%f</td></tr>
<tr><td>买入数额</td><td>%f</td>
<td>卖出总资金</td><td>%f</td></tr>
<tr><td>卖出均价</td><td>%f</td>
<td>卖出数额</td><td>%f</td></tr>
</table><br></br><table border="1">
<tr><th colspan="4">Gold States</th></tr>
<tr><td>银行买入价</td><td>%f</td>
<td>银行卖出价</td><td>%f</td></tr>
<tr><td>中间价</td><td>%f</td>
<td>最高中间价</td><td>%f</td></tr>
<tr><td>最低中间价</td><td>%f</td>
<td>PDN</td><td>%f</td></tr>
<tr><td>N</td><td>%f</td>
<td>TR</td><td>%f</td></tr>
</table>`, timenow,values[0],values[1],values[2],values[3],values[4],values[5],values[6],values[7],values[8],values[9],values[10],values[11],values[12],values[13],values[14],values[15])
    table1 += fmt.Sprintf(`<br></br><table border="1">
<tr><th colspan="5">Action Log</th></tr>
<tr><td>买/卖</td>
<td>操作类型</td>
<td>操作价格</td>
<td>操作数目</td>
<td>操作盈利</td>
<td>目前仓位</td></tr>
<tr><td>%d</td>
<td>%d</td>
<td>%f</td>
<td>%f</td>
<td>%f</td>
<td>%f</td></tr>
<tr><td>%d</td>
<td>%d</td>
<td>%f</td>
<td>%f</td>
<td>%f</td>
<td>%f</td></tr>
</table>(注：0表示没有操作，1表示买入，2表示卖出。加仓0，建仓1，盈利清仓2，止损清仓3。)`,ints[0],ints[1],values[16],values[17],values[18],values[19],ints[2],ints[3],values[20],values[21],values[22],values[23])
fmt.Println(values)
fmt.Println(ints)
return table1;
}


