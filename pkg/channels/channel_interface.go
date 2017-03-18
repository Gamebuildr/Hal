package channels

// Channel is the interface for creating new listening channels
type Channel interface {
	CreateChannel()
	StartChannelListener()
}
