package models

type Flashcard struct {
    ID       int    `json:"id"`
    Question string `json:"question"`
    Answer   string `json:"answer"`
}
