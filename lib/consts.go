package lib

import "errors"

type transactionCode string

const (
	transactionCodePutJob     = "02,001"
	transactionCodeGetJob     = "02,051"
	transactionCodeRunCommand = "01,000"
)

func (s Signal) is(o Signal) bool {
	return s.name == o.name
}

var (
	SOH = Signal{[]byte{0x01}, "SOH"} // Start Of Heading: denotes the start of the message heading
	STX = Signal{[]byte{0x02}, "STX"} // Start Of Text: denotes end of heading and beginning of data
	ETX = Signal{[]byte{0x03}, "ETX"} // End Of Text: signals that all payload data has been sent
	EOT = Signal{[]byte{0x04}, "EOT"} // End Of Transmission: indicates the end of transmission
	ENQ = Signal{[]byte{0x05}, "ENQ"} // Enquiry: requests a response from the receiving station
	DLE = Signal{[]byte{0x10}, "DLE"} // Data Link Escape: modifies the meaning of a subsequent char
	NAK = Signal{[]byte{0x15}, "NAK"} // Negative Acknowledge: indicates improper communication
	ETB = Signal{[]byte{0x17}, "ETB"} // End of Transmission Block: used in place of ETX to indicate

	ACK0 = Signal{append(DLE.raw, 0x30), "ACK0"}      // Even acknowledgment
	ACK1 = Signal{append(DLE.raw, 0x31), "ACK1"}      // Odd acknowledgment
	WACK = Signal{append(DLE.raw, 0x6b), "WACK"}      // Wait acknowledgement
	RVI  = Signal{append(DLE.raw, 0x7c), "RVI"}       // Reverse interrupt
	TTD  = Signal{append(STX.raw, ENQ.raw...), "TTD"} // Temporary transmission delay

	TRANSACTIONS = map[string]string{
		// how we issue commands to the robot
		"01,000": "command from remote computer",

		// job and special system files - transmission
		"02,001": "put *.JBI",       // indicates independent job data
		"02,002": "put *.JBR",       // indicates related (master) job data
		"02,011": "put WEAV.DAT",    // weave data
		"02,012": "put TOOL.DAT",    // tool data
		"02,013": "put UFRAME.DAT",  // user coordinate data
		"02,014": "put ABSWELD.DAT", // welder condition data undocumented
		"02,015": "put CV.DAT",      // conveyer condition data
		"02,016": "put SENSOR.DAT",  // locus correction condition data
		"02,017": "put COMARC2.DAT", // com-arc 2 condition data
		"02,018": "put PC1PC2.DAT",  // phase comprehension data
		"02,020": "put POSOUT.DAT",  // unknown, undocumented
		"02,022": "put RECIPRO.DAT", // unknown,  undocumented
		"02,023": "put PALACT.DAT",  // palletizing action data, undocumented
		"02,030": "put SYSTEM.DAT",  // system data

		// job and special system files - request
		"02,051": "get *.JBI",       // independent job data
		"02,052": "get *.JBR",       // related (master) job data
		"02,061": "get WEAV.DAT",    // weave data
		"02,062": "get TOOL.DAT",    // tool data
		"02,063": "get UFRAME.DAT",  // user coordinate data
		"02,064": "get ABSWELD.DAT", // welder condition data, undoc"d
		"02,065": "get CV.DAT",      // conveyer condition
		"02,066": "get SENSOR.DAT",  // locus correction condition data
		"02,067": "get COMARC2.DAT", // COM-ARC2 condition data
		"02,068": "get PC1PC2.DAT",  // phase comprehension data
		"02,070": "get POSOUT.DAT",  // unknown, undocumented
		"02,072": "get RECIPRO.DAT", // unknown, undocumented
		"02,073": "get PALACT.DAT",  // palletizing action data, undoc"d.
		"02,080": "get SYSTEM.DAT",  // system data

		// variable data - transmission
		"03,001": "put byte",
		"03,002": "put integer",
		"03,003": "put double",
		"03,004": "put real",
		"03,005": "put position (pulse data)",
		"03,006": "put position (rectangular data)",
		"03,007": "put external axis (pulse data)",
		"03,008": "put external axis (rectangular data)",

		// variable data - request
		"03,051": "get byte",
		"03,052": "get integer",
		"03,053": "get double",
		"03,054": "get real",
		"03,055": "get position (pulse data)",
		"03,056": "get position (rectangular data)",
		"03,057": "get external axis (pulse data)",
		"03,058": "get external axis (rectangular data)",

		// job execution response
		"90,000": "0000 or a 4 digit error code, response to a command",
		"90,001": "data response, variable number of digits/data sent as csv",
	}

	ERRORS = map[string]ercError{
		// 1xxx - command test
		"1010": {"1010", "command failure"},
		"1011": {"1011", "command operand number failure"},
		"1012": {"1012", "command operand value excessive"},
		"1013": {"1013", "command operand length failure"},

		// 2xxx - command execution mode error
		"2010": {"2010", "during robot operation"},
		"2020": {"2020", "during T-PENDANT"},
		"2030": {"2030", "during panel HOLD"},
		"2040": {"2040", "during external HOLD"},
		"2050": {"2050", "during command HOLD"},
		"2060": {"2060", "during error alarm"},
		"2070": {"2070", "in servo OFF or stopping by a panel HOLD"},

		// 3xxx - command execution error
		"3010": {"3010", "servo power on"},
		"3040": {"3040", "set home position"},
		"3070": {"3070", "current position is not input"},
		"3080": {"3080", "END command of job (except master job)"},

		// 4xxx - job registration error
		"4010": {"4010", "shortage of memory capacity (job registration)"},
		"4012": {"4012", "shortage of memory capacity (position data registration)"},
		"4020": {"4020", "job edit prohibit"},
		"4030": {"4030", "job of same name exists"},
		"4040": {"4040", "no desired job"},
		"4060": {"4060", "set execution"},
		"4120": {"4120", "position data broken"},
		"4130": {"4130", "no position data"},
		"4150": {"4150", "END command of job (except master job)"},
		"4170": {"4170", "instruction data broken"},
		"4190": {"4190", "unsuitable characters in job name exist"},
		"4200": {"4200", "unsuitable characters in job name exist"},
		"4230": {"4230", "instructions which cannot be used by this system exist"},

		// 5xxx - file text error
		"5110": {"5110", "instruction syntax error"},
		"5120": {"5120", "position data fault"},
		"5130": {"5130", "neither NOP or END exists"},
		"5170": {"5170", "format error"},
		"5180": {"5180", "data number is inadequate"},
		"5200": {"5200", "data range exceeded"},
	}

	ErrTimedout = errors.New("timed out")
)
