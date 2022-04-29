package main

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10
	serialBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	serialMax   int64 = -1 ^ (-1 << serialBits)
	timeShift   uint8 = workerBits + serialBits
	workerShift uint8 = serialBits
	epoch       int64 = 1650399601 // 2022.0420.4:20:1
)

/*
將 -1 左移 10 位和 12 位後在將兩個數值做 XOR 就會得到 workerMax跟 serialMax。
取得 ID 的時後不再同一個時間點就會重製 serial 並更新 timestamp
若在同一時間點的話會一直累加 serial 直到超出最大數量(4092)。
*/
type Worker struct {
	mulock    sync.Mutex
	timestamp int64
	workId    int64
	serial    int64
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("Worker ID is over maxium :" + strconv.FormatInt(workerMax, 10))
	}
	return &Worker{
		timestamp: 0,
		workId:    workerId,
		serial:    0,
	}, nil
}

func (w *Worker) Generate() int64 {
	w.mulock.Lock()
	defer w.mulock.Lock()

	now := time.Now().UnixNano() / 1000000
	if w.timestamp == now {
		w.serial++
		if w.serial > serialMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		w.serial = 0
		w.timestamp = now
	}

	ID := int64((now-epoch)<<timeShift | (w.workId << workerShift) | w.serial)
	return ID
}
