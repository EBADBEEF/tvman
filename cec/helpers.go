package cec

import (
	"fmt"
	"strings"
)

var cecMsgToString = map[uint8]string{
	CEC_MSG_ACTIVE_SOURCE:                  "active_source",
	CEC_MSG_IMAGE_VIEW_ON:                  "image_view_on",
	CEC_MSG_TEXT_VIEW_ON:                   "text_view_on",
	CEC_MSG_INACTIVE_SOURCE:                "inactive_source",
	CEC_MSG_REQUEST_ACTIVE_SOURCE:          "request_active_source",
	CEC_MSG_ROUTING_CHANGE:                 "routing_change",
	CEC_MSG_ROUTING_INFORMATION:            "routing_information",
	CEC_MSG_SET_STREAM_PATH:                "set_stream_path",
	CEC_MSG_STANDBY:                        "standby",
	CEC_MSG_RECORD_OFF:                     "record_off",
	CEC_MSG_RECORD_ON:                      "record_on",
	CEC_MSG_RECORD_STATUS:                  "record_status",
	CEC_MSG_RECORD_TV_SCREEN:               "record_tv_screen",
	CEC_MSG_CLEAR_ANALOGUE_TIMER:           "clear_analogue_timer",
	CEC_MSG_CLEAR_DIGITAL_TIMER:            "clear_digital_timer",
	CEC_MSG_CLEAR_EXT_TIMER:                "clear_ext_timer",
	CEC_MSG_SET_ANALOGUE_TIMER:             "set_analogue_timer",
	CEC_MSG_SET_DIGITAL_TIMER:              "set_digital_timer",
	CEC_MSG_SET_EXT_TIMER:                  "set_ext_timer",
	CEC_MSG_SET_TIMER_PROGRAM_TITLE:        "set_timer_program_title",
	CEC_MSG_TIMER_CLEARED_STATUS:           "timer_cleared_status",
	CEC_MSG_TIMER_STATUS:                   "timer_status",
	CEC_MSG_CEC_VERSION:                    "cec_version",
	CEC_MSG_GET_CEC_VERSION:                "get_cec_version",
	CEC_MSG_GIVE_PHYSICAL_ADDR:             "give_physical_addr",
	CEC_MSG_GET_MENU_LANGUAGE:              "get_menu_language",
	CEC_MSG_REPORT_PHYSICAL_ADDR:           "report_physical_addr",
	CEC_MSG_SET_MENU_LANGUAGE:              "set_menu_language",
	CEC_MSG_REPORT_FEATURES:                "report_features",
	CEC_MSG_GIVE_FEATURES:                  "give_features",
	CEC_MSG_DECK_CONTROL:                   "deck_control",
	CEC_MSG_DECK_STATUS:                    "deck_status",
	CEC_MSG_GIVE_DECK_STATUS:               "give_deck_status",
	CEC_MSG_PLAY:                           "play",
	CEC_MSG_GIVE_TUNER_DEVICE_STATUS:       "give_tuner_device_status",
	CEC_MSG_SELECT_ANALOGUE_SERVICE:        "select_analogue_service",
	CEC_MSG_SELECT_DIGITAL_SERVICE:         "select_digital_service",
	CEC_MSG_TUNER_DEVICE_STATUS:            "tuner_device_status",
	CEC_MSG_TUNER_STEP_DECREMENT:           "tuner_step_decrement",
	CEC_MSG_TUNER_STEP_INCREMENT:           "tuner_step_increment",
	CEC_MSG_DEVICE_VENDOR_ID:               "device_vendor_id",
	CEC_MSG_GIVE_DEVICE_VENDOR_ID:          "give_device_vendor_id",
	CEC_MSG_VENDOR_COMMAND:                 "vendor_command",
	CEC_MSG_VENDOR_COMMAND_WITH_ID:         "vendor_command_with_id",
	CEC_MSG_VENDOR_REMOTE_BUTTON_DOWN:      "vendor_remote_button_down",
	CEC_MSG_VENDOR_REMOTE_BUTTON_UP:        "vendor_remote_button_up",
	CEC_MSG_SET_OSD_STRING:                 "set_osd_string",
	CEC_MSG_GIVE_OSD_NAME:                  "give_osd_name",
	CEC_MSG_SET_OSD_NAME:                   "set_osd_name",
	CEC_MSG_MENU_REQUEST:                   "menu_request",
	CEC_MSG_MENU_STATUS:                    "menu_status",
	CEC_MSG_USER_CONTROL_PRESSED:           "user_control_pressed",
	CEC_MSG_USER_CONTROL_RELEASED:          "user_control_released",
	CEC_MSG_GIVE_DEVICE_POWER_STATUS:       "give_device_power_status",
	CEC_MSG_REPORT_POWER_STATUS:            "report_power_status",
	CEC_MSG_FEATURE_ABORT:                  "feature_abort",
	CEC_MSG_ABORT:                          "abort",
	CEC_MSG_GIVE_AUDIO_STATUS:              "give_audio_status",
	CEC_MSG_GIVE_SYSTEM_AUDIO_MODE_STATUS:  "give_system_audio_mode_status",
	CEC_MSG_REPORT_AUDIO_STATUS:            "report_audio_status",
	CEC_MSG_REPORT_SHORT_AUDIO_DESCRIPTOR:  "report_short_audio_descriptor",
	CEC_MSG_REQUEST_SHORT_AUDIO_DESCRIPTOR: "request_short_audio_descriptor",
	CEC_MSG_SET_SYSTEM_AUDIO_MODE:          "set_system_audio_mode",
	CEC_MSG_SYSTEM_AUDIO_MODE_REQUEST:      "system_audio_mode_request",
	CEC_MSG_SYSTEM_AUDIO_MODE_STATUS:       "system_audio_mode_status",
	CEC_MSG_SET_AUDIO_VOLUME_LEVEL:         "set_audio_volume_level",
	CEC_MSG_SET_AUDIO_RATE:                 "set_audio_rate",
	CEC_MSG_INITIATE_ARC:                   "initiate_arc",
	CEC_MSG_REPORT_ARC_INITIATED:           "report_arc_initiated",
	CEC_MSG_REPORT_ARC_TERMINATED:          "report_arc_terminated",
	CEC_MSG_REQUEST_ARC_INITIATION:         "request_arc_initiation",
	CEC_MSG_REQUEST_ARC_TERMINATION:        "request_arc_termination",
	CEC_MSG_TERMINATE_ARC:                  "terminate_arc",
	CEC_MSG_REQUEST_CURRENT_LATENCY:        "request_current_latency",
	CEC_MSG_REPORT_CURRENT_LATENCY:         "report_current_latency",
	CEC_MSG_CDC_MESSAGE:                    "cdc_message",
}

var cecUiCmdToString = map[uint8]string{
	CEC_OP_UI_CMD_SELECT:                       "select",
	CEC_OP_UI_CMD_UP:                           "up",
	CEC_OP_UI_CMD_DOWN:                         "down",
	CEC_OP_UI_CMD_LEFT:                         "left",
	CEC_OP_UI_CMD_RIGHT:                        "right",
	CEC_OP_UI_CMD_RIGHT_UP:                     "right_up",
	CEC_OP_UI_CMD_RIGHT_DOWN:                   "right_down",
	CEC_OP_UI_CMD_LEFT_UP:                      "left_up",
	CEC_OP_UI_CMD_LEFT_DOWN:                    "left_down",
	CEC_OP_UI_CMD_DEVICE_ROOT_MENU:             "device_root_menu",
	CEC_OP_UI_CMD_DEVICE_SETUP_MENU:            "device_setup_menu",
	CEC_OP_UI_CMD_CONTENTS_MENU:                "contents_menu",
	CEC_OP_UI_CMD_FAVORITE_MENU:                "favorite_menu",
	CEC_OP_UI_CMD_BACK:                         "back",
	CEC_OP_UI_CMD_MEDIA_TOP_MENU:               "media_top_menu",
	CEC_OP_UI_CMD_MEDIA_CONTEXT_SENSITIVE_MENU: "media_context_sensitive_menu",
	CEC_OP_UI_CMD_NUMBER_ENTRY_MODE:            "number_entry_mode",
	CEC_OP_UI_CMD_NUMBER_11:                    "number_11",
	CEC_OP_UI_CMD_NUMBER_12:                    "number_12",
	CEC_OP_UI_CMD_NUMBER_0_OR_NUMBER_10:        "number_0_or_number_10",
	CEC_OP_UI_CMD_NUMBER_1:                     "number_1",
	CEC_OP_UI_CMD_NUMBER_2:                     "number_2",
	CEC_OP_UI_CMD_NUMBER_3:                     "number_3",
	CEC_OP_UI_CMD_NUMBER_4:                     "number_4",
	CEC_OP_UI_CMD_NUMBER_5:                     "number_5",
	CEC_OP_UI_CMD_NUMBER_6:                     "number_6",
	CEC_OP_UI_CMD_NUMBER_7:                     "number_7",
	CEC_OP_UI_CMD_NUMBER_8:                     "number_8",
	CEC_OP_UI_CMD_NUMBER_9:                     "number_9",
	CEC_OP_UI_CMD_DOT:                          "dot",
	CEC_OP_UI_CMD_ENTER:                        "enter",
	CEC_OP_UI_CMD_CLEAR:                        "clear",
	CEC_OP_UI_CMD_NEXT_FAVORITE:                "next_favorite",
	CEC_OP_UI_CMD_CHANNEL_UP:                   "channel_up",
	CEC_OP_UI_CMD_CHANNEL_DOWN:                 "channel_down",
	CEC_OP_UI_CMD_PREVIOUS_CHANNEL:             "previous_channel",
	CEC_OP_UI_CMD_SOUND_SELECT:                 "sound_select",
	CEC_OP_UI_CMD_INPUT_SELECT:                 "input_select",
	CEC_OP_UI_CMD_DISPLAY_INFORMATION:          "display_information",
	CEC_OP_UI_CMD_HELP:                         "help",
	CEC_OP_UI_CMD_PAGE_UP:                      "page_up",
	CEC_OP_UI_CMD_PAGE_DOWN:                    "page_down",
	CEC_OP_UI_CMD_POWER:                        "power",
	CEC_OP_UI_CMD_VOLUME_UP:                    "volume_up",
	CEC_OP_UI_CMD_VOLUME_DOWN:                  "volume_down",
	CEC_OP_UI_CMD_MUTE:                         "mute",
	CEC_OP_UI_CMD_PLAY:                         "play",
	CEC_OP_UI_CMD_STOP:                         "stop",
	CEC_OP_UI_CMD_PAUSE:                        "pause",
	CEC_OP_UI_CMD_RECORD:                       "record",
	CEC_OP_UI_CMD_REWIND:                       "rewind",
	CEC_OP_UI_CMD_FAST_FORWARD:                 "fast_forward",
	CEC_OP_UI_CMD_EJECT:                        "eject",
	CEC_OP_UI_CMD_SKIP_FORWARD:                 "skip_forward",
	CEC_OP_UI_CMD_SKIP_BACKWARD:                "skip_backward",
	CEC_OP_UI_CMD_STOP_RECORD:                  "stop_record",
	CEC_OP_UI_CMD_PAUSE_RECORD:                 "pause_record",
	CEC_OP_UI_CMD_ANGLE:                        "angle",
	CEC_OP_UI_CMD_SUB_PICTURE:                  "sub_picture",
	CEC_OP_UI_CMD_VIDEO_ON_DEMAND:              "video_on_demand",
	CEC_OP_UI_CMD_ELECTRONIC_PROGRAM_GUIDE:     "electronic_program_guide",
	CEC_OP_UI_CMD_TIMER_PROGRAMMING:            "timer_programming",
	CEC_OP_UI_CMD_INITIAL_CONFIGURATION:        "initial_configuration",
	CEC_OP_UI_CMD_SELECT_BROADCAST_TYPE:        "select_broadcast_type",
	CEC_OP_UI_CMD_SELECT_SOUND_PRESENTATION:    "select_sound_presentation",
	CEC_OP_UI_CMD_AUDIO_DESCRIPTION:            "audio_description",
	CEC_OP_UI_CMD_INTERNET:                     "internet",
	CEC_OP_UI_CMD_3D_MODE:                      "3d_mode",
	CEC_OP_UI_CMD_PLAY_FUNCTION:                "play_function",
	CEC_OP_UI_CMD_PAUSE_PLAY_FUNCTION:          "pause_play_function",
	CEC_OP_UI_CMD_RECORD_FUNCTION:              "record_function",
	CEC_OP_UI_CMD_PAUSE_RECORD_FUNCTION:        "pause_record_function",
	CEC_OP_UI_CMD_STOP_FUNCTION:                "stop_function",
	CEC_OP_UI_CMD_MUTE_FUNCTION:                "mute_function",
	CEC_OP_UI_CMD_RESTORE_VOLUME_FUNCTION:      "restore_volume_function",
	CEC_OP_UI_CMD_TUNE_FUNCTION:                "tune_function",
	CEC_OP_UI_CMD_SELECT_MEDIA_FUNCTION:        "select_media_function",
	CEC_OP_UI_CMD_SELECT_AV_INPUT_FUNCTION:     "select_av_input_function",
	CEC_OP_UI_CMD_SELECT_AUDIO_INPUT_FUNCTION:  "select_audio_input_function",
	CEC_OP_UI_CMD_POWER_TOGGLE_FUNCTION:        "power_toggle_function",
	CEC_OP_UI_CMD_POWER_OFF_FUNCTION:           "power_off_function",
	CEC_OP_UI_CMD_POWER_ON_FUNCTION:            "power_on_function",
	CEC_OP_UI_CMD_F1_BLUE:                      "f1_blue",
	CEC_OP_UI_CMD_F2_RED:                       "f2_red",
	CEC_OP_UI_CMD_F3_GREEN:                     "f3_green",
	CEC_OP_UI_CMD_F4_YELLOW:                    "f4_yellow",
	CEC_OP_UI_CMD_F5:                           "f5",
	CEC_OP_UI_CMD_DATA:                         "data",
}

