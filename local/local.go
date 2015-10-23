package local

type Local struct {
	store map[string][]byte
}

func New() Local {
	return Local{store: map[string][]byte{}}
}

func (l Local) Get(key string) []byte {
	return l.store[key]
}

func (l Local) Put(key string, data []byte) error {
	l.store[key] = data
	return nil
}
