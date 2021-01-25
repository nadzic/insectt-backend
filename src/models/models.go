package models

// Measurement schema of the raw_data_t table
type Measurement struct {
    ID                  int64    `json:"id"`
    MeasuredAt          string   `json:"measured_at"`
    SignalTypeID        int      `json:"signal_type_id"`
    SignalValue         float32  `json:"signal_value"`
    MeasurementID       int8     `json:"measurement_id"`
}