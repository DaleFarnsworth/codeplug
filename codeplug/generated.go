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

// Record types
const (
	RtBasicInformation_md380 RecordType = "BasicInformation"
	RtBasicInformation_md390 RecordType = "BasicInformation"
	RtBasicInformation_md40  RecordType = "BasicInformation"
	RtChannels_md380         RecordType = "Channels"
	RtChannels_md40          RecordType = "Channels"
	RtContacts               RecordType = "Contacts"
	RtGPSSystems             RecordType = "GPSSystems"
	RtGeneralSettings_md380  RecordType = "GeneralSettings"
	RtGeneralSettings_md40   RecordType = "GeneralSettings"
	RtGroupLists             RecordType = "GroupLists"
	RtScanLists_md380        RecordType = "ScanLists"
	RtScanLists_md40         RecordType = "ScanLists"
	RtTextMessages           RecordType = "TextMessages"
	RtZones_md380            RecordType = "Zones"
	RtZones_md40             RecordType = "Zones"
)

// Field types
const (
	FtBiCpsVersion                FieldType = "CpsVersion"
	FtBiFrequencyRange_md380      FieldType = "FrequencyRange"
	FtBiFrequencyRange_md390      FieldType = "FrequencyRange"
	FtBiFrequencyRange_md40       FieldType = "FrequencyRange"
	FtBiHighFrequency             FieldType = "HighFrequency"
	FtBiLastProgrammedTime        FieldType = "LastProgrammedTime"
	FtBiLowFrequency              FieldType = "LowFrequency"
	FtBiModel                     FieldType = "Model"
	FtBiNewFilename_md380         FieldType = "NewFilename"
	FtBiNewFilename_md390         FieldType = "NewFilename"
	FtBiNewFilename_md40          FieldType = "NewFilename"
	FtCiAdmitCriteria             FieldType = "AdmitCriteria"
	FtCiAllowTalkaround           FieldType = "AllowTalkaround"
	FtCiAutoscan                  FieldType = "Autoscan"
	FtCiBandwidth                 FieldType = "Bandwidth"
	FtCiChannelMode               FieldType = "ChannelMode"
	FtCiColorCode                 FieldType = "ColorCode"
	FtCiCompressedUdpDataHeader   FieldType = "CompressedUdpDataHeader"
	FtCiContactName               FieldType = "ContactName"
	FtCiCtcssDecode               FieldType = "CtcssDecode"
	FtCiCtcssEncode               FieldType = "CtcssEncode"
	FtCiDataCallConfirmed         FieldType = "DataCallConfirmed"
	FtCiDecode1                   FieldType = "Decode1"
	FtCiDecode2                   FieldType = "Decode2"
	FtCiDecode3                   FieldType = "Decode3"
	FtCiDecode4                   FieldType = "Decode4"
	FtCiDecode5                   FieldType = "Decode5"
	FtCiDecode6                   FieldType = "Decode6"
	FtCiDecode7                   FieldType = "Decode7"
	FtCiDecode8                   FieldType = "Decode8"
	FtCiDisplayPTTID              FieldType = "DisplayPTTID"
	FtCiEmergencyAlarmAck         FieldType = "EmergencyAlarmAck"
	FtCiGPSSystem                 FieldType = "GPSSystem"
	FtCiGroupList                 FieldType = "GroupList"
	FtCiLoneWorker                FieldType = "LoneWorker"
	FtCiName                      FieldType = "Name"
	FtCiPower                     FieldType = "Power"
	FtCiPrivacy                   FieldType = "Privacy"
	FtCiPrivacyNumber             FieldType = "PrivacyNumber"
	FtCiPrivateCallConfirmed      FieldType = "PrivateCallConfirmed"
	FtCiQtReverse                 FieldType = "QtReverse"
	FtCiReceiveGPSInfo            FieldType = "ReceiveGPSInfo"
	FtCiRepeaterSlot              FieldType = "RepeaterSlot"
	FtCiReverseBurst              FieldType = "ReverseBurst"
	FtCiRxFrequency               FieldType = "RxFrequency"
	FtCiRxOnly                    FieldType = "RxOnly"
	FtCiRxRefFrequency            FieldType = "RxRefFrequency"
	FtCiRxSignallingSystem        FieldType = "RxSignallingSystem"
	FtCiScanList_md380            FieldType = "ScanList"
	FtCiScanList_md40             FieldType = "ScanList"
	FtCiSendGPSInfo               FieldType = "SendGPSInfo"
	FtCiSquelch                   FieldType = "Squelch"
	FtCiTot                       FieldType = "Tot"
	FtCiTotRekeyDelay             FieldType = "TotRekeyDelay"
	FtCiTxFrequency               FieldType = "TxFrequency"
	FtCiTxRefFrequency            FieldType = "TxRefFrequency"
	FtCiTxSignallingSystem        FieldType = "TxSignallingSystem"
	FtCiVox                       FieldType = "Vox"
	FtDcCallID                    FieldType = "CallID"
	FtDcCallReceiveTone           FieldType = "CallReceiveTone"
	FtDcCallType                  FieldType = "CallType"
	FtDcName                      FieldType = "Name"
	FtGlContact                   FieldType = "Contact"
	FtGlName                      FieldType = "Name"
	FtGpDestinationID             FieldType = "DestinationID"
	FtGpGPSDefaultReportInterval  FieldType = "GPSDefaultReportInterval"
	FtGpGPSRevertChannel          FieldType = "GPSRevertChannel"
	FtGsBacklightColor            FieldType = "BacklightColor"
	FtGsBacklightTime             FieldType = "BacklightTime"
	FtGsCallAlertToneDuration     FieldType = "CallAlertToneDuration"
	FtGsChFreeIndicationTone      FieldType = "ChFreeIndicationTone"
	FtGsDisableAllLeds            FieldType = "DisableAllLeds"
	FtGsDisableAllTones           FieldType = "DisableAllTones"
	FtGsFreqChannelMode           FieldType = "FreqChannelMode"
	FtGsGroupCallHangTime         FieldType = "GroupCallHangTime"
	FtGsIntroScreen               FieldType = "IntroScreen"
	FtGsIntroScreenLine1          FieldType = "IntroScreenLine1"
	FtGsIntroScreenLine2          FieldType = "IntroScreenLine2"
	FtGsLockUnlock                FieldType = "LockUnlock"
	FtGsLoneWorkerReminderTime    FieldType = "LoneWorkerReminderTime"
	FtGsLoneWorkerResponseTime    FieldType = "LoneWorkerResponseTime"
	FtGsMode                      FieldType = "Mode"
	FtGsModeSelect                FieldType = "ModeSelect"
	FtGsMonitorType               FieldType = "MonitorType"
	FtGsPcProgPassword            FieldType = "PcProgPassword"
	FtGsPowerOnPassword           FieldType = "PowerOnPassword"
	FtGsPrivateCallHangTime       FieldType = "PrivateCallHangTime"
	FtGsPwAndLockEnable           FieldType = "PwAndLockEnable"
	FtGsRadioID                   FieldType = "RadioID"
	FtGsRadioName                 FieldType = "RadioName"
	FtGsRadioProgPassword         FieldType = "RadioProgPassword"
	FtGsRxLowBatteryInterval      FieldType = "RxLowBatteryInterval"
	FtGsSaveModeReceive           FieldType = "SaveModeReceive"
	FtGsSavePreamble              FieldType = "SavePreamble"
	FtGsScanAnalogHangTime        FieldType = "ScanAnalogHangTime"
	FtGsScanDigitalHangTime       FieldType = "ScanDigitalHangTime"
	FtGsSetKeypadLockTime         FieldType = "SetKeypadLockTime"
	FtGsTalkPermitTone            FieldType = "TalkPermitTone"
	FtGsTxPreambleDuration        FieldType = "TxPreambleDuration"
	FtGsVoxSensitivity            FieldType = "VoxSensitivity"
	FtSlChannel_md380             FieldType = "Channel"
	FtSlChannel_md40              FieldType = "Channel"
	FtSlName                      FieldType = "Name"
	FtSlPriorityChannel1_md380    FieldType = "PriorityChannel1"
	FtSlPriorityChannel1_md40     FieldType = "PriorityChannel1"
	FtSlPriorityChannel2_md380    FieldType = "PriorityChannel2"
	FtSlPriorityChannel2_md40     FieldType = "PriorityChannel2"
	FtSlPrioritySampleTime        FieldType = "PrioritySampleTime"
	FtSlSignallingHoldTime        FieldType = "SignallingHoldTime"
	FtSlTxDesignatedChannel_md380 FieldType = "TxDesignatedChannel"
	FtSlTxDesignatedChannel_md40  FieldType = "TxDesignatedChannel"
	FtTmTextMessage               FieldType = "TextMessage"
	FtZiChannel_md380             FieldType = "Channel"
	FtZiChannel_md40              FieldType = "Channel"
	FtZiName                      FieldType = "Name"
)

