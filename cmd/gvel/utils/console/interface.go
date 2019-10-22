package console

type Prompt interface {
	RequestPassphrase() string
}
