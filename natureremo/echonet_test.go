package natureremo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEchonet_CalcCumulativePower01(t *testing.T) {
	currentTime := time.Now()

	echonet := NewEchonet()
	echonet.SetValue("d3", "00000001", currentTime)
	echonet.SetValue("d7", "06", currentTime)
	echonet.SetValue("e0", "0000e242", currentTime)
	echonet.SetValue("e1", "01", currentTime)
	echonet.SetValue("e3", "00000008", currentTime)
	echonet.SetValue("e7", "00000110", currentTime)

	power, powerTime, err := echonet.CalcCumulativePower()
	assert.Equal(t, 5.2, power)
	assert.Equal(t, currentTime, powerTime)
	assert.NoError(t, err)
}

func TestEchonet_CalcInstantaneousPower01(t *testing.T) {
	currentTime := time.Now()

	echonet := NewEchonet()
	echonet.SetValue("d3", "00000001", currentTime)
	echonet.SetValue("d7", "06", currentTime)
	echonet.SetValue("e0", "0000e242", currentTime)
	echonet.SetValue("e1", "01", currentTime)
	echonet.SetValue("e3", "00000008", currentTime)
	echonet.SetValue("e7", "00000110", currentTime)

	power, powerTime, err := echonet.CalcInstantaneousPower()
	assert.Equal(t, 272.0, power)
	assert.Equal(t, currentTime, powerTime)
	assert.NoError(t, err)
}
