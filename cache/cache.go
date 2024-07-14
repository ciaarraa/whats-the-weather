package cache

import (
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"time"

	"git.mills.io/prologic/bitcask"
	"github.com/google/uuid"
)

type Cache struct {
	location    string
	cacheFolder string
	database    *bitcask.Bitcask
}

func (cache *Cache) open() {
	db, err := bitcask.Open(cache.location)
	if err != nil {
		fmt.Print(err)
	}
	cache.database = db
}

func (cache *Cache) Reopen() {
	err := cache.database.Reopen()
	if err != nil {
		fmt.Print(err)
	}
}

func (cache *Cache) close() {
	cache.database.Sync()
	cache.database.Close()
}

func NewCache(cacheLocation string, cacheFolder string) *Cache {
	checkCacheFolder(cacheFolder)
	return &Cache{location: cacheLocation, cacheFolder: cacheFolder}
}

func (cache *Cache) Add(object []byte, key string) {
	cache.open()
	defer cache.close()
	hashKey := getHashKey(key)
	if cache.database.Has([]byte(hashKey)) {
		cache.close()
		return
	}
	id := uuid.New().String()
	if _, err := os.Stat(cache.cacheFolder + "/" + id); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(cache.cacheFolder + "/" + id)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(cache.cacheFolder+"/"+id, []byte(object), 0644)
		fmt.Println()
		if err != nil {
			fmt.Println(err)
		}
	}
	err := cache.database.Put([]byte(hashKey), []byte(id))
	if err != nil {
		fmt.Print(err)
	}
}
func (cache *Cache) AddWithTTL(object []byte, key string, ttl time.Duration) {
	cache.open()
	defer cache.close()
	hashKey := getHashKey(key)
	if cache.database.Has([]byte(hashKey)) {
		cache.close()
		return
	}
	id := uuid.New().String()
	if _, err := os.Stat(cache.cacheFolder + "/" + id); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(cache.cacheFolder + "/" + id)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(cache.cacheFolder+"/"+id, []byte(object), 0644)
		fmt.Println()
		if err != nil {
			fmt.Println(err)
		}
	}
	err := cache.database.PutWithTTL([]byte(hashKey), []byte(id), ttl)
	if err != nil {
		fmt.Print(err)
	}
}

func (cache *Cache) Clean() {
	err := os.RemoveAll(cache.cacheFolder)
	if err != nil {
		log.Println(err)
	}
	err = os.RemoveAll(cache.location)
	if err != nil {
		log.Println(err)
	}
}

func (cache *Cache) Get(key string) (string, error) {
	cache.open()
	defer cache.close()
	hashKey := getHashKey(key)
	val, err := cache.database.Get([]byte(hashKey))
	if err == bitcask.ErrKeyNotFound || err == bitcask.ErrKeyExpired {
		return "", err
	}
	object, err := os.ReadFile(cache.cacheFolder + "/" + string(val))
	if err != nil {
		log.Fatalln(err)
	}
	return string(object), nil
}

func checkCacheFolder(cacheFolder string) {
	if _, err := os.Stat(cacheFolder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(cacheFolder, os.ModePerm)
		if err != nil {
			fmt.Println("An error has occured: ", err)
		}
	}
}

func getHashKey(key string) string {
	keyHash := fnv.New64()
	keyHash.Write([]byte(key))
	return fmt.Sprintf("%v", keyHash.Sum64())
}

func keyInCache(db *bitcask.Bitcask, key string) bool {
	key = getHashKey(key)
	val := db.Has([]byte(key))
	return val
}
