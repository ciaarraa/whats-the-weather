package cache

import (
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"os"

	"git.mills.io/prologic/bitcask"
	"github.com/google/uuid"
)

type Cache struct {
	location string
	database *bitcask.Bitcask
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

func NewCache() *Cache {
	checkCacheFolder()
	return &Cache{location: "tmp/db"}
}

func (cache *Cache) Add(object []byte, key string) {
	cache.open()
	hashKey := getHashKey(key)
	if cache.database.Has([]byte(hashKey)) {
		cache.close()
		return
	}
	id := uuid.New().String()
	if _, err := os.Stat(".cache/" + id); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(".cache/" + id)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(".cache/"+id, []byte(object), 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

	err := cache.database.Put([]byte(hashKey), []byte(id))
	if err != nil {
		fmt.Print(err)
	}
	cache.close()
}
func (cache *Cache) Fold(f func(key []byte) error) {
	cache.database.Fold(f)
}
func (cache *Cache) Get(key string) string {
	cache.open()
	hashKey := getHashKey(key)
	if cache.database.Has([]byte(hashKey)) {
		val, err := cache.database.Get([]byte(hashKey))
		if err != nil {
			log.Fatalln(err)
		}
		object, err := os.ReadFile(".cache/" + string(val))
		if err != nil {
			log.Fatalln(err)
		}
		return string(object)
	}
	cache.close()
	return ""
}

func checkCacheFolder() {
	if _, err := os.Stat(".cache"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(".cache", os.ModePerm)
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