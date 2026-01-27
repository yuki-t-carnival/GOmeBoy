package apu

import (
	"encoding/binary"
	"math"
	"time"
)

const SampleRate float64 = 44100.0
const CyclesPerSample float64 = 4194304.0 / SampleRate
const SamplesPerLengthTimerTick float64 = SampleRate / 256.0
const SamplesPerEnvelopeTick float64 = SampleRate / 64.0
const SamplesPerSweepTick float64 = SampleRate / 128.0

var b [8]byte // Declared here for optimization.

type APU struct {
	AudioStream *AudioStream

	cycles float64

	// For debug
	debugStrings              []string
	timeOfDebugStringsCreated time.Time

	// Global control registers
	nr52 byte
	nr51 byte
	nr50 byte

	// Sound channel 1
	nr10                         byte
	nr11                         byte
	nr12                         byte
	nr13                         byte
	nr14                         byte
	ch1Vol                       byte
	ch1SampleCountForLengthTimer float64
	ch1SampleCountForEnvelope    float64
	ch1LengthTimer               int
	ch1Phase                     float64
	ch1SampleCountForSweep       float64
	ch1Period                    uint16

	// Sound channel 2
	nr21                         byte
	nr22                         byte
	nr23                         byte
	nr24                         byte
	ch2Vol                       byte
	ch2SampleCountForLengthTimer float64
	ch2SampleCountForEnvelope    float64
	ch2LengthTimer               int
	ch2Phase                     float64
	ch2Period                    uint16

	// Sound channel 3
	waveRAM                      [16]byte
	nr30                         byte
	nr31                         byte
	nr32                         byte
	nr33                         byte
	nr34                         byte
	ch3SampleCountForLengthTimer float64
	ch3LengthTimer               int
	ch3Phase                     float64
	indexWaveRAM                 int

	// Sound channel 4
	nr41                         byte
	nr42                         byte
	nr43                         byte
	nr44                         byte
	ch4Vol                       byte
	ch4SampleCountForLengthTimer float64
	ch4SampleCountForEnvelope    float64
	ch4SampleCountForLFSR        float64
	ch4LengthTimer               int
	lfsr                         uint16

	distanceThreshold    int
	samplingAcceleration float64
}

func NewAPU() *APU {
	bufferMilliSecond := float64(120)
	bufferSize := int(SampleRate * 8 * bufferMilliSecond / 1000)
	a := &APU{
		AudioStream:          NewAudioStream(bufferSize),
		lfsr:                 0x7FFF,
		samplingAcceleration: 1.0,
	}

	a.distanceThreshold = int(float64(len(a.AudioStream.buffer)) * 0.40)
	return a
}

// The Step adjusts the sampling interval depending on the current buffer size.
func (a *APU) Step(cpuCycles int) {
	targetAcceleration := 1.0
	w := a.AudioStream.w
	r := a.AudioStream.r
	t := a.distanceThreshold
	if w-r < -t {
		targetAcceleration = 1.5
	} else if w-r > t {
		targetAcceleration = 0.5
	}
	a.samplingAcceleration += (targetAcceleration - a.samplingAcceleration) * 0.01
	a.cycles += float64(cpuCycles) * a.samplingAcceleration
	for a.cycles >= CyclesPerSample {
		a.cycles -= CyclesPerSample
		b := a.generateSample()
		a.AudioStream.write(b)
	}
}

func (a *APU) ReadWaveRAM(addr uint16) byte {
	return a.waveRAM[addr]
}

func (a *APU) WriteWaveRAM(addr uint16, val byte) {
	a.waveRAM[addr] = val
}

func (a *APU) GetAPUInfo() []string {
	if time.Since(a.timeOfDebugStringsCreated).Milliseconds() >= 500 {
		a.debugStrings = []string{}
		//a.debugStrings = append(a.debugStrings, "BUF STATUS")
		a.timeOfDebugStringsCreated = time.Now()
	}
	return a.debugStrings
}

func (a *APU) generateSample() [8]byte {
	sr1, sl1 := a.generateSquareChannel(1)
	sr2, sl2 := a.generateSquareChannel(2)
	sr3, sl3 := a.generateWaveChannel()
	sr4, sl4 := a.generateNoiseChannel()
	/*sr1 = 0
	sl1 = 0
	sr2 = 0
	sl2 = 0
	sr3 = 0
	sl3 = 0
	sr4 = 0
	sl4 = 0 */
	sampleR := float32((sr1 + sr2 + sr3 + sr4) / 4.0)
	sampleL := float32((sl1 + sl2 + sl3 + sl4) / 4.0)
	binary.LittleEndian.PutUint32(b[0:4], math.Float32bits(sampleL))
	binary.LittleEndian.PutUint32(b[4:8], math.Float32bits(sampleR))
	return b
}

