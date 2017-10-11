// Copyright 2017 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of Codeplug.
//
// Codeplug is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU Lesser General Public
// License as published by the Free Software Foundation.
//
// Codeplug is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Codeplug.  If not, see <http://www.gnu.org/licenses/>.

// Package codeplug implements access to MD380-style codeplug files.
// It can read/update/write both .rdt files and .bin files.
package codeplug

//go:generate genCodeplugInfo codeplugs.json

// Codeplug types
const (
	CtMd380 CodeplugType = "md380"
)

// Record types
const (
	RtChannelInformation RecordType = "ChannelInformation"
	RtDigitalContacts    RecordType = "DigitalContacts"
	RtGeneralSettings    RecordType = "GeneralSettings"
	RtGroupList          RecordType = "GroupList"
	RtRdtHeader          RecordType = "RdtHeader"
	RtScanList           RecordType = "ScanList"
	RtTextMessage        RecordType = "TextMessage"
	RtZoneInformation    RecordType = "ZoneInformation"
)

// Field types
const (
	FtAdmitCriteria           FieldType = "AdmitCriteria"
	FtAllowTalkaround         FieldType = "AllowTalkaround"
	FtAutoscan                FieldType = "Autoscan"
	FtBandwidth               FieldType = "Bandwidth"
	FtCallAlertToneDuration   FieldType = "CallAlertToneDuration"
	FtCallID                  FieldType = "CallID"
	FtCallReceiveTone         FieldType = "CallReceiveTone"
	FtCallType                FieldType = "CallType"
	FtChFreeIndicationTone    FieldType = "ChFreeIndicationTone"
	FtChannelMember           FieldType = "ChannelMember"
	FtChannelMode             FieldType = "ChannelMode"
	FtChannelName             FieldType = "ChannelName"
	FtColorCode               FieldType = "ColorCode"
	FtCompressedUdpDataHeader FieldType = "CompressedUdpDataHeader"
	FtContactMember           FieldType = "ContactMember"
	FtContactName             FieldType = "ContactName"
	FtCtcssDecode             FieldType = "CtcssDecode"
	FtCtcssEncode             FieldType = "CtcssEncode"
	FtDataCallConfirmed       FieldType = "DataCallConfirmed"
	FtDecode1                 FieldType = "Decode1"
	FtDecode2                 FieldType = "Decode2"
	FtDecode3                 FieldType = "Decode3"
	FtDecode4                 FieldType = "Decode4"
	FtDecode5                 FieldType = "Decode5"
	FtDecode6                 FieldType = "Decode6"
	FtDecode7                 FieldType = "Decode7"
	FtDecode8                 FieldType = "Decode8"
	FtDisableAllLeds          FieldType = "DisableAllLeds"
	FtDisableAllTones         FieldType = "DisableAllTones"
	FtDisplayPTTID            FieldType = "DisplayPTTID"
	FtEmergencyAlarmAck       FieldType = "EmergencyAlarmAck"
	FtGroupCallHangTime       FieldType = "GroupCallHangTime"
	FtGroupList               FieldType = "GroupList"
	FtHighFrequency           FieldType = "HighFrequency"
	FtIntroScreen             FieldType = "IntroScreen"
	FtIntroScreenLine1        FieldType = "IntroScreenLine1"
	FtIntroScreenLine2        FieldType = "IntroScreenLine2"
	FtLoneWorker              FieldType = "LoneWorker"
	FtLoneWorkerReminderTime  FieldType = "LoneWorkerReminderTime"
	FtLoneWorkerResponseTime  FieldType = "LoneWorkerResponseTime"
	FtLowFrequency            FieldType = "LowFrequency"
	FtMode                    FieldType = "Mode"
	FtMonitorType             FieldType = "MonitorType"
	FtName                    FieldType = "Name"
	FtPcProgPw                FieldType = "PcProgPw"
	FtPower                   FieldType = "Power"
	FtPowerOnPassword         FieldType = "PowerOnPassword"
	FtPriorityChannel1        FieldType = "PriorityChannel1"
	FtPriorityChannel2        FieldType = "PriorityChannel2"
	FtPrioritySampleTime      FieldType = "PrioritySampleTime"
	FtPrivacy                 FieldType = "Privacy"
	FtPrivacyNumber           FieldType = "PrivacyNumber"
	FtPrivateCallConfirmed    FieldType = "PrivateCallConfirmed"
	FtPrivateCallHangTime     FieldType = "PrivateCallHangTime"
	FtPwAndLockEnable         FieldType = "PwAndLockEnable"
	FtQtReverse               FieldType = "QtReverse"
	FtRadioID                 FieldType = "RadioID"
	FtRadioName               FieldType = "RadioName"
	FtRadioProgPw             FieldType = "RadioProgPw"
	FtRepeaterSlot            FieldType = "RepeaterSlot"
	FtReverseBurst            FieldType = "ReverseBurst"
	FtRxFrequency             FieldType = "RxFrequency"
	FtRxLowBatteryInterval    FieldType = "RxLowBatteryInterval"
	FtRxOnly                  FieldType = "RxOnly"
	FtRxRefFrequency          FieldType = "RxRefFrequency"
	FtRxSignallingSystem      FieldType = "RxSignallingSystem"
	FtSaveModeReceive         FieldType = "SaveModeReceive"
	FtSavePreamble            FieldType = "SavePreamble"
	FtScanAnalogHangTime      FieldType = "ScanAnalogHangTime"
	FtScanDigitalHangTime     FieldType = "ScanDigitalHangTime"
	FtScanList                FieldType = "ScanList"
	FtSetKeypadLockTime       FieldType = "SetKeypadLockTime"
	FtSignallingHoldTime      FieldType = "SignallingHoldTime"
	FtSquelch                 FieldType = "Squelch"
	FtTalkPermitTone          FieldType = "TalkPermitTone"
	FtTextMessage             FieldType = "TextMessage"
	FtTot                     FieldType = "Tot"
	FtTotRekeyDelay           FieldType = "TotRekeyDelay"
	FtTxDesignatedChannel     FieldType = "TxDesignatedChannel"
	FtTxFrequency             FieldType = "TxFrequency"
	FtTxPreambleDuration      FieldType = "TxPreambleDuration"
	FtTxRefFrequency          FieldType = "TxRefFrequency"
	FtTxSignallingSystem      FieldType = "TxSignallingSystem"
	FtVox                     FieldType = "Vox"
	FtVoxSensitivity          FieldType = "VoxSensitivity"
)

