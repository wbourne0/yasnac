package lib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"go.bug.st/serial"
)

type Signal struct {
	raw  []byte
	name string
}

func debug(f string, args ...any) {
	// no-op for now
	//
	return

	fmt.Printf(f, args...)
}

func (s Signal) string() string {
	return string(s.raw)
}

func (s Signal) equals(dat []byte) bool {
	return string(s.raw) == string(dat)
}

func (s Signal) matches(buf []byte) bool {
	return len(buf) >= len(s.raw) && string(s.raw) == string(buf[:len(s.raw)])
}

const KiB = 1024

var inbuf = [4 * KiB]byte{}

func getChecksum(bytes []byte) (sum uint16) {
	for _, b := range bytes[1:] {
		sum += uint16(b)
	}

	return
}

func getChecksumBytes(in []byte) (out []byte) {
	sum := getChecksum(in)

	out = make([]byte, 2)
	binary.LittleEndian.PutUint16(out, sum)
	return
}

type message struct {
	header []byte
	body   []byte
	footer
}
type footer struct {
	sig      Signal
	checksum uint16
}

const maxChunkSize = 256

type ercError struct {
	code string
	desc string
}

func (e ercError) Error() string {
	return fmt.Sprintf("error #%s: %s", e.code, e.desc)
}

type ERCMachine struct {
	port serial.Port

	isAck1 bool
}

type Position struct {
	// Coordinates
	X, Y, Z float64
	// Rotation (degrees)
	TX, TY, TZ float64
}

func (e *ERCMachine) write(buf []byte) (int, error) {
	debug("raw_write %d bytes: %#v\n", len(buf), string(buf))

	return e.port.Write(buf)
}

func (e *ERCMachine) readRaw() (buf []byte, err error) {
	var size int

	e.port.GetModemStatusBits()

	e.port.SetReadTimeout(time.Second * 5) // arbitrary value
	if size, err = e.port.Read(inbuf[:]); err != nil {
		return
	}

	if size == 0 {
		return nil, errors.New("timed out while reading")
	}

	e.port.SetReadTimeout(100 * time.Millisecond)

	for {
		var read int

		if read, err = e.port.Read(inbuf[size:]); err != nil {
			return
		}

		if read == 0 {
			break
		}

		size += read
	}

	buf = make([]byte, size)
	copy(buf, inbuf[:size])

	debug("raw_read: %#v\n", string(buf))
	return
}

func NewERCMachine(port string) (e *ERCMachine, err error) {
	e = new(ERCMachine)

	e.port, err = serial.Open(port, &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.EvenParity,
		StopBits: serial.OneStopBit,
	})

	return
}

func (e *ERCMachine) jwait() error {
	// for {

	// 	fmt.Printf("%#v\n", string(buf))

	// 	if string(buf) == "1011\r" {
	// 		return nil
	// 	}
	// }
	// 
	
		_, err := e.RunCommand("JWAIT -1")

		return err
}

func (e *ERCMachine) MoveTo(pos Position, speed float64) (err error) {
// 	const movjobfmt = `/JOB
// //NAME __MOV
// //POS
// ///NPOS 1,0,0,0
// ///TOOL 0
// ///RECTAN
// ///RCONF 0,0,0,0,0
// C000=%s
// //INST
// ///DATE 1984/12/34 56:78
// ///ATTR 0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0
// ///FRAME BASE
// NOP
// MOVL C000 V=%f CONT
// END
// `

// 	str := fmt.Sprintf(movjobfmt, pos.String(), speed)
// 	fmt.Println(str)

// 	if err = e.WriteFile("__MOV", bytes.NewBufferString(strings.Replace(str, "\n", "\r\n", -1))); err != nil {
// 		return
// 	}

// 	if _, err = e.RunCommand("START __MOV"); err != nil {
// 		return
// 	}

// 	time.Sleep(5 * time.Second)

	if _, err = e.RunCommand(fmt.Sprintf("MOVL 0,%f,0,%s,0,0,0,0,0,0,0,0", speed, pos.String())); err != nil {
		return
	}

	err = e.jwait()

	return
}

func (e *ERCMachine) signal(s Signal) (err error) {
	_, err = e.write(s.raw)
	return
}

func (e *ERCMachine) GetPosition() (pos Position, err error) {
	var buf []byte
	if buf, err = e.RunCommand("RPOS"); err != nil {
		return
	}

	pos.parse(string(buf))
	return
}

func mustParseFloat(str string) float64 {
	f64, err := strconv.ParseFloat(str, 64)

	if err != nil {
		panic(err)
	}

	return f64
}

func (p *Position) parse(str string) {
	// From section 6.9.8.2.4 of yasnac erc communications
	// http://spaz.org/~jake/robot/479236-17-Communications.pdf
	// The response data is in the following format:
	// Data(1): X axis (mm)
	// Data(2): Y axis (mm)
	// Data(3): Z axis (mm)
	// Data(4): List angle TX (°)
	// Data(5): List angle TY (°)
	// Data(6): List angle TZ (°)
	// Data(7): Type 1 ("0" flip; "1" no flip)
	// Data(8): Type 2 ("0" upper arm; "1" lower arm)
	// Data(9): Type 3 ("0" front; "1" back)
	// Data(10): Pulse number of 7th axis (Traverse axis indicated in mm)
	// Data(11): Pulse number of 8th axis (Traverse axis indicated in mm)
	// Data(12): Pulse number of 9th axis (Traverse axis indicated in mm)
	// Data(13): Pulse number of 10th axis
	// Data(14): Pulse number of 11th axis
	// Data(15): Pulse number of 12th axis
	split := strings.Split(string(str[:len(str)-1]), ",")

	p.X = mustParseFloat(split[0])
	p.Y = mustParseFloat(split[1])
	p.Z = mustParseFloat(split[2])
	p.TX = mustParseFloat(split[3])
	p.TY = mustParseFloat(split[4])
	p.TZ = mustParseFloat(split[5])
}

