package handlers

import "time"

type DataHandler interface {
	Write(measurement string, fields map[string]interface{}, tags map[string]string, timestamp *time.Time) error
	Close()
}
