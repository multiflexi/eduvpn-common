package failover

import (
	"context"
	"fmt"
	"time"

	"github.com/eduvpn/eduvpn-common/internal/log"
)

type sender interface {
	Read(deadline time.Time) error
	Send(seq int) error
}

// The DroppedConMon is a connection monitor that checks for an increase in rx bytes in certain intervals
type DroppedConMon struct {
	// pInterval means how the interval in which to send pings
	pInterval time.Duration
	// pDropped means how many pings need to be send before checking if the connection is dropped
	pDropped int
	// The function that reads Rx bytes
	// If this function returns an error, the monitor exits
	readRxBytes func() (int64, error)
	// newPinger creates a new pinger
	// This gets used in the tests to mock the Ping sender interface
	newPinger func(gateway string, mtu int) (sender, error)
}

// NewDroppedMonitor creates a new failover monitor
// `pingInterval` is the interval in which to send pings
// `pDropped` is how many pings we need to send before we deem it is dropped
// `readRxBytes` is a function that gets the rx bytes from the client
func NewDroppedMonitor(pingInterval time.Duration, pDropped int, readRxBytes func() (int64, error)) *DroppedConMon {
	return &DroppedConMon{pInterval: pingInterval, pDropped: pDropped, readRxBytes: readRxBytes, newPinger: func(gateway string, mtu int) (sender, error) {
		return NewPinger(gateway, mtu)
	}}
}

// Dropped checks whether or not the connection is 'dropped'
// In other words, it checks if rx bytes has increased
func (m *DroppedConMon) dropped(startBytes int64) (bool, error) {
	b, err := m.readRxBytes()
	if err != nil {
		return false, err
	}
	log.Logger.Debugf("[Failover] Alive check, current Rx bytes: %d, start Rx bytes: %d", b, startBytes)
	return b <= startBytes, nil
}

// Start starts ticking every ping interval and check if the connection is dropped or alive
// This does not check Rx bytes every tick, but rather when pAlive or pDropped is reached
// It returns an error if there was an invalid input or a ping was failed to be sent
func (m *DroppedConMon) Start(ctx context.Context, gateway string, mtuSize int) (bool, error) {
	if mtuSize < mtuOverhead {
		return false, fmt.Errorf("invalid MTU size given, MTU has to be at least: %v bytes", mtuOverhead)
	}

	// Create a ping struct with our mtu size
	p, err := m.newPinger(gateway, mtuSize)
	if err != nil {
		return false, err
	}

	// Read the start Rx bytes
	b, err := m.readRxBytes()
	if err != nil {
		return false, err
	}

	// Send a ping and wait for max 2 seconds
	// If we have then increased Rx bytes we return early
	if err = p.Send(1); err != nil {
		log.Logger.Debugf("[Failover] First ping failed, exiting...")
		return false, err
	}
	log.Logger.Debugf("[Failover] Now we are doing alive check")

	// Read the pong, if we got the echo reply then everything is fine, early return
	if err = p.Read(time.Now().Add(m.pInterval)); err == nil {
		log.Logger.Debugf("[Failover] Got early pong, exiting...")
		return false, err
	}
	log.Logger.Debugf("[Failover] Error reading pong: %v", err)

	// Create a new ticker that executes our ping function every 'interval' seconds
	// It starts immediately and stops when we reach the end
	ticker := time.NewTicker(m.pInterval)
	defer ticker.Stop()

	// Otherwise send n pings, without waiting for pong and then check if dropped
	log.Logger.Debugf("[Failover] Starting by sending pings and not waiting for pong...")
	// Loop until the max drop counter
	// We begin with 2 as this is used as the sequence number for ping
	// and we have already sent a ping
	for s := 2; s <= m.pDropped; s++ {
		log.Logger.Debugf("[Failover] Sending ping: %d, with size: %d", s, mtuSize)
		// Send a ping and return if an error occurs
		if err := p.Send(s); err != nil {
			log.Logger.Debugf("[Failover] A ping failed, exiting...")
			return false, err
		}
		// Wait for the next tick to continue
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return false, fmt.Errorf("failover was stopped with error: %w", context.Canceled)
		}
	}

	// Dropped check if we have not returned early
	return m.dropped(b)
}
