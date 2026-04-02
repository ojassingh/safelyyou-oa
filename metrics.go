package main

import "time"

func calculateStats(heartbeats []time.Time, uploadTimes []int64) (GetDeviceStatsResponse, error) {
	if len(heartbeats) == 0 || len(uploadTimes) == 0 {
		return GetDeviceStatsResponse{}, ErrNoStatsAvailable
	}

	firstMinute := heartbeats[0].UTC().Unix() / 60
	lastMinute := firstMinute

	for _, heartbeat := range heartbeats[1:] {
		minute := heartbeat.UTC().Unix() / 60
		if minute < firstMinute {
			firstMinute = minute
		}
		if minute > lastMinute {
			lastMinute = minute
		}
	}

	totalMinutes := lastMinute - firstMinute
	uptime := 100.0
	if totalMinutes > 0 {
		uptime = (float64(len(heartbeats)) / float64(totalMinutes)) * 100
	}

	var totalUploadTime int64
	for _, uploadTime := range uploadTimes {
		totalUploadTime += uploadTime
	}

	return GetDeviceStatsResponse{
		AvgUploadTime: time.Duration(totalUploadTime / int64(len(uploadTimes))).String(),
		Uptime:        uptime,
	}, nil
}