// The value types a field may contain
const (
	VtCallID          ValueType = "callID"
	VtCtcssDcs        ValueType = "ctcssDcs"
	VtFrequency       ValueType = "frequency"
	VtIStrings        ValueType = "iStrings"
	VtIndexedStrings  ValueType = "indexedStrings"
	VtIntroLine       ValueType = "introLine"
	VtListIndex       ValueType = "listIndex"
	VtMemberListIndex ValueType = "memberListIndex"
	VtName            ValueType = "name"
	VtOffOn           ValueType = "offOn"
	VtOnOff           ValueType = "onOff"
	VtPcPassword      ValueType = "pcPassword"
	VtPrivacyNumber   ValueType = "privacyNumber"
	VtRadioName       ValueType = "radioName"
	VtRadioPassword   ValueType = "radioPassword"
	VtRhFrequency     ValueType = "rhFrequency"
	VtSpan            ValueType = "span"
	VtTextMessage     ValueType = "textMessage"
)

// newValue returns a new value of the given ValueType
func newValue(vt ValueType) value {
	switch vt {
	case VtCallID:
		return new(callID)
	case VtCtcssDcs:
		return new(ctcssDcs)
	case VtFrequency:
		return new(frequency)
	case VtIStrings:
		return new(iStrings)
	case VtIndexedStrings:
		return new(indexedStrings)
	case VtIntroLine:
		return new(introLine)
	case VtListIndex:
		return new(listIndex)
	case VtMemberListIndex:
		return new(memberListIndex)
	case VtName:
		return new(name)
	case VtOffOn:
		return new(offOn)
	case VtOnOff:
		return new(onOff)
	case VtPcPassword:
		return new(pcPassword)
	case VtPrivacyNumber:
		return new(privacyNumber)
	case VtRadioName:
		return new(radioName)
	case VtRadioPassword:
		return new(radioPassword)
	case VtRhFrequency:
		return new(rhFrequency)
	case VtSpan:
		return new(span)
	case VtTextMessage:
		return new(textMessage)
	}

	return nil
}

