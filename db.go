package main
import (
    "database/sql"
    _"github.com/go-sql-driver/mysql"
    "log"
    "io/ioutil"
    "strings"
 //   "strconv"
    //"time"
	"fmt"
    "os"
)

var dbaddr string



func readGoldDay(num int) ([]Gold,error){//here suggest the num is 55 for the 55 day's high or low
    golds := []Gold{}
    db, err := sql.Open("mysql",dbaddr)
    if err!= nil{
        return golds,err
    }
    defer db.Close()
    err = db.Ping()
    if err!= nil{
        return golds,err
    }
    searchGED := fmt.Sprintf("SELECT * FROM gold.goldeveryday order by getDate DESC LIMIT %d",num)
    rows,err := db.Query(searchGED)
    if err!= nil{
        return golds,err
    }
    for rows.Next(){
        tmp := new(Gold)
        err := rows.Scan(&tmp.date,&tmp.price,&tmp.open,&tmp.high,&tmp.low)
        if err != nil{
            return golds,err
        }
        golds = append(golds,*tmp)
    }
    return golds,nil
}

func writeLog(act string,action,actionType int,price,amount,earn,watermark float64)error{
    db, err := sql.Open("mysql",dbaddr)
    if err!= nil{
        return err
    }
    defer db.Close()
    err = db.Ping()
    if err!= nil{
        return err
    }
    insertIntoLog := "insert into actionLog(ActionTime,action,price,amount,actionType,earn,watermark)values(?,?,?,?,?,?,?)"
    stmt,err := db.Prepare(insertIntoLog)
    if err!= nil{
        return err
    }
    _,err = stmt.Exec(act,action,price,amount,actionType,earn,watermark)
    return err
}

func writeGoldMin(g Gold) error{
    db, err := sql.Open("mysql",dbaddr)
    if err!= nil{
        return err
    }
    defer db.Close()
    err = db.Ping()
    if err!= nil{
        return err
    }
    insertIntoGEM := "insert into GoldEveryMin(GetTime,price,buy,sell,highMid,lowMid)values(?,?,?,?,?,?)" 
    stmt,err := db.Prepare(insertIntoGEM)
    if err!=nil{
        return err
    }
    _,err = stmt.Exec(g.date,g.price,g.buy,g.sell,g.high,g.low)
    return err
}

func importDailyGold(g []Gold){
    db, err := sql.Open("mysql",dbaddr)
    if err!= nil{
        log.Println(err)
        return
    }
    defer db.Close()
    err = db.Ping()
    if err!= nil{
        log.Println(err)
        return
    }
    insertIntoGED := "insert into GoldEveryDay(getDate,price,open,high,low)values(?,?,?,?,?)"

    stmt,err := db.Prepare(insertIntoGED)
    if err != nil{
        log.Println(err)
        return
    }
    for i:= len(g)-1;i>=0;i--{
    g[i].date = strings.Replace(g[i].date,"年","-",-1)
    g[i].date = strings.Replace(g[i].date,"月","-",-1)
    g[i].date = strings.Replace(g[i].date,"日","",-1)
    _,err := stmt.Exec(g[i].date,g[i].price,g[i].open,g[i].high,g[i].low)
    if err!=nil{
        log.Println(err)
        continue
    }
}
    log.Println("Job done.")
}

func initDB() error{

    file,err := os.Open("f:/addr.txt")
    if err!=nil{
        return err
    }
    defer file.Close()
    str,_ := ioutil.ReadAll(file)
    dbaddr = string(str)//need to read file for safety
  return nil
}