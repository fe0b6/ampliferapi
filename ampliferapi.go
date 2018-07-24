package ampliferapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const host = "https://amplifr.com"

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

func (api *API) rq(link string, params map[string]string) (res rqAns, err error) {
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

	log.Println(string(content))

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

	return
}
