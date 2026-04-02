package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func loadDeviceIDs(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	deviceIDs := make([]string, 0, len(rows))
	for i, row := range rows {
		if i == 0 || len(row) == 0 {
			continue
		}

		deviceID := strings.TrimSpace(row[0])
		if deviceID != "" {
			deviceIDs = append(deviceIDs, deviceID)
		}
	}

	return deviceIDs, nil
}
