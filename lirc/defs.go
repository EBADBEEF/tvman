package lirc

import (
	"tvman/ioctl"
	"unsafe"
)

/*
#define LIRC_SPACE(val) (((val) & LIRC_VALUE_MASK) | LIRC_MODE2_SPACE)
#define LIRC_PULSE(val) (((val) & LIRC_VALUE_MASK) | LIRC_MODE2_PULSE)
#define LIRC_FREQUENCY(val) (((val) & LIRC_VALUE_MASK) | LIRC_MODE2_FREQUENCY)
#define LIRC_TIMEOUT(val) (((val) & LIRC_VALUE_MASK) | LIRC_MODE2_TIMEOUT)
#define LIRC_OVERFLOW(val) (((val) & LIRC_VALUE_MASK) | LIRC_MODE2_OVERFLOW)
#define LIRC_VALUE(val) ((val)&LIRC_VALUE_MASK)
#define LIRC_MODE2(val) ((val)&LIRC_MODE2_MASK)
#define LIRC_IS_SPACE(val) (LIRC_MODE2(val) == LIRC_MODE2_SPACE)
#define LIRC_IS_PULSE(val) (LIRC_MODE2(val) == LIRC_MODE2_PULSE)
#define LIRC_IS_FREQUENCY(val) (LIRC_MODE2(val) == LIRC_MODE2_FREQUENCY)
#define LIRC_IS_TIMEOUT(val) (LIRC_MODE2(val) == LIRC_MODE2_TIMEOUT)
#define LIRC_IS_OVERFLOW(val) (LIRC_MODE2(val) == LIRC_MODE2_OVERFLOW)
#define LIRC_CAN_SEND(x) ((x)&LIRC_CAN_SEND_MASK)
#define LIRC_CAN_REC(x) ((x)&LIRC_CAN_REC_MASK)
*/

const (
	PULSE_BIT                      = 0x01000000
	PULSE_MASK                     = 0x00FFFFFF
	LIRC_MODE2_SPACE               = 0x00000000
	LIRC_MODE2_PULSE               = 0x01000000
	LIRC_MODE2_FREQUENCY           = 0x02000000
	LIRC_MODE2_TIMEOUT             = 0x03000000
	LIRC_MODE2_OVERFLOW            = 0x04000000
	LIRC_VALUE_MASK                = 0x00FFFFFF
	LIRC_MODE2_MASK                = 0xFF000000
	LIRC_MODE_RAW                  = 0x00000001
	LIRC_MODE_PULSE                = 0x00000002
	LIRC_MODE_MODE2                = 0x00000004
	LIRC_MODE_SCANCODE             = 0x00000008
	LIRC_MODE_LIRCCODE             = 0x00000010
	LIRC_CAN_SEND_RAW              = LIRC_MODE_RAW
	LIRC_CAN_SEND_PULSE            = LIRC_MODE_PULSE
	LIRC_CAN_SEND_MODE2            = LIRC_MODE_MODE2
	LIRC_CAN_SEND_LIRCCODE         = LIRC_MODE_LIRCCODE
	LIRC_CAN_SEND_MASK             = 0x0000003f
	LIRC_CAN_SET_SEND_CARRIER      = 0x00000100
	LIRC_CAN_SET_SEND_DUTY_CYCLE   = 0x00000200
	LIRC_CAN_SET_TRANSMITTER_MASK  = 0x00000400
	LIRC_CAN_REC_RAW               = LIRC_MODE_RAW << 16
	LIRC_CAN_REC_PULSE             = LIRC_MODE_PULSE << 16
	LIRC_CAN_REC_MODE2             = LIRC_MODE_MODE2 << 16
	LIRC_CAN_REC_SCANCODE          = LIRC_MODE_SCANCODE << 16
	LIRC_CAN_REC_LIRCCODE          = LIRC_MODE_LIRCCODE << 16
	LIRC_CAN_REC_MASK              = LIRC_CAN_SEND_MASK << 16
	LIRC_CAN_SET_REC_CARRIER       = (LIRC_CAN_SET_SEND_CARRIER << 16)
	LIRC_CAN_SET_REC_CARRIER_RANGE = 0x80000000
	LIRC_CAN_GET_REC_RESOLUTION    = 0x20000000
	LIRC_CAN_SET_REC_TIMEOUT       = 0x10000000
	LIRC_CAN_MEASURE_CARRIER       = 0x02000000
	LIRC_CAN_USE_WIDEBAND_RECEIVER = 0x04000000
)

//const (
//	LIRC_CAN_SET_REC_FILTER	= 0
//	LIRC_CAN_NOTIFY_DECODE = 0
//)

