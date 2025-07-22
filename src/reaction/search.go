package reaction

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"telegram-listener/helper"
)

type Post struct {
	ID     int
	Title  string
	Slug   string
	URL    string
	Poster string
	Year   string
}

type SearchInlineMenuRow struct {
	Row []struct {
		Title string `json:"title"`
		Value string `json:"value"`
	} `json:"row"`
}

func (s *Service) searchPosts(apiURL, appURL, query string) (inlineMenu string, err error) {
	posts := []Post{}
	limit := 5

	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("q", query)
	u := fmt.Sprintf(apiURL, params.Encode())
	body, err := helper.GetURL(u)
	if err == nil {
		err = json.Unmarshal(body, &posts)
		if err != nil {
			log.Printf("api fetched but cant be unmarshalled: %s", err)
		}
	}
	if len(posts) > limit {
		posts = posts[:limit] // limit the number of posts
	}
	if len(posts) > 0 {
		menu := []SearchInlineMenuRow{}
		for _, post := range posts {
			row := SearchInlineMenuRow{}
			row.Row = append(row.Row, struct {
				Title string `json:"title"`
				Value string `json:"value"`
			}{})
			row.Row[0].Title = fmt.Sprintf("%s %s", post.Title, post.Year)
			row.Row[0].Value = fmt.Sprintf("%s%s", appURL, post.Slug)
			menu = append(menu, row)
		}
		inlineMenuBytes, err := json.Marshal(menu)
		if err != nil {
			log.Printf("Failed to marshal inline menu: %v", err)
			return "", err
		}
		inlineMenu = string(inlineMenuBytes)
	}
	//log.Println("searchPosts:", apiURL, "query:", query, "found:", len(posts), "inlineMenu:", inlineMenu)
	return
}
