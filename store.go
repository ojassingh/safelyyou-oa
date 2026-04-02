package main

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrDeviceNotFound   = errors.New("device not found")
	ErrNoStatsAvailable = errors.New("device stats not available")
)

type DeviceData struct {
	Heartbeats  []time.Time
	UploadTimes []int64
}

type Store struct {
	mu      sync.RWMutex
	devices map[string]*DeviceData
}

func NewStore(deviceIDs []string) *Store {
	devices := make(map[string]*DeviceData, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		devices[deviceID] = &DeviceData{}
	}

	return &Store{devices: devices}
}

func (s *Store) AddHeartbeat(deviceID string, sentAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	device, ok := s.devices[deviceID]
	if !ok {
		return ErrDeviceNotFound
	}

	device.Heartbeats = append(device.Heartbeats, sentAt.UTC())
	return nil
}

func (s *Store) AddUploadStat(deviceID string, uploadTime int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	device, ok := s.devices[deviceID]
	if !ok {
		return ErrDeviceNotFound
	}

	device.UploadTimes = append(device.UploadTimes, uploadTime)
	return nil
}

func (s *Store) GetStats(deviceID string) (GetDeviceStatsResponse, error) {
	s.mu.RLock()
	device, ok := s.devices[deviceID]
	if !ok {
		s.mu.RUnlock()
		return GetDeviceStatsResponse{}, ErrDeviceNotFound
	}

	heartbeats := append([]time.Time(nil), device.Heartbeats...)
	uploadTimes := append([]int64(nil), device.UploadTimes...)
	s.mu.RUnlock()

	return calculateStats(heartbeats, uploadTimes)
}