var Aliases = map[string]CecMsg{
	"av-input-0":               {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x00}},
	"av-input-1":               {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x01}},
	"av-input-2":               {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x02}},
	"av-input-3":               {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x03}},
	"av-input-4":               {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x04}},
	"av-input-5":               {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x05}},
	"av-input-6":               {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x06}},
	"av-input-7":               {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x07}},
	"av-input-8":               {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x08}},
	"av-input-9":               {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x09}},
	"av-input-10":              {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x0a}},
	"av-input-11":              {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x0b}},
	"av-input-12":              {Len: 4, Message: CecMsgMessage{0x00, 0x44, 0x69, 0x0c}},
	"ui-cmd-power-off":         {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x6c}},
	"ui-cmd-power-on":          {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x6d}},
	"ui-cmd-volume-up":         {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x41}},
	"ui-cmd-volume-down":       {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x42}},
	"ui-cmd-select":            {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x00}},
	"ui-cmd-up":                {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x01}},
	"ui-cmd-down":              {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x02}},
	"ui-cmd-left":              {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x03}},
	"ui-cmd-right":             {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x04}},
	"ui-cmd-back":              {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x0d}},
	"ui-cmd-enter":             {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x2b}},
	"ui-cmd-root-menu":         {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x09}},
	"ui-cmd-exit":              {Len: 3, Message: CecMsgMessage{0x00, 0x44, 0x0d}},
	"ui-cmd-release":           {Len: 2, Message: CecMsgMessage{0x00, 0x45}},
	"standby":                  {Len: 2, Message: CecMsgMessage{0x00, 0x36}},
	"give-osd-name":            {Len: 2, Message: CecMsgMessage{0x00, 0x46}},
	"give-device-power-status": {Len: 2, Message: CecMsgMessage{0x00, CEC_MSG_GIVE_DEVICE_POWER_STATUS}},
}

