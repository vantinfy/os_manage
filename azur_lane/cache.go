package azur_lane

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"sync"
)

// Cache 实现缓存和持久化功能
type Cache struct {
	mu   sync.RWMutex
	Data map[string]*Ship
	file string
}

// NewCache 创建一个新的缓存实例
func NewCache(file string) (*Cache, error) {
	cache := &Cache{
		Data: make(map[string]*Ship),
		file: file,
	}

	// 尝试加载文件中的数据
	if err := cache.load(); err != nil {
		return nil, err
	}
	return cache, nil
}

// Save 将缓存数据持久化到文件
func (c *Cache) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 将缓存数据编码为 JSON
	data, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		return err
	}
	// 单独保存一份json格式的
	_ = os.WriteFile("ships.json", data, 0644)

	// 压缩
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	_, err = gzipWriter.Write(data)
	if err != nil {
		return err
	}
	gzipWriter.Close()

	// 将数据写入文件
	return os.WriteFile(c.file, compressedData.Bytes(), 0644)
}

// Load 从文件加载数据
func (c *Cache) load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.file)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果文件不存在，说明缓存文件为空，直接返回
			return nil
		}
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// 读取解压后的数据
	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, gzipReader)
	if err != nil {
		return err
	}

	// 解码 JSON 数据到缓存中
	return json.Unmarshal(decompressedData.Bytes(), &c.Data)
}

// Set 将 ship 存入缓存
func (c *Cache) Set(name string, ship *Ship) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Data[name] = ship
}

// Get 从缓存中获取 ship
func (c *Cache) Get(name string) (*Ship, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ship, exists := c.Data[name]
	return ship, exists
}