// The value types a field may contain
const (
	VtBiFilename      ValueType = "biFilename"
	VtBiFrequency     ValueType = "biFrequency"
	VtCallID          ValueType = "callID"
	VtCallType        ValueType = "callType"
	VtCpsVersion      ValueType = "cpsVersion"
	VtCtcssDcs        ValueType = "ctcssDcs"
	VtFrequency       ValueType = "frequency"
	VtGpsListIndex    ValueType = "gpsListIndex"
	VtIStrings        ValueType = "iStrings"
	VtIndexedStrings  ValueType = "indexedStrings"
	VtIntroLine       ValueType = "introLine"
	VtListIndex       ValueType = "listIndex"
	VtMemberListIndex ValueType = "memberListIndex"
	VtModel           ValueType = "model"
	VtName            ValueType = "name"
	VtOffOn           ValueType = "offOn"
	VtOnOff           ValueType = "onOff"
	VtPcPassword      ValueType = "pcPassword"
	VtPrivacyNumber   ValueType = "privacyNumber"
	VtRadioName       ValueType = "radioName"
	VtRadioPassword   ValueType = "radioPassword"
	VtSpan            ValueType = "span"
	VtTextMessage     ValueType = "textMessage"
	VtTimeStamp       ValueType = "timeStamp"
	VtUniqueName      ValueType = "uniqueName"
)

// newValue returns a new value of the given ValueType
func newValue(vt ValueType) value {
	switch vt {
	case VtBiFilename:
		return new(biFilename)
	case VtBiFrequency:
		return new(biFrequency)
	case VtCallID:
		return new(callID)
	case VtCallType:
		return new(callType)
	case VtCpsVersion:
		return new(cpsVersion)
	case VtCtcssDcs:
		return new(ctcssDcs)
	case VtFrequency:
		return new(frequency)
	case VtGpsListIndex:
		return new(gpsListIndex)
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
	case VtModel:
		return new(model)
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
	case VtSpan:
		return new(span)
	case VtTextMessage:
		return new(textMessage)
	case VtTimeStamp:
		return new(timeStamp)
	case VtUniqueName:
		return new(uniqueName)
	}

	return nil
}

var codeplugInfos = []*CodeplugInfo{
	&cpMd380,
	&cpMd390,
	&cpMd40,
}

var cpMd380 = CodeplugInfo{
	Type: "md380",
	Models: []string{
		"MD380",
		"DR780",
	},
	Ext:       "rdt",
	RdtSize:   262709,
	BinSize:   262144,
	BinOffset: 549,
	RecordInfos: []*recordInfo{
		&riBasicInformation_md380,
		&riGeneralSettings_md380,
		&riContacts,
		&riGroupLists,
		&riZones_md380,
		&riScanLists_md380,
		&riChannels_md380,
		&riGPSSystems,
	},
}

var cpMd390 = CodeplugInfo{
	Type: "md390",
	Models: []string{
		"MD390",
	},
	Ext:       "rdt",
	RdtSize:   262709,
	BinSize:   262144,
	BinOffset: 549,
	RecordInfos: []*recordInfo{
		&riBasicInformation_md390,
		&riGeneralSettings_md380,
		&riContacts,
		&riGroupLists,
		&riZones_md380,
		&riScanLists_md380,
		&riChannels_md380,
		&riGPSSystems,
	},
}

var cpMd40 = CodeplugInfo{
	Type: "md40",
	Models: []string{
		"DJ-MD40",
	},
	Ext:       "rdt",
	RdtSize:   262709,
	BinSize:   262144,
	BinOffset: 549,
	RecordInfos: []*recordInfo{
		&riBasicInformation_md40,
		&riGeneralSettings_md40,
		&riContacts,
		&riGroupLists,
		&riZones_md40,
		&riScanLists_md40,
		&riChannels_md40,
	},
}

var riBasicInformation_md380 = recordInfo{
	rType:    RtBasicInformation_md380,
	typeName: "Basic Information",
	max:      1,
	offset:   0,
	size:     8805,
	fieldInfos: []*fieldInfo{
		&fiBiModel,
		&fiBiFrequencyRange_md380,
		&fiBiNewFilename_md380,
		&fiBiLowFrequency,
		&fiBiHighFrequency,
		&fiBiLastProgrammedTime,
		&fiBiCpsVersion,
	},
}

var riBasicInformation_md390 = recordInfo{
	rType:    RtBasicInformation_md390,
	typeName: "Basic Information",
	max:      1,
	offset:   0,
	size:     8805,
	fieldInfos: []*fieldInfo{
		&fiBiModel,
		&fiBiFrequencyRange_md390,
		&fiBiNewFilename_md390,
		&fiBiLowFrequency,
		&fiBiHighFrequency,
		&fiBiLastProgrammedTime,
		&fiBiCpsVersion,
	},
}