func stringifyFloat(val float64) string {
	str := fmt.Sprintf("%f", val)

	str = strings.Trim(str, "0")

	if str == "." {
		return "0"
	}

	return strings.TrimRight(str, ".")
}

func (p Position) String() string {
	return fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s",
		stringifyFloat(p.X),
		stringifyFloat(p.Y),
		stringifyFloat(p.Z),
		stringifyFloat(p.TX),
		stringifyFloat(p.TY),
		stringifyFloat(p.TZ),
	)
	// return fmt.Sprintf("V=10 0,10,0,35.695,1123.442,402.783,175.91,-70.69,1,0,0,0,0,0,0,0,0")
	// return fmt.Sprintf("V=10 35.695,1123.442,402.783,175.91,-70.69,1 CONT")
}

func (e *ERCMachine) currentAck() Signal {
	if e.isAck1 {
		e.isAck1 = false
		return ACK1
	}

	e.isAck1 = true
	return ACK0
}

func (e *ERCMachine) sendAck() error {
	return e.signal(e.currentAck())
}

func (e *ERCMachine) sendEOT() (err error) {
	err = e.signal(EOT)

	e.isAck1 = false
	return
}

func (e *ERCMachine) checkSignal(s Signal) (bool, error) {
	buf, err := e.readRaw()

	if err != nil {
		return false, err
	}

	if !s.equals(buf) {
		return false, nil
	}

	return true, nil
}

func (e *ERCMachine) expectSignal(s Signal) error {
	didReceive, err := e.checkSignal(s)

	if err != nil {
		return err
	}

	if !didReceive {
		return fmt.Errorf("invalid transaction: expected signal %s", s.name)
	}

	return nil
}

func (e *ERCMachine) expectACK() error {
	return e.expectSignal(e.currentAck())
}

func (e *ERCMachine) sendHandshake() (err error) {
	if err := e.signal(ENQ); err != nil {
		return err
	}

	return e.expectACK()
}

func (e *ERCMachine) expectHandshake() (err error) {
	err = e.expectSignal(ENQ)

	if err != ErrTimedout {
		e.sendAck()
	} else {
		fmt.Println("timedout")
	}
	return
}

func (e *ERCMachine) confirmedWrite(bl []byte) (err error) {
	for didReceiveAck := false; !didReceiveAck; {
		e.write(bl)

		if didReceiveAck, err = e.checkSignal(e.currentAck()); err != nil {
			return
		}
	}

	return
}

func (e *ERCMachine) message(head, name string) (err error) {
	if err = e.sendHandshake(); err != nil {
		return
	}

	var enc *encoder

	if enc, err = newEncoder(e, []byte(head), []byte(name+"\r")); err != nil {
		return
	}

	if err = enc.Close(); err != nil {
		return
	}

	return e.sendEOT()
}

type LineSkipper struct {
	w       io.Writer
	didRead bool
}

func (l *LineSkipper) Write(buf []byte) (n int, err error) {
	if l.didRead {
		return l.w.Write(buf)
	}

	var val byte

	for n, val = range buf {
		if val == '\r' {
			l.didRead = true
			break
		}
	}

	n++

	if !l.didRead {
		return n, err
	}

	var r int

	r, err = l.w.Write(buf[n:])
	n += r

	return
}

func (e *ERCMachine) ReadFileTo(file string, out io.Writer) (err error) {
	if err = e.message(string(transactionCodeGetJob), file); err != nil {
		return
	}

	if err = e.expectHandshake(); err != nil {
		return
	}

	ls := LineSkipper{w: out}

	_, err = e.decodeInto(&ls)
	return
}

func (e *ERCMachine) readMessage() (message, error) {
	var b bytes.Buffer

	header, err := e.decodeInto(&b)

	return message{
		body:   b.Bytes(),
		header: header,
	}, err
}

func (e *ERCMachine) expectExecutionResponse() (msg message, err error) {
	if err = e.expectHandshake(); err != nil {
		return
	}

	if msg, err = e.readMessage(); err != nil {
		return
	}

	err = e.expectEOT()

	return

}

func (e *ERCMachine) RunCommand(command string) (out []byte, err error) {
	if err = e.message(transactionCodeRunCommand, command+"\r"); err != nil {
		return
	}

	var msg message

	if msg, err = e.expectExecutionResponse(); err != nil {
		return
	}

	out = msg.body

	return
}

func (e *ERCMachine) WriteFile(name string, in io.Reader) (err error) {
	// if err = e.message(string(transactionCodePutJob), name); err != nil {
	// 	return
	// }

	if err = e.sendHandshake(); err != nil {
		return
	}

	var enc *encoder
	if enc, err = newEncoder(e, []byte(transactionCodePutJob), []byte(name+"\r")); err != nil {
		return
	}

	if _, err = io.Copy(enc, in); err != nil {
		return
	}

	if err = enc.Close(); err != nil {
		return
	}

	if err = e.sendEOT(); err != nil {
		return
	}

	_, err = e.expectExecutionResponse()
	return
}

func (e *ERCMachine) handleEOT() {
	e.isAck1 = false
}

func (e *ERCMachine) expectEOT() error {
	return e.expectSignal(EOT)
}
