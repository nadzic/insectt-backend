package websocket

import (
	"fmt"
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"
)

var (
	upgrader  = websocket.Upgrader{}
)

// Measurement model
type Measurement struct {
	ID                  int64    `json:"id"`
	MeasuredAt          string   `json:"measured_at"`
	SignalTypeID        int      `json:"signal_type_id"`
	SignalValue         float32  `json:"signal_value"`
	MeasurementID       int8     `json:"measurement_id"`
}

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
			log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
			panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
			panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// CreateMeasurement create a measurement in the postgres db
func CreateMeasurement(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty measurement of type Measurement
	var measurement Measurement

	// decode the json request to measurement
	err := json.NewDecoder(r.Body).Decode(&measurement)

	if err != nil {
			log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert measurement function and pass the measurement
	insertID := insertMeasurement(measurement)

	// format a response object
	res := response{
			ID:      insertID,
			Message: "measurement created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// insert one measurement in the DB
func insertMeasurement(measurement Measurement) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning measurement id will return the id of the inserted measurement
	sqlStatement := `INSERT INTO raw_data_t (measured_at, signal_type_id, signal_value, measurement_id) VALUES ($1, $2, $3, $4) RETURNING id`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, measurement.MeasuredAt, measurement.SignalTypeID, measurement.SignalValue, measurement.MeasurementID).Scan(&id)

	if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// DbHandler export handling requests from frontend
func DbHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

