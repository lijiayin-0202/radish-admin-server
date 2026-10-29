package global

import (
	"errors"
	"sync"
	"time"
)

const (
	// snowflakeEpoch 雪花算法起始时间戳：2023-01-01 00:00:00 +0800
	snowflakeEpoch = int64(1672502400000)
	// workerBits 机器ID占用位数，取值范围 0~1023
	workerBits = 10
	// seqBits 毫秒内序列号占用位数，取值范围 0~4095
	seqBits = 12
)

const (
	// workerIDMax 机器ID最大值 1023
	workerIDMax = int64(-1) ^ (int64(-1) << workerBits)
	// seqMax 毫秒内序列号最大值 4095
	seqMax = int64(-1) ^ (int64(-1) << seqBits)
	// timeShift 时间戳左移位数 22
	timeShift = workerBits + seqBits
	// workerShift 机器ID左移位数 12
	workerShift = seqBits
)

// SnowFlake 标准雪花算法ID生成器，加锁保证并发安全
type SnowFlake struct {
	sync.Mutex
	timestamp int64 // 上一次生成ID的时间戳（毫秒）
	workerID  int64 // 机器ID
	sequence  int64 // 同一毫秒内的序列号
}

// Snow 全局雪花ID生成器，由 initialize.InitSnowflake 初始化
var Snow *SnowFlake

// NewSnowFlake 创建雪花ID生成器，workerID 取值范围 0~1023
func NewSnowFlake(workerID int64) (*SnowFlake, error) {
	if workerID < 0 || workerID > workerIDMax {
		return nil, errors.New("workerID 超出范围，取值应为 0~1023")
	}
	return &SnowFlake{
		timestamp: 0,
		workerID:  workerID,
		sequence:  0,
	}, nil
}

// NextID 生成下一个雪花ID
func (s *SnowFlake) NextID() int64 {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixNano() / int64(time.Millisecond)
	if s.timestamp == now {
		// 同一毫秒内序列号自增，溢出后阻塞等待下一毫秒
		s.sequence = (s.sequence + 1) & seqMax
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		s.sequence = 0
	}
	s.timestamp = now

	return (now-snowflakeEpoch)<<timeShift | s.workerID<<workerShift | s.sequence
}
