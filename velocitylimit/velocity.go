package velocity

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type parameters struct {
	MaxLoadPerDay  float64
	MaxTrxPerDay   int
	MaxLoadPerWeek int
}

// Array of transactions
type transactions struct {
	Transactions []transaction
}

// Transaction def
type transaction struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	LoadAmount string `json:"load_amount"`
	Time       string `json:"time"`
}

type transactionKey1 struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
}

type transactionKey struct {
	CustomerID string `json:"customer_id"`
	Time       string `json:"time"`
}

type response struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
}

var params = parameters{ // b == Student{"Bob", 0}
	MaxLoadPerDay:  5000.00,
	MaxTrxPerDay:   3,
	MaxLoadPerWeek: 20000.00,
}

var mTrx = make(map[transactionKey1]int)        // record of transactions
var mTrxDay = make(map[transactionKey]int)      // a struct map counting the number of times the customer has made a transaction within a day
var mTrxLoad = make(map[transactionKey]float64) // a struct map reducing load per customer within a day
var mTrxWeek = make(map[transactionKey]int)     // a struct map reducing load per customer within a day

//Limit - @trx  transaction
func Limit(trx string) string {
	var rtnJSON []byte
	var acceptTrx = false // transaction flag

	byt := []byte(trx)                                        // decode JSON to Go values
	transaction := transaction{}                              // variable to store decoded data
	if err := json.Unmarshal(byt, &transaction); err != nil { // actual decode and error check
		panic(err)
	}

	// convert into Go time
	t, _ := time.Parse(time.RFC3339, transaction.Time)
	day := t.Day()
	month := t.Month()
	intMonth := int(month)
	//fmt.Println((intMonth))
	year := t.Year()
	strDate := []string{strconv.Itoa(year), strconv.Itoa(intMonth), strconv.Itoa(day)}
	strings.Join(strDate, "")
	timestamp := strings.Join(strDate, "") // concatenate timestamp to be used as part of map key

	// get week number
	yearWeek, week := t.ISOWeek()
	strDateWeek := []string{strconv.Itoa(yearWeek), strconv.Itoa(week)}
	strings.Join(strDateWeek, "")
	timestampWeek := strings.Join(strDateWeek, "")
	//fmt.Println(yearWeek, week)

	// convert load amount to float
	loadAmount := strings.Replace(transaction.LoadAmount, "$", "", 1)
	var floatLoadAmount1 float64
	if floatLoadAmount, err := strconv.ParseFloat(loadAmount, 32); err == nil {
		floatLoadAmount1 = floatLoadAmount
	}

	// ignore duplicate transaction id's with the same customer
	if mTrx[transactionKey1{ID: transaction.ID, CustomerID: transaction.CustomerID}] == 0 {
		mTrx[transactionKey1{ID: transaction.ID, CustomerID: transaction.CustomerID}]++
		acceptTrx = true
	} else {
		return ""
	}

	// check mTrxLoad for amounts < 5000.00 load amount per day, check mTrxDay < 3, check mTrxWeek for amounts < 20000.00
	switch {
	case mTrxLoad[transactionKey{CustomerID: transaction.CustomerID, Time: timestamp}]+floatLoadAmount1 <= params.MaxLoadPerDay && mTrxDay[transactionKey{CustomerID: transaction.CustomerID, Time: timestamp}] <= params.MaxTrxPerDay && mTrxWeek[transactionKey{CustomerID: transaction.CustomerID, Time: timestamp}] <= params.MaxLoadPerWeek:
		// add load amount to hashmap by CustomerID and timestamp to count max load per day
		mTrxLoad[transactionKey{CustomerID: transaction.CustomerID, Time: timestamp}] += floatLoadAmount1
		//fmt.Println("mTrxLoad", mTrxLoad)

		mTrxWeek[transactionKey{CustomerID: transaction.CustomerID, Time: timestampWeek}]++
		//fmt.Println("mTrxWeek", mTrxWeek)

		// add transaction to hashmap by CustomerID and timestamp to count number of transactions per day
		mTrxDay[transactionKey{CustomerID: transaction.CustomerID, Time: timestamp}]++
		//fmt.Println("mTrxDay", mTrxDay)

		acceptTrx = true

	default:
		acceptTrx = false
	}

	data := &response{
		ID:         transaction.ID,
		CustomerID: transaction.CustomerID,
		Accepted:   acceptTrx,
	}

	rtnJSON, _ = json.Marshal(data)

	//fmt.Println(string(rtnJSON))

	return string(rtnJSON)
}
