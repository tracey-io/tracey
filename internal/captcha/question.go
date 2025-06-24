package captcha

import (
	"github.com/google/uuid"
	"math/rand"
)

const DefaultQuestionPrompt = "Pick the word that doesn't belong:"

type Categories map[string][]string

var DefaultCategories = Categories{
	"animals": {
		"Cat", "Dog", "Rabbit", "Horse", "Lion", "Tiger", "Elephant", "Zebra", "Bear",
		"Giraffe", "Kangaroo", "Fox", "Wolf", "Monkey", "Deer", "Sheep", "Cow", "Goat",
		"Pig", "Camel", "Cheetah", "Leopard", "Panda", "Otter", "Squirrel", "Raccoon",
		"Hedgehog", "Moose", "Sloth", "Bat", "Donkey", "Hyena", "Llama", "Crocodile",
		"Alligator", "Frog", "Toad", "Turtle", "Lizard", "Snake", "Whale", "Dolphin",
		"Shark", "Octopus", "Crab", "Lobster", "Penguin", "Owl", "Eagle", "Parrot",
		"Swan", "Duck", "Goose", "Rooster", "Hen", "Turkey", "Peacock", "Flamingo",
	},
	"colors": {
		"Red", "Blue", "Green", "Yellow", "Orange", "Purple", "Pink", "Brown", "Black", "White",
		"Gray", "Cyan", "Magenta", "Beige", "Turquoise", "Lavender", "Maroon", "Navy", "Teal",
		"Gold", "Silver", "Bronze", "Ivory", "Coral", "Olive", "Indigo", "Mint", "Peach",
	},
	"fruits": {
		"Apple", "Banana", "Grapes", "Mango", "Pineapple", "Peach", "Strawberry", "Cherry",
		"Watermelon", "Papaya", "Kiwi", "Orange", "Lemon", "Lime", "Blueberry", "Raspberry",
		"Coconut", "Pomegranate", "Fig", "Guava", "Lychee", "Passionfruit", "Dragonfruit",
		"Avocado", "Cranberry", "Tangerine", "Cantaloupe", "Plum", "Apricot", "Date",
	},
	"vehicles": {
		"Car", "Bike", "Bus", "Truck", "Scooter", "Van", "Train", "Airplane", "Boat",
		"Helicopter", "Tram", "Submarine", "Taxi", "Pickup", "SUV", "Motorcycle", "Jet",
		"Yacht", "Ferry", "Tank", "Rickshaw", "Skateboard", "Rollerblades", "Cruise", "Spaceship",
	},
	"countries": {
		"France", "Brazil", "India", "Japan", "Canada", "Germany", "Australia", "Mexico", "Italy",
		"Russia", "China", "South Korea", "Spain", "Argentina", "Sweden", "Norway", "Finland",
		"Egypt", "South Africa", "Thailand", "Vietnam", "Turkey", "New Zealand", "Indonesia",
		"Philippines", "Greece", "Netherlands", "Poland", "Ukraine", "Switzerland", "Portugal",
	},
	"tools": {
		"Hammer", "Wrench", "Screwdriver", "Drill", "Pliers", "Saw", "Chisel", "Tape Measure",
		"Level", "Utility Knife", "Mallet", "Socket Wrench", "Allen Key", "Clamp", "File",
		"Hacksaw", "Chainsaw", "Sander", "Trowel", "Crowbar", "Stud Finder", "Voltage Tester",
		"Caulking Gun", "Nail Gun", "Ladder", "Workbench", "Vice", "Wire Cutter",
	},
	"languages": {
		"English", "Spanish", "French", "Chinese", "German", "Arabic", "Russian", "Portuguese",
		"Hindi", "Japanese", "Korean", "Italian", "Dutch", "Greek", "Turkish", "Swedish",
		"Polish", "Hebrew", "Thai", "Vietnamese", "Romanian", "Czech", "Finnish", "Indonesian",
		"Malay", "Bengali", "Tamil", "Telugu", "Ukrainian", "Persian",
	},
	"sports": {
		"Soccer", "Tennis", "Basketball", "Hockey", "Golf", "Baseball", "Cricket", "Rugby",
		"Swimming", "Volleyball", "Skating", "Boxing", "Wrestling", "Table Tennis", "Badminton",
		"Karate", "Judo", "Skiing", "Snowboarding", "Surfing", "Cycling", "Rowing", "Fencing",
		"Handball", "Archery", "Skateboarding", "Gymnastics", "Equestrian", "Diving",
	},
	"shapes": {
		"Circle", "Square", "Triangle", "Rectangle", "Oval", "Hexagon", "Pentagon", "Octagon",
		"Star", "Heart", "Diamond", "Trapezoid", "Parallelogram", "Crescent", "Cross", "Arrow",
	},
	"clothing": {
		"Shirt", "Pants", "Jacket", "Socks", "Shoes", "Hat", "Scarf", "Gloves", "Dress", "Shorts",
		"Sweater", "Blazer", "Boots", "Sandals", "Belt", "Skirt", "Tie", "Hoodie", "Coat",
		"Raincoat", "Cap", "Mittens", "Undershirt", "Jeans",
	},
	"planets": {
		"Earth", "Mars", "Jupiter", "Venus", "Saturn", "Mercury", "Uranus", "Neptune",
	},
}

var DefaultQuestionConfig = &QuestionConfig{
	Categories: DefaultCategories,
	Prompt:     DefaultQuestionPrompt,
}

type QuestionConfig struct {
	Categories Categories
	Prompt     string
}
type QuestionManager struct {
	QuestionConfig *QuestionConfig
}

func NewQuestionManager(questionConfig *QuestionConfig) *QuestionManager {
	return &QuestionManager{
		QuestionConfig: questionConfig,
	}
}

type Question struct {
	ID      string   `json:"id"`
	Prompt  string   `json:"prompt"`
	Options []string `json:"options"`
	Answer  string   `json:"-"`
}

func (q *QuestionManager) Generate() *Question {
	categoryKeys := make([]string, 0, len(q.QuestionConfig.Categories))
	for k := range q.QuestionConfig.Categories {
		categoryKeys = append(categoryKeys, k)
	}

	mainCategory := categoryKeys[rand.Intn(len(categoryKeys))]
	mainWords := q.QuestionConfig.Categories[mainCategory]

	rand.Shuffle(len(mainWords), func(i, j int) {
		mainWords[i], mainWords[j] = mainWords[j], mainWords[i]
	})

	selected := mainWords[:3]

	var oddWord string
	for {
		otherCategory := categoryKeys[rand.Intn(len(categoryKeys))]
		if otherCategory != mainCategory {
			others := q.QuestionConfig.Categories[otherCategory]
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
		Prompt:  q.QuestionConfig.Prompt,
		Options: options,
		Answer:  oddWord,
	}
}
