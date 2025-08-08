package doc

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// (placeholder) later you can add:
// type UserStore interface { Ensure(id, name string) error }
