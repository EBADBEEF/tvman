package cec

import (
	"tvman/ioctl"
	"unsafe"
)

const (
	CEC_MAX_MSG_SIZE                            = 16
	CEC_MSG_FL_REPLY_TO_FOLLOWERS               = (1 << 0)
	CEC_MSG_FL_RAW                              = (1 << 1)
	CEC_TX_STATUS_OK                            = (1 << 0)
	CEC_TX_STATUS_ARB_LOST                      = (1 << 1)
	CEC_TX_STATUS_NACK                          = (1 << 2)
	CEC_TX_STATUS_LOW_DRIVE                     = (1 << 3)
	CEC_TX_STATUS_ERROR                         = (1 << 4)
	CEC_TX_STATUS_MAX_RETRIES                   = (1 << 5)
	CEC_TX_STATUS_ABORTED                       = (1 << 6)
	CEC_TX_STATUS_TIMEOUT                       = (1 << 7)
	CEC_RX_STATUS_OK                            = (1 << 0)
	CEC_RX_STATUS_TIMEOUT                       = (1 << 1)
	CEC_RX_STATUS_FEATURE_ABORT                 = (1 << 2)
	CEC_RX_STATUS_ABORTED                       = (1 << 3)
	CEC_LOG_ADDR_INVALID                        = 0xff
	CEC_PHYS_ADDR_INVALID                       = 0xffff
	CEC_MAX_LOG_ADDRS                           = 4
	CEC_LOG_ADDR_TV                             = 0
	CEC_LOG_ADDR_RECORD_1                       = 1
	CEC_LOG_ADDR_RECORD_2                       = 2
	CEC_LOG_ADDR_TUNER_1                        = 3
	CEC_LOG_ADDR_PLAYBACK_1                     = 4
	CEC_LOG_ADDR_AUDIOSYSTEM                    = 5
	CEC_LOG_ADDR_TUNER_2                        = 6
	CEC_LOG_ADDR_TUNER_3                        = 7
	CEC_LOG_ADDR_PLAYBACK_2                     = 8
	CEC_LOG_ADDR_RECORD_3                       = 9
	CEC_LOG_ADDR_TUNER_4                        = 10
	CEC_LOG_ADDR_PLAYBACK_3                     = 11
	CEC_LOG_ADDR_BACKUP_1                       = 12
	CEC_LOG_ADDR_BACKUP_2                       = 13
	CEC_LOG_ADDR_SPECIFIC                       = 14
	CEC_LOG_ADDR_UNREGISTERED                   = 15 /* as initiator address */
	CEC_LOG_ADDR_BROADCAST                      = 15 /* as destination address */
	CEC_LOG_ADDR_TYPE_TV                        = 0
	CEC_LOG_ADDR_TYPE_RECORD                    = 1
	CEC_LOG_ADDR_TYPE_TUNER                     = 2
	CEC_LOG_ADDR_TYPE_PLAYBACK                  = 3
	CEC_LOG_ADDR_TYPE_AUDIOSYSTEM               = 4
	CEC_LOG_ADDR_TYPE_SPECIFIC                  = 5
	CEC_LOG_ADDR_TYPE_UNREGISTERED              = 6
	CEC_LOG_ADDR_MASK_TV                        = (1 << CEC_LOG_ADDR_TV)
	CEC_LOG_ADDR_MASK_RECORD                    = ((1 << CEC_LOG_ADDR_RECORD_1) | (1 << CEC_LOG_ADDR_RECORD_2) | (1 << CEC_LOG_ADDR_RECORD_3))
	CEC_LOG_ADDR_MASK_TUNER                     = ((1 << CEC_LOG_ADDR_TUNER_1) | (1 << CEC_LOG_ADDR_TUNER_2) | (1 << CEC_LOG_ADDR_TUNER_3) | (1 << CEC_LOG_ADDR_TUNER_4))
	CEC_LOG_ADDR_MASK_PLAYBACK                  = ((1 << CEC_LOG_ADDR_PLAYBACK_1) | (1 << CEC_LOG_ADDR_PLAYBACK_2) | (1 << CEC_LOG_ADDR_PLAYBACK_3))
	CEC_LOG_ADDR_MASK_AUDIOSYSTEM               = (1 << CEC_LOG_ADDR_AUDIOSYSTEM)
	CEC_LOG_ADDR_MASK_BACKUP                    = ((1 << CEC_LOG_ADDR_BACKUP_1) | (1 << CEC_LOG_ADDR_BACKUP_2))
	CEC_LOG_ADDR_MASK_SPECIFIC                  = (1 << CEC_LOG_ADDR_SPECIFIC)
	CEC_LOG_ADDR_MASK_UNREGISTERED              = (1 << CEC_LOG_ADDR_UNREGISTERED)
	CEC_VENDOR_ID_NONE                          = 0xffffffff
	CEC_MODE_NO_INITIATOR                       = (0x0 << 0)
	CEC_MODE_INITIATOR                          = (0x1 << 0)
	CEC_MODE_EXCL_INITIATOR                     = (0x2 << 0)
	CEC_MODE_INITIATOR_MSK                      = 0x0f
	CEC_MODE_NO_FOLLOWER                        = (0x0 << 4)
	CEC_MODE_FOLLOWER                           = (0x1 << 4)
	CEC_MODE_EXCL_FOLLOWER                      = (0x2 << 4)
	CEC_MODE_EXCL_FOLLOWER_PASSTHRU             = (0x3 << 4)
	CEC_MODE_MONITOR_PIN                        = (0xd << 4)
	CEC_MODE_MONITOR                            = (0xe << 4)
	CEC_MODE_MONITOR_ALL                        = (0xf << 4)
	CEC_MODE_FOLLOWER_MSK                       = 0xf0
	CEC_CAP_PHYS_ADDR                           = (1 << 0)
	CEC_CAP_LOG_ADDRS                           = (1 << 1)
	CEC_CAP_TRANSMIT                            = (1 << 2)
	CEC_CAP_PASSTHROUGH                         = (1 << 3)
	CEC_CAP_RC                                  = (1 << 4)
	CEC_CAP_MONITOR_ALL                         = (1 << 5)
	CEC_CAP_NEEDS_HPD                           = (1 << 6)
	CEC_CAP_MONITOR_PIN                         = (1 << 7)
	CEC_CAP_CONNECTOR_INFO                      = (1 << 8)
	CEC_LOG_ADDRS_FL_ALLOW_UNREG_FALLBACK       = (1 << 0)
	CEC_LOG_ADDRS_FL_ALLOW_RC_PASSTHRU          = (1 << 1)
	CEC_LOG_ADDRS_FL_CDC_ONLY                   = (1 << 2)
	CEC_MSG_ACTIVE_SOURCE                       = 0x82
	CEC_MSG_IMAGE_VIEW_ON                       = 0x04
	CEC_MSG_TEXT_VIEW_ON                        = 0x0d
	CEC_MSG_INACTIVE_SOURCE                     = 0x9d
	CEC_MSG_REQUEST_ACTIVE_SOURCE               = 0x85
	CEC_MSG_ROUTING_CHANGE                      = 0x80
	CEC_MSG_ROUTING_INFORMATION                 = 0x81
	CEC_MSG_SET_STREAM_PATH                     = 0x86
	CEC_MSG_STANDBY                             = 0x36
	CEC_MSG_RECORD_OFF                          = 0x0b
	CEC_MSG_RECORD_ON                           = 0x09
	CEC_OP_RECORD_SRC_OWN                       = 1
	CEC_OP_RECORD_SRC_DIGITAL                   = 2
	CEC_OP_RECORD_SRC_ANALOG                    = 3
	CEC_OP_RECORD_SRC_EXT_PLUG                  = 4
	CEC_OP_RECORD_SRC_EXT_PHYS_ADDR             = 5
	CEC_OP_SERVICE_ID_METHOD_BY_DIG_ID          = 0
	CEC_OP_SERVICE_ID_METHOD_BY_CHANNEL         = 1
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_ARIB_GEN    = 0x00
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_ATSC_GEN    = 0x01
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_DVB_GEN     = 0x02
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_ARIB_BS     = 0x08
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_ARIB_CS     = 0x09
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_ARIB_T      = 0x0a
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_ATSC_CABLE  = 0x10
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_ATSC_SAT    = 0x11
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_ATSC_T      = 0x12
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_DVB_C       = 0x18
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_DVB_S       = 0x19
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_DVB_S2      = 0x1a
	CEC_OP_DIG_SERVICE_BCAST_SYSTEM_DVB_T       = 0x1b
	CEC_OP_ANA_BCAST_TYPE_CABLE                 = 0
	CEC_OP_ANA_BCAST_TYPE_SATELLITE             = 1
	CEC_OP_ANA_BCAST_TYPE_TERRESTRIAL           = 2
	CEC_OP_BCAST_SYSTEM_PAL_BG                  = 0x00
	CEC_OP_BCAST_SYSTEM_SECAM_LQ                = 0x01 /* SECAM L' */
	CEC_OP_BCAST_SYSTEM_PAL_M                   = 0x02
	CEC_OP_BCAST_SYSTEM_NTSC_M                  = 0x03
	CEC_OP_BCAST_SYSTEM_PAL_I                   = 0x04
	CEC_OP_BCAST_SYSTEM_SECAM_DK                = 0x05
	CEC_OP_BCAST_SYSTEM_SECAM_BG                = 0x06
	CEC_OP_BCAST_SYSTEM_SECAM_L                 = 0x07
	CEC_OP_BCAST_SYSTEM_PAL_DK                  = 0x08
	CEC_OP_BCAST_SYSTEM_OTHER                   = 0x1f
	CEC_OP_CHANNEL_NUMBER_FMT_1_PART            = 0x01
	CEC_OP_CHANNEL_NUMBER_FMT_2_PART            = 0x02
	CEC_MSG_RECORD_STATUS                       = 0x0a
	CEC_OP_RECORD_STATUS_CUR_SRC                = 0x01
	CEC_OP_RECORD_STATUS_DIG_SERVICE            = 0x02
	CEC_OP_RECORD_STATUS_ANA_SERVICE            = 0x03
	CEC_OP_RECORD_STATUS_EXT_INPUT              = 0x04
	CEC_OP_RECORD_STATUS_NO_DIG_SERVICE         = 0x05
	CEC_OP_RECORD_STATUS_NO_ANA_SERVICE         = 0x06
	CEC_OP_RECORD_STATUS_NO_SERVICE             = 0x07
	CEC_OP_RECORD_STATUS_INVALID_EXT_PLUG       = 0x09
	CEC_OP_RECORD_STATUS_INVALID_EXT_PHYS_ADDR  = 0x0a
	CEC_OP_RECORD_STATUS_UNSUP_CA               = 0x0b
	CEC_OP_RECORD_STATUS_NO_CA_ENTITLEMENTS     = 0x0c
	CEC_OP_RECORD_STATUS_CANT_COPY_SRC          = 0x0d
	CEC_OP_RECORD_STATUS_NO_MORE_COPIES         = 0x0e
	CEC_OP_RECORD_STATUS_NO_MEDIA               = 0x10
	CEC_OP_RECORD_STATUS_PLAYING                = 0x11
	CEC_OP_RECORD_STATUS_ALREADY_RECORDING      = 0x12
	CEC_OP_RECORD_STATUS_MEDIA_PROT             = 0x13
	CEC_OP_RECORD_STATUS_NO_SIGNAL              = 0x14
	CEC_OP_RECORD_STATUS_MEDIA_PROBLEM          = 0x15
	CEC_OP_RECORD_STATUS_NO_SPACE               = 0x16
	CEC_OP_RECORD_STATUS_PARENTAL_LOCK          = 0x17
	CEC_OP_RECORD_STATUS_TERMINATED_OK          = 0x1a
	CEC_OP_RECORD_STATUS_ALREADY_TERM           = 0x1b
	CEC_OP_RECORD_STATUS_OTHER                  = 0x1f
	CEC_MSG_RECORD_TV_SCREEN                    = 0x0f
	CEC_MSG_CLEAR_ANALOGUE_TIMER                = 0x33
	CEC_OP_REC_SEQ_SUNDAY                       = 0x01
	CEC_OP_REC_SEQ_MONDAY                       = 0x02
	CEC_OP_REC_SEQ_TUESDAY                      = 0x04
	CEC_OP_REC_SEQ_WEDNESDAY                    = 0x08
	CEC_OP_REC_SEQ_THURSDAY                     = 0x10
	CEC_OP_REC_SEQ_FRIDAY                       = 0x20
	CEC_OP_REC_SEQ_SATURDAY                     = 0x40
	CEC_OP_REC_SEQ_ONCE_ONLY                    = 0x00
	CEC_MSG_CLEAR_DIGITAL_TIMER                 = 0x99
	CEC_MSG_CLEAR_EXT_TIMER                     = 0xa1
	CEC_OP_EXT_SRC_PLUG                         = 0x04
	CEC_OP_EXT_SRC_PHYS_ADDR                    = 0x05
	CEC_MSG_SET_ANALOGUE_TIMER                  = 0x34
	CEC_MSG_SET_DIGITAL_TIMER                   = 0x97
	CEC_MSG_SET_EXT_TIMER                       = 0xa2
	CEC_MSG_SET_TIMER_PROGRAM_TITLE             = 0x67
	CEC_MSG_TIMER_CLEARED_STATUS                = 0x43
	CEC_OP_TIMER_CLR_STAT_RECORDING             = 0x00
	CEC_OP_TIMER_CLR_STAT_NO_MATCHING           = 0x01
	CEC_OP_TIMER_CLR_STAT_NO_INFO               = 0x02
	CEC_OP_TIMER_CLR_STAT_CLEARED               = 0x80
	CEC_MSG_TIMER_STATUS                        = 0x35
	CEC_OP_TIMER_OVERLAP_WARNING_NO_OVERLAP     = 0
	CEC_OP_TIMER_OVERLAP_WARNING_OVERLAP        = 1
	CEC_OP_MEDIA_INFO_UNPROT_MEDIA              = 0
	CEC_OP_MEDIA_INFO_PROT_MEDIA                = 1
	CEC_OP_MEDIA_INFO_NO_MEDIA                  = 2
	CEC_OP_PROG_IND_NOT_PROGRAMMED              = 0
	CEC_OP_PROG_IND_PROGRAMMED                  = 1
	CEC_OP_PROG_INFO_ENOUGH_SPACE               = 0x08
	CEC_OP_PROG_INFO_NOT_ENOUGH_SPACE           = 0x09
	CEC_OP_PROG_INFO_MIGHT_NOT_BE_ENOUGH_SPACE  = 0x0b
	CEC_OP_PROG_INFO_NONE_AVAILABLE             = 0x0a
	CEC_OP_PROG_ERROR_NO_FREE_TIMER             = 0x01
	CEC_OP_PROG_ERROR_DATE_OUT_OF_RANGE         = 0x02
	CEC_OP_PROG_ERROR_REC_SEQ_ERROR             = 0x03
	CEC_OP_PROG_ERROR_INV_EXT_PLUG              = 0x04
	CEC_OP_PROG_ERROR_INV_EXT_PHYS_ADDR         = 0x05
	CEC_OP_PROG_ERROR_CA_UNSUPP                 = 0x06
	CEC_OP_PROG_ERROR_INSUF_CA_ENTITLEMENTS     = 0x07
	CEC_OP_PROG_ERROR_RESOLUTION_UNSUPP         = 0x08
	CEC_OP_PROG_ERROR_PARENTAL_LOCK             = 0x09
	CEC_OP_PROG_ERROR_CLOCK_FAILURE             = 0x0a
	CEC_OP_PROG_ERROR_DUPLICATE                 = 0x0e
	CEC_MSG_CEC_VERSION                         = 0x9e
	CEC_OP_CEC_VERSION_1_3A                     = 4
	CEC_OP_CEC_VERSION_1_4                      = 5
	CEC_OP_CEC_VERSION_2_0                      = 6
	CEC_MSG_GET_CEC_VERSION                     = 0x9f
	CEC_MSG_GIVE_PHYSICAL_ADDR                  = 0x83
	CEC_MSG_GET_MENU_LANGUAGE                   = 0x91
	CEC_MSG_REPORT_PHYSICAL_ADDR                = 0x84
	CEC_OP_PRIM_DEVTYPE_TV                      = 0
	CEC_OP_PRIM_DEVTYPE_RECORD                  = 1
	CEC_OP_PRIM_DEVTYPE_TUNER                   = 3
	CEC_OP_PRIM_DEVTYPE_PLAYBACK                = 4
	CEC_OP_PRIM_DEVTYPE_AUDIOSYSTEM             = 5
	CEC_OP_PRIM_DEVTYPE_SWITCH                  = 6
	CEC_OP_PRIM_DEVTYPE_PROCESSOR               = 7
	CEC_MSG_SET_MENU_LANGUAGE                   = 0x32
	CEC_MSG_REPORT_FEATURES                     = 0xa6 /* HDMI 2.0 */
	CEC_OP_ALL_DEVTYPE_TV                       = 0x80
	CEC_OP_ALL_DEVTYPE_RECORD                   = 0x40
	CEC_OP_ALL_DEVTYPE_TUNER                    = 0x20
	CEC_OP_ALL_DEVTYPE_PLAYBACK                 = 0x10
	CEC_OP_ALL_DEVTYPE_AUDIOSYSTEM              = 0x08
	CEC_OP_ALL_DEVTYPE_SWITCH                   = 0x04
	CEC_OP_FEAT_EXT                             = 0x80 /* Extension bit */
	CEC_OP_FEAT_RC_TV_PROFILE_NONE              = 0x00
	CEC_OP_FEAT_RC_TV_PROFILE_1                 = 0x02
	CEC_OP_FEAT_RC_TV_PROFILE_2                 = 0x06
	CEC_OP_FEAT_RC_TV_PROFILE_3                 = 0x0a
	CEC_OP_FEAT_RC_TV_PROFILE_4                 = 0x0e
	CEC_OP_FEAT_RC_SRC_HAS_DEV_ROOT_MENU        = 0x50
	CEC_OP_FEAT_RC_SRC_HAS_DEV_SETUP_MENU       = 0x48
	CEC_OP_FEAT_RC_SRC_HAS_CONTENTS_MENU        = 0x44
	CEC_OP_FEAT_RC_SRC_HAS_MEDIA_TOP_MENU       = 0x42
	CEC_OP_FEAT_RC_SRC_HAS_MEDIA_CONTEXT_MENU   = 0x41
	CEC_OP_FEAT_DEV_HAS_RECORD_TV_SCREEN        = 0x40
	CEC_OP_FEAT_DEV_HAS_SET_OSD_STRING          = 0x20
	CEC_OP_FEAT_DEV_HAS_DECK_CONTROL            = 0x10
	CEC_OP_FEAT_DEV_HAS_SET_AUDIO_RATE          = 0x08
	CEC_OP_FEAT_DEV_SINK_HAS_ARC_TX             = 0x04
	CEC_OP_FEAT_DEV_SOURCE_HAS_ARC_RX           = 0x02
	CEC_OP_FEAT_DEV_HAS_SET_AUDIO_VOLUME_LEVEL  = 0x01
	CEC_MSG_GIVE_FEATURES                       = 0xa5 /* HDMI 2.0 */
	CEC_MSG_DECK_CONTROL                        = 0x42
	CEC_OP_DECK_CTL_MODE_SKIP_FWD               = 1
	CEC_OP_DECK_CTL_MODE_SKIP_REV               = 2
	CEC_OP_DECK_CTL_MODE_STOP                   = 3
	CEC_OP_DECK_CTL_MODE_EJECT                  = 4
	CEC_MSG_DECK_STATUS                         = 0x1b
	CEC_OP_DECK_INFO_PLAY                       = 0x11
	CEC_OP_DECK_INFO_RECORD                     = 0x12
	CEC_OP_DECK_INFO_PLAY_REV                   = 0x13
	CEC_OP_DECK_INFO_STILL                      = 0x14
	CEC_OP_DECK_INFO_SLOW                       = 0x15
	CEC_OP_DECK_INFO_SLOW_REV                   = 0x16
	CEC_OP_DECK_INFO_FAST_FWD                   = 0x17
	CEC_OP_DECK_INFO_FAST_REV                   = 0x18
	CEC_OP_DECK_INFO_NO_MEDIA                   = 0x19
	CEC_OP_DECK_INFO_STOP                       = 0x1a
	CEC_OP_DECK_INFO_SKIP_FWD                   = 0x1b
	CEC_OP_DECK_INFO_SKIP_REV                   = 0x1c
	CEC_OP_DECK_INFO_INDEX_SEARCH_FWD           = 0x1d
	CEC_OP_DECK_INFO_INDEX_SEARCH_REV           = 0x1e
	CEC_OP_DECK_INFO_OTHER                      = 0x1f
	CEC_MSG_GIVE_DECK_STATUS                    = 0x1a
	CEC_OP_STATUS_REQ_ON                        = 1
	CEC_OP_STATUS_REQ_OFF                       = 2
	CEC_OP_STATUS_REQ_ONCE                      = 3
	CEC_MSG_PLAY                                = 0x41
	CEC_OP_PLAY_MODE_PLAY_FWD                   = 0x24
	CEC_OP_PLAY_MODE_PLAY_REV                   = 0x20
	CEC_OP_PLAY_MODE_PLAY_STILL                 = 0x25
	CEC_OP_PLAY_MODE_PLAY_FAST_FWD_MIN          = 0x05
	CEC_OP_PLAY_MODE_PLAY_FAST_FWD_MED          = 0x06
	CEC_OP_PLAY_MODE_PLAY_FAST_FWD_MAX          = 0x07
	CEC_OP_PLAY_MODE_PLAY_FAST_REV_MIN          = 0x09
	CEC_OP_PLAY_MODE_PLAY_FAST_REV_MED          = 0x0a
	CEC_OP_PLAY_MODE_PLAY_FAST_REV_MAX          = 0x0b
	CEC_OP_PLAY_MODE_PLAY_SLOW_FWD_MIN          = 0x15
	CEC_OP_PLAY_MODE_PLAY_SLOW_FWD_MED          = 0x16
	CEC_OP_PLAY_MODE_PLAY_SLOW_FWD_MAX          = 0x17
	CEC_OP_PLAY_MODE_PLAY_SLOW_REV_MIN          = 0x19
	CEC_OP_PLAY_MODE_PLAY_SLOW_REV_MED          = 0x1a
	CEC_OP_PLAY_MODE_PLAY_SLOW_REV_MAX          = 0x1b
	CEC_MSG_GIVE_TUNER_DEVICE_STATUS            = 0x08
	CEC_MSG_SELECT_ANALOGUE_SERVICE             = 0x92
	CEC_MSG_SELECT_DIGITAL_SERVICE              = 0x93
	CEC_MSG_TUNER_DEVICE_STATUS                 = 0x07
	CEC_OP_REC_FLAG_NOT_USED                    = 0
	CEC_OP_REC_FLAG_USED                        = 1
	CEC_OP_TUNER_DISPLAY_INFO_DIGITAL           = 0
	CEC_OP_TUNER_DISPLAY_INFO_NONE              = 1
	CEC_OP_TUNER_DISPLAY_INFO_ANALOGUE          = 2
	CEC_MSG_TUNER_STEP_DECREMENT                = 0x06
	CEC_MSG_TUNER_STEP_INCREMENT                = 0x05
	CEC_MSG_DEVICE_VENDOR_ID                    = 0x87
	CEC_MSG_GIVE_DEVICE_VENDOR_ID               = 0x8c
	CEC_MSG_VENDOR_COMMAND                      = 0x89
	CEC_MSG_VENDOR_COMMAND_WITH_ID              = 0xa0
	CEC_MSG_VENDOR_REMOTE_BUTTON_DOWN           = 0x8a
	CEC_MSG_VENDOR_REMOTE_BUTTON_UP             = 0x8b
	CEC_MSG_SET_OSD_STRING                      = 0x64
	CEC_OP_DISP_CTL_DEFAULT                     = 0x00
	CEC_OP_DISP_CTL_UNTIL_CLEARED               = 0x40
	CEC_OP_DISP_CTL_CLEAR                       = 0x80
	CEC_MSG_GIVE_OSD_NAME                       = 0x46
	CEC_MSG_SET_OSD_NAME                        = 0x47
	CEC_MSG_MENU_REQUEST                        = 0x8d
	CEC_OP_MENU_REQUEST_ACTIVATE                = 0x00
	CEC_OP_MENU_REQUEST_DEACTIVATE              = 0x01
	CEC_OP_MENU_REQUEST_QUERY                   = 0x02
	CEC_MSG_MENU_STATUS                         = 0x8e
	CEC_OP_MENU_STATE_ACTIVATED                 = 0x00
	CEC_OP_MENU_STATE_DEACTIVATED               = 0x01
	CEC_MSG_USER_CONTROL_PRESSED                = 0x44
	CEC_OP_UI_CMD_SELECT                        = 0x00
	CEC_OP_UI_CMD_UP                            = 0x01
	CEC_OP_UI_CMD_DOWN                          = 0x02
	CEC_OP_UI_CMD_LEFT                          = 0x03
	CEC_OP_UI_CMD_RIGHT                         = 0x04
	CEC_OP_UI_CMD_RIGHT_UP                      = 0x05
	CEC_OP_UI_CMD_RIGHT_DOWN                    = 0x06
	CEC_OP_UI_CMD_LEFT_UP                       = 0x07
	CEC_OP_UI_CMD_LEFT_DOWN                     = 0x08
	CEC_OP_UI_CMD_DEVICE_ROOT_MENU              = 0x09
	CEC_OP_UI_CMD_DEVICE_SETUP_MENU             = 0x0a
	CEC_OP_UI_CMD_CONTENTS_MENU                 = 0x0b
	CEC_OP_UI_CMD_FAVORITE_MENU                 = 0x0c
	CEC_OP_UI_CMD_BACK                          = 0x0d
	CEC_OP_UI_CMD_MEDIA_TOP_MENU                = 0x10
	CEC_OP_UI_CMD_MEDIA_CONTEXT_SENSITIVE_MENU  = 0x11
	CEC_OP_UI_CMD_NUMBER_ENTRY_MODE             = 0x1d
	CEC_OP_UI_CMD_NUMBER_11                     = 0x1e
	CEC_OP_UI_CMD_NUMBER_12                     = 0x1f
	CEC_OP_UI_CMD_NUMBER_0_OR_NUMBER_10         = 0x20
	CEC_OP_UI_CMD_NUMBER_1                      = 0x21
	CEC_OP_UI_CMD_NUMBER_2                      = 0x22
	CEC_OP_UI_CMD_NUMBER_3                      = 0x23
	CEC_OP_UI_CMD_NUMBER_4                      = 0x24
	CEC_OP_UI_CMD_NUMBER_5                      = 0x25
	CEC_OP_UI_CMD_NUMBER_6                      = 0x26
	CEC_OP_UI_CMD_NUMBER_7                      = 0x27
	CEC_OP_UI_CMD_NUMBER_8                      = 0x28
	CEC_OP_UI_CMD_NUMBER_9                      = 0x29
	CEC_OP_UI_CMD_DOT                           = 0x2a
	CEC_OP_UI_CMD_ENTER                         = 0x2b
	CEC_OP_UI_CMD_CLEAR                         = 0x2c
	CEC_OP_UI_CMD_NEXT_FAVORITE                 = 0x2f
	CEC_OP_UI_CMD_CHANNEL_UP                    = 0x30
	CEC_OP_UI_CMD_CHANNEL_DOWN                  = 0x31
	CEC_OP_UI_CMD_PREVIOUS_CHANNEL              = 0x32
	CEC_OP_UI_CMD_SOUND_SELECT                  = 0x33
	CEC_OP_UI_CMD_INPUT_SELECT                  = 0x34
	CEC_OP_UI_CMD_DISPLAY_INFORMATION           = 0x35
	CEC_OP_UI_CMD_HELP                          = 0x36
	CEC_OP_UI_CMD_PAGE_UP                       = 0x37
	CEC_OP_UI_CMD_PAGE_DOWN                     = 0x38
	CEC_OP_UI_CMD_POWER                         = 0x40
	CEC_OP_UI_CMD_VOLUME_UP                     = 0x41
	CEC_OP_UI_CMD_VOLUME_DOWN                   = 0x42
	CEC_OP_UI_CMD_MUTE                          = 0x43
	CEC_OP_UI_CMD_PLAY                          = 0x44
	CEC_OP_UI_CMD_STOP                          = 0x45
	CEC_OP_UI_CMD_PAUSE                         = 0x46
	CEC_OP_UI_CMD_RECORD                        = 0x47
	CEC_OP_UI_CMD_REWIND                        = 0x48
	CEC_OP_UI_CMD_FAST_FORWARD                  = 0x49
	CEC_OP_UI_CMD_EJECT                         = 0x4a
	CEC_OP_UI_CMD_SKIP_FORWARD                  = 0x4b
	CEC_OP_UI_CMD_SKIP_BACKWARD                 = 0x4c
	CEC_OP_UI_CMD_STOP_RECORD                   = 0x4d
	CEC_OP_UI_CMD_PAUSE_RECORD                  = 0x4e
	CEC_OP_UI_CMD_ANGLE                         = 0x50
	CEC_OP_UI_CMD_SUB_PICTURE                   = 0x51
	CEC_OP_UI_CMD_VIDEO_ON_DEMAND               = 0x52
	CEC_OP_UI_CMD_ELECTRONIC_PROGRAM_GUIDE      = 0x53
	CEC_OP_UI_CMD_TIMER_PROGRAMMING             = 0x54
	CEC_OP_UI_CMD_INITIAL_CONFIGURATION         = 0x55
	CEC_OP_UI_CMD_SELECT_BROADCAST_TYPE         = 0x56
	CEC_OP_UI_CMD_SELECT_SOUND_PRESENTATION     = 0x57
	CEC_OP_UI_CMD_AUDIO_DESCRIPTION             = 0x58
	CEC_OP_UI_CMD_INTERNET                      = 0x59
	CEC_OP_UI_CMD_3D_MODE                       = 0x5a
	CEC_OP_UI_CMD_PLAY_FUNCTION                 = 0x60
	CEC_OP_UI_CMD_PAUSE_PLAY_FUNCTION           = 0x61
	CEC_OP_UI_CMD_RECORD_FUNCTION               = 0x62
	CEC_OP_UI_CMD_PAUSE_RECORD_FUNCTION         = 0x63
	CEC_OP_UI_CMD_STOP_FUNCTION                 = 0x64
	CEC_OP_UI_CMD_MUTE_FUNCTION                 = 0x65
	CEC_OP_UI_CMD_RESTORE_VOLUME_FUNCTION       = 0x66
	CEC_OP_UI_CMD_TUNE_FUNCTION                 = 0x67
	CEC_OP_UI_CMD_SELECT_MEDIA_FUNCTION         = 0x68
	CEC_OP_UI_CMD_SELECT_AV_INPUT_FUNCTION      = 0x69
	CEC_OP_UI_CMD_SELECT_AUDIO_INPUT_FUNCTION   = 0x6a
	CEC_OP_UI_CMD_POWER_TOGGLE_FUNCTION         = 0x6b
	CEC_OP_UI_CMD_POWER_OFF_FUNCTION            = 0x6c
	CEC_OP_UI_CMD_POWER_ON_FUNCTION             = 0x6d
	CEC_OP_UI_CMD_F1_BLUE                       = 0x71
	CEC_OP_UI_CMD_F2_RED                        = 0x72
	CEC_OP_UI_CMD_F3_GREEN                      = 0x73
	CEC_OP_UI_CMD_F4_YELLOW                     = 0x74
	CEC_OP_UI_CMD_F5                            = 0x75
	CEC_OP_UI_CMD_DATA                          = 0x76
	CEC_OP_UI_BCAST_TYPE_TOGGLE_ALL             = 0x00
	CEC_OP_UI_BCAST_TYPE_TOGGLE_DIG_ANA         = 0x01
	CEC_OP_UI_BCAST_TYPE_ANALOGUE               = 0x10
	CEC_OP_UI_BCAST_TYPE_ANALOGUE_T             = 0x20
	CEC_OP_UI_BCAST_TYPE_ANALOGUE_CABLE         = 0x30
	CEC_OP_UI_BCAST_TYPE_ANALOGUE_SAT           = 0x40
	CEC_OP_UI_BCAST_TYPE_DIGITAL                = 0x50
	CEC_OP_UI_BCAST_TYPE_DIGITAL_T              = 0x60
	CEC_OP_UI_BCAST_TYPE_DIGITAL_CABLE          = 0x70
	CEC_OP_UI_BCAST_TYPE_DIGITAL_SAT            = 0x80
	CEC_OP_UI_BCAST_TYPE_DIGITAL_COM_SAT        = 0x90
	CEC_OP_UI_BCAST_TYPE_DIGITAL_COM_SAT2       = 0x91
	CEC_OP_UI_BCAST_TYPE_IP                     = 0xa0
	CEC_OP_UI_SND_PRES_CTL_DUAL_MONO            = 0x10
	CEC_OP_UI_SND_PRES_CTL_KARAOKE              = 0x20
	CEC_OP_UI_SND_PRES_CTL_DOWNMIX              = 0x80
	CEC_OP_UI_SND_PRES_CTL_REVERB               = 0x90
	CEC_OP_UI_SND_PRES_CTL_EQUALIZER            = 0xa0
	CEC_OP_UI_SND_PRES_CTL_BASS_UP              = 0xb1
	CEC_OP_UI_SND_PRES_CTL_BASS_NEUTRAL         = 0xb2
	CEC_OP_UI_SND_PRES_CTL_BASS_DOWN            = 0xb3
	CEC_OP_UI_SND_PRES_CTL_TREBLE_UP            = 0xc1
	CEC_OP_UI_SND_PRES_CTL_TREBLE_NEUTRAL       = 0xc2
	CEC_OP_UI_SND_PRES_CTL_TREBLE_DOWN          = 0xc3
	CEC_MSG_USER_CONTROL_RELEASED               = 0x45
	CEC_MSG_GIVE_DEVICE_POWER_STATUS            = 0x8f
	CEC_MSG_REPORT_POWER_STATUS                 = 0x90
	CEC_OP_POWER_STATUS_ON                      = 0
	CEC_OP_POWER_STATUS_STANDBY                 = 1
	CEC_OP_POWER_STATUS_TO_ON                   = 2
	CEC_OP_POWER_STATUS_TO_STANDBY              = 3
	CEC_MSG_FEATURE_ABORT                       = 0x00
	CEC_OP_ABORT_UNRECOGNIZED_OP                = 0
	CEC_OP_ABORT_INCORRECT_MODE                 = 1
	CEC_OP_ABORT_NO_SOURCE                      = 2
	CEC_OP_ABORT_INVALID_OP                     = 3
	CEC_OP_ABORT_REFUSED                        = 4
	CEC_OP_ABORT_UNDETERMINED                   = 5
	CEC_MSG_ABORT                               = 0xff
	CEC_MSG_GIVE_AUDIO_STATUS                   = 0x71
	CEC_MSG_GIVE_SYSTEM_AUDIO_MODE_STATUS       = 0x7d
	CEC_MSG_REPORT_AUDIO_STATUS                 = 0x7a
	CEC_OP_AUD_MUTE_STATUS_OFF                  = 0
	CEC_OP_AUD_MUTE_STATUS_ON                   = 1
	CEC_MSG_REPORT_SHORT_AUDIO_DESCRIPTOR       = 0xa3
	CEC_MSG_REQUEST_SHORT_AUDIO_DESCRIPTOR      = 0xa4
	CEC_MSG_SET_SYSTEM_AUDIO_MODE               = 0x72
	CEC_OP_SYS_AUD_STATUS_OFF                   = 0
	CEC_OP_SYS_AUD_STATUS_ON                    = 1
	CEC_MSG_SYSTEM_AUDIO_MODE_REQUEST           = 0x70
	CEC_MSG_SYSTEM_AUDIO_MODE_STATUS            = 0x7e
	CEC_OP_AUD_FMT_ID_CEA861                    = 0
	CEC_OP_AUD_FMT_ID_CEA861_CXT                = 1
	CEC_MSG_SET_AUDIO_VOLUME_LEVEL              = 0x73
	CEC_MSG_SET_AUDIO_RATE                      = 0x9a
	CEC_OP_AUD_RATE_OFF                         = 0
	CEC_OP_AUD_RATE_WIDE_STD                    = 1
	CEC_OP_AUD_RATE_WIDE_FAST                   = 2
	CEC_OP_AUD_RATE_WIDE_SLOW                   = 3
	CEC_OP_AUD_RATE_NARROW_STD                  = 4
	CEC_OP_AUD_RATE_NARROW_FAST                 = 5
	CEC_OP_AUD_RATE_NARROW_SLOW                 = 6
	CEC_MSG_INITIATE_ARC                        = 0xc0
	CEC_MSG_REPORT_ARC_INITIATED                = 0xc1
	CEC_MSG_REPORT_ARC_TERMINATED               = 0xc2
	CEC_MSG_REQUEST_ARC_INITIATION              = 0xc3
	CEC_MSG_REQUEST_ARC_TERMINATION             = 0xc4
	CEC_MSG_TERMINATE_ARC                       = 0xc5
	CEC_MSG_REQUEST_CURRENT_LATENCY             = 0xa7
	CEC_MSG_REPORT_CURRENT_LATENCY              = 0xa8
	CEC_OP_LOW_LATENCY_MODE_OFF                 = 0
	CEC_OP_LOW_LATENCY_MODE_ON                  = 1
	CEC_OP_AUD_OUT_COMPENSATED_NA               = 0
	CEC_OP_AUD_OUT_COMPENSATED_DELAY            = 1
	CEC_OP_AUD_OUT_COMPENSATED_NO_DELAY         = 2
	CEC_OP_AUD_OUT_COMPENSATED_PARTIAL_DELAY    = 3
	CEC_MSG_CDC_MESSAGE                         = 0xf8
	CEC_MSG_CDC_HEC_INQUIRE_STATE               = 0x00
	CEC_MSG_CDC_HEC_REPORT_STATE                = 0x01
	CEC_OP_HEC_FUNC_STATE_NOT_SUPPORTED         = 0
	CEC_OP_HEC_FUNC_STATE_INACTIVE              = 1
	CEC_OP_HEC_FUNC_STATE_ACTIVE                = 2
	CEC_OP_HEC_FUNC_STATE_ACTIVATION_FIELD      = 3
	CEC_OP_HOST_FUNC_STATE_NOT_SUPPORTED        = 0
	CEC_OP_HOST_FUNC_STATE_INACTIVE             = 1
	CEC_OP_HOST_FUNC_STATE_ACTIVE               = 2
	CEC_OP_ENC_FUNC_STATE_EXT_CON_NOT_SUPPORTED = 0
	CEC_OP_ENC_FUNC_STATE_EXT_CON_INACTIVE      = 1
	CEC_OP_ENC_FUNC_STATE_EXT_CON_ACTIVE        = 2
	CEC_OP_CDC_ERROR_CODE_NONE                  = 0
	CEC_OP_CDC_ERROR_CODE_CAP_UNSUPPORTED       = 1
	CEC_OP_CDC_ERROR_CODE_WRONG_STATE           = 2
	CEC_OP_CDC_ERROR_CODE_OTHER                 = 3
	CEC_OP_HEC_SUPPORT_NO                       = 0
	CEC_OP_HEC_SUPPORT_YES                      = 1
	CEC_OP_HEC_ACTIVATION_ON                    = 0
	CEC_OP_HEC_ACTIVATION_OFF                   = 1
	CEC_MSG_CDC_HEC_SET_STATE_ADJACENT          = 0x02
	CEC_MSG_CDC_HEC_SET_STATE                   = 0x03
	CEC_OP_HEC_SET_STATE_DEACTIVATE             = 0
	CEC_OP_HEC_SET_STATE_ACTIVATE               = 1
	CEC_MSG_CDC_HEC_REQUEST_DEACTIVATION        = 0x04
	CEC_MSG_CDC_HEC_NOTIFY_ALIVE                = 0x05
	CEC_MSG_CDC_HEC_DISCOVER                    = 0x06
	CEC_MSG_CDC_HPD_SET_STATE                   = 0x10
	CEC_OP_HPD_STATE_CP_EDID_DISABLE            = 0
	CEC_OP_HPD_STATE_CP_EDID_ENABLE             = 1
	CEC_OP_HPD_STATE_CP_EDID_DISABLE_ENABLE     = 2
	CEC_OP_HPD_STATE_EDID_DISABLE               = 3
	CEC_OP_HPD_STATE_EDID_ENABLE                = 4
	CEC_OP_HPD_STATE_EDID_DISABLE_ENABLE        = 5
	CEC_MSG_CDC_HPD_REPORT_STATE                = 0x11
	CEC_OP_HPD_ERROR_NONE                       = 0
	CEC_OP_HPD_ERROR_INITIATOR_NOT_CAPABLE      = 1
	CEC_OP_HPD_ERROR_INITIATOR_WRONG_STATE      = 2
	CEC_OP_HPD_ERROR_OTHER                      = 3
	CEC_OP_HPD_ERROR_NONE_NO_VIDEO              = 4
)

