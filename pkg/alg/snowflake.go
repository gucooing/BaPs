package alg

import (
	"sync"
	"time"
)

// 雪花算法的基本实现

// snowflake ID 是一个64位的int数据 由四部分组成
// A-B-C-D
// A 1位 最高位不使用
// B 41位 时间戳
// C 10位 节点ID
// D 12位 毫秒内序列号
// 时间戳 节点ID 毫秒内序列号 位数可按需调整

const (
	workerBits  uint8 = 10                      // 节点ID位数 2^10=1024
	numberBits  uint8 = 12                      // 毫秒内序列号位数 2^12=4096
	workerMax   int64 = -1 ^ (-1 << workerBits) // 节点ID最大值
	numberMax   int64 = -1 ^ (-1 << numberBits) // 毫秒内序列号最大值
	timeShift         = workerBits + numberBits // 时间戳向左偏移量
	workerShift       = numberBits              // 节点ID向左偏移量
	/*
	* 原始算法使用41位字节作为时间戳数值
	* 大约68年也就是2038年就会用完
	* 这里做个偏移以增加可用时间
	* !!!这个一旦定义且开始生成ID后千万不要改了!!!
	* 不然可能会生成相同的ID
	 */
	epoch int64 = 1657148827000 // 2022-07-07 07:07:07
)

type SnowflakeWorker struct {
	lock      sync.Mutex // 互斥锁
	timestamp int64      // 记录时间戳
	workerId  int64      // 节点ID
	number    int64      // 当前毫秒内已经生成的ID序列号 从0开始累加
}

func NewSnowflakeWorker(workerId int64) *SnowflakeWorker {
	if workerId < 0 || workerId > workerMax {
		// worker id error
		return nil
	}
	worker := &SnowflakeWorker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}
	return worker
}

func (s *SnowflakeWorker) GenId() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 当前毫秒时间戳
	now := time.Now().UnixNano() / 1e6
	if s.timestamp > now {
		// 发生了时钟回拨
		if s.timestamp-now > 1000 {
			// 时钟回拨太严重
			return -1
		}
		for now <= s.timestamp {
			// 自旋等待当前时间超过上一次ID生成的时间
			now = time.Now().UnixNano() / 1e6
		}
	}
	if s.timestamp == now {
		s.number++
		if s.number > numberMax {
			// 当前毫秒内生成ID数量超过限制
			for now <= s.timestamp {
				// 自旋等待
				now = time.Now().UnixNano() / 1e6
			}
		}
	}
	if s.timestamp < now {
		// 新的毫秒到来重置序列号和时间戳
		s.number = 0
		s.timestamp = now
	}
	// 生成ID
	id := (now-epoch)<<timeShift | (s.workerId << workerShift) | (s.number)
	return id
}
