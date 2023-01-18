package lib

import (
	"encoding/binary"
	"errors"
	"io"
)

const chunkSize = 256
const headerSize = 8
const checksumSize = 2
const footerSize = checksumSize + 1 // +1 for ETX/ETB
const headerBufSize = chunkSize + headerSize + checksumSize
const bodyBufSize = chunkSize + footerSize + 1 // +1 for STX

type encoder struct {
	buf []byte
	*ERCMachine
	isNameChunk bool
}

func newEncoder(m *ERCMachine, head, name []byte) (e *encoder, err error) {
	e = &encoder{
		buf:        make([]byte, 0, headerBufSize),
		ERCMachine: m,
	}

	e.buf = append(e.buf, SOH.raw...)
	e.buf = append(e.buf, head...)
	e.buf = append(e.buf, STX.raw...)

	if name != nil {
		e.buf = append(e.buf, name...)
		e.isNameChunk = true
	}

	return
}

func (e *encoder) writeChunk(hasNext bool) (err error) {
	if hasNext {
		e.buf = append(e.buf, ETB.raw...)
	} else {
		e.buf = append(e.buf, ETX.raw...)
	}

	e.buf = binary.LittleEndian.AppendUint16(e.buf, getChecksum(e.buf))

	_, err = e.write(e.buf)

	if cap(e.buf) != bodyBufSize {
		e.buf = make([]byte, 0, bodyBufSize)
		e.buf = append(e.buf, STX.raw...)
	} else {
		e.buf = e.buf[:0]
		e.buf = append(e.buf, STX.raw...)
	}

	if err != nil {
		return
	}

	return e.expectACK()
}

func (e *encoder) Write(dat []byte) (n int, err error) {
	n = len(dat)
	if e.isNameChunk && len(dat) > 0 {
		if err = e.writeChunk(true); err != nil {
			return
		}

		e.isNameChunk = false
	}

	for len(dat) > 0 {
		if len(e.buf) == cap(e.buf)-footerSize {
			if err = e.writeChunk(true); err != nil {
				return
			}
		}

		rem := cap(e.buf) - len(e.buf) - footerSize

		if rem > len(dat) {
			rem = len(dat)
		}
		e.buf = append(e.buf, dat[:rem]...)
		dat = dat[rem:]
	}

	return
}

func (e *encoder) Close() (err error) {
	if len(e.buf) == 1 {
		return errors.New("empty buf; this shouldn't happen")
	}

	return e.writeChunk(false)
}

type decoder struct {
	*ERCMachine
	to io.Writer
	// lastSig       Signal
	header []byte
	footer footer
}

func (d *decoder) readBlock() (err error) {
	var (
		block              []byte
		bodyStart, bodyEnd int
	)

	block, err = d.readRaw()

	if SOH.matches(block) {
		bodyStart = 8
		d.header = block[1:7]
	} else if STX.matches(block) {
		bodyStart = 1
		d.header = block[:1]
	} else {
		err = errors.New("invalid message")
		return
	}

	for i := bodyStart; i < len(block) && i < maxChunkSize+bodyStart+1; i++ {
		if ETB.equals(block[i : i+1]) {
			bodyEnd = i
			d.footer.sig = ETB
			break
		}

		if ETX.equals(block[i : i+1]) {
			bodyEnd = i
			d.footer.sig = ETX
			break
		}
	}

	if bodyEnd == 0 {
		return errors.New("invalid body")
	}

	if _, err = d.to.Write(block[bodyStart:bodyEnd]); err != nil {
		return
	}

	footer := block[bodyEnd:]

	checksum := binary.LittleEndian.Uint16(footer[1:])

	if checksum != getChecksum(block[:bodyEnd+1]) {
		err = errors.New("invalid checksum")
	}

	d.sendAck()

	return
}

func (m *ERCMachine) decodeInto(to io.Writer) (header []byte, err error) {
	dec := &decoder{ERCMachine: m, to: to}

	if err = dec.readBlock(); err != nil {
		return
	}

	for dec.footer.sig.is(ETB) {
		if err = dec.readBlock(); err != nil {
			return
		}
	}

	header = dec.header

	return
}
