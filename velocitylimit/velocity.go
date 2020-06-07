package velocity

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// parameter struct for business logic
type parameters struct {
	MaxLoadPerDay  float64
	MaxTrxPerDay   int
	MaxLoadPerWeek int
}

// transaction struct for converting to usable Go values
type transaction struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	LoadAmount string `json:"load_amount"`
	Time       string `json:"time"`
}

// struct used as a key for mapping
type transactionKey1 struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
}

// struct used as a key for mapping
type transactionKey struct {
	CustomerID string `json:"customer_id"`
	Time       string `json:"time"`
}

// response struct for encoding text back to JSON
type response struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
}

// parameter for business logic
var params = parameters{
	MaxLoadPerDay:  5000.00,
	MaxTrxPerDay:   3,
	MaxLoadPerWeek: 20000.00,
}

var mTrx = make(map[transactionKey1]int)        // map with key transactionKey to count record of transactions and ignore duplicates
var mTrxDay = make(map[transactionKey]int)      // map with key transactionKey to count the number of times the customer has made a transaction within a day
var mTrxLoad = make(map[transactionKey]float64) // map with key transactionKey to count the load per customer within a day
var mTrxWeek = make(map[transactionKey]int)     // map with key transactionKey to count the load per customer within a week

// Limit - Analyze the velocity limit of a Customer based on the Transaction ID and Customer ID
//  @trx transaction with properties:
//  ID         string
//	CustomerID string
//	LoadAmount string
//  Time       string
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
	year := t.Year()

	// get date in yyyymd - month and day are not 0 padded
	strDate := []string{strconv.Itoa(year), strconv.Itoa(intMonth), strconv.Itoa(day)}
	strings.Join(strDate, "")
	timestamp := strings.Join(strDate, "") // concatenate timestamp to be used as part of map key

	// get week number
	yearWeek, week := t.ISOWeek()
	strDateWeek := []string{strconv.Itoa(yearWeek), strconv.Itoa(week)}
	strings.Join(strDateWeek, "")
	timestampWeek := strings.Join(strDateWeek, "")

	// convert load amount to float
	loadAmount := strings.Replace(transaction.LoadAmount, "$", "", 1)
	var floatLoadAmount1 float64
	if floatLoadAmount, err := strconv.ParseFloat(loadAmount, 32); err == nil {
		floatLoadAmount1 = floatLoadAmount
	}

	// ignore duplicate transaction ID's for the same customer
	if mTrx[transactionKey1{ID: transaction.ID, CustomerID: transaction.CustomerID}] == 0 {
		mTrx[transactionKey1{ID: transaction.ID, CustomerID: transaction.CustomerID}]++
		acceptTrx = true
	} else {
		return ""
	}

	// check mTrxLoad for amounts < 5000.00 load amount per day, check mTrxDay < 3, check mTrxWeek for amounts < 20000.00
	switch {
	case mTrxLoad[transactionKey{CustomerID: transaction.CustomerID, Time: timestamp}]+floatLoadAmount1 <= params.MaxLoadPerDay && mTrxDay[transactionKey{CustomerID: transaction.CustomerID, Time: timestamp}] <= params.MaxTrxPerDay && mTrxWeek[transactionKey{CustomerID: transaction.CustomerID, Time: timestamp}] <= params.MaxLoadPerWeek:
		// add load amount to bucket by CustomerID and timestamp to count max load per day
		mTrxLoad[transactionKey{CustomerID: transaction.CustomerID, Time: timestamp}] += floatLoadAmount1

		// add load amount to bucket by CustomerID and timestamp to count max load per week
		mTrxWeek[transactionKey{CustomerID: transaction.CustomerID, Time: timestampWeek}]++

		// add transaction to bucket by CustomerID and timestamp to count number of transactions per day
		mTrxDay[transactionKey{CustomerID: transaction.CustomerID, Time: timestamp}]++

		acceptTrx = true // flag transaction as accepted

	default:
		acceptTrx = false // flag transaction as not accepted
	}

	data := &response{
		ID:         transaction.ID,
		CustomerID: transaction.CustomerID,
		Accepted:   acceptTrx,
	}

	rtnJSON, _ = json.Marshal(data)

	return string(rtnJSON)
}