var riBasicInformation_md40 = recordInfo{
	rType:    RtBasicInformation_md40,
	typeName: "Basic Information",
	max:      1,
	offset:   0,
	size:     8805,
	fieldInfos: []*fieldInfo{
		&fiBiModel,
		&fiBiFrequencyRange_md40,
		&fiBiNewFilename_md40,
		&fiBiLowFrequency,
		&fiBiHighFrequency,
		&fiBiLastProgrammedTime,
		&fiBiCpsVersion,
	},
}

var riChannels_md380 = recordInfo{
	rType:    RtChannels_md380,
	typeName: "Channels",
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
	fieldInfos: []*fieldInfo{
		&fiCiName,
		&fiCiRxFrequency,
		&fiCiTxFrequency,
		&fiCiChannelMode,
		&fiCiBandwidth,
		&fiCiScanList_md380,
		&fiCiSquelch,
		&fiCiRxRefFrequency,
		&fiCiTxRefFrequency,
		&fiCiTot,
		&fiCiTotRekeyDelay,
		&fiCiPower,
		&fiCiAdmitCriteria,
		&fiCiAutoscan,
		&fiCiRxOnly,
		&fiCiLoneWorker,
		&fiCiVox,
		&fiCiAllowTalkaround,
		&fiCiPrivateCallConfirmed,
		&fiCiEmergencyAlarmAck,
		&fiCiDataCallConfirmed,
		&fiCiCompressedUdpDataHeader,
		&fiCiContactName,
		&fiCiGroupList,
		&fiCiColorCode,
		&fiCiRepeaterSlot,
		&fiCiPrivacy,
		&fiCiPrivacyNumber,
		&fiCiDisplayPTTID,
		&fiCiCtcssEncode,
		&fiCiTxSignallingSystem,
		&fiCiQtReverse,
		&fiCiReverseBurst,
		&fiCiCtcssDecode,
		&fiCiRxSignallingSystem,
		&fiCiDecode1,
		&fiCiDecode2,
		&fiCiDecode3,
		&fiCiDecode4,
		&fiCiDecode5,
		&fiCiDecode6,
		&fiCiDecode7,
		&fiCiDecode8,
		&fiCiReceiveGPSInfo,
		&fiCiSendGPSInfo,
		&fiCiGPSSystem,
	},
}

var riChannels_md40 = recordInfo{
	rType:    RtChannels_md40,
	typeName: "Channels",
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
	fieldInfos: []*fieldInfo{
		&fiCiName,
		&fiCiRxFrequency,
		&fiCiTxFrequency,
		&fiCiChannelMode,
		&fiCiBandwidth,
		&fiCiScanList_md40,
		&fiCiSquelch,
		&fiCiRxRefFrequency,
		&fiCiTxRefFrequency,
		&fiCiTot,
		&fiCiTotRekeyDelay,
		&fiCiPower,
		&fiCiAdmitCriteria,
		&fiCiAutoscan,
		&fiCiRxOnly,
		&fiCiLoneWorker,
		&fiCiVox,
		&fiCiAllowTalkaround,
		&fiCiPrivateCallConfirmed,
		&fiCiEmergencyAlarmAck,
		&fiCiDataCallConfirmed,
		&fiCiCompressedUdpDataHeader,
		&fiCiContactName,
		&fiCiGroupList,
		&fiCiColorCode,
		&fiCiRepeaterSlot,
		&fiCiPrivacy,
		&fiCiPrivacyNumber,
		&fiCiDisplayPTTID,
		&fiCiCtcssEncode,
		&fiCiTxSignallingSystem,
		&fiCiQtReverse,
		&fiCiReverseBurst,
		&fiCiCtcssDecode,
		&fiCiRxSignallingSystem,
		&fiCiDecode1,
		&fiCiDecode2,
		&fiCiDecode3,
		&fiCiDecode4,
		&fiCiDecode5,
		&fiCiDecode6,
		&fiCiDecode7,
		&fiCiDecode8,
	},
}

var riContacts = recordInfo{
	rType:    RtContacts,
	typeName: "Contacts",
	max:      1000,
	offset:   24997,
	size:     36,
	delDescs: []delDesc{
		delDesc{
			offset: 3,
			size:   1,
			value:  192,
		},
	},
	fieldInfos: []*fieldInfo{
		&fiDcName,
		&fiDcCallID,
		&fiDcCallType,
		&fiDcCallReceiveTone,
	},
}

var riGPSSystems = recordInfo{
	rType:      RtGPSSystems,
	typeName:   "GPS Systems",
	max:        16,
	offset:     257637,
	size:       16,
	namePrefix: "GPS ",
	fieldInfos: []*fieldInfo{
		&fiGpGPSRevertChannel,
		&fiGpGPSDefaultReportInterval,
		&fiGpDestinationID,
	},
}

var riGeneralSettings_md380 = recordInfo{
	rType:    RtGeneralSettings_md380,
	typeName: "General Settings",
	max:      1,
	offset:   8805,
	size:     144,
	fieldInfos: []*fieldInfo{
		&fiGsRadioName,
		&fiGsRadioID,
		&fiGsIntroScreen,
		&fiGsIntroScreenLine1,
		&fiGsIntroScreenLine2,
		&fiGsSavePreamble,
		&fiGsSaveModeReceive,
		&fiGsDisableAllTones,
		&fiGsChFreeIndicationTone,
		&fiGsTalkPermitTone,
		&fiGsCallAlertToneDuration,
		&fiGsScanDigitalHangTime,
		&fiGsScanAnalogHangTime,
		&fiGsLoneWorkerResponseTime,
		&fiGsLoneWorkerReminderTime,
		&fiGsPwAndLockEnable,
		&fiGsPowerOnPassword,
		&fiGsMonitorType,
		&fiGsVoxSensitivity,
		&fiGsTxPreambleDuration,
		&fiGsRxLowBatteryInterval,
		&fiGsPcProgPassword,
		&fiGsRadioProgPassword,
		&fiGsMode,
		&fiGsBacklightTime,
		&fiGsSetKeypadLockTime,
		&fiGsDisableAllLeds,
		&fiGsGroupCallHangTime,
		&fiGsPrivateCallHangTime,
	},
}

var riGeneralSettings_md40 = recordInfo{
	rType:    RtGeneralSettings_md40,
	typeName: "General Settings",
	max:      1,
	offset:   8805,
	size:     144,
	fieldInfos: []*fieldInfo{
		&fiGsRadioName,
		&fiGsRadioID,
		&fiGsIntroScreen,
		&fiGsIntroScreenLine1,
		&fiGsIntroScreenLine2,
		&fiGsSavePreamble,
		&fiGsSaveModeReceive,
		&fiGsDisableAllTones,
		&fiGsChFreeIndicationTone,
		&fiGsTalkPermitTone,
		&fiGsCallAlertToneDuration,
		&fiGsScanDigitalHangTime,
		&fiGsScanAnalogHangTime,
		&fiGsLoneWorkerResponseTime,
		&fiGsLoneWorkerReminderTime,
		&fiGsPwAndLockEnable,
		&fiGsPowerOnPassword,
		&fiGsMonitorType,
		&fiGsVoxSensitivity,
		&fiGsTxPreambleDuration,
		&fiGsRxLowBatteryInterval,
		&fiGsPcProgPassword,
		&fiGsRadioProgPassword,
		&fiGsMode,
		&fiGsFreqChannelMode,
		&fiGsBacklightColor,
		&fiGsModeSelect,
		&fiGsLockUnlock,
		&fiGsSetKeypadLockTime,
		&fiGsDisableAllLeds,
		&fiGsGroupCallHangTime,
		&fiGsPrivateCallHangTime,
	},
}

