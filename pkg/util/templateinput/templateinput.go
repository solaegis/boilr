package templateinput

func init() {
	UserInput.Init()
	UsedKeys.Init()
}

type set map[string]interface{}

var UserInput set

func (s *set) Init() {
	*s = make(map[string]interface{})
}

// store keys used during boilr execution
type s map[string]interface{}

var UsedKeys s

func (s *s) Init() {
	*s = make(map[string]interface{})
}
