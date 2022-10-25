package glice

type Editors []*Editor
type EditorMap map[string]*Editor
type Editor struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
	Added string `yaml:"added"`
}

var (
	defaultName  = "Username Goes Here"
	defaultEmail = "email-alias@example.com"
)

func NewEditor(ua UserAdapter) *Editor {
	e := &Editor{
		Name:  ua.GetName(),
		Email: ua.GetEmail(),
		Added: Timestamp(),
	}
	e.ID = e.GetID()
	return e
}

func (em EditorMap) ToEditors() Editors {
	editors := make(Editors, len(em))
	index := 0
	for _, ed := range em {
		editors[index] = ed
		index++
	}
	return editors
}

func (e *Editor) GetID() string {
	e.ensureID()
	return e.ID
}

func (e *Editor) ensureID() {
	if e.ID == "" {
		e.ID = UpToN(e.Email, '@', 1)
	}
}
