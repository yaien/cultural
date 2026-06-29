package label

type Preset struct {
	Key     string
	Name    string
	Options func() ([]Option, error)
	Load    func(params ...string) (any, error)
}

type Option struct {
	Value string
	Label string
}
