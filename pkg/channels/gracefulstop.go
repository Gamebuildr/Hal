package channels

import (
	"os"

	"os/signal"
	"syscall"

	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
	"github.com/braintree/manners"
)

// OSSigChannel listens to os system events
type OSSigChannel struct {
	Server    *manners.GracefulServer
	Log       logger.Log
	osChannel chan os.Signal
}

// CreateChannel setups a new channel to listen to OS events
func (channel *OSSigChannel) CreateChannel() {
	newChannel := make(chan os.Signal)
	channel.osChannel = newChannel
}

// StartChannelListener listens and handles SIGTERM / SIGINT events
func (channel *OSSigChannel) StartChannelListener() {
	signal.Notify(channel.osChannel, syscall.SIGTERM)
	signal.Notify(channel.osChannel, syscall.SIGINT)

	go func() {
		<-channel.osChannel

		channel.Log.Info("Stopping Hal Client gracefully")
		channel.Server.Close()
	}()
}