var riGroupLists = recordInfo{
	rType:    RtGroupLists,
	typeName: "Group Lists",
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
	fieldInfos: []*fieldInfo{
		&fiGlName,
		&fiGlContact,
	},
}

var riScanLists_md380 = recordInfo{
	rType:    RtScanLists_md380,
	typeName: "Scan Lists",
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
	fieldInfos: []*fieldInfo{
		&fiSlName,
		&fiSlPriorityChannel1_md380,
		&fiSlPriorityChannel2_md380,
		&fiSlTxDesignatedChannel_md380,
		&fiSlSignallingHoldTime,
		&fiSlPrioritySampleTime,
		&fiSlChannel_md380,
	},
}

var riScanLists_md40 = recordInfo{
	rType:    RtScanLists_md40,
	typeName: "Scan Lists",
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
	fieldInfos: []*fieldInfo{
		&fiSlName,
		&fiSlPriorityChannel1_md40,
		&fiSlPriorityChannel2_md40,
		&fiSlTxDesignatedChannel_md40,
		&fiSlSignallingHoldTime,
		&fiSlPrioritySampleTime,
		&fiSlChannel_md40,
	},
}

var riTextMessages = recordInfo{
	rType:    RtTextMessages,
	typeName: "Text Messages",
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
	fieldInfos: []*fieldInfo{
		&fiTmTextMessage,
	},
}

var riZones_md380 = recordInfo{
	rType:    RtZones_md380,
	typeName: "Zones",
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
	fieldInfos: []*fieldInfo{
		&fiZiName,
		&fiZiChannel_md380,
	},
}

var riZones_md40 = recordInfo{
	rType:    RtZones_md40,
	typeName: "Zones",
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
	fieldInfos: []*fieldInfo{
		&fiZiName,
		&fiZiChannel_md40,
	},
}

var fiBiCpsVersion = fieldInfo{
	fType:     FtBiCpsVersion,
	typeName:  "CPS Version",
	max:       1,
	bitOffset: 69992,
	bitSize:   32,
	valueType: VtCpsVersion,
}

var fiBiFrequencyRange_md380 = fieldInfo{
	fType:     FtBiFrequencyRange_md380,
	typeName:  "Frequency Range",
	max:       1,
	bitOffset: 2480,
	bitSize:   8,
	valueType: VtIStrings,
	strings: &[]string{
		"136-174 MHz",
		"350-400 MHz",
		"400-480 MHz",
		"450-520 MHz",
	},
}

var fiBiFrequencyRange_md390 = fieldInfo{
	fType:     FtBiFrequencyRange_md390,
	typeName:  "Frequency Range",
	max:       1,
	bitOffset: 2480,
	bitSize:   8,
	valueType: VtIStrings,
	strings: &[]string{
		"136-174 MHz",
		"350-400 MHz",
		"400-480 MHz",
		"450-520 MHz",
	},
}

var fiBiFrequencyRange_md40 = fieldInfo{
	fType:     FtBiFrequencyRange_md40,
	typeName:  "Frequency Range",
	max:       1,
	bitOffset: 2480,
	bitSize:   8,
	valueType: VtIStrings,
	strings: &[]string{
		"MD40 400-480 MHz",
		"MD40HT 420-450 MHz",
		"MD40HE 430-440 MHz",
		"MD40T 400-480 MHz",
	},
}

var fiBiHighFrequency = fieldInfo{
	fType:     FtBiHighFrequency,
	typeName:  "High Frequency",
	max:       1,
	bitOffset: 2520,
	bitSize:   16,
	valueType: VtBiFrequency,
}

var fiBiLastProgrammedTime = fieldInfo{
	fType:     FtBiLastProgrammedTime,
	typeName:  "Last Programmed Time",
	max:       1,
	bitOffset: 69936,
	bitSize:   56,
	valueType: VtTimeStamp,
}

var fiBiLowFrequency = fieldInfo{
	fType:     FtBiLowFrequency,
	typeName:  "Low Frequency",
	max:       1,
	bitOffset: 2504,
	bitSize:   16,
	valueType: VtBiFrequency,
}

var fiBiModel = fieldInfo{
	fType:     FtBiModel,
	typeName:  "Model Name",
	max:       1,
	bitOffset: 2344,
	bitSize:   64,
	valueType: VtModel,
}

var fiBiNewFilename_md380 = fieldInfo{
	fType:     FtBiNewFilename_md380,
	typeName:  "Codeplug Model Filename",
	max:       1,
	bitOffset: 2480,
	bitSize:   8,
	valueType: VtBiFilename,
	strings: &[]string{
		"md380_136-174.rdt",
		"md380_350-400.rdt",
		"md380_400-480.rdt",
		"md380_450-520.rdt",
	},
}

var fiBiNewFilename_md390 = fieldInfo{
	fType:     FtBiNewFilename_md390,
	typeName:  "Codeplug Model Filename",
	max:       1,
	bitOffset: 2480,
	bitSize:   8,
	valueType: VtBiFilename,
	strings: &[]string{
		"md390_136-174.rdt",
		"md390_350-400.rdt",
		"md390_400-480.rdt",
		"md390_450-520.rdt",
	},
}

var fiBiNewFilename_md40 = fieldInfo{
	fType:     FtBiNewFilename_md40,
	typeName:  "Codeplug Model Filename",
	max:       1,
	bitOffset: 2480,
	bitSize:   8,
	valueType: VtBiFilename,
	strings: &[]string{
		"md40_400-480.rdt",
		"md40ht_420-450.rdt",
		"md40he_430-440.rdt",
		"md40t_400-480.rdt",
	},
}

var fiCiAdmitCriteria = fieldInfo{
	fType:     FtCiAdmitCriteria,
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
}

var fiCiAllowTalkaround = fieldInfo{
	fType:     FtCiAllowTalkaround,
	typeName:  "Allow Talkaround",
	max:       1,
	bitOffset: 15,
	bitSize:   1,
	valueType: VtOffOn,
}

var fiCiAutoscan = fieldInfo{
	fType:     FtCiAutoscan,
	typeName:  "Autoscan",
	max:       1,
	bitOffset: 3,
	bitSize:   1,
	valueType: VtOffOn,
}

var fiCiBandwidth = fieldInfo{
	fType:     FtCiBandwidth,
	typeName:  "Bandwidth",
	max:       1,
	bitOffset: 4,
	bitSize:   1,
	valueType: VtIStrings,
	strings: &[]string{
		"12.5",
		"25",
	},
}