var helperIsVerbose = false

func SetVerbose(v bool) {
	helperIsVerbose = v
}

func (cm *CecMsg) String() string {
	return cm.StringHelper(helperIsVerbose)
}

func (cm *CecMsg) StringHelper(verbose bool) string {
	var s strings.Builder
	from := (cm.Message[0] >> 4) & 0xf
	to := (cm.Message[0] >> 0) & 0xf
	if cm.TxTimestamp != 0 {
		s.WriteString("tx ")
	} else if cm.RxTimestamp != 0 {
		s.WriteString("rx ")
	}
	s.WriteString(fmt.Sprintf("%d->%d", from, to))
	if cm.Len > 1 {
		s.WriteRune(' ')
		pretty, found := cecMsgToString[cm.Message[1]]
		if !found {
			pretty = fmt.Sprintf("%02x", cm.Message[1])
		}
		s.WriteString(pretty)
		if cm.Len == 2 {
			// do nothing
		} else if cm.Message[1] == CEC_MSG_USER_CONTROL_PRESSED {
			s.WriteRune(' ')
			s.WriteString(cecUiCmdToString[cm.Message[2]])
			if cm.Len > 3 && cm.Message[2] == CEC_OP_UI_CMD_SELECT_AV_INPUT_FUNCTION {
				s.WriteString(fmt.Sprintf(" %d", cm.Message[3]))
			}
		} else if cm.Message[1] == CEC_MSG_SET_OSD_NAME {
			s.WriteRune(' ')
			name := string(cm.Message[2:cm.Len])
			s.WriteString(name)
		} else if cm.Message[1] == CEC_MSG_REPORT_POWER_STATUS {
			switch cm.Message[2] {
			case CEC_OP_POWER_STATUS_ON:
				s.WriteString(" on")
			case CEC_OP_POWER_STATUS_STANDBY:
				s.WriteString(" standby")
			case CEC_OP_POWER_STATUS_TO_ON:
				s.WriteString(" to-on")
			case CEC_OP_POWER_STATUS_TO_STANDBY:
				s.WriteString(" to-standby")
			default:
				s.WriteString(" unknown")
			}
		}
	}
	if verbose {
		s.WriteString(" (")
		for i := 0; i < int(cm.Len); i++ {
			s.WriteString(fmt.Sprintf("%02x", cm.Message[i]))
			if i < int(cm.Len)-1 {
				s.WriteRune(':')
			}
		}
		s.WriteString(")")
	}
	if cm.TxTimestamp != 0 {
		s.WriteString(" | ")
		if (cm.TxStatus & CEC_TX_STATUS_OK) != 0 {
			s.WriteString("ok ")
		}
		if (cm.TxStatus & CEC_TX_STATUS_ARB_LOST) != 0 {
			s.WriteString("arb ")
		}
		if (cm.TxStatus & CEC_TX_STATUS_NACK) != 0 {
			s.WriteString("nack ")
		}
		if (cm.TxStatus & CEC_TX_STATUS_LOW_DRIVE) != 0 {
			s.WriteString("low ")
		}
		if (cm.TxStatus & CEC_TX_STATUS_ERROR) != 0 {
			s.WriteString("err ")
		}
		if (cm.TxStatus & CEC_TX_STATUS_MAX_RETRIES) != 0 {
			s.WriteString("rtr ")
		}
		if (cm.TxStatus & CEC_TX_STATUS_ABORTED) != 0 {
			s.WriteString("abrt ")
		}
		if (cm.TxStatus & CEC_TX_STATUS_TIMEOUT) != 0 {
			s.WriteString("time ")
		}
	} else if cm.RxTimestamp != 0 {
		s.WriteString(" | ")
		if (cm.RxStatus & CEC_RX_STATUS_OK) != 0 {
			s.WriteString("ok ")
		}
		if (cm.RxStatus & CEC_RX_STATUS_TIMEOUT) != 0 {
			s.WriteString("time ")
		}
		if (cm.RxStatus & CEC_RX_STATUS_FEATURE_ABORT) != 0 {
			s.WriteString("feat ")
		}
		if (cm.RxStatus & CEC_RX_STATUS_ABORTED) != 0 {
			s.WriteString("abrt ")
		}
	}
	return s.String()
}