// Channel 1 & 2.
func (a *APU) generateSquareChannel(ch byte) (float64, float64) {
	var nrX1 byte
	var nrX2 byte
	//var nrX3 byte
	var nrX4 byte
	var chXPhase *float64
	var chXSampleCountForLengthTimer *float64
	var chXLengthTimer *int
	var chXSampleCountForEnvelope *float64
	var chXVol *byte
	var chXPeriod *uint16
	switch ch {
	case 1:
		nrX1 = a.nr11
		nrX2 = a.nr12
		//nrX3 = a.nr13
		nrX4 = a.nr14
		chXPhase = &a.ch1Phase
		chXSampleCountForLengthTimer = &a.ch1SampleCountForLengthTimer
		chXLengthTimer = &a.ch1LengthTimer
		chXSampleCountForEnvelope = &a.ch1SampleCountForEnvelope
		chXVol = &a.ch1Vol
		chXPeriod = &a.ch1Period
	case 2:
		nrX1 = a.nr21
		nrX2 = a.nr22
		//nrX3 = a.nr23
		nrX4 = a.nr24
		chXPhase = &a.ch2Phase
		chXSampleCountForLengthTimer = &a.ch2SampleCountForLengthTimer
		chXLengthTimer = &a.ch2LengthTimer
		chXSampleCountForEnvelope = &a.ch2SampleCountForEnvelope
		chXVol = &a.ch2Vol
		chXPeriod = &a.ch2Period
	default:
		panic("")
	}

	// Sweep function is only available on Channel 1.
	if ch == 1 {
		pace := a.nr10 & 0x70 >> 4
		if pace != 0 {
			sweepItaration := SamplesPerSweepTick * float64(pace)
			a.ch1SampleCountForSweep++
			for a.ch1SampleCountForSweep >= sweepItaration {
				a.ch1SampleCountForSweep -= sweepItaration
				individualStep := int(a.nr10 & 0x07)
				delta := int(*chXPeriod) >> individualStep
				direction := a.nr10 & (1 << 3) >> 3
				if direction == 1 {
					delta = -delta
				}
				tmpPeriod := int(*chXPeriod) + delta
				if tmpPeriod < 0 || tmpPeriod > 0x7FF {
					a.nr52 = a.nr52 & 0xFE
				} else {
					*chXPeriod = uint16(tmpPeriod)
				}
			}
		}
	}

	freqency := 131072.0 / (2048.0 - float64(*chXPeriod))
	*chXPhase += freqency / float64(SampleRate)
	for *chXPhase >= 1.0 {
		*chXPhase -= 1.0
	}
	var dutyRatio float64
	waveDuty := nrX1 >> 6
	switch waveDuty {
	case 0:
		dutyRatio = 0.125
	case 1:
		dutyRatio = 0.25
	case 2:
		dutyRatio = 0.5
	case 3:
		dutyRatio = 0.75
	}

	a.execLengthTimer(ch, nrX4, chXLengthTimer, chXSampleCountForLengthTimer)
	a.execEnvelope(chXVol, nrX2, chXSampleCountForEnvelope)
	volumeR, volumeL := a.execMixing(ch, float64(*chXVol)/15.0)

	var sampleR, sampleL float64
	if *chXPhase < dutyRatio {
		sampleR = +volumeR
		sampleL = +volumeL
	} else {
		sampleR = -volumeR
		sampleL = -volumeL
	}
	return sampleR, sampleL
}

// Channel 3.
func (a *APU) generateWaveChannel() (float64, float64) {
	period := uint16(a.nr34&0x07)<<8 | uint16(a.nr33)
	freqency := (65536.0 / (2048.0 - float64(period))) * 32

	a.ch3Phase += freqency / float64(SampleRate)
	for a.ch3Phase >= 1.0 {
		a.ch3Phase -= 1.0
		a.indexWaveRAM = (a.indexWaveRAM + 1) % 32
	}

	// ram[0]hi, ram[0]lo, ram[1]hi...
	wave := a.waveRAM[a.indexWaveRAM/2]
	if a.indexWaveRAM%2 == 0 {
		wave >>= 4
	} else {
		wave &= 0x0F
	}

	outputLevel := a.nr32 & 0x60 >> 5
	switch outputLevel {
	case 0:
		wave = 0
	case 2:
		wave >>= 1
	case 3:
		wave >>= 2
	}

	a.execLengthTimer(3, a.nr34, &a.ch3LengthTimer, &a.ch3SampleCountForLengthTimer)
	volumeR, volumeL := a.execMixing(3, 1.0)

	isDACOn := a.nr30&0x80 != 0
	if !isDACOn {
		volumeR = 0
		volumeL = 0
	}

	var sampleR, sampleL float64
	sampleR = (float64(wave)/7.5 - 1.0) * volumeR // ram[x]hi/lo = 0~15 to -1~+1.
	sampleL = (float64(wave)/7.5 - 1.0) * volumeL // ram[x]hi/lo = 0~15 to -1~+1.
	return sampleR, sampleL
}