var fiCiChannelMode = fieldInfo{
	fType:     FtCiChannelMode,
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
}

var fiCiColorCode = fieldInfo{
	fType:     FtCiColorCode,
	typeName:  "Color Code",
	max:       1,
	bitOffset: 8,
	bitSize:   4,
	valueType: VtSpan,
	span: &Span{
		min: 0,
		max: 15,
	},
	enabler: FtCiChannelMode,
}

var fiCiCompressedUdpDataHeader = fieldInfo{
	fType:     FtCiCompressedUdpDataHeader,
	typeName:  "Compressed UDP Data Header",
	max:       1,
	bitOffset: 25,
	bitSize:   1,
	valueType: VtOffOn,
	enabler:   FtCiChannelMode,
}

var fiCiContactName = fieldInfo{
	fType:        FtCiContactName,
	typeName:     "Contact Name",
	max:          1,
	bitOffset:    48,
	bitSize:      16,
	valueType:    VtListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "None"},
	},
	listRecordType: RtContacts,
	enabler:        FtCiChannelMode,
}

var fiCiCtcssDecode = fieldInfo{
	fType:        FtCiCtcssDecode,
	typeName:     "CTCSS/DCS Decode",
	max:          1,
	bitOffset:    192,
	bitSize:      16,
	valueType:    VtCtcssDcs,
	defaultValue: "None",
	disabler:     FtCiChannelMode,
}

var fiCiCtcssEncode = fieldInfo{
	fType:         FtCiCtcssEncode,
	typeName:      "CTCSS/DCS Encode",
	max:           1,
	bitOffset:     208,
	bitSize:       16,
	valueType:     VtCtcssDcs,
	defaultValue:  "None",
	enablingValue: "None",
	disabler:      FtCiChannelMode,
}

var fiCiDataCallConfirmed = fieldInfo{
	fType:     FtCiDataCallConfirmed,
	typeName:  "Data Call Confirmed",
	max:       1,
	bitOffset: 16,
	bitSize:   1,
	valueType: VtOffOn,
	enabler:   FtCiChannelMode,
}

var fiCiDecode1 = fieldInfo{
	fType:     FtCiDecode1,
	typeName:  "Decode 1",
	max:       1,
	bitOffset: 112,
	bitSize:   1,
	valueType: VtOffOn,
	disabler:  FtCiRxSignallingSystem,
}

var fiCiDecode2 = fieldInfo{
	fType:     FtCiDecode2,
	typeName:  "Decode 2",
	max:       1,
	bitOffset: 113,
	bitSize:   1,
	valueType: VtOffOn,
	disabler:  FtCiRxSignallingSystem,
}

var fiCiDecode3 = fieldInfo{
	fType:     FtCiDecode3,
	typeName:  "Decode 3",
	max:       1,
	bitOffset: 114,
	bitSize:   1,
	valueType: VtOffOn,
	disabler:  FtCiRxSignallingSystem,
}

var fiCiDecode4 = fieldInfo{
	fType:     FtCiDecode4,
	typeName:  "Decode 4",
	max:       1,
	bitOffset: 115,
	bitSize:   1,
	valueType: VtOffOn,
	disabler:  FtCiRxSignallingSystem,
}

var fiCiDecode5 = fieldInfo{
	fType:     FtCiDecode5,
	typeName:  "Decode 5",
	max:       1,
	bitOffset: 116,
	bitSize:   1,
	valueType: VtOffOn,
	disabler:  FtCiRxSignallingSystem,
}

var fiCiDecode6 = fieldInfo{
	fType:     FtCiDecode6,
	typeName:  "Decode 6",
	max:       1,
	bitOffset: 117,
	bitSize:   1,
	valueType: VtOffOn,
	disabler:  FtCiRxSignallingSystem,
}

var fiCiDecode7 = fieldInfo{
	fType:     FtCiDecode7,
	typeName:  "Decode 7",
	max:       1,
	bitOffset: 118,
	bitSize:   1,
	valueType: VtOffOn,
	disabler:  FtCiRxSignallingSystem,
}

var fiCiDecode8 = fieldInfo{
	fType:     FtCiDecode8,
	typeName:  "Decode 8",
	max:       1,
	bitOffset: 119,
	bitSize:   1,
	valueType: VtOffOn,
	disabler:  FtCiRxSignallingSystem,
}

var fiCiDisplayPTTID = fieldInfo{
	fType:     FtCiDisplayPTTID,
	typeName:  "Display PTT ID",
	max:       1,
	bitOffset: 24,
	bitSize:   1,
	valueType: VtOnOff,
	disabler:  FtCiChannelMode,
}

var fiCiEmergencyAlarmAck = fieldInfo{
	fType:     FtCiEmergencyAlarmAck,
	typeName:  "Emergency Alarm Ack",
	max:       1,
	bitOffset: 28,
	bitSize:   1,
	valueType: VtOffOn,
	enabler:   FtCiChannelMode,
}

var fiCiGPSSystem = fieldInfo{
	fType:     FtCiGPSSystem,
	typeName:  "GPS System",
	max:       1,
	bitOffset: 107,
	bitSize:   5,
	valueType: VtSpan,
	span: &Span{
		min:       0,
		max:       16,
		minString: "None",
	},
}

var fiCiGroupList = fieldInfo{
	fType:        FtCiGroupList,
	typeName:     "Group List",
	max:          1,
	bitOffset:    96,
	bitSize:      8,
	valueType:    VtListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "None"},
	},
	listRecordType: RtGroupLists,
	enabler:        FtCiChannelMode,
}

var fiCiLoneWorker = fieldInfo{
	fType:     FtCiLoneWorker,
	typeName:  "Lone Worker",
	max:       1,
	bitOffset: 0,
	bitSize:   1,
	valueType: VtOffOn,
}

var fiCiName = fieldInfo{
	fType:     FtCiName,
	typeName:  "Channel Name",
	max:       1,
	bitOffset: 256,
	bitSize:   256,
	valueType: VtUniqueName,
}

var fiCiPower = fieldInfo{
	fType:     FtCiPower,
	typeName:  "Power",
	max:       1,
	bitOffset: 34,
	bitSize:   1,
	valueType: VtIStrings,
	strings: &[]string{
		"Low",
		"High",
	},
}

var fiCiPrivacy = fieldInfo{
	fType:        FtCiPrivacy,
	typeName:     "Privacy",
	max:          1,
	bitOffset:    18,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "None",
	strings: &[]string{
		"None",
		"Basic",
		"Enhanced",
	},
	enablingValue: "None",
	enabler:       FtCiChannelMode,
}

var fiCiPrivacyNumber = fieldInfo{
	fType:        FtCiPrivacyNumber,
	typeName:     "Privacy Number",
	max:          1,
	bitOffset:    20,
	bitSize:      4,
	valueType:    VtPrivacyNumber,
	defaultValue: "0",
	span: &Span{
		min: 0,
		max: 15,
	},
	disabler: FtCiPrivacy,
}

