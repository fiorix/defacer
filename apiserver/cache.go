package apiserver

import (
	"image"
	"sync"
	"time"
)

// imageCache holds resized overlay images indexed by their width and height.
type imageCache struct {
	sync.Mutex
	m    map[image.Point]*imageCacheItem
	once sync.Once
}

type imageCacheItem struct {
	Time  time.Time
	Image image.Image
}

var defaultImageCache = &imageCache{
	m: make(map[image.Point]*imageCacheItem),
}

func (ic *imageCache) Get(size *image.Point) image.Image {
	ic.once.Do(func() { go ic.flush() })
	ic.Lock()
	defer ic.Unlock()
	item := ic.m[*size]
	if item == nil {
		return nil
	}
	item.Time = time.Now()
	defaceImageCache.WithLabelValues("hits").Inc()
	return item.Image
}

func (ic *imageCache) Set(size *image.Point, m image.Image) {
	ic.Lock()
	ic.m[*size] = &imageCacheItem{Time: time.Now(), Image: m}
	ic.Unlock()
	defaceImageCache.WithLabelValues("items").Inc()
}

// flush runs every 5s to evict items that are inactive for up to 5min.
func (ic *imageCache) flush() {
	for range time.Tick(5 * time.Second) {
		ic.Lock()
		for k, v := range ic.m {
			if time.Since(v.Time) > 5*time.Minute {
				delete(ic.m, k)
			}
		}
		ic.Unlock()
	}
}
