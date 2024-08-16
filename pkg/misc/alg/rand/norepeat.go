package rand

import (
	"errors"
	"math"
)

//	快速不重复随机数
//	uint32  取值范围[0～4,294,967,295]
//	100,000,000 个数的内存开销约 400MB
type NoRepeatRand struct {
	Index  uint32   // 随机标号,标记还剩下多少个可选数
	Seed   uint32   // 随机种子
	Count  uint32   // 总数
	_rand  Rand     // 随机算法
	_array []uint32 // 随机数组
}

// [start,end]
func NewNoRepeatRand(start, end, seed uint32) *NoRepeatRand {
	if start > end {
		start, end = end, start
	}

	index := end - start

	r := &NoRepeatRand{index, seed, index + 1, Rand{seed}, make([]uint32, index+1)}

	for i := uint32(0); i < index+1; i++ {
		r._array[i] = start + i
	}

	return r
}

// [start,end]   skip 略过多少个
func NewNoRepeatRandEx(start, end, seed, skip uint32) *NoRepeatRand {
	r := NewNoRepeatRand(start, end, seed)

	if skip > 0 {
		count := uint32(math.Abs(float64(end-start))) + 1
		if skip > count {
			skip = count
		}
		r.Random(skip)
	}

	return r
}

// 获取剩余随机个数
func (r *NoRepeatRand) GetRemainCount() uint32 {
	return r.Index + 1
}

// 获取随机数池总数
func (r *NoRepeatRand) GetCount() uint32 {
	return r.Count
}

// 随机 count 个数
func (r *NoRepeatRand) Random(count uint32) ([]uint32, error) {
	if count > r.Index+1 || uint32(len(r._array)) < r.Index+1 {
		return nil, errors.New("not enough random numbers")
	}

	out := make([]uint32, count)

	for i := uint32(0); i < count; i++ {
		index := r._rand.RandAB(0, r.Index) // [A,B] 从剩下的随机数里生成

		out[i] = r._array[index]

		r._array[index] = r._array[r.Index]

		r.Index--
	}

	// 减少内存开销
	r._array = r._array[:r.Index+1]

	return out, nil
}

// 添加新的随机数，已经被使用过的随机数重复利用
func (r *NoRepeatRand) Append(array []uint32) {
	r._array = append(r._array, array...)
	r.Index += uint32(len(array))
}