type CecMsgMessage [CEC_MAX_MSG_SIZE]uint8

type CecMsg struct {
	TxTimestamp            uint64
	RxTimestamp            uint64
	Len                    uint32
	Timeout                uint32
	Sequence               uint32
	Flags                  uint32
	Message                [CEC_MAX_MSG_SIZE]uint8
	Reply                  uint8
	RxStatus               uint8
	TxStatus               uint8
	TxArbitrationLostCount uint8
	TxNackCount            uint8
	TxLowDriveCount        uint8
	TxErrorCount           uint8
}

type CecCaps struct {
	Driver             [32]byte
	Name               [32]byte
	AvailableLogAddres uint32
	Capabilities       uint32
	Version            uint32
}

type CecDrmConnectorInfo struct {
	CardNumber  uint32
	ConnectorId uint32
}

const (
	CEC_CONNECTOR_TYPE_NO_CONNECTOR = 0
	CEC_CONNECTOR_TYPE_DRM          = 1
)

type CecConnectorInfo struct {
	Type uint32
	Raw  [16]uint32 /* CecDrmConnectorInfo */
}

type CecOsdName [15]byte

type CecLogAddrs struct {
	LogAddr           [CEC_MAX_LOG_ADDRS]uint8
	LogAddrMask       uint16
	CecVersion        uint8
	NumLogAddrs       uint8
	VendorId          uint32
	Flags             uint32
	OsdName           CecOsdName
	PrimaryDeviceType [CEC_MAX_LOG_ADDRS]uint8
	LogAddrType       [CEC_MAX_LOG_ADDRS]uint8
	/* CEC 2.0 */
	AllDeviceTypes [CEC_MAX_LOG_ADDRS]uint8
	Features       [CEC_MAX_LOG_ADDRS][12]uint8
}

