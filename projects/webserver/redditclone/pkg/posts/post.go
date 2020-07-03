package posts

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"redditclone/pkg/users"
)

type Vote struct {
	Author users.UserIdType `json:"user"`
	Vote   VoteType         `json:"vote"`
}

type Comment struct {
	ID      uint32      `json:"id"`
	Author  *users.User `json:"author"`
	Body    string      `json:"body"`
	Created time.Time   `json:"created"`
}

type Post struct {
	ID              PostIdType   `json:"id"`
	Author          *users.User  `json:"author"`
	PType           PostType     `json:"type"`
	Category        CategoryType `json:"category"`
	Title           string       `json:"title"`
	Text            string       `json:"text"`
	Votes           []*Vote      `json:"votes"`
	Url             string       `json:"url"`
	ratingUpCount   uint
	ratingDownCount uint
	ratingUpPercent uint       `json:"upvotePercentage"`
	Comments        []*Comment `json:"comments"`
	Score           uint64     `json:"score"`
	Views           int64      `json:"views"`
	Created         time.Time  `json:"created"`
}

type PostRep struct {
	data   map[PostIdType]*Post
	lastId uint64
	mutexx *sync.RWMutex
}

type (
	VoteType     int8
	PostType     int8
	CategoryType int8
	PostIdType   uint64
)

var (
	nameToType map[string]CategoryType = map[string]CategoryType{
		"all":         AllCategory,
		"music":       MusicCategory,
		"funny":       FunnyCategory,
		"videos":      VideosCategory,
		"programming": ProgrammingCategory,
		"news":        NewsCategory,
		"fashion":     FashionCategory,
	}
	typeToName map[CategoryType]string = map[CategoryType]string{
		AllCategory:         "all",
		MusicCategory:       "music",
		FunnyCategory:       "funny",
		VideosCategory:      "videos",
		ProgrammingCategory: "programming",
		NewsCategory:        "news",
		FashionCategory:     "fashion",
	}
	ErrNoPost    = errors.New("No post found")
	ErrNoComment = errors.New("No comment found")
)

const (
	RatingUp   = 1
	RatingDown = -1

	UnrecognizedPost = 0
	TextPost         = 1
	LinkPost         = 2

	UnrecognizedCategory = 0
	AllCategory          = 1
	MusicCategory        = 2
	FunnyCategory        = 3
	VideosCategory       = 4
	ProgrammingCategory  = 5
	NewsCategory         = 6
	FashionCategory      = 7
)

func (post *PostType) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "link":
		*post = LinkPost
	case "text":
		*post = TextPost
	default:
		*post = UnrecognizedPost
	}
	return nil
}

func (post PostType) MarshalText() ([]byte, error) {
	var data string
	switch post {
	case LinkPost:
		data = "link"
	case TextPost:
		data = "text"
	default:
		data = ""
	}
	return []byte(data), nil
}

func (category *CategoryType) UnmarshalText(text []byte) error {
	*category = NameToType(string(text))
	return nil
}

func (category CategoryType) MarshalText() ([]byte, error) {
	var data string = TypeToName(category)
	return []byte(data), nil
}

func NameToType(name string) CategoryType {
	category, ok := nameToType[name]

	if !ok {
		return UnrecognizedCategory
	}
	return category
}

func TypeToName(cat CategoryType) string {
	return typeToName[cat]
}

func (post *PostIdType) UnmarshalText(text []byte) error {
	data, _ := strconv.Atoi(string(text))
	*post = PostIdType(data)
	return nil
}

func (post PostIdType) MarshalText() ([]byte, error) {
	var data string = fmt.Sprintf(`%d`, uint32(post))
	return []byte(data), nil
}

func NewPost() *Post {
	return &Post{Votes: []*Vote{}, Comments: []*Comment{}}
}

func NewPostRep() *PostRep {
	return &PostRep{
		data:   make(map[PostIdType]*Post, 0),
		mutexx: &sync.RWMutex{},
	}
}

func (post *Post) GetCommentIndex(id uint32) (int, bool) {
	for index, comment := range post.Comments {
		if comment.ID == id {
			return index, true
		}
	}
	return 0, false
}

func (post *Post) GetComment(id uint32) (*Comment, bool) {
	for _, comment := range post.Comments {
		if comment.ID == id {
			return comment, true
		}
	}
	return nil, false
}

func (post *Post) AddComment(comment *Comment) (*Comment, error) {
	var lastId uint32 = 1

	for _, comm := range post.Comments {
		if lastId < comm.ID {
			lastId = comm.ID
		}
	}

	comment.ID = lastId
	post.Comments = append(post.Comments, comment)
	return comment, nil
}