// Channel 4
func (a *APU) generateNoiseChannel() (float64, float64) {
	clockShift := float64(a.nr43 >> 4)
	clockDivider := float64(a.nr43 & 0x07)
	if clockDivider == 0 {
		clockDivider = 0.5
	}
	lfsrWidth := a.nr43 & (1 << 3) >> 3 // 0:15bit   1:7bit
	a.ch4SampleCountForLFSR++
	// Set a min value for when LFSR clock > Sample rate.
	samplesPerLFSRClock := max(1.0, SampleRate/(262144.0/(clockDivider*math.Pow(2, clockShift))))
	for a.ch4SampleCountForLFSR >= samplesPerLFSRClock {
		a.ch4SampleCountForLFSR -= samplesPerLFSRClock
		xor := a.lfsr&1 ^ a.lfsr&2>>1
		a.lfsr = xor<<15 | a.lfsr&^(1<<15)
		if lfsrWidth == 1 {
			a.lfsr = xor<<7 | a.lfsr&^(1<<7)
		}
		a.lfsr >>= 1
	}

	a.execLengthTimer(4, a.nr44, &a.ch4LengthTimer, &a.ch4SampleCountForLengthTimer)
	a.execEnvelope(&a.ch4Vol, a.nr42, &a.ch4SampleCountForEnvelope)
	volumeR, volumeL := a.execMixing(4, float64(a.ch4Vol)/15.0)

	var sampleR, sampleL float64
	if a.lfsr&1 == 0 {
		sampleR = +volumeR
		sampleL = +volumeL
	} else {
		sampleR = -volumeR
		sampleL = -volumeL
	}
	return sampleR, sampleL
}

func (a *APU) execMixing(ch byte, chVolume float64) (float64, float64) {
	isAudioOn := a.nr52&(1<<7) != 0
	isChannelOn := a.nr52&(1<<(ch-1)) != 0
	if !isAudioOn || !isChannelOn {
		return 0, 0
	}
	panR := float64((a.nr51 & (1 << (ch - 1))) >> (ch - 1)) // 0 or 1
	panL := float64((a.nr51 & (1 << (ch + 3))) >> (ch + 3)) // 0 or 1
	masterR := float64(a.nr50&0x07) / 7.0
	masterL := float64(a.nr50&0x70>>4) / 7.0
	volumeR := panR * masterR * chVolume
	volumeL := panL * masterL * chVolume
	return volumeR, volumeL
}

func (a *APU) execLengthTimer(ch, controlRegister byte, lengthTimer *int, sampleCounter *float64) {
	maxLength := 64
	if ch == 3 {
		maxLength = 256
	}
	isChannelOn := a.nr52&(1<<(ch-1)) != 0
	isLengthEnabled := controlRegister&(1<<6) != 0
	if isChannelOn && isLengthEnabled {
		(*sampleCounter)++
		for *sampleCounter >= SamplesPerLengthTimerTick {
			*sampleCounter -= SamplesPerLengthTimerTick
			if *lengthTimer >= maxLength {
				a.nr52 &^= 1 << (ch - 1) // Disable chX.
			} else {
				(*lengthTimer)++
			}
		}
	}
}

func (a *APU) execEnvelope(chVolume *byte, envelopeRegister byte, sampleCounter *float64) {
	envelopePeriod := envelopeRegister & 0x07
	if envelopePeriod != 0 {
		(*sampleCounter)++
		for *sampleCounter >= SamplesPerEnvelopeTick*float64(envelopePeriod) {
			*sampleCounter -= SamplesPerEnvelopeTick * float64(envelopePeriod)
			isEnvelopeDirectionUp := envelopeRegister&(1<<3) != 0
			if isEnvelopeDirectionUp {
				if *chVolume < 15 {
					(*chVolume)++
				}
			} else {
				if *chVolume > 0 {
					(*chVolume)--
				}
			}
		}
	}
}