const (
	CEC_EVENT_STATE_CHANGE      = 1
	CEC_EVENT_LOST_MSGS         = 2
	CEC_EVENT_PIN_CEC_LOW       = 3
	CEC_EVENT_PIN_CEC_HIGH      = 4
	CEC_EVENT_PIN_HPD_LOW       = 5
	CEC_EVENT_PIN_HPD_HIGH      = 6
	CEC_EVENT_PIN_5V_LOW        = 7
	CEC_EVENT_PIN_5V_HIGH       = 8
	CEC_EVENT_FL_INITIAL_STATE  = (1 << 0)
	CEC_EVENT_FL_DROPPED_EVENTS = (1 << 1)
)

type CecEvent struct {
	Timestamp uint64
	Event     uint32
	Flags     uint32
	Raw       [16]uint32 /* CecEventStateChange, CecEventLostMsgs */
}

var (
	CEC_ADAP_G_CAPS           = uintptr(ioctl.IOWR('a', 0, unsafe.Sizeof(CecCaps{})))
	CEC_ADAP_G_PHYS_ADDR      = uintptr(ioctl.IOR('a', 1, unsafe.Sizeof(uint16(0))))
	CEC_ADAP_S_PHYS_ADDR      = uintptr(ioctl.IOW('a', 2, unsafe.Sizeof(uint16(0))))
	CEC_ADAP_G_LOG_ADDRS      = uintptr(ioctl.IOR('a', 3, unsafe.Sizeof(CecLogAddrs{})))
	CEC_ADAP_S_LOG_ADDRS      = uintptr(ioctl.IOWR('a', 4, unsafe.Sizeof(CecLogAddrs{})))
	CEC_TRANSMIT              = uintptr(ioctl.IOWR('a', 5, unsafe.Sizeof(CecMsg{})))
	CEC_RECEIVE               = uintptr(ioctl.IOWR('a', 6, unsafe.Sizeof(CecMsg{})))
	CEC_DQEVENT               = uintptr(ioctl.IOWR('a', 7, unsafe.Sizeof(CecEvent{})))
	CEC_G_MODE                = uintptr(ioctl.IOR('a', 8, unsafe.Sizeof(uint32(0))))
	CEC_S_MODE                = uintptr(ioctl.IOW('a', 9, unsafe.Sizeof(uint32(0))))
	CEC_ADAP_G_CONNECTOR_INFO = uintptr(ioctl.IOR('a', 10, unsafe.Sizeof(CecConnectorInfo{})))
)
