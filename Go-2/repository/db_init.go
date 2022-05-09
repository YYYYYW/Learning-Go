package repository

import (
	"bufio"
	"encoding/json"
	"os"
)

var (
	topicIndexMap map[int64]*Topic
	postIndexMap  map[int64][]*Post
)

func Init(filepath string) error {
	if err := initTopicIndexMap(filepath); err != nil {
		return err
	}
	if err := initPostIndexMap(filepath); err != nil {
		return err
	}
	return nil
}

func initTopicIndexMap(filepath string) error {
	open, err := os.Open(filepath + "topic")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(open)
	topicTmpMap := make(map[int64]*Topic)
	for scanner.Scan() {
		text := scanner.Text()
		var topic Topic
		if err := json.Unmarshal([]byte(text), &topic); err != nil {
			return err
		}
		topicTmpMap[topic.Id] = &topic
	}
	topicIndexMap = topicTmpMap
	return nil
}

func initPostIndexMap(filepath string) error {
	open, err := os.Open(filepath + "post")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(open)
	postTmpMap := make(map[int64][]*Post)
	for scanner.Scan() {
		text := scanner.Text()
		var post Post
		println(text)
		if err := json.Unmarshal([]byte(text), &post); err != nil {
			return err
		}
		posts, ok := postTmpMap[post.ParentId]
		if !ok {
			postTmpMap[post.ParentId] = []*Post{&post}
			continue
		}
		posts = append(posts, &post)
		postTmpMap[post.ParentId] = posts
	}
	postIndexMap = postTmpMap
	return nil
}