var fiCiPrivateCallConfirmed = fieldInfo{
	fType:     FtCiPrivateCallConfirmed,
	typeName:  "Private Call Confimed",
	max:       1,
	bitOffset: 17,
	bitSize:   1,
	valueType: VtOffOn,
	enabler:   FtCiChannelMode,
}

var fiCiQtReverse = fieldInfo{
	fType:     FtCiQtReverse,
	typeName:  "QT Reverse",
	max:       1,
	bitOffset: 36,
	bitSize:   1,
	valueType: VtIStrings,
	strings: &[]string{
		"180",
		"120",
	},
	disabler: FtCiCtcssEncode,
}

var fiCiReceiveGPSInfo = fieldInfo{
	fType:     FtCiReceiveGPSInfo,
	typeName:  "Receive GPS Info",
	max:       1,
	bitOffset: 254,
	bitSize:   1,
	valueType: VtOnOff,
}

var fiCiRepeaterSlot = fieldInfo{
	fType:        FtCiRepeaterSlot,
	typeName:     "Repeater Slot",
	max:          1,
	bitOffset:    12,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "1",
	strings: &[]string{
		"",
		"1",
		"2",
	},
	enabler: FtCiChannelMode,
}

var fiCiReverseBurst = fieldInfo{
	fType:     FtCiReverseBurst,
	typeName:  "Reverse Burst/Turn Off Code",
	max:       1,
	bitOffset: 37,
	bitSize:   1,
	valueType: VtOffOn,
	disabler:  FtCiCtcssEncode,
}

var fiCiRxFrequency = fieldInfo{
	fType:     FtCiRxFrequency,
	typeName:  "Rx Frequency (MHz)",
	max:       1,
	bitOffset: 128,
	bitSize:   32,
	valueType: VtFrequency,
}

var fiCiRxOnly = fieldInfo{
	fType:     FtCiRxOnly,
	typeName:  "Rx Only",
	max:       1,
	bitOffset: 14,
	bitSize:   1,
	valueType: VtOffOn,
}

var fiCiRxRefFrequency = fieldInfo{
	fType:     FtCiRxRefFrequency,
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
}

var fiCiRxSignallingSystem = fieldInfo{
	fType:        FtCiRxSignallingSystem,
	typeName:     "Rx Signaling System",
	max:          1,
	bitOffset:    229,
	bitSize:      3,
	valueType:    VtIStrings,
	defaultValue: "Off",
	strings: &[]string{
		"Off",
		"DTMF-1",
		"DTMF-2",
		"DTMF-3",
		"DTMF-4",
	},
	enablingValue: "Off",
	disabler:      FtCiChannelMode,
}

var fiCiScanList_md380 = fieldInfo{
	fType:     FtCiScanList_md380,
	typeName:  "Scan List",
	max:       1,
	bitOffset: 88,
	bitSize:   8,
	valueType: VtListIndex,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "None"},
	},
	listRecordType: RtScanLists_md380,
}

var fiCiScanList_md40 = fieldInfo{
	fType:     FtCiScanList_md40,
	typeName:  "Scan List",
	max:       1,
	bitOffset: 88,
	bitSize:   8,
	valueType: VtListIndex,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "None"},
	},
	listRecordType: RtScanLists_md40,
}

var fiCiSendGPSInfo = fieldInfo{
	fType:     FtCiSendGPSInfo,
	typeName:  "Send GPS Info",
	max:       1,
	bitOffset: 255,
	bitSize:   1,
	valueType: VtOnOff,
}

var fiCiSquelch = fieldInfo{
	fType:     FtCiSquelch,
	typeName:  "Squelch",
	max:       1,
	bitOffset: 2,
	bitSize:   1,
	valueType: VtIStrings,
	strings: &[]string{
		"Tight",
		"Normal",
	},
}

var fiCiTot = fieldInfo{
	fType:     FtCiTot,
	typeName:  "TOT (S)",
	max:       1,
	bitOffset: 66,
	bitSize:   6,
	valueType: VtSpan,
	span: &Span{
		min:       0,
		max:       63,
		scale:     15,
		minString: "Infinite",
	},
}

var fiCiTotRekeyDelay = fieldInfo{
	fType:     FtCiTotRekeyDelay,
	typeName:  "TOT Rekey Delay (S)",
	max:       1,
	bitOffset: 72,
	bitSize:   8,
	valueType: VtSpan,
	span: &Span{
		min: 0,
		max: 255,
	},
}

var fiCiTxFrequency = fieldInfo{
	fType:     FtCiTxFrequency,
	typeName:  "Tx Frequency (MHz)",
	max:       1,
	bitOffset: 160,
	bitSize:   32,
	valueType: VtFrequency,
}

var fiCiTxRefFrequency = fieldInfo{
	fType:     FtCiTxRefFrequency,
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
}

var fiCiTxSignallingSystem = fieldInfo{
	fType:        FtCiTxSignallingSystem,
	typeName:     "Tx Signaling System",
	max:          1,
	bitOffset:    237,
	bitSize:      3,
	valueType:    VtIStrings,
	defaultValue: "Off",
	strings: &[]string{
		"Off",
		"DTMF-1",
		"DTMF-2",
		"DTMF-3",
		"DTMF-4",
	},
	disabler: FtCiChannelMode,
}

var fiCiVox = fieldInfo{
	fType:     FtCiVox,
	typeName:  "VOX",
	max:       1,
	bitOffset: 35,
	bitSize:   1,
	valueType: VtOffOn,
}

var fiDcCallID = fieldInfo{
	fType:     FtDcCallID,
	typeName:  "Call ID",
	max:       1,
	bitOffset: 0,
	bitSize:   24,
	valueType: VtCallID,
}

var fiDcCallReceiveTone = fieldInfo{
	fType:     FtDcCallReceiveTone,
	typeName:  "Call Receive Tone",
	max:       1,
	bitOffset: 26,
	bitSize:   1,
	valueType: VtIStrings,
	strings: &[]string{
		"No",
		"Yes",
	},
}

var fiDcCallType = fieldInfo{
	fType:     FtDcCallType,
	typeName:  "Call Type",
	max:       1,
	bitOffset: 30,
	bitSize:   2,
	valueType: VtCallType,
	strings: &[]string{
		"",
		"Group",
		"Private",
		"All",
	},
}

var fiDcName = fieldInfo{
	fType:     FtDcName,
	typeName:  "Contact Name",
	max:       1,
	bitOffset: 32,
	bitSize:   256,
	valueType: VtName,
}

var fiGlContact = fieldInfo{
	fType:          FtGlContact,
	typeName:       "Contacts",
	max:            32,
	bitOffset:      256,
	bitSize:        16,
	valueType:      VtListIndex,
	listRecordType: RtContacts,
}

var fiGlName = fieldInfo{
	fType:     FtGlName,
	typeName:  "Group List Name",
	max:       1,
	bitOffset: 0,
	bitSize:   256,
	valueType: VtUniqueName,
}

