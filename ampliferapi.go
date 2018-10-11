package ampliferapi

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/fe0b6/tools"
	"github.com/syndtr/goleveldb/leveldb"
)

const host = "https://amplifr.com"

var (
	cacheHandler *leveldb.DB
	cacheTime    int
)

// InitCache - инициализация кэша
func InitCache(h *leveldb.DB, ex int) {
	cacheHandler = h
	cacheTime = ex
}

// GetProjects - Получаем список проектов
func (api *API) GetProjects() (ans GetProjectsAns, err error) {

	res, err := api.rq("/api/v1/projects", map[string]string{})
	if err != nil {
		log.Println("[error]", err)
		return
	}

	err = json.Unmarshal(res.Result, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// GetProjectPosts - Получаем список постов проекта
func (api *API) GetProjectPosts(projectID int64, params map[string]string) (ans GetProjectPostsAns, err error) {

	res, err := api.rq(fmt.Sprintf("/api/v1/projects/%d/posts", projectID), params)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	err = json.Unmarshal(res.Result, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// GetProjectPostStats - Получаем статистику поста
func (api *API) GetProjectPostStats(projectID int64, postID int, params map[string]string) (ans GetProjectPostStatsAns, err error) {

	res, err := api.rq(fmt.Sprintf("/api/v1/projects/%d/stats/%d", projectID, postID), params)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	err = json.Unmarshal(res.Result, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// GetProjectStats - Получаем статистику проекта
func (api *API) GetProjectStats(projectID int64, params map[string]string) (ans GetProjectStatsAns, err error) {

	res, err := api.rq(fmt.Sprintf("/api/v1/projects/%d/stats", projectID), params)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	err = json.Unmarshal(res.Result, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// GetProjectStatsByPost - Получаем статистику постов проекта
func (api *API) GetProjectStatsByPost(projectID int64, params map[string]string) (ans GetProjectStatsAns, err error) {

	res, err := api.rq(fmt.Sprintf("/api/v1/projects/%d/stats/by_posts", projectID), params)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	err = json.Unmarshal(res.Result, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

func (api *API) rq(link string, params map[string]string) (res rqAns, err error) {
	var cacheKey []byte
	if cacheTime > 0 {
		cacheKey = getCacheKey(link, params)
	}

	// Чекаем кэш
	if len(cacheKey) > 0 {
		var b []byte
		b, err = getCache(cacheKey)
		if err != nil && err.Error() != "leveldb: not found" {
			log.Println("[error]", err)
		} else if len(b) > 0 {
			tools.FromGob(&res, b)
			return
		}
	}

	q := url.Values{}
	q.Add("access_token", api.AccessToken)
	for k, v := range params {
		q.Add(k, v)
	}

	resp, err := http.Get(host + link + "?" + q.Encode())
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println("[error]", err)
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Status: %s, Code: %d", resp.Status, resp.StatusCode)
		log.Println("[error]", err)
		return
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	err = json.Unmarshal(content, &res)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	if !res.OK {
		err = errors.New("bad ans")
		log.Println("[error]", string(content))
		return
	}

	// Если ответ хороший - кэшируем его
	if len(cacheKey) > 0 {
		err2 := setCache(cacheKey, tools.ToGob(res))
		if err2 != nil {
			log.Println("[error]", err2)
		}
	}

	return
}

func getCacheKey(link string, params map[string]string) (key []byte) {
	h := sha512.New()

	_, err := h.Write([]byte(link))
	if err != nil {
		log.Println("[error]", err)
		return
	}

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		_, err := h.Write([]byte(k + "=" + params[k]))
		if err != nil {
			log.Println("[error]", err)
			return
		}
	}

	return h.Sum(nil)
}

func getCache(cacheKey []byte) (ans []byte, err error) {
	var bt []byte
	bt, err = cacheHandler.Get(cacheKey, nil)
	if err != nil {
		if err.Error() != "leveldb: not found" {
			log.Println("[error]", err)
		}
		return
	}

	var obj cacheObj
	tools.FromGob(&obj, bt)

	if obj.Expire.Before(time.Now()) {
		cacheHandler.Delete(cacheKey, nil)
		return
	}

	ans = obj.Data
	return
}

func setCache(key []byte, b []byte) (err error) {
	err = cacheHandler.Put(key, tools.ToGob(cacheObj{
		Data:   b,
		Expire: time.Now().Add(time.Duration(cacheTime) * time.Second),
	}), nil)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}