const (
	LIRC_SCANCODE_FLAG_TOGGLE = 1
	LIRC_SCANCODE_FLAG_REPEAT = 2
)

type Scancode struct {
	Timestamp uint64
	Flags     uint16
	RcProto   uint16
	Keycode   uint32
	Scancode  uint64
}

const (
	RC_PROTO_UNKNOWN   = 0
	RC_PROTO_OTHER     = 1
	RC_PROTO_RC5       = 2
	RC_PROTO_RC5X_20   = 3
	RC_PROTO_RC5_SZ    = 4
	RC_PROTO_JVC       = 5
	RC_PROTO_SONY12    = 6
	RC_PROTO_SONY15    = 7
	RC_PROTO_SONY20    = 8
	RC_PROTO_NEC       = 9
	RC_PROTO_NECX      = 10
	RC_PROTO_NEC32     = 11
	RC_PROTO_SANYO     = 12
	RC_PROTO_MCIR2_KBD = 13
	RC_PROTO_MCIR2_MSE = 14
	RC_PROTO_RC6_0     = 15
	RC_PROTO_RC6_6A_20 = 16
	RC_PROTO_RC6_6A_24 = 17
	RC_PROTO_RC6_6A_32 = 18
	RC_PROTO_RC6_MCE   = 19
	RC_PROTO_SHARP     = 20
	RC_PROTO_XMP       = 21
	RC_PROTO_CEC       = 22
	RC_PROTO_IMON      = 23
	RC_PROTO_RCMM12    = 24
	RC_PROTO_RCMM24    = 25
	RC_PROTO_RCMM32    = 26
	RC_PROTO_XBOX_DVD  = 27
	RC_PROTO_MAX       = RC_PROTO_XBOX_DVD
)

var (
	LIRC_GET_FEATURES             = uintptr(ioctl.IOR('i', 0x00000000, unsafe.Sizeof(uint32(0))))
	LIRC_GET_SEND_MODE            = uintptr(ioctl.IOR('i', 0x00000001, unsafe.Sizeof(uint32(0))))
	LIRC_GET_REC_MODE             = uintptr(ioctl.IOR('i', 0x00000002, unsafe.Sizeof(uint32(0))))
	LIRC_GET_REC_RESOLUTION       = uintptr(ioctl.IOR('i', 0x00000007, unsafe.Sizeof(uint32(0))))
	LIRC_GET_MIN_TIMEOUT          = uintptr(ioctl.IOR('i', 0x00000008, unsafe.Sizeof(uint32(0))))
	LIRC_GET_MAX_TIMEOUT          = uintptr(ioctl.IOR('i', 0x00000009, unsafe.Sizeof(uint32(0))))
	LIRC_GET_LENGTH               = uintptr(ioctl.IOR('i', 0x0000000f, unsafe.Sizeof(uint32(0))))
	LIRC_SET_SEND_MODE            = uintptr(ioctl.IOW('i', 0x00000011, unsafe.Sizeof(uint32(0))))
	LIRC_SET_REC_MODE             = uintptr(ioctl.IOW('i', 0x00000012, unsafe.Sizeof(uint32(0))))
	LIRC_SET_SEND_CARRIER         = uintptr(ioctl.IOW('i', 0x00000013, unsafe.Sizeof(uint32(0))))
	LIRC_SET_REC_CARRIER          = uintptr(ioctl.IOW('i', 0x00000014, unsafe.Sizeof(uint32(0))))
	LIRC_SET_SEND_DUTY_CYCLE      = uintptr(ioctl.IOW('i', 0x00000015, unsafe.Sizeof(uint32(0))))
	LIRC_SET_TRANSMITTER_MASK     = uintptr(ioctl.IOW('i', 0x00000017, unsafe.Sizeof(uint32(0))))
	LIRC_SET_REC_TIMEOUT          = uintptr(ioctl.IOW('i', 0x00000018, unsafe.Sizeof(uint32(0))))
	LIRC_SET_REC_TIMEOUT_REPORTS  = uintptr(ioctl.IOW('i', 0x00000019, unsafe.Sizeof(uint32(0))))
	LIRC_SET_MEASURE_CARRIER_MODE = uintptr(ioctl.IOW('i', 0x0000001d, unsafe.Sizeof(uint32(0))))
	LIRC_SET_REC_CARRIER_RANGE    = uintptr(ioctl.IOW('i', 0x0000001f, unsafe.Sizeof(uint32(0))))
	LIRC_SET_WIDEBAND_RECEIVER    = uintptr(ioctl.IOW('i', 0x00000023, unsafe.Sizeof(uint32(0))))
	LIRC_GET_REC_TIMEOUT          = uintptr(ioctl.IOR('i', 0x00000024, unsafe.Sizeof(uint32(0))))
)