var fiGpDestinationID = fieldInfo{
	fType:        FtGpDestinationID,
	typeName:     "Destination ID",
	max:          1,
	bitOffset:    32,
	bitSize:      16,
	valueType:    VtGpsListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "None"},
	},
	listRecordType: RtContacts,
}

var fiGpGPSDefaultReportInterval = fieldInfo{
	fType:     FtGpGPSDefaultReportInterval,
	typeName:  "GPS Default Report Interval (S)",
	max:       1,
	bitOffset: 16,
	bitSize:   8,
	valueType: VtSpan,
	span: &Span{
		min:       0,
		max:       240,
		scale:     30,
		minString: "Off",
	},
}

var fiGpGPSRevertChannel = fieldInfo{
	fType:     FtGpGPSRevertChannel,
	typeName:  "GPS Revert Channel",
	max:       1,
	bitOffset: 0,
	bitSize:   16,
	valueType: VtGpsListIndex,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Current Channel"},
	},
	listRecordType: RtChannels_md380,
}

var fiGsBacklightColor = fieldInfo{
	fType:     FtGsBacklightColor,
	typeName:  "Backlight Color",
	max:       1,
	bitOffset: 542,
	bitSize:   2,
	valueType: VtIStrings,
	strings: &[]string{
		"Off",
		"Orange",
		"White",
		"Sakura",
	},
}

var fiGsBacklightTime = fieldInfo{
	fType:     FtGsBacklightTime,
	typeName:  "Backlight Time (S)",
	max:       1,
	bitOffset: 686,
	bitSize:   2,
	valueType: VtIStrings,
	strings: &[]string{
		"Always",
		"5",
		"10",
		"15",
	},
}

var fiGsCallAlertToneDuration = fieldInfo{
	fType:     FtGsCallAlertToneDuration,
	typeName:  "Call Alert Tone Duration (S)",
	max:       1,
	bitOffset: 632,
	bitSize:   8,
	valueType: VtSpan,
	span: &Span{
		min:       0,
		max:       240,
		scale:     5,
		minString: "Continue",
	},
}

var fiGsChFreeIndicationTone = fieldInfo{
	fType:     FtGsChFreeIndicationTone,
	typeName:  "Channel Free Indication Tone",
	max:       1,
	bitOffset: 523,
	bitSize:   1,
	valueType: VtOnOff,
}

var fiGsDisableAllLeds = fieldInfo{
	fType:     FtGsDisableAllLeds,
	typeName:  "Disable All LEDS",
	max:       1,
	bitOffset: 517,
	bitSize:   1,
	valueType: VtOnOff,
}

var fiGsDisableAllTones = fieldInfo{
	fType:     FtGsDisableAllTones,
	typeName:  "Disable All Tones",
	max:       1,
	bitOffset: 525,
	bitSize:   1,
	valueType: VtOnOff,
}

var fiGsFreqChannelMode = fieldInfo{
	fType:     FtGsFreqChannelMode,
	typeName:  "Freq/Channel Mode",
	max:       1,
	bitOffset: 540,
	bitSize:   1,
	valueType: VtIStrings,
	strings: &[]string{
		"Frequency",
		"Channel",
	},
	enablingValue: "Frequency",
}

var fiGsGroupCallHangTime = fieldInfo{
	fType:     FtGsGroupCallHangTime,
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
}

var fiGsIntroScreen = fieldInfo{
	fType:     FtGsIntroScreen,
	typeName:  "Intro Screen",
	max:       1,
	bitOffset: 531,
	bitSize:   1,
	valueType: VtIStrings,
	strings: &[]string{
		"Character String",
		"Picture",
	},
}

var fiGsIntroScreenLine1 = fieldInfo{
	fType:     FtGsIntroScreenLine1,
	typeName:  "Intro Screen Line 1",
	max:       1,
	bitOffset: 0,
	bitSize:   160,
	valueType: VtIntroLine,
}

var fiGsIntroScreenLine2 = fieldInfo{
	fType:     FtGsIntroScreenLine2,
	typeName:  "Intro Screen Line 2",
	max:       1,
	bitOffset: 160,
	bitSize:   160,
	valueType: VtIntroLine,
}

var fiGsLockUnlock = fieldInfo{
	fType:     FtGsLockUnlock,
	typeName:  "Lock/Unlock",
	max:       1,
	bitOffset: 539,
	bitSize:   1,
	valueType: VtIStrings,
	strings: &[]string{
		"Unlock",
		"Lock",
	},
	disabler: FtGsFreqChannelMode,
}

var fiGsLoneWorkerReminderTime = fieldInfo{
	fType:     FtGsLoneWorkerReminderTime,
	typeName:  "Lone Worker Reminder Time (S)",
	max:       1,
	bitOffset: 648,
	bitSize:   8,
	valueType: VtSpan,
	span: &Span{
		min: 1,
		max: 255,
	},
}

var fiGsLoneWorkerResponseTime = fieldInfo{
	fType:     FtGsLoneWorkerResponseTime,
	typeName:  "Lone Worker Response Time (min)",
	max:       1,
	bitOffset: 640,
	bitSize:   8,
	valueType: VtSpan,
	span: &Span{
		min: 1,
		max: 255,
	},
}

var fiGsMode = fieldInfo{
	fType:     FtGsMode,
	typeName:  "Mode",
	max:       1,
	bitOffset: 696,
	bitSize:   8,
	valueType: VtIndexedStrings,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Memory"},
		IndexedString{255, "Channel"},
	},
}

var fiGsModeSelect = fieldInfo{
	fType:     FtGsModeSelect,
	typeName:  "Mode Select",
	max:       1,
	bitOffset: 541,
	bitSize:   1,
	valueType: VtIStrings,
	strings: &[]string{
		"VFO",
		"Memory",
	},
	enabler: FtGsFreqChannelMode,
}

var fiGsMonitorType = fieldInfo{
	fType:     FtGsMonitorType,
	typeName:  "Monitor Type",
	max:       1,
	bitOffset: 515,
	bitSize:   1,
	valueType: VtIStrings,
	strings: &[]string{
		"Silent",
		"Open Squelch",
	},
}

var fiGsPcProgPassword = fieldInfo{
	fType:     FtGsPcProgPassword,
	typeName:  "PC Programming Password",
	max:       1,
	bitOffset: 768,
	bitSize:   64,
	valueType: VtPcPassword,
}

var fiGsPowerOnPassword = fieldInfo{
	fType:        FtGsPowerOnPassword,
	typeName:     "Power On Password",
	max:          1,
	bitOffset:    704,
	bitSize:      32,
	valueType:    VtRadioPassword,
	defaultValue: "00000000",
	enabler:      FtGsPwAndLockEnable,
}

var fiGsPrivateCallHangTime = fieldInfo{
	fType:     FtGsPrivateCallHangTime,
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
}

