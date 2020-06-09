package soundtouch

func (s *Speaker) AddUpdateHandler(uhc UpdateHandlerConfig) {
	s.RemoveUpdateHandler("NotConfigured")
	if uhc.Name == "" || s.Name() == uhc.Name {
		s.UpdateHandlers = append(s.UpdateHandlers, uhc)
	}
}

func (s *Speaker) RemoveUpdateHandler(name string) {
	var newHandler []UpdateHandlerConfig

	for _, suhc := range s.UpdateHandlers {
		if suhc.Name != "NotConfigured" {
			newHandler = append(newHandler, suhc)
		}
	}
	s.UpdateHandlers = newHandler
}

// HasUpdateHandler returns true if speaker has an UpdateHandler named name. False otherwise
func (s *Speaker) HasUpdateHandler(name string) bool {
	for _, suhc := range s.UpdateHandlers {
		if suhc.Name == name {
			return true
		}
	}
	return false
}

func (s *Speaker) Handle(msgChan chan *Update) {
	for _, uh := range s.UpdateHandlers {
		uh.UpdateHandler.Handle(msgChan, *s)
		if uh.Terminate {
			return
		}
	}
}

// UpdateHandlerFunc turns a function with the right signature into a update handler
type UpdateHandlerFunc func(msgChan chan *Update, speaker Speaker)

// Handle executing the request and returning a response
func (fn UpdateHandlerFunc) Handle(msgChan chan *Update, speaker Speaker) {
	fn(msgChan, speaker)
}

// UpdateHandler interface for that can handle valid update params
type UpdateHandler interface {
	Handle(msgChan chan *Update, speaker Speaker)
}

// UpdateHandlerConfig describes an UpdateHandler. It has a
// Name to be able to remove again
// UpdateHandler the function
// Terminate indicates whether this is the last handler to be called
type UpdateHandlerConfig struct {
	Name          string
	UpdateHandler UpdateHandler
	Terminate     bool
}
