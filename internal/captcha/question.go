package captcha

import (
	"github.com/google/uuid"
	"math/rand"
)

const defaultPrompt = "Pick the word that doesn't belong:"

type Categories map[string][]string

var DefaultCategories = Categories{
	"animals":   {"Cat", "Dog", "Rabbit", "Horse", "Lion", "Tiger", "Elephant", "Zebra", "Bear"},
	"colors":    {"Red", "Blue", "Green", "Yellow", "Orange", "Purple", "Pink"},
	"fruits":    {"Apple", "Banana", "Grapes", "Mango", "Pineapple", "Peach"},
	"vehicles":  {"Car", "Bike", "Bus", "Truck", "Scooter", "Van", "Train"},
	"countries": {"France", "Brazil", "India", "Japan", "Canada", "Germany"},
	"tools":     {"Hammer", "Wrench", "Screwdriver", "Drill", "Pliers"},
	"languages": {"English", "Spanish", "French", "Chinese", "German", "Arabic"},
	"sports":    {"Soccer", "Tennis", "Basketball", "Hockey", "Golf"},
	"shapes":    {"Circle", "Square", "Triangle", "Rectangle", "Oval"},
	"clothing":  {"Shirt", "Pants", "Jacket", "Socks", "Shoes"},
	"planets":   {"Earth", "Mars", "Jupiter", "Venus", "Saturn"},
}

type QuestionManager struct {
	Categories Categories
}

func NewQuestionManager(categories Categories) *QuestionManager {
	return &QuestionManager{Categories: categories}
}

type Question struct {
	ID      string   `json:"id"`
	Prompt  string   `json:"prompt"`
	Options []string `json:"options"`
	Answer  string   `json:"-"`
}

func (q *QuestionManager) Generate() *Question {
	categoryKeys := make([]string, 0, len(q.Categories))
	for k := range q.Categories {
		categoryKeys = append(categoryKeys, k)
	}

	mainCategory := categoryKeys[rand.Intn(len(categoryKeys))]
	mainWords := q.Categories[mainCategory]

	rand.Shuffle(len(mainWords), func(i, j int) {
		mainWords[i], mainWords[j] = mainWords[j], mainWords[i]
	})

	selected := mainWords[:3]

	var oddWord string
	for {
		otherCategory := categoryKeys[rand.Intn(len(categoryKeys))]
		if otherCategory != mainCategory {
			others := q.Categories[otherCategory]
			oddWord = others[rand.Intn(len(others))]
			break
		}
	}

	options := append([]string{}, selected...)
	options = append(options, oddWord)
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return &Question{
		ID:      uuid.NewString(),
		Prompt:  defaultPrompt,
		Options: options,
		Answer:  oddWord,
	}
}