var fiGsPwAndLockEnable = fieldInfo{
	fType:         FtGsPwAndLockEnable,
	typeName:      "Password And Lock Enable",
	max:           1,
	bitOffset:     522,
	bitSize:       1,
	valueType:     VtOnOff,
	enablingValue: "On",
}

var fiGsRadioID = fieldInfo{
	fType:     FtGsRadioID,
	typeName:  "Radio ID",
	max:       1,
	bitOffset: 544,
	bitSize:   24,
	valueType: VtCallID,
}

var fiGsRadioName = fieldInfo{
	fType:     FtGsRadioName,
	typeName:  "Radio Name",
	max:       1,
	bitOffset: 896,
	bitSize:   256,
	valueType: VtRadioName,
}

var fiGsRadioProgPassword = fieldInfo{
	fType:     FtGsRadioProgPassword,
	typeName:  "Radio Programming Password",
	max:       1,
	bitOffset: 736,
	bitSize:   32,
	valueType: VtRadioPassword,
}

var fiGsRxLowBatteryInterval = fieldInfo{
	fType:     FtGsRxLowBatteryInterval,
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
}

var fiGsSaveModeReceive = fieldInfo{
	fType:     FtGsSaveModeReceive,
	typeName:  "Save Mode Receive",
	max:       1,
	bitOffset: 526,
	bitSize:   1,
	valueType: VtOffOn,
}

var fiGsSavePreamble = fieldInfo{
	fType:     FtGsSavePreamble,
	typeName:  "Save Preamble",
	max:       1,
	bitOffset: 527,
	bitSize:   1,
	valueType: VtOffOn,
}

var fiGsScanAnalogHangTime = fieldInfo{
	fType:     FtGsScanAnalogHangTime,
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
}

var fiGsScanDigitalHangTime = fieldInfo{
	fType:     FtGsScanDigitalHangTime,
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
}

var fiGsSetKeypadLockTime = fieldInfo{
	fType:     FtGsSetKeypadLockTime,
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
}

var fiGsTalkPermitTone = fieldInfo{
	fType:     FtGsTalkPermitTone,
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
}

var fiGsTxPreambleDuration = fieldInfo{
	fType:     FtGsTxPreambleDuration,
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
}

var fiGsVoxSensitivity = fieldInfo{
	fType:     FtGsVoxSensitivity,
	typeName:  "VOX Sensitivity",
	max:       1,
	bitOffset: 600,
	bitSize:   8,
	valueType: VtSpan,
	span: &Span{
		min: 1,
		max: 10,
	},
}

var fiSlChannel_md380 = fieldInfo{
	fType:          FtSlChannel_md380,
	typeName:       "Channels",
	max:            31,
	bitOffset:      336,
	bitSize:        16,
	valueType:      VtListIndex,
	listRecordType: RtChannels_md380,
}

var fiSlChannel_md40 = fieldInfo{
	fType:          FtSlChannel_md40,
	typeName:       "Channels",
	max:            31,
	bitOffset:      336,
	bitSize:        16,
	valueType:      VtListIndex,
	listRecordType: RtChannels_md40,
}

var fiSlName = fieldInfo{
	fType:     FtSlName,
	typeName:  "Scan List Name",
	max:       1,
	bitOffset: 0,
	bitSize:   256,
	valueType: VtUniqueName,
}

var fiSlPriorityChannel1_md380 = fieldInfo{
	fType:     FtSlPriorityChannel1_md380,
	typeName:  "Priority Channel 1",
	max:       1,
	bitOffset: 256,
	bitSize:   16,
	valueType: VtMemberListIndex,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Selected"},
		IndexedString{65535, "None"},
	},
	listRecordType: RtChannels_md380,
	enablingValue:  "None",
}

var fiSlPriorityChannel1_md40 = fieldInfo{
	fType:     FtSlPriorityChannel1_md40,
	typeName:  "Priority Channel 1",
	max:       1,
	bitOffset: 256,
	bitSize:   16,
	valueType: VtMemberListIndex,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Selected"},
		IndexedString{65535, "None"},
	},
	listRecordType: RtChannels_md40,
	enablingValue:  "None",
}

var fiSlPriorityChannel2_md380 = fieldInfo{
	fType:        FtSlPriorityChannel2_md380,
	typeName:     "Priority Channel 2",
	max:          1,
	bitOffset:    272,
	bitSize:      16,
	valueType:    VtMemberListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Selected"},
		IndexedString{65535, "None"},
	},
	listRecordType: RtChannels_md380,
	disabler:       FtSlPriorityChannel1_md380,
}

var fiSlPriorityChannel2_md40 = fieldInfo{
	fType:        FtSlPriorityChannel2_md40,
	typeName:     "Priority Channel 2",
	max:          1,
	bitOffset:    272,
	bitSize:      16,
	valueType:    VtMemberListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Selected"},
		IndexedString{65535, "None"},
	},
	listRecordType: RtChannels_md40,
	disabler:       FtSlPriorityChannel1_md40,
}

var fiSlPrioritySampleTime = fieldInfo{
	fType:     FtSlPrioritySampleTime,
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
}

var fiSlSignallingHoldTime = fieldInfo{
	fType:     FtSlSignallingHoldTime,
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
}

var fiSlTxDesignatedChannel_md380 = fieldInfo{
	fType:     FtSlTxDesignatedChannel_md380,
	typeName:  "Tx Designated Channel",
	max:       1,
	bitOffset: 288,
	bitSize:   16,
	valueType: VtListIndex,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Selected"},
		IndexedString{65535, "Last Active Channel"},
	},
	listRecordType: RtChannels_md380,
}

var fiSlTxDesignatedChannel_md40 = fieldInfo{
	fType:     FtSlTxDesignatedChannel_md40,
	typeName:  "Tx Designated Channel",
	max:       1,
	bitOffset: 288,
	bitSize:   16,
	valueType: VtListIndex,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Selected"},
		IndexedString{65535, "Last Active Channel"},
	},
	listRecordType: RtChannels_md40,
}

var fiTmTextMessage = fieldInfo{
	fType:     FtTmTextMessage,
	typeName:  "Message",
	max:       1,
	bitOffset: 0,
	bitSize:   2304,
	valueType: VtTextMessage,
}

var fiZiChannel_md380 = fieldInfo{
	fType:          FtZiChannel_md380,
	typeName:       "Channels",
	max:            16,
	bitOffset:      256,
	bitSize:        16,
	valueType:      VtListIndex,
	listRecordType: RtChannels_md380,
}

var fiZiChannel_md40 = fieldInfo{
	fType:          FtZiChannel_md40,
	typeName:       "Channels",
	max:            64,
	bitOffset:      256,
	bitSize:        16,
	valueType:      VtListIndex,
	listRecordType: RtChannels_md40,
	extOffset:      201253,
	extSize:        224,
	extIndex:       16,
}

var fiZiName = fieldInfo{
	fType:     FtZiName,
	typeName:  "Zone Name",
	max:       1,
	bitOffset: 0,
	bitSize:   256,
	valueType: VtUniqueName,
}

//go:generate genCodeplugInfo
