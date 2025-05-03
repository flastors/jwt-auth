package empty

type EmptyMailer struct {
}

func NewUserEmail() *EmptyMailer {
	return &EmptyMailer{}
}

func (m *EmptyMailer) Send(to, subject, body string) error {
	return nil
}
