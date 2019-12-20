package ratelimiter

import "time"

//漏桶模式限流

type LeakyBucket struct {
	Size uint16
	Fill float64
	LeakInterval time.Duration
	LastUpdate time.Time
	Now func()time.Time
}

func NewLeakyBucket(size uint16, leakInterval time.Duration) *LeakyBucket {
	bucket := LeakyBucket{
		Size:         size,
		Fill:         0,
		LeakInterval: leakInterval,
		Now:          time.Now,
		LastUpdate:   time.Now(),
	}
	return &bucket
}

func (b *LeakyBucket) updateFill() {
	now := b.Now()
	if b.Fill > 0 {
		elapsed := now.Sub(b.LastUpdate)

		b.Fill -= float64(elapsed) / float64(b.LeakInterval)
		if b.Fill < 0 {
			b.Fill = 0
		}
	}
	b.LastUpdate = now
}

func (b *LeakyBucket) Pour(amount uint16) bool {
	b.updateFill()

	var newfill = b.Fill + float64(amount)

	if newfill > float64(b.Size) {
		return false
	}

	b.Fill = newfill

	return true
}
// The time at which this bucket will be completely drained
func (b *LeakyBucket) DrainedAt() time.Time {
	return b.LastUpdate.Add(time.Duration(b.Fill * float64(b.LeakInterval)))
}
// The duration until this bucket is completely drained
func (b *LeakyBucket) TimeToDrain() time.Duration {
	return b.DrainedAt().Sub(b.Now())
}

func (b *LeakyBucket) TimeSinceLastUpdate() time.Duration {
	return b.Now().Sub(b.LastUpdate)
}

type LeakyBucketSer struct {
	Size         uint16
	Fill         float64
	LeakInterval time.Duration // time.Duration for 1 unit of size to leak
	Lastupdate   time.Time
}

func (b *LeakyBucket) Serialise() *LeakyBucketSer {
	bucket := LeakyBucketSer{
		Size:         b.Size,
		Fill:         b.Fill,
		LeakInterval: b.LeakInterval,
		Lastupdate:   b.LastUpdate,
	}

	return &bucket
}

func (b *LeakyBucketSer) DeSerialise() *LeakyBucket {
	bucket := LeakyBucket{
		Size:         b.Size,
		Fill:         b.Fill,
		LeakInterval: b.LeakInterval,
		LastUpdate:   b.Lastupdate,
		Now:          time.Now,
	}

	return &bucket
}