func (post *Post) UpdateComment(id uint32, comment *Comment) (*Comment, error) {
	index, ok := post.GetCommentIndex(id)

	if !ok {
		return nil, ErrNoComment
	}

	comment.ID = post.Comments[index].ID
	post.Comments[index] = comment
	return comment, nil
}

func (post *Post) DeleteComment(id uint32) (*Comment, error) {
	index, ok := post.GetCommentIndex(id)

	if !ok {
		return nil, ErrNoComment
	}

	comment := post.Comments[index]
	post.Comments = removeComment(post.Comments, index)
	return comment, nil
}

func (pr *PostRep) DeleteComment(id PostIdType, cid uint32) (*Comment, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	post, err := pr.getPost(id)

	if err != nil {
		return nil, err
	}
	return post.DeleteComment(cid)
}

func (pr *PostRep) AddComment(id PostIdType, comment *Comment) (*Comment, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	post, err := pr.getPost(id)

	if err != nil {
		return nil, err
	}
	return post.AddComment(comment)
}

func (pr *PostRep) GetComment(id PostIdType, cid uint32) (*Comment, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	post, err := pr.getPost(id)

	if err != nil {
		return nil, err
	}

	comment, has := post.GetComment(cid)

	if !has {
		return nil, ErrNoComment
	}
	return comment, nil
}

func (pr *PostRep) UpdateComment(id PostIdType, cid uint32, comment *Comment) (*Comment, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	post, err := pr.getPost(id)

	if err != nil {
		return nil, err
	}
	return post.UpdateComment(cid, comment)
}

func removeComment(s []*Comment, i int) []*Comment {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (pr *PostRep) GetCategoryPosts(category CategoryType) ([]*Post, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	posts := make([]*Post, 0)

	for _, p := range pr.data {
		if p.Category == category {
			posts = append(posts, p)
		}
	}
	return posts, nil
}

func (pr *PostRep) GetUserPosts(uid users.UserIdType) ([]*Post, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	posts := make([]*Post, 0)

	for _, p := range pr.data {
		if p.Author.ID == uid {
			posts = append(posts, p)
		}
	}
	return posts, nil
}

func (pr *PostRep) DeletePost(id PostIdType) (*Post, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	post, err := pr.getPost(id)

	if err != nil {
		return post, err
	}

	delete(pr.data, id)
	return post, nil
}

func (post *Post) GetVote(id users.UserIdType) (*Vote, bool) {
	for _, v := range post.Votes {
		if v.Author == id {
			return v, true
		}
	}
	return nil, false
}

func (post *Post) AddVote(v *Vote) (*Vote, error) {
	if vote, ok := post.GetVote(v.Author); ok {
		if vote.Vote == RatingUp {
			post.ratingUpCount--
		} else {
			post.ratingDownCount--
		}
		vote.Vote = v.Vote
	} else {
		post.Votes = append(post.Votes, v)
	}

	if v.Vote == RatingUp {
		post.ratingUpCount++
	} else {
		post.ratingDownCount++
	}

	post.ratingUpPercent = post.ratingUpCount / (post.ratingUpCount + post.ratingDownCount)
	return v, nil
}

func (pr *PostRep) Vote(id PostIdType, uid users.UserIdType, vote VoteType) (*Post, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	post, err := pr.getPost(id)

	if err != nil {
		return nil, err
	}

	post.AddVote(&Vote{Author: uid, Vote: vote})
	return post, nil
}

func (rep *PostRep) getNewId() PostIdType {
	rep.lastId++
	return PostIdType(rep.lastId)
}

func (pr *PostRep) getPost(id PostIdType) (*Post, error) {
	post, ok := pr.data[id]

	if !ok {
		return nil, ErrNoPost
	}
	return post, nil
}

func (pr *PostRep) Get(id PostIdType) (*Post, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	return pr.getPost(id)
}

func (pr *PostRep) GetAll() ([]*Post, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	posts := make([]*Post, 0, len(pr.data))

	for _, p := range pr.data {
		posts = append(posts, p)
	}
	return posts, nil
}

func (pr *PostRep) Add(post *Post) (*Post, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	post.ID = pr.getNewId()
	pr.data[post.ID] = post
	return post, nil
}

func (pr *PostRep) Update(id PostIdType, newPost *Post) (*Post, error) {
	pr.mutexx.Lock()
	defer pr.mutexx.Unlock()
	post, err := pr.getPost(id)

	if err != nil {
		return post, err
	}

	newPost.ID = post.ID
	pr.data[id] = newPost
	return newPost, nil
}