// Codeplug types and their records, fields, with offsets, sizes, etc.
var cpTypes = map[CodeplugType][]rInfo{
	CtMd380: []rInfo{
		rInfo{
			rType:    RtRdtHeader,
			typeName: "Rdt Header",
			max:      1,
			offset:   0,
			size:     549,
			fInfos: []fInfo{
				fInfo{
					fType:     FtLowFrequency,
					typeName:  "Low Frequency",
					max:       1,
					bitOffset: 2504,
					bitSize:   16,
					valueType: VtRhFrequency,
				},
				fInfo{
					fType:     FtHighFrequency,
					typeName:  "High Frequency",
					max:       1,
					bitOffset: 2520,
					bitSize:   16,
					valueType: VtRhFrequency,
				},
			},
		},
		rInfo{
			rType:    RtGeneralSettings,
			typeName: "General Settings",
			max:      1,
			offset:   8805,
			size:     144,
			fInfos: []fInfo{
				fInfo{
					fType:     FtIntroScreenLine1,
					typeName:  "Intro Screen Line 1",
					max:       1,
					bitOffset: 0,
					bitSize:   160,
					valueType: VtIntroLine,
				},
				fInfo{
					fType:     FtIntroScreenLine2,
					typeName:  "Intro Screen Line 2",
					max:       1,
					bitOffset: 160,
					bitSize:   160,
					valueType: VtIntroLine,
				},
				fInfo{
					fType:     FtMonitorType,
					typeName:  "Monitor Type",
					max:       1,
					bitOffset: 515,
					bitSize:   1,
					valueType: VtIStrings,
					strings: &[]string{
						"Silent",
						"Open Squelch",
					},
				},
				fInfo{
					fType:     FtDisableAllLeds,
					typeName:  "Disable All LEDS",
					max:       1,
					bitOffset: 517,
					bitSize:   1,
					valueType: VtOnOff,
				},
				fInfo{
					fType:     FtTalkPermitTone,
					typeName:  "Talk Permit Tone",
					max:       1,
					bitOffset: 520,
					bitSize:   2,
					valueType: VtIStrings,
					strings: &[]string{
						"None",
						"Digital",
						"Analog",
						"Digital and Analog",
					},
				},
				fInfo{
					fType:         FtPwAndLockEnable,
					typeName:      "Password And Lock Enable",
					max:           1,
					bitOffset:     522,
					bitSize:       1,
					valueType:     VtOnOff,
					enablingValue: "On",
				},
				fInfo{
					fType:     FtChFreeIndicationTone,
					typeName:  "Channel Free Indication Tone",
					max:       1,
					bitOffset: 523,
					bitSize:   1,
					valueType: VtOnOff,
				},
				fInfo{
					fType:     FtDisableAllTones,
					typeName:  "Disable All Tones",
					max:       1,
					bitOffset: 525,
					bitSize:   1,
					valueType: VtOnOff,
				},
				fInfo{
					fType:     FtSaveModeReceive,
					typeName:  "Save Mode Receive",
					max:       1,
					bitOffset: 526,
					bitSize:   1,
					valueType: VtOffOn,
				},
				fInfo{
					fType:     FtSavePreamble,
					typeName:  "Save Preamble",
					max:       1,
					bitOffset: 527,
					bitSize:   1,
					valueType: VtOffOn,
				},
				fInfo{
					fType:     FtIntroScreen,
					typeName:  "Intro Screen",
					max:       1,
					bitOffset: 531,
					bitSize:   1,
					valueType: VtIStrings,
					strings: &[]string{
						"Character string",
						"Picture",
					},
				},
				fInfo{
					fType:     FtRadioID,
					typeName:  "Radio ID",
					max:       1,
					bitOffset: 544,
					bitSize:   24,
					valueType: VtCallID,
				},
				fInfo{
					fType:     FtTxPreambleDuration,
					typeName:  "Tx Preamble Duration (mS)",
					max:       1,
					bitOffset: 576,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min:   0,
						max:   144,
						scale: 60,
					},
				},
				fInfo{
					fType:     FtGroupCallHangTime,
					typeName:  "Group Call Hang Time (mS)",
					max:       1,
					bitOffset: 584,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min:      0,
						max:      70,
						scale:    100,
						interval: 5,
					},
				},
				fInfo{
					fType:     FtPrivateCallHangTime,
					typeName:  "Private Call Hang Time (mS)",
					max:       1,
					bitOffset: 592,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min:      0,
						max:      70,
						scale:    100,
						interval: 5,
					},
				},
				fInfo{
					fType:     FtVoxSensitivity,
					typeName:  "VOX Sensitivity",
					max:       1,
					bitOffset: 600,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min: 1,
						max: 10,
					},
				},
				fInfo{
					fType:     FtRxLowBatteryInterval,
					typeName:  "Rx Low Battery Interval (S)",
					max:       1,
					bitOffset: 624,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min:   0,
						max:   127,
						scale: 5,
					},
				},
				fInfo{
					fType:     FtCallAlertToneDuration,
					typeName:  "Call Alert Tone Duration (S)",
					max:       1,
					bitOffset: 632,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min:       0,
						max:       240,
						scale:     5,
						interval:  1,
						minString: "Continue",
					},
				},
				fInfo{
					fType:     FtLoneWorkerResponseTime,
					typeName:  "Lone Worker Response Time (min)",
					max:       1,
					bitOffset: 640,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min: 1,
						max: 255,
					},
				},
				fInfo{
					fType:     FtLoneWorkerReminderTime,
					typeName:  "Lone Worker Reminder Time (S)",
					max:       1,
					bitOffset: 648,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min: 1,
						max: 255,
					},
				},
				fInfo{
					fType:     FtScanDigitalHangTime,
					typeName:  "Scan Digital Hang Time (mS)",
					max:       1,
					bitOffset: 664,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min:      5,
						max:      100,
						scale:    100,
						interval: 5,
					},
				},
				fInfo{
					fType:     FtScanAnalogHangTime,
					typeName:  "Scan Analog Hang Time (mS)",
					max:       1,
					bitOffset: 672,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min:      5,
						max:      100,
						scale:    100,
						interval: 5,
					},
				},
				fInfo{
					fType:     FtSetKeypadLockTime,
					typeName:  "Set Keypad Lock Time (S)",
					max:       1,
					bitOffset: 688,
					bitSize:   8,
					valueType: VtIndexedStrings,
					indexedStrings: &[]IndexedString{
						IndexedString{255, "Manual"},
						IndexedString{5, "5"},
						IndexedString{10, "10"},
						IndexedString{15, "15"},
					},
				},
				fInfo{
					fType:     FtMode,
					typeName:  "Mode",
					max:       1,
					bitOffset: 696,
					bitSize:   8,
					valueType: VtIndexedStrings,
					indexedStrings: &[]IndexedString{
						IndexedString{0, "Memory"},
						IndexedString{255, "Channel"},
					},
				},
				fInfo{
					fType:     FtPowerOnPassword,
					typeName:  "Power On Password",
					max:       1,
					bitOffset: 704,
					bitSize:   32,
					valueType: VtRadioPassword,
					enabler:   FtPwAndLockEnable,
				},
				fInfo{
					fType:     FtRadioProgPw,
					typeName:  "Radio Programming Password",
					max:       1,
					bitOffset: 736,
					bitSize:   32,
					valueType: VtRadioPassword,
				},
				fInfo{
					fType:     FtPcProgPw,
					typeName:  "PC Programming Password",
					max:       1,
					bitOffset: 768,
					bitSize:   64,
					valueType: VtPcPassword,
				},
				fInfo{
					fType:     FtRadioName,
					typeName:  "Radio Name",
					max:       1,
					bitOffset: 896,
					bitSize:   256,
					valueType: VtRadioName,
				},
			},
		},
		rInfo{
			rType:    RtTextMessage,
			typeName: "Text Message",
			max:      50,
			offset:   9125,
			size:     288,
			delDescs: []delDesc{
				delDesc{
					offset: 0,
					size:   8,
					value:  0,
				},
			},
			fInfos: []fInfo{
				fInfo{
					fType:     FtTextMessage,
					typeName:  "Message",
					max:       1,
					bitOffset: 0,
					bitSize:   2304,
					valueType: VtTextMessage,
				},
			},
		},
		rInfo{
			rType:    RtDigitalContacts,
			typeName: "Digital Contacts",
			max:      1000,
			offset:   24997,
			size:     36,
			delDescs: []delDesc{
				delDesc{
					offset: 0,
					size:   3,
					value:  255,
				},
				delDesc{
					offset: 4,
					size:   2,
					value:  0,
				},
				delDesc{
					offset: 4,
					size:   16,
					value:  0,
				},
			},
			fInfos: []fInfo{
				fInfo{
					fType:     FtCallID,
					typeName:  "Call ID",
					max:       1,
					bitOffset: 0,
					bitSize:   24,
					valueType: VtCallID,
				},
				fInfo{
					fType:     FtCallReceiveTone,
					typeName:  "Call Receive Tone",
					max:       1,
					bitOffset: 26,
					bitSize:   1,
					valueType: VtIStrings,
					strings: &[]string{
						"No",
						"Yes",
					},
				},
				fInfo{
					fType:     FtCallType,
					typeName:  "Call Type",
					max:       1,
					bitOffset: 30,
					bitSize:   2,
					valueType: VtIStrings,
					strings: &[]string{
						"",
						"Group",
						"Private",
						"All",
					},
				},
				fInfo{
					fType:     FtContactName,
					typeName:  "Contact Name",
					max:       1,
					bitOffset: 32,
					bitSize:   256,
					valueType: VtName,
				},
			},
		},
		rInfo{
			rType:    RtGroupList,
			typeName: "Digital Rx Group List",
			max:      250,
			offset:   60997,
			size:     96,
			delDescs: []delDesc{
				delDesc{
					offset: 0,
					size:   1,
					value:  0,
				},
			},
			fInfos: []fInfo{
				fInfo{
					fType:     FtName,
					typeName:  "Group List Name",
					max:       1,
					bitOffset: 0,
					bitSize:   256,
					valueType: VtName,
				},
				fInfo{
					fType:          FtContactMember,
					typeName:       "Contact Member",
					max:            32,
					bitOffset:      256,
					bitSize:        16,
					valueType:      VtListIndex,
					listRecordType: RtDigitalContacts,
				},
			},
		},
		rInfo{
			rType:    RtZoneInformation,
			typeName: "Zone Information",
			max:      250,
			offset:   84997,
			size:     64,
			delDescs: []delDesc{
				delDesc{
					offset: 0,
					size:   1,
					value:  0,
				},
			},
			fInfos: []fInfo{
				fInfo{
					fType:     FtName,
					typeName:  "Zone Name",
					max:       1,
					bitOffset: 0,
					bitSize:   256,
					valueType: VtName,
				},
				fInfo{
					fType:          FtChannelMember,
					typeName:       "Channel Member",
					max:            16,
					bitOffset:      256,
					bitSize:        16,
					valueType:      VtListIndex,
					listRecordType: RtChannelInformation,
				},
			},
		},
		rInfo{
			rType:    RtScanList,
			typeName: "Scan List",
			max:      250,
			offset:   100997,
			size:     104,
			delDescs: []delDesc{
				delDesc{
					offset: 0,
					size:   1,
					value:  0,
				},
			},
			fInfos: []fInfo{
				fInfo{
					fType:     FtName,
					typeName:  "Scan List Name",
					max:       1,
					bitOffset: 0,
					bitSize:   256,
					valueType: VtName,
				},
				fInfo{
					fType:     FtPriorityChannel1,
					typeName:  "Priority Channel 1",
					max:       1,
					bitOffset: 256,
					bitSize:   16,
					valueType: VtMemberListIndex,
					indexedStrings: &[]IndexedString{
						IndexedString{0, "Selected"},
						IndexedString{65535, "None"},
					},
					listRecordType: RtChannelInformation,
					enablingValue:  "None",
				},
				fInfo{
					fType:     FtPriorityChannel2,
					typeName:  "Priority Channel 2",
					max:       1,
					bitOffset: 272,
					bitSize:   16,
					valueType: VtMemberListIndex,
					indexedStrings: &[]IndexedString{
						IndexedString{0, "Selected"},
						IndexedString{65535, "None"},
					},
					listRecordType: RtChannelInformation,
					disabler:       FtPriorityChannel1,
				},
				fInfo{
					fType:     FtTxDesignatedChannel,
					typeName:  "Tx Designated Channel",
					max:       1,
					bitOffset: 288,
					bitSize:   16,
					valueType: VtListIndex,
					indexedStrings: &[]IndexedString{
						IndexedString{0, "Selected"},
						IndexedString{65535, "Last Active Channel"},
					},
					listRecordType: RtChannelInformation,
				},
				fInfo{
					fType:     FtSignallingHoldTime,
					typeName:  "Signalling Hold Time (mS)",
					max:       1,
					bitOffset: 312,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min:   2,
						max:   255,
						scale: 25,
					},
				},
				fInfo{
					fType:     FtPrioritySampleTime,
					typeName:  "Priority Sample Time (mS)",
					max:       1,
					bitOffset: 320,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min:   3,
						max:   31,
						scale: 250,
					},
				},
				fInfo{
					fType:          FtChannelMember,
					typeName:       "Channel Member",
					max:            31,
					bitOffset:      336,
					bitSize:        16,
					valueType:      VtListIndex,
					listRecordType: RtChannelInformation,
				},
			},
		},
		rInfo{
			rType:    RtChannelInformation,
			typeName: "Channel Information",
			max:      1000,
			offset:   127013,
			size:     64,
			delDescs: []delDesc{
				delDesc{
					offset: 16,
					size:   1,
					value:  255,
				},
			},
			fInfos: []fInfo{
				fInfo{
					fType:     FtLoneWorker,
					typeName:  "Lone Worker",
					max:       1,
					bitOffset: 0,
					bitSize:   1,
					valueType: VtOffOn,
				},
				fInfo{
					fType:     FtSquelch,
					typeName:  "Squelch",
					max:       1,
					bitOffset: 2,
					bitSize:   1,
					valueType: VtIStrings,
					strings: &[]string{
						"Tight",
						"Normal",
					},
				},
				fInfo{
					fType:     FtAutoscan,
					typeName:  "Autoscan",
					max:       1,
					bitOffset: 3,
					bitSize:   1,
					valueType: VtOffOn,
				},
				fInfo{
					fType:     FtBandwidth,
					typeName:  "Bandwidth",
					max:       1,
					bitOffset: 4,
					bitSize:   1,
					valueType: VtIStrings,
					strings: &[]string{
						"12.5",
						"25",
					},
				},
				fInfo{
					fType:     FtChannelMode,
					typeName:  "Channel Mode",
					max:       1,
					bitOffset: 6,
					bitSize:   2,
					valueType: VtIStrings,
					strings: &[]string{
						"",
						"Analog",
						"Digital",
					},
					enablingValue: "Digital",
				},
				fInfo{
					fType:     FtColorCode,
					typeName:  "Color Code",
					max:       1,
					bitOffset: 8,
					bitSize:   4,
					valueType: VtSpan,
					span: &Span{
						min: 0,
						max: 15,
					},
					enabler: FtChannelMode,
				},
				fInfo{
					fType:     FtRepeaterSlot,
					typeName:  "Repeater Slot",
					max:       1,
					bitOffset: 12,
					bitSize:   2,
					valueType: VtIStrings,
					strings: &[]string{
						"",
						"1",
						"2",
					},
					enabler: FtChannelMode,
				},
				fInfo{
					fType:     FtRxOnly,
					typeName:  "Rx Only",
					max:       1,
					bitOffset: 14,
					bitSize:   1,
					valueType: VtOffOn,
				},
				fInfo{
					fType:     FtAllowTalkaround,
					typeName:  "Allow Talkaround",
					max:       1,
					bitOffset: 15,
					bitSize:   1,
					valueType: VtOffOn,
				},
				fInfo{
					fType:     FtDataCallConfirmed,
					typeName:  "Data Call Confirmed",
					max:       1,
					bitOffset: 16,
					bitSize:   1,
					valueType: VtOffOn,
					enabler:   FtChannelMode,
				},
				fInfo{
					fType:     FtPrivateCallConfirmed,
					typeName:  "Private Call Confimed",
					max:       1,
					bitOffset: 17,
					bitSize:   1,
					valueType: VtOffOn,
					enabler:   FtChannelMode,
				},
				fInfo{
					fType:     FtPrivacy,
					typeName:  "Privacy",
					max:       1,
					bitOffset: 18,
					bitSize:   2,
					valueType: VtIStrings,
					strings: &[]string{
						"None",
						"Basic",
						"Enhanced",
					},
					enablingValue: "None",
					enabler:       FtChannelMode,
				},
				fInfo{
					fType:     FtPrivacyNumber,
					typeName:  "Privacy Number",
					max:       1,
					bitOffset: 20,
					bitSize:   4,
					valueType: VtPrivacyNumber,
					span: &Span{
						min: 0,
						max: 15,
					},
					disabler: FtPrivacy,
				},
				fInfo{
					fType:     FtDisplayPTTID,
					typeName:  "Display PTT ID",
					max:       1,
					bitOffset: 24,
					bitSize:   1,
					valueType: VtOnOff,
					disabler:  FtChannelMode,
				},
				fInfo{
					fType:     FtCompressedUdpDataHeader,
					typeName:  "Compressed UDP Data Header",
					max:       1,
					bitOffset: 25,
					bitSize:   1,
					valueType: VtOffOn,
					enabler:   FtChannelMode,
				},
				fInfo{
					fType:     FtEmergencyAlarmAck,
					typeName:  "Emergency Alarm Ack",
					max:       1,
					bitOffset: 28,
					bitSize:   1,
					valueType: VtOffOn,
					enabler:   FtChannelMode,
				},
				fInfo{
					fType:     FtRxRefFrequency,
					typeName:  "Rx Ref Frequency",
					max:       1,
					bitOffset: 30,
					bitSize:   2,
					valueType: VtIStrings,
					strings: &[]string{
						"Low",
						"Medium",
						"High",
					},
				},
				fInfo{
					fType:     FtAdmitCriteria,
					typeName:  "Admit Criteria",
					max:       1,
					bitOffset: 32,
					bitSize:   2,
					valueType: VtIStrings,
					strings: &[]string{
						"Always",
						"Channel free",
						"CTCSS/DCS",
						"Color code",
					},
				},
				fInfo{
					fType:     FtPower,
					typeName:  "Power",
					max:       1,
					bitOffset: 34,
					bitSize:   1,
					valueType: VtIStrings,
					strings: &[]string{
						"Low",
						"High",
					},
				},
				fInfo{
					fType:     FtVox,
					typeName:  "VOX",
					max:       1,
					bitOffset: 35,
					bitSize:   1,
					valueType: VtOffOn,
				},
				fInfo{
					fType:     FtQtReverse,
					typeName:  "QT Reverse",
					max:       1,
					bitOffset: 36,
					bitSize:   1,
					valueType: VtIStrings,
					strings: &[]string{
						"180",
						"120",
					},
					disabler: FtCtcssEncode,
				},
				fInfo{
					fType:     FtReverseBurst,
					typeName:  "Reverse Burst/Turn Off Code",
					max:       1,
					bitOffset: 37,
					bitSize:   1,
					valueType: VtOffOn,
					disabler:  FtCtcssEncode,
				},
				fInfo{
					fType:     FtTxRefFrequency,
					typeName:  "Tx Ref Frequency",
					max:       1,
					bitOffset: 38,
					bitSize:   2,
					valueType: VtIStrings,
					strings: &[]string{
						"Low",
						"Medium",
						"High",
					},
				},
				fInfo{
					fType:     FtContactName,
					typeName:  "Contact Name",
					max:       1,
					bitOffset: 48,
					bitSize:   16,
					valueType: VtListIndex,
					indexedStrings: &[]IndexedString{
						IndexedString{0, "None"},
					},
					listRecordType: RtDigitalContacts,
					enabler:        FtChannelMode,
				},
				fInfo{
					fType:     FtTot,
					typeName:  "TOT (S)",
					max:       1,
					bitOffset: 66,
					bitSize:   6,
					valueType: VtSpan,
					span: &Span{
						min:       0,
						max:       63,
						scale:     15,
						interval:  1,
						minString: "Infinite",
					},
				},
				fInfo{
					fType:     FtTotRekeyDelay,
					typeName:  "TOT Rekey Delay (S)",
					max:       1,
					bitOffset: 72,
					bitSize:   8,
					valueType: VtSpan,
					span: &Span{
						min: 0,
						max: 255,
					},
				},
				fInfo{
					fType:     FtScanList,
					typeName:  "Scan List",
					max:       1,
					bitOffset: 88,
					bitSize:   8,
					valueType: VtListIndex,
					indexedStrings: &[]IndexedString{
						IndexedString{0, "None"},
					},
					listRecordType: RtScanList,
				},
				fInfo{
					fType:     FtGroupList,
					typeName:  "Group List",
					max:       1,
					bitOffset: 96,
					bitSize:   8,
					valueType: VtListIndex,
					indexedStrings: &[]IndexedString{
						IndexedString{0, "None"},
					},
					listRecordType: RtGroupList,
					enabler:        FtChannelMode,
				},
				fInfo{
					fType:     FtDecode1,
					typeName:  "Decode 1",
					max:       1,
					bitOffset: 112,
					bitSize:   1,
					valueType: VtOffOn,
					disabler:  FtRxSignallingSystem,
				},
				fInfo{
					fType:     FtDecode2,
					typeName:  "Decode 2",
					max:       1,
					bitOffset: 113,
					bitSize:   1,
					valueType: VtOffOn,
					disabler:  FtRxSignallingSystem,
				},
				fInfo{
					fType:     FtDecode3,
					typeName:  "Decode 3",
					max:       1,
					bitOffset: 114,
					bitSize:   1,
					valueType: VtOffOn,
					disabler:  FtRxSignallingSystem,
				},
				fInfo{
					fType:     FtDecode4,
					typeName:  "Decode 4",
					max:       1,
					bitOffset: 115,
					bitSize:   1,
					valueType: VtOffOn,
					disabler:  FtRxSignallingSystem,
				},
				fInfo{
					fType:     FtDecode5,
					typeName:  "Decode 5",
					max:       1,
					bitOffset: 116,
					bitSize:   1,
					valueType: VtOffOn,
					disabler:  FtRxSignallingSystem,
				},
				fInfo{
					fType:     FtDecode6,
					typeName:  "Decode 6",
					max:       1,
					bitOffset: 117,
					bitSize:   1,
					valueType: VtOffOn,
					disabler:  FtRxSignallingSystem,
				},
				fInfo{
					fType:     FtDecode7,
					typeName:  "Decode 7",
					max:       1,
					bitOffset: 118,
					bitSize:   1,
					valueType: VtOffOn,
					disabler:  FtRxSignallingSystem,
				},
				fInfo{
					fType:     FtDecode8,
					typeName:  "Decode 8",
					max:       1,
					bitOffset: 119,
					bitSize:   1,
					valueType: VtOffOn,
					disabler:  FtRxSignallingSystem,
				},
				fInfo{
					fType:     FtRxFrequency,
					typeName:  "Rx Frequency (MHz)",
					max:       1,
					bitOffset: 128,
					bitSize:   32,
					valueType: VtFrequency,
				},
				fInfo{
					fType:     FtTxFrequency,
					typeName:  "Tx Frequency (MHz)",
					max:       1,
					bitOffset: 160,
					bitSize:   32,
					valueType: VtFrequency,
				},
				fInfo{
					fType:     FtCtcssDecode,
					typeName:  "CTCSS/DCS Decode",
					max:       1,
					bitOffset: 192,
					bitSize:   16,
					valueType: VtCtcssDcs,
					disabler:  FtChannelMode,
				},
				fInfo{
					fType:         FtCtcssEncode,
					typeName:      "CTCSS/DCS Encode",
					max:           1,
					bitOffset:     208,
					bitSize:       16,
					valueType:     VtCtcssDcs,
					enablingValue: "None",
					disabler:      FtChannelMode,
				},
				fInfo{
					fType:     FtRxSignallingSystem,
					typeName:  "Rx Signaling System",
					max:       1,
					bitOffset: 229,
					bitSize:   3,
					valueType: VtIStrings,
					strings: &[]string{
						"Off",
						"DTMF-1",
						"DTMF-2",
						"DTMF-3",
						"DTMF-4",
					},
					enablingValue: "Off",
					disabler:      FtChannelMode,
				},
				fInfo{
					fType:     FtTxSignallingSystem,
					typeName:  "Tx Signaling System",
					max:       1,
					bitOffset: 237,
					bitSize:   3,
					valueType: VtIStrings,
					strings: &[]string{
						"Off",
						"DTMF-1",
						"DTMF-2",
						"DTMF-3",
						"DTMF-4",
					},
					disabler: FtChannelMode,
				},
				fInfo{
					fType:     FtChannelName,
					typeName:  "Channel Name",
					max:       1,
					bitOffset: 256,
					bitSize:   256,
					valueType: VtName,
				},
			},
		},
	},
}

//go:generate genCodeplugInfo
