package natureremo

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// reference: https://developer.nature.global/docs/how-to-calculate-energy-data-from-smart-meter-values/

type Echonet struct {
	plus uint64 //積算電力量計測値(正方向)
	mins uint64 //積算電力量計測値(逆方向)
	coef uint64 //係数
	unit uint64 //積算電力量単位
	digt uint64 //積算電力量有効桁数
	inst int64  //瞬時電力計測値

	plusTime time.Time
	minsTime time.Time
	coefTime time.Time
	unitTime time.Time
	digtTime time.Time
	instTime time.Time
}

func NewEchonet() *Echonet {
	return &Echonet{
		plus: 0,
		mins: 0,
		coef: 1,
		unit: 0,
		digt: 1,
		inst: 0,
	}
}

func (e *Echonet) SetValue(epc string, value string, time time.Time) error {
	var err error

	parseUint := func(ref *uint64, value string, bits int) error {
		*ref, err = strconv.ParseUint(value, 16, bits)
		return err
	}
	parseInt := func(ref *int64, value string, bits int) error {
		*ref, err = strconv.ParseInt(value, 16, bits)
		return err
	}

	switch epc {
	case "e0":
		err = parseUint(&e.plus, value, 32)
		e.plusTime = time
	case "e3":
		err = parseUint(&e.mins, value, 32)
		e.minsTime = time
	case "d3":
		err = parseUint(&e.coef, value, 32)
		e.coefTime = time
	case "e1":
		err = parseUint(&e.unit, value, 8)
		e.unitTime = time
	case "d7":
		err = parseUint(&e.digt, value, 8)
		e.digtTime = time
	case "e7":
		err = parseInt(&e.inst, value, 32)
		e.instTime = time
	default:
		return fmt.Errorf("unsupport epc=%s", epc)
	}

	return err
}

func (e *Echonet) CalcCumulativePower() (float64, time.Time, error) {
	maxvalue := math.Pow10(int(e.digt))

	coef := float64(e.coef)

	unit := float64(0)
	switch e.unit {
	case 0x00:
		unit = 1.0
	case 0x01:
		unit = 0.1
	case 0x02:
		unit = 0.01
	case 0x03:
		unit = 0.001
	case 0x04:
		unit = 0.0001
	case 0x0A:
		unit = 10.0
	case 0x0B:
		unit = 100.0
	case 0x0C:
		unit = 1000.0
	case 0x0D:
		unit = 10000.0
	default:
		return 0, e.plusTime, fmt.Errorf("unexpect effective digits: e1=0x%x", e.unit)
	}

	value := float64(e.plus)
	value += maxvalue
	value -= float64(e.mins)
	if value > maxvalue {
		value -= maxvalue
	}

	value *= unit * coef
	return value, e.plusTime, nil
}

func (e *Echonet) CalcInstantaneousPower() (float64, time.Time, error) {
	return float64(e.inst), e.instTime, nil
}
