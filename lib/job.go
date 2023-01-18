package lib

import (
	"bytes"
	"fmt"
	"time"
)

type Job struct {
	Name string
	Date time.Time
	Pos  []Position

	Steps []string
}

func (j *Job) Dout(pin int, on bool) {
	var i int
	if on {
		i = 1
	}
	j.Steps = append(j.Steps, fmt.Sprintf("DOUT OT#%d %d", pin, i))
}

func (j *Job) Pause() {
	j.Steps = append(j.Steps, "PAUSE")
}

func (j *Job) MoveTo(p Position, speed float64) bool {

	j.Steps = append(j.Steps, fmt.Sprintf("MOVL C%03d V=%s CONT\r\nCWAIT", len(j.Pos), stringifyFloat(speed)))
	j.Pos = append(j.Pos, p)

	if len(j.Pos) == 998 {
		return true
	}

	return false
}

func (j *Job) Buffer() (buf bytes.Buffer) {
	buf.WriteString("/JOB\r\n")
	buf.WriteString(fmt.Sprintf("//NAME %s\r\n", j.Name))
	buf.WriteString("//POS\r\n")
	buf.WriteString(fmt.Sprintf("///NPOS %d,0,0,0\r\n", len(j.Pos)))
	buf.WriteString("///TOOL 0\r\n")
	buf.WriteString("///RECTAN\r\n")
	buf.WriteString("///RCONF 0,0,0,0,0\r\n")
	for idx, pos := range j.Pos {
		buf.WriteString(fmt.Sprintf("C%03d=%s\r\n", idx, pos.String()))
	}
	buf.WriteString("//INST\r\n")
	buf.WriteString(fmt.Sprintf("///DATE %s\r\n", j.Date.Format("2006/01/02 03:04")))
	buf.WriteString("///ATTR 0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0\r\n")
	buf.WriteString("///FRAME BASE\r\n")
	buf.WriteString("NOP\r\n")



	for _, st := range j.Steps {
		buf.WriteString(st + "\r\n")
	}

	buf.WriteString("END\r\n")

	return
}
