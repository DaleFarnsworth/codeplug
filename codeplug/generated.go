// Copyright 2017-2019 Dale Farnsworth. All rights reserved.

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
	RtBasicInformation_uv380 RecordType = "BasicInformation"
	RtButtonDefinitions      RecordType = "ButtonDefinitions"
	RtChannels_md380         RecordType = "Channels"
	RtChannels_md40          RecordType = "Channels"
	RtChannels_uv380         RecordType = "Channels"
	RtContacts               RecordType = "Contacts"
	RtContacts_uv380         RecordType = "Contacts"
	RtGPSSystems             RecordType = "GPSSystems"
	RtGeneralSettings_md2017 RecordType = "GeneralSettings"
	RtGeneralSettings_md380  RecordType = "GeneralSettings"
	RtGeneralSettings_md40   RecordType = "GeneralSettings"
	RtGeneralSettings_uv380  RecordType = "GeneralSettings"
	RtGroupLists             RecordType = "GroupLists"
	RtMenuItems              RecordType = "MenuItems"
	RtNumberKey              RecordType = "NumberKey"
	RtOneTouch               RecordType = "OneTouch"
	RtPrivacySettings        RecordType = "PrivacySettings"
	RtRadioButtons           RecordType = "RadioButtons"
	RtRadioButtons2          RecordType = "RadioButtons2"
	RtRadioButtons_md2017    RecordType = "RadioButtons"
	RtRadioButtons_md40      RecordType = "RadioButtons"
	RtRadioButtons_uv380     RecordType = "RadioButtons"
	RtScanLists_md380        RecordType = "ScanLists"
	RtScanLists_md40         RecordType = "ScanLists"
	RtScanLists_uv380        RecordType = "ScanLists"
	RtTextMessages           RecordType = "TextMessages"
	RtZones_md380            RecordType = "Zones"
	RtZones_md40             RecordType = "Zones"
	RtZones_uv380            RecordType = "Zones"
)

// Field types
const (
	FtBdLongPressDuration         FieldType = "LongPressDuration"
	FtBiCpsVersion                FieldType = "CpsVersion"
	FtBiFrequencyRangeA           FieldType = "FrequencyRangeA"
	FtBiFrequencyRangeB           FieldType = "FrequencyRangeB"
	FtBiFrequencyRange_md380      FieldType = "FrequencyRange"
	FtBiFrequencyRange_md390      FieldType = "FrequencyRange"
	FtBiFrequencyRange_md40       FieldType = "FrequencyRange"
	FtBiHighFrequency             FieldType = "HighFrequency"
	FtBiHighFrequencyA            FieldType = "HighFrequencyA"
	FtBiHighFrequencyB            FieldType = "HighFrequencyB"
	FtBiLastProgrammedTime        FieldType = "LastProgrammedTime"
	FtBiLowFrequency              FieldType = "LowFrequency"
	FtBiLowFrequencyA             FieldType = "LowFrequencyA"
	FtBiLowFrequencyB             FieldType = "LowFrequencyB"
	FtBiModel                     FieldType = "Model"
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
	FtCiDCDMSwitch                FieldType = "DCDMSwitch"
	FtCiDQTTurnoffFreq            FieldType = "DQTTurnoffFreq"
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
	FtCiEmergencySystem           FieldType = "EmergencySystem"
	FtCiGPSSystem                 FieldType = "GPSSystem"
	FtCiGroupList                 FieldType = "GroupList"
	FtCiInCallCriteria            FieldType = "InCallCriteria"
	FtCiLeaderMS                  FieldType = "LeaderMS"
	FtCiLoneWorker                FieldType = "LoneWorker"
	FtCiName                      FieldType = "Name"
	FtCiPower                     FieldType = "Power"
	FtCiPower_uv380               FieldType = "Power"
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
	FtCiSquelch_uv380             FieldType = "Squelch"
	FtCiTalkaround                FieldType = "Talkaround"
	FtCiTot                       FieldType = "Tot"
	FtCiTotRekeyDelay             FieldType = "TotRekeyDelay"
	FtCiTxFrequencyOffset         FieldType = "TxFrequencyOffset"
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
	FtGsCHVoiceAnnouncement       FieldType = "CHVoiceAnnouncement"
	FtGsCallAlertToneDuration     FieldType = "CallAlertToneDuration"
	FtGsChFreeIndicationTone      FieldType = "ChFreeIndicationTone"
	FtGsChannelsHangTime          FieldType = "ChannelsHangTime"
	FtGsDisableAllLeds            FieldType = "DisableAllLeds"
	FtGsDisableAllTones           FieldType = "DisableAllTones"
	FtGsEditRadioID               FieldType = "EditRadioID"
	FtGsEnableContactsCSV         FieldType = "EnableContactsCSV"
	FtGsFreqChannelMode           FieldType = "FreqChannelMode"
	FtGsFreqChannelMode_uv380     FieldType = "FreqChannelMode"
	FtGsGroupCallHangTime         FieldType = "GroupCallHangTime"
	FtGsGroupCallMatch            FieldType = "GroupCallMatch"
	FtGsIntroScreen               FieldType = "IntroScreen"
	FtGsIntroScreenLine1          FieldType = "IntroScreenLine1"
	FtGsIntroScreenLine2          FieldType = "IntroScreenLine2"
	FtGsKeypadTones               FieldType = "KeypadTones"
	FtGsLockUnlock                FieldType = "LockUnlock"
	FtGsLoneWorkerReminderTime    FieldType = "LoneWorkerReminderTime"
	FtGsLoneWorkerResponseTime    FieldType = "LoneWorkerResponseTime"
	FtGsMenuControl               FieldType = "MenuControl"
	FtGsMicLevel                  FieldType = "MicLevel"
	FtGsModeSelect                FieldType = "ModeSelect"
	FtGsModeSelectA               FieldType = "ModeSelectA"
	FtGsModeSelectB               FieldType = "ModeSelectB"
	FtGsMonitorType               FieldType = "MonitorType"
	FtGsPcProgPassword            FieldType = "PcProgPassword"
	FtGsPowerOnPassword           FieldType = "PowerOnPassword"
	FtGsPrivateCallHangTime       FieldType = "PrivateCallHangTime"
	FtGsPrivateCallMatch          FieldType = "PrivateCallMatch"
	FtGsPublicZone                FieldType = "PublicZone"
	FtGsPwAndLockEnable           FieldType = "PwAndLockEnable"
	FtGsRadioID                   FieldType = "RadioID"
	FtGsRadioID1                  FieldType = "RadioID1"
	FtGsRadioID2                  FieldType = "RadioID2"
	FtGsRadioID3                  FieldType = "RadioID3"
	FtGsRadioName                 FieldType = "RadioName"
	FtGsRadioProgPassword         FieldType = "RadioProgPassword"
	FtGsRxLowBatteryInterval      FieldType = "RxLowBatteryInterval"
	FtGsSaveModeReceive           FieldType = "SaveModeReceive"
	FtGsSavePreamble              FieldType = "SavePreamble"
	FtGsScanAnalogHangTime        FieldType = "ScanAnalogHangTime"
	FtGsScanDigitalHangTime       FieldType = "ScanDigitalHangTime"
	FtGsSetKeypadLockTime         FieldType = "SetKeypadLockTime"
	FtGsTalkPermitTone            FieldType = "TalkPermitTone"
	FtGsTimeZone                  FieldType = "TimeZone"
	FtGsTwoChannel                FieldType = "TwoChannel"
	FtGsTxMode                    FieldType = "TxMode"
	FtGsTxPreambleDuration        FieldType = "TxPreambleDuration"
	FtGsVoxSensitivity            FieldType = "VoxSensitivity"
	FtMiAnswered                  FieldType = "Answered"
	FtMiBacklight                 FieldType = "Backlight"
	FtMiCallAlert                 FieldType = "CallAlert"
	FtMiDisplayMode               FieldType = "DisplayMode"
	FtMiEdit                      FieldType = "Edit"
	FtMiEditList                  FieldType = "EditList"
	FtMiGps                       FieldType = "Gps"
	FtMiHangTime                  FieldType = "HangTime"
	FtMiIntroScreen               FieldType = "IntroScreen"
	FtMiKeyboardLock              FieldType = "KeyboardLock"
	FtMiLedIndicator              FieldType = "LedIndicator"
	FtMiManualDial                FieldType = "ManualDial"
	FtMiMissed                    FieldType = "Missed"
	FtMiOutgoingRadio             FieldType = "OutgoingRadio"
	FtMiPasswordAndLock           FieldType = "PasswordAndLock"
	FtMiPower                     FieldType = "Power"
	FtMiProgramKey                FieldType = "ProgramKey"
	FtMiProgramRadio              FieldType = "ProgramRadio"
	FtMiRadioCheck                FieldType = "RadioCheck"
	FtMiRadioDisable              FieldType = "RadioDisable"
	FtMiRadioEnable               FieldType = "RadioEnable"
	FtMiRemoteMonitor             FieldType = "RemoteMonitor"
	FtMiScan                      FieldType = "Scan"
	FtMiSquelch                   FieldType = "Squelch"
	FtMiTalkaround                FieldType = "Talkaround"
	FtMiTextMessage               FieldType = "TextMessage"
	FtMiToneOrAlert               FieldType = "ToneOrAlert"
	FtMiVox                       FieldType = "Vox"
	FtNkContact                   FieldType = "Contact"
	FtOtCall                      FieldType = "Call"
	FtOtCallType                  FieldType = "CallType"
	FtOtDTMF                      FieldType = "DTMF"
	FtOtEncode                    FieldType = "Encode"
	FtOtMode                      FieldType = "Mode"
	FtOtTextMessage               FieldType = "TextMessage"
	FtPsBasicKey                  FieldType = "BasicKey"
	FtPsEnhancedKey               FieldType = "EnhancedKey"
	FtRbButton                    FieldType = "Button"
	FtRbButton_md40               FieldType = "Button"
	FtRbButton_uv380              FieldType = "Button"
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
	FtZiChannelA_uv380            FieldType = "ChannelA"
	FtZiChannelB_uv380            FieldType = "ChannelB"
	FtZiChannel_md380             FieldType = "Channel"
	FtZiChannel_md40              FieldType = "Channel"
	FtZiName                      FieldType = "Name"
)

// The value types a field may contain
const (
	VtAscii              ValueType = "ascii"
	VtBandwidth          ValueType = "bandwidth"
	VtBiFrequency        ValueType = "biFrequency"
	VtCallID             ValueType = "callID"
	VtCallType           ValueType = "callType"
	VtContactListIndex   ValueType = "contactListIndex"
	VtContactName        ValueType = "contactName"
	VtCpsVersion         ValueType = "cpsVersion"
	VtCtcssDcs           ValueType = "ctcssDcs"
	VtDerefListIndex     ValueType = "derefListIndex"
	VtFrequency          ValueType = "frequency"
	VtFrequencyOffset    ValueType = "frequencyOffset"
	VtGpsListIndex       ValueType = "gpsListIndex"
	VtGpsReportInterval  ValueType = "gpsReportInterval"
	VtHexadecimal32      ValueType = "hexadecimal32"
	VtHexadecimal4       ValueType = "hexadecimal4"
	VtIStrings           ValueType = "iStrings"
	VtIndexedStrings     ValueType = "indexedStrings"
	VtIntroLine          ValueType = "introLine"
	VtListIndex          ValueType = "listIndex"
	VtMemberListIndex    ValueType = "memberListIndex"
	VtName               ValueType = "name"
	VtNkContactListIndex ValueType = "nkContactListIndex"
	VtOffOn              ValueType = "offOn"
	VtOnOff              ValueType = "onOff"
	VtPcPassword         ValueType = "pcPassword"
	VtPrivacyNumber      ValueType = "privacyNumber"
	VtRadioButton        ValueType = "radioButton"
	VtRadioName          ValueType = "radioName"
	VtRadioPassword      ValueType = "radioPassword"
	VtRadioProgPassword  ValueType = "radioProgPassword"
	VtSpan               ValueType = "span"
	VtSpanList           ValueType = "spanList"
	VtTextMessage        ValueType = "textMessage"
	VtTimeStamp          ValueType = "timeStamp"
)

// newValue returns a new value of the given ValueType
func newValue(vt ValueType) value {
	switch vt {
	case VtAscii:
		return new(ascii)
	case VtBandwidth:
		return new(bandwidth)
	case VtBiFrequency:
		return new(biFrequency)
	case VtCallID:
		return new(callID)
	case VtCallType:
		return new(callType)
	case VtContactListIndex:
		return new(contactListIndex)
	case VtContactName:
		return new(contactName)
	case VtCpsVersion:
		return new(cpsVersion)
	case VtCtcssDcs:
		return new(ctcssDcs)
	case VtDerefListIndex:
		return new(derefListIndex)
	case VtFrequency:
		return new(frequency)
	case VtFrequencyOffset:
		return new(frequencyOffset)
	case VtGpsListIndex:
		return new(gpsListIndex)
	case VtGpsReportInterval:
		return new(gpsReportInterval)
	case VtHexadecimal32:
		return new(hexadecimal32)
	case VtHexadecimal4:
		return new(hexadecimal4)
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
	case VtNkContactListIndex:
		return new(nkContactListIndex)
	case VtOffOn:
		return new(offOn)
	case VtOnOff:
		return new(onOff)
	case VtPcPassword:
		return new(pcPassword)
	case VtPrivacyNumber:
		return new(privacyNumber)
	case VtRadioButton:
		return new(radioButton)
	case VtRadioName:
		return new(radioName)
	case VtRadioPassword:
		return new(radioPassword)
	case VtRadioProgPassword:
		return new(radioProgPassword)
	case VtSpan:
		return new(span)
	case VtSpanList:
		return new(spanList)
	case VtTextMessage:
		return new(textMessage)
	case VtTimeStamp:
		return new(timeStamp)
	}

	return nil
}

var codeplugInfos = []*CodeplugInfo{
	&cpMD380,
	&cpRT3,
	&cpMD390,
	&cpRT3G,
	&cpDJMD40,
	&cpMDUV380,
	&cpMDUV390,
	&cpRT3S,
	&cpMD2017,
	&cpRT82,
}

var cpMD380 = CodeplugInfo{
	Type: "MD-380",
	Models: []string{
		"MD380",
		"DR780",
	},
	Ext:           "rdt",
	RdtSize:       262709,
	HeaderSize:    549,
	TrailerOffset: 262693,
	TrailerSize:   16,
	RecordInfos: []*recordInfo{
		&riBasicInformation_md380,
		&riGeneralSettings_md380,
		&riMenuItems,
		&riRadioButtons,
		&riButtonDefinitions,
		&riOneTouch,
		&riNumberKey,
		&riTextMessages,
		&riPrivacySettings,
		&riContacts,
		&riGroupLists,
		&riZones_md380,
		&riScanLists_md380,
		&riChannels_md380,
		&riGPSSystems,
	},
}

var cpRT3 = CodeplugInfo{
	Type: "RT3",
	Models: []string{
		"DR780",
	},
	Ext:           "rdt",
	RdtSize:       262709,
	HeaderSize:    549,
	TrailerOffset: 262693,
	TrailerSize:   16,
	RecordInfos: []*recordInfo{
		&riBasicInformation_md380,
		&riGeneralSettings_md380,
		&riMenuItems,
		&riRadioButtons,
		&riButtonDefinitions,
		&riOneTouch,
		&riNumberKey,
		&riTextMessages,
		&riPrivacySettings,
		&riContacts,
		&riGroupLists,
		&riZones_md380,
		&riScanLists_md380,
		&riChannels_md380,
		&riGPSSystems,
	},
}

var cpMD390 = CodeplugInfo{
	Type: "MD-390",
	Models: []string{
		"MD390",
	},
	Ext:           "rdt",
	RdtSize:       262709,
	HeaderSize:    549,
	TrailerOffset: 262693,
	TrailerSize:   16,
	RecordInfos: []*recordInfo{
		&riBasicInformation_md390,
		&riGeneralSettings_md380,
		&riMenuItems,
		&riRadioButtons,
		&riButtonDefinitions,
		&riOneTouch,
		&riNumberKey,
		&riTextMessages,
		&riPrivacySettings,
		&riContacts,
		&riGroupLists,
		&riZones_md380,
		&riScanLists_md380,
		&riChannels_md380,
		&riGPSSystems,
	},
}

var cpRT3G = CodeplugInfo{
	Type: "RT3-G",
	Models: []string{
		"MD390",
	},
	Ext:           "rdt",
	RdtSize:       262709,
	HeaderSize:    549,
	TrailerOffset: 262693,
	TrailerSize:   16,
	RecordInfos: []*recordInfo{
		&riBasicInformation_md390,
		&riGeneralSettings_md380,
		&riMenuItems,
		&riRadioButtons,
		&riButtonDefinitions,
		&riOneTouch,
		&riNumberKey,
		&riTextMessages,
		&riPrivacySettings,
		&riContacts,
		&riGroupLists,
		&riZones_md380,
		&riScanLists_md380,
		&riChannels_md380,
		&riGPSSystems,
	},
}

var cpDJMD40 = CodeplugInfo{
	Type: "DJ-MD40",
	Models: []string{
		"DJ-MD40",
	},
	Ext:           "rdt",
	RdtSize:       262709,
	HeaderSize:    549,
	TrailerOffset: 262693,
	TrailerSize:   16,
	RecordInfos: []*recordInfo{
		&riBasicInformation_md40,
		&riGeneralSettings_md40,
		&riMenuItems,
		&riRadioButtons_md40,
		&riRadioButtons2,
		&riButtonDefinitions,
		&riOneTouch,
		&riNumberKey,
		&riTextMessages,
		&riPrivacySettings,
		&riContacts,
		&riGroupLists,
		&riZones_md40,
		&riScanLists_md40,
		&riChannels_md40,
	},
}

var cpMDUV380 = CodeplugInfo{
	Type: "MD-UV380",
	Models: []string{
		"MD-UV380",
	},
	Ext:           "rdt",
	RdtSize:       852533,
	HeaderSize:    549,
	TrailerOffset: 262693,
	TrailerSize:   16,
	RecordInfos: []*recordInfo{
		&riBasicInformation_uv380,
		&riGeneralSettings_uv380,
		&riMenuItems,
		&riRadioButtons_uv380,
		&riButtonDefinitions,
		&riOneTouch,
		&riNumberKey,
		&riTextMessages,
		&riPrivacySettings,
		&riContacts_uv380,
		&riGroupLists,
		&riZones_uv380,
		&riScanLists_uv380,
		&riChannels_uv380,
		&riGPSSystems,
	},
}

var cpMDUV390 = CodeplugInfo{
	Type: "MD-UV390",
	Models: []string{
		"MD-UV390",
	},
	Ext:           "rdt",
	RdtSize:       852533,
	HeaderSize:    549,
	TrailerOffset: 262693,
	TrailerSize:   16,
	RecordInfos: []*recordInfo{
		&riBasicInformation_uv380,
		&riGeneralSettings_uv380,
		&riMenuItems,
		&riRadioButtons_uv380,
		&riButtonDefinitions,
		&riOneTouch,
		&riNumberKey,
		&riTextMessages,
		&riPrivacySettings,
		&riContacts_uv380,
		&riGroupLists,
		&riZones_uv380,
		&riScanLists_uv380,
		&riChannels_uv380,
		&riGPSSystems,
	},
}

var cpRT3S = CodeplugInfo{
	Type: "RT3S",
	Models: []string{
		"MD-UV390",
	},
	Ext:           "rdt",
	RdtSize:       852533,
	HeaderSize:    549,
	TrailerOffset: 262693,
	TrailerSize:   16,
	RecordInfos: []*recordInfo{
		&riBasicInformation_uv380,
		&riGeneralSettings_uv380,
		&riMenuItems,
		&riRadioButtons_uv380,
		&riButtonDefinitions,
		&riOneTouch,
		&riNumberKey,
		&riTextMessages,
		&riPrivacySettings,
		&riContacts_uv380,
		&riGroupLists,
		&riZones_uv380,
		&riScanLists_uv380,
		&riChannels_uv380,
		&riGPSSystems,
	},
}

var cpMD2017 = CodeplugInfo{
	Type: "MD-2017",
	Models: []string{
		"2017",
	},
	Ext:           "rdt",
	RdtSize:       852533,
	HeaderSize:    549,
	TrailerOffset: 262693,
	TrailerSize:   16,
	RecordInfos: []*recordInfo{
		&riBasicInformation_uv380,
		&riGeneralSettings_md2017,
		&riMenuItems,
		&riRadioButtons_md2017,
		&riButtonDefinitions,
		&riOneTouch,
		&riNumberKey,
		&riTextMessages,
		&riPrivacySettings,
		&riContacts_uv380,
		&riGroupLists,
		&riZones_uv380,
		&riScanLists_uv380,
		&riChannels_uv380,
		&riGPSSystems,
	},
}

var cpRT82 = CodeplugInfo{
	Type: "RT82",
	Models: []string{
		"2017",
	},
	Ext:           "rdt",
	RdtSize:       852533,
	HeaderSize:    549,
	TrailerOffset: 262693,
	TrailerSize:   16,
	RecordInfos: []*recordInfo{
		&riBasicInformation_uv380,
		&riGeneralSettings_md2017,
		&riMenuItems,
		&riRadioButtons_md2017,
		&riButtonDefinitions,
		&riOneTouch,
		&riNumberKey,
		&riTextMessages,
		&riPrivacySettings,
		&riContacts_uv380,
		&riGroupLists,
		&riZones_uv380,
		&riScanLists_uv380,
		&riChannels_uv380,
		&riGPSSystems,
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
		&fiBiLowFrequency,
		&fiBiHighFrequency,
		&fiBiLastProgrammedTime,
		&fiBiCpsVersion,
	},
}

var riBasicInformation_uv380 = recordInfo{
	rType:    RtBasicInformation_uv380,
	typeName: "Basic Information",
	max:      1,
	offset:   0,
	size:     8805,
	fieldInfos: []*fieldInfo{
		&fiBiModel,
		&fiBiFrequencyRangeA,
		&fiBiFrequencyRangeB,
		&fiBiLowFrequencyA,
		&fiBiHighFrequencyA,
		&fiBiLowFrequencyB,
		&fiBiHighFrequencyB,
		&fiBiLastProgrammedTime,
		&fiBiCpsVersion,
	},
}

var riButtonDefinitions = recordInfo{
	rType:    RtButtonDefinitions,
	typeName: "Button Definitions",
	max:      1,
	offset:   9014,
	size:     1,
	fieldInfos: []*fieldInfo{
		&fiBdLongPressDuration,
	},
}

var riChannels_md380 = recordInfo{
	rType:    RtChannels_md380,
	typeName: "Channels",
	max:      1000,
	offset:   127013,
	size:     64,
	delDesc: &delDesc{
		offset: 16,
		size:   1,
		value:  255,
	},
	fieldInfos: []*fieldInfo{
		&fiCiName,
		&fiCiRxFrequency,
		&fiCiTxFrequencyOffset,
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
		&fiCiTalkaround,
		&fiCiEmergencySystem,
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
		&fiCiInCallCriteria,
	},
}

var riChannels_md40 = recordInfo{
	rType:    RtChannels_md40,
	typeName: "Channels",
	max:      1000,
	offset:   127013,
	size:     64,
	delDesc: &delDesc{
		offset: 16,
		size:   1,
		value:  255,
	},
	fieldInfos: []*fieldInfo{
		&fiCiName,
		&fiCiRxFrequency,
		&fiCiTxFrequencyOffset,
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
		&fiCiTalkaround,
		&fiCiEmergencySystem,
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

var riChannels_uv380 = recordInfo{
	rType:    RtChannels_uv380,
	typeName: "Channels",
	max:      3000,
	offset:   262709,
	size:     64,
	delDesc: &delDesc{
		offset: 16,
		size:   1,
		value:  255,
	},
	fieldInfos: []*fieldInfo{
		&fiCiName,
		&fiCiRxFrequency,
		&fiCiTxFrequencyOffset,
		&fiCiChannelMode,
		&fiCiBandwidth,
		&fiCiScanList_md40,
		&fiCiSquelch_uv380,
		&fiCiRxRefFrequency,
		&fiCiTxRefFrequency,
		&fiCiTot,
		&fiCiTotRekeyDelay,
		&fiCiPower_uv380,
		&fiCiAdmitCriteria,
		&fiCiAutoscan,
		&fiCiRxOnly,
		&fiCiLoneWorker,
		&fiCiVox,
		&fiCiAllowTalkaround,
		&fiCiSendGPSInfo,
		&fiCiReceiveGPSInfo,
		&fiCiPrivateCallConfirmed,
		&fiCiTalkaround,
		&fiCiEmergencyAlarmAck,
		&fiCiDataCallConfirmed,
		&fiCiDCDMSwitch,
		&fiCiLeaderMS,
		&fiCiEmergencySystem,
		&fiCiContactName,
		&fiCiGroupList,
		&fiCiColorCode,
		&fiCiRepeaterSlot,
		&fiCiPrivacy,
		&fiCiPrivacyNumber,
		&fiCiGPSSystem,
		&fiCiInCallCriteria,
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
	delDesc: &delDesc{
		offset: 3,
		size:   1,
		value:  192,
	},
	fieldInfos: []*fieldInfo{
		&fiDcName,
		&fiDcCallID,
		&fiDcCallType,
		&fiDcCallReceiveTone,
	},
}

var riContacts_uv380 = recordInfo{
	rType:    RtContacts_uv380,
	typeName: "Contacts",
	max:      10000,
	offset:   459317,
	size:     36,
	delDesc: &delDesc{
		offset: 3,
		size:   1,
		value:  192,
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

var riGeneralSettings_md2017 = recordInfo{
	rType:    RtGeneralSettings_md2017,
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
		&fiGsCHVoiceAnnouncement,
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
		&fiGsChannelsHangTime,
		&fiGsPcProgPassword,
		&fiGsRadioProgPassword,
		&fiGsSetKeypadLockTime,
		&fiGsFreqChannelMode_uv380,
		&fiGsModeSelectA,
		&fiGsModeSelectB,
		&fiGsTimeZone,
		&fiGsBacklightTime,
		&fiGsDisableAllLeds,
		&fiGsGroupCallMatch,
		&fiGsPrivateCallMatch,
		&fiGsGroupCallHangTime,
		&fiGsPrivateCallHangTime,
		&fiGsRadioID1,
		&fiGsRadioID2,
		&fiGsRadioID3,
		&fiGsMicLevel,
		&fiGsTxMode,
		&fiGsEditRadioID,
		&fiGsPublicZone,
		&fiGsEnableContactsCSV,
		&fiGsMenuControl,
		&fiGsTwoChannel,
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
		&fiGsKeypadTones,
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
		&fiGsKeypadTones,
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

var riGeneralSettings_uv380 = recordInfo{
	rType:    RtGeneralSettings_uv380,
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
		&fiGsCHVoiceAnnouncement,
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
		&fiGsChannelsHangTime,
		&fiGsPcProgPassword,
		&fiGsRadioProgPassword,
		&fiGsSetKeypadLockTime,
		&fiGsFreqChannelMode_uv380,
		&fiGsModeSelectA,
		&fiGsModeSelectB,
		&fiGsTimeZone,
		&fiGsBacklightTime,
		&fiGsDisableAllLeds,
		&fiGsGroupCallMatch,
		&fiGsPrivateCallMatch,
		&fiGsGroupCallHangTime,
		&fiGsPrivateCallHangTime,
		&fiGsRadioID1,
		&fiGsRadioID2,
		&fiGsRadioID3,
		&fiGsMicLevel,
		&fiGsTxMode,
		&fiGsEditRadioID,
		&fiGsPublicZone,
		&fiGsEnableContactsCSV,
	},
}

var riGroupLists = recordInfo{
	rType:    RtGroupLists,
	typeName: "RX Group Lists",
	max:      250,
	offset:   60997,
	size:     96,
	delDesc: &delDesc{
		offset: 0,
		size:   1,
		value:  0,
	},
	fieldInfos: []*fieldInfo{
		&fiGlName,
		&fiGlContact,
	},
}

var riMenuItems = recordInfo{
	rType:    RtMenuItems,
	typeName: "Menu Items",
	max:      1,
	offset:   8981,
	size:     5,
	fieldInfos: []*fieldInfo{
		&fiMiHangTime,
		&fiMiRadioDisable,
		&fiMiRadioEnable,
		&fiMiRemoteMonitor,
		&fiMiRadioCheck,
		&fiMiManualDial,
		&fiMiEdit,
		&fiMiCallAlert,
		&fiMiTextMessage,
		&fiMiToneOrAlert,
		&fiMiTalkaround,
		&fiMiOutgoingRadio,
		&fiMiAnswered,
		&fiMiMissed,
		&fiMiEditList,
		&fiMiScan,
		&fiMiProgramKey,
		&fiMiVox,
		&fiMiSquelch,
		&fiMiLedIndicator,
		&fiMiKeyboardLock,
		&fiMiIntroScreen,
		&fiMiBacklight,
		&fiMiPower,
		&fiMiGps,
		&fiMiProgramRadio,
		&fiMiDisplayMode,
		&fiMiPasswordAndLock,
	},
}

var riNumberKey = recordInfo{
	rType:      RtNumberKey,
	typeName:   "Number Key",
	max:        10,
	offset:     9041,
	size:       2,
	namePrefix: "Number Key ",
	fieldInfos: []*fieldInfo{
		&fiNkContact,
	},
}

var riOneTouch = recordInfo{
	rType:    RtOneTouch,
	typeName: "One Touch",
	max:      6,
	offset:   9017,
	size:     4,
	fieldInfos: []*fieldInfo{
		&fiOtMode,
		&fiOtCallType,
		&fiOtDTMF,
		&fiOtTextMessage,
		&fiOtEncode,
		&fiOtCall,
	},
}

var riPrivacySettings = recordInfo{
	rType:    RtPrivacySettings,
	typeName: "Privacy Settings",
	max:      1,
	offset:   23525,
	size:     176,
	fieldInfos: []*fieldInfo{
		&fiPsEnhancedKey,
		&fiPsBasicKey,
	},
}

var riRadioButtons = recordInfo{
	rType:    RtRadioButtons,
	typeName: "Radio Buttons",
	max:      4,
	offset:   8999,
	size:     1,
	names: []string{
		"Side Button 1 Short Press",
		"Side Button 1 Long Press",
		"Side Button 2 Short Press",
		"Side Button 2 Long Press",
	},
	fieldInfos: []*fieldInfo{
		&fiRbButton,
	},
}

var riRadioButtons2 = recordInfo{
	rType:    RtRadioButtons2,
	typeName: "Radio Buttons",
	max:      12,
	offset:   9061,
	size:     1,
	names: []string{
		"Func + Key 0",
		"Func + Key 1",
		"Func + Key 2",
		"Func + Key 3",
		"Func + Key 4",
		"Func + Key 5",
		"Func + Key 6",
		"Func + Key 7",
		"Func + Key 8",
		"Func + Key 9",
		"Func + Key #",
		"Func + Key *",
	},
	fieldInfos: []*fieldInfo{
		&fiRbButton_md40,
	},
}

var riRadioButtons_md2017 = recordInfo{
	rType:    RtRadioButtons_md2017,
	typeName: "Radio Buttons",
	max:      12,
	offset:   8997,
	size:     1,
	names: []string{
		"Top Button 1 Short Press",
		"Top Button 1 Long Press",
		"Side Button 1 Short Press",
		"Side Button 1 Long Press",
		"Side Button 2 Short Press",
		"Side Button 2 Long Press",
		"Side Button 3 Short Press",
		"Side Button 3 Long Press",
		"P1 Short Press",
		"P1 Long Press",
		"P2 Short Press",
		"P2 Long Press",
	},
	fieldInfos: []*fieldInfo{
		&fiRbButton_uv380,
	},
}

var riRadioButtons_md40 = recordInfo{
	rType:    RtRadioButtons_md40,
	typeName: "Radio Buttons",
	max:      6,
	offset:   8997,
	size:     1,
	names: []string{
		"Top Button 1 Short Press",
		"Top Button 1 Long Press",
		"Side Button 1 Short Press",
		"Side Button 1 Long Press",
		"Side Button 2 Short Press",
		"Side Button 2 Long Press",
	},
	fieldInfos: []*fieldInfo{
		&fiRbButton_md40,
	},
}

var riRadioButtons_uv380 = recordInfo{
	rType:    RtRadioButtons_uv380,
	typeName: "Radio Buttons",
	max:      4,
	offset:   8999,
	size:     1,
	names: []string{
		"Side Button 1 Short Press",
		"Side Button 1 Long Press",
		"Side Button 2 Short Press",
		"Side Button 2 Long Press",
	},
	fieldInfos: []*fieldInfo{
		&fiRbButton_uv380,
	},
}

var riScanLists_md380 = recordInfo{
	rType:    RtScanLists_md380,
	typeName: "Scan Lists",
	max:      250,
	offset:   100997,
	size:     104,
	delDesc: &delDesc{
		offset: 0,
		size:   1,
		value:  0,
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
	delDesc: &delDesc{
		offset: 0,
		size:   1,
		value:  0,
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

var riScanLists_uv380 = recordInfo{
	rType:    RtScanLists_uv380,
	typeName: "Scan Lists",
	max:      250,
	offset:   100997,
	size:     104,
	delDesc: &delDesc{
		offset: 0,
		size:   1,
		value:  0,
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
	rType:      RtTextMessages,
	typeName:   "Text Messages",
	max:        50,
	offset:     9125,
	size:       288,
	namePrefix: "Message ",
	delDesc: &delDesc{
		offset: 0,
		size:   8,
		value:  0,
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
	delDesc: &delDesc{
		offset: 0,
		size:   1,
		value:  0,
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
	delDesc: &delDesc{
		offset: 0,
		size:   1,
		value:  0,
	},
	fieldInfos: []*fieldInfo{
		&fiZiName,
		&fiZiChannel_md40,
	},
}

var riZones_uv380 = recordInfo{
	rType:    RtZones_uv380,
	typeName: "Zones",
	max:      250,
	offset:   84997,
	size:     64,
	delDesc: &delDesc{
		offset: 0,
		size:   1,
		value:  0,
	},
	fieldInfos: []*fieldInfo{
		&fiZiName,
		&fiZiChannelA_uv380,
		&fiZiChannelB_uv380,
	},
}

var fiBdLongPressDuration = fieldInfo{
	fType:        FtBdLongPressDuration,
	typeName:     "Long Press Duration (mS)",
	max:          1,
	bitOffset:    0,
	bitSize:      8,
	valueType:    VtSpan,
	defaultValue: "1000",
	span: &Span{
		min:   4,
		max:   15,
		scale: 250,
	},
}

var fiBiCpsVersion = fieldInfo{
	fType:        FtBiCpsVersion,
	typeName:     "CPS Version",
	max:          1,
	bitOffset:    69992,
	bitSize:      32,
	valueType:    VtCpsVersion,
	defaultValue: "1001",
}

var fiBiFrequencyRangeA = fieldInfo{
	fType:     FtBiFrequencyRangeA,
	typeName:  "Frequency Range A (MHz)",
	max:       1,
	bitOffset: 2480,
	bitSize:   8,
	valueType: VtIStrings,
	strings: &[]string{
		"136-174",
		"350-400",
		"400-480",
		"450-520",
	},
}

var fiBiFrequencyRangeB = fieldInfo{
	fType:     FtBiFrequencyRangeB,
	typeName:  "Frequency Range B (MHz)",
	max:       1,
	bitOffset: 2488,
	bitSize:   8,
	valueType: VtIStrings,
	strings: &[]string{
		"136-174",
		"350-400",
		"400-480",
		"450-520",
	},
}

var fiBiFrequencyRange_md380 = fieldInfo{
	fType:     FtBiFrequencyRange_md380,
	typeName:  "Frequency Range (MHz)",
	max:       1,
	bitOffset: 2480,
	bitSize:   8,
	valueType: VtIStrings,
	strings: &[]string{
		"136-174",
		"350-400",
		"400-480",
		"450-520",
	},
}

var fiBiFrequencyRange_md390 = fieldInfo{
	fType:     FtBiFrequencyRange_md390,
	typeName:  "Frequency Range (MHz)",
	max:       1,
	bitOffset: 2480,
	bitSize:   8,
	valueType: VtIStrings,
	strings: &[]string{
		"136-174",
		"350-400",
		"400-480",
		"450-520",
	},
}

var fiBiFrequencyRange_md40 = fieldInfo{
	fType:     FtBiFrequencyRange_md40,
	typeName:  "Frequency Range (MHz)",
	max:       1,
	bitOffset: 2480,
	bitSize:   8,
	valueType: VtIStrings,
	strings: &[]string{
		"MD40 400-480",
		"MD40HT 420-450",
		"MD40HE 430-440",
		"MD40T 400-480",
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

var fiBiHighFrequencyA = fieldInfo{
	fType:     FtBiHighFrequencyA,
	typeName:  "High Frequency A",
	max:       1,
	bitOffset: 2520,
	bitSize:   16,
	valueType: VtBiFrequency,
}

var fiBiHighFrequencyB = fieldInfo{
	fType:     FtBiHighFrequencyB,
	typeName:  "High Frequency B",
	max:       1,
	bitOffset: 2552,
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

var fiBiLowFrequencyA = fieldInfo{
	fType:     FtBiLowFrequencyA,
	typeName:  "Low Frequency A",
	max:       1,
	bitOffset: 2504,
	bitSize:   16,
	valueType: VtBiFrequency,
}

var fiBiLowFrequencyB = fieldInfo{
	fType:     FtBiLowFrequencyB,
	typeName:  "Low Frequency B",
	max:       1,
	bitOffset: 2536,
	bitSize:   16,
	valueType: VtBiFrequency,
}

var fiBiModel = fieldInfo{
	fType:     FtBiModel,
	typeName:  "Model Name",
	max:       1,
	bitOffset: 2344,
	bitSize:   64,
	valueType: VtAscii,
}

var fiCiAdmitCriteria = fieldInfo{
	fType:        FtCiAdmitCriteria,
	typeName:     "Admit Criteria",
	max:          1,
	bitOffset:    32,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "Always",
	strings: &[]string{
		"Always",
		"Channel free",
		"CTCSS/DCS",
		"Color code",
	},
}

var fiCiAllowTalkaround = fieldInfo{
	fType:        FtCiAllowTalkaround,
	typeName:     "Allow Talkaround",
	max:          1,
	bitOffset:    15,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiTxFrequencyOffset,
	enablers: []enabler{
		enabler{
			value:  "+0.00000",
			enable: false,
		},
	},
}

var fiCiAutoscan = fieldInfo{
	fType:        FtCiAutoscan,
	typeName:     "Autoscan",
	max:          1,
	bitOffset:    3,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
}

var fiCiBandwidth = fieldInfo{
	fType:        FtCiBandwidth,
	typeName:     "Bandwidth (KHz)",
	max:          1,
	bitOffset:    4,
	bitSize:      2,
	valueType:    VtBandwidth,
	defaultValue: "12.5",
	strings: &[]string{
		"12.5",
		"20",
		"25",
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: false,
		},
	},
}

var fiCiChannelMode = fieldInfo{
	fType:        FtCiChannelMode,
	typeName:     "Channel Mode",
	max:          1,
	bitOffset:    6,
	bitSize:      2,
	valueType:    VtIndexedStrings,
	defaultValue: "Analog",
	indexedStrings: &[]IndexedString{
		IndexedString{1, "Analog"},
		IndexedString{2, "Digital"},
	},
	enables: []FieldType{
		FtCiBandwidth,
		FtCiColorCode,
		FtCiCompressedUdpDataHeader,
		FtCiContactName,
		FtCiCtcssDecode,
		FtCiCtcssEncode,
		FtCiDCDMSwitch,
		FtCiDQTTurnoffFreq,
		FtCiDataCallConfirmed,
		FtCiDisplayPTTID,
		FtCiEmergencyAlarmAck,
		FtCiEmergencySystem,
		FtCiGPSSystem,
		FtCiGroupList,
		FtCiInCallCriteria,
		FtCiPrivacy,
		FtCiPrivateCallConfirmed,
		FtCiRepeaterSlot,
		FtCiRxSignallingSystem,
		FtCiTxSignallingSystem,
	},
}

var fiCiColorCode = fieldInfo{
	fType:        FtCiColorCode,
	typeName:     "Color Code",
	max:          1,
	bitOffset:    8,
	bitSize:      4,
	valueType:    VtSpanList,
	defaultValue: "1",
	span: &Span{
		min: 0,
		max: 15,
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiCompressedUdpDataHeader = fieldInfo{
	fType:        FtCiCompressedUdpDataHeader,
	typeName:     "Compressed UDP Data Header",
	max:          1,
	bitOffset:    25,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
	enablerType:  FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiContactName = fieldInfo{
	fType:        FtCiContactName,
	typeName:     "Contact Name",
	max:          1,
	bitOffset:    48,
	bitSize:      16,
	valueType:    VtContactListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "None"},
	},
	listRecordType: RtContacts,
	enablerType:    FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiCtcssDecode = fieldInfo{
	fType:        FtCiCtcssDecode,
	typeName:     "CTCSS/DCS Decode",
	max:          1,
	bitOffset:    192,
	bitSize:      16,
	valueType:    VtCtcssDcs,
	defaultValue: "None",
	enablerType:  FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: false,
		},
	},
}

var fiCiCtcssEncode = fieldInfo{
	fType:        FtCiCtcssEncode,
	typeName:     "CTCSS/DCS Encode",
	max:          1,
	bitOffset:    208,
	bitSize:      16,
	valueType:    VtCtcssDcs,
	defaultValue: "None",
	enables: []FieldType{
		FtCiQtReverse,
		FtCiReverseBurst,
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: false,
		},
	},
}

var fiCiDCDMSwitch = fieldInfo{
	fType:        FtCiDCDMSwitch,
	typeName:     "DCDM Switch",
	max:          1,
	bitOffset:    252,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
	enables: []FieldType{
		FtCiLeaderMS,
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiDQTTurnoffFreq = fieldInfo{
	fType:        FtCiDQTTurnoffFreq,
	typeName:     "Non-QT/DQT Turn-off Freq",
	max:          1,
	bitOffset:    40,
	bitSize:      2,
	valueType:    VtIndexedStrings,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "259.2 Hz"},
		IndexedString{2, "55.2 Hz"},
		IndexedString{3, "None"},
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: false,
		},
	},
}

var fiCiDataCallConfirmed = fieldInfo{
	fType:        FtCiDataCallConfirmed,
	typeName:     "Data Call Confirmed",
	max:          1,
	bitOffset:    16,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiDecode1 = fieldInfo{
	fType:        FtCiDecode1,
	typeName:     "Decode 1",
	max:          1,
	bitOffset:    112,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiRxSignallingSystem,
	enablers: []enabler{
		enabler{
			value:  "Off",
			enable: false,
		},
	},
}

var fiCiDecode2 = fieldInfo{
	fType:        FtCiDecode2,
	typeName:     "Decode 2",
	max:          1,
	bitOffset:    113,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiRxSignallingSystem,
	enablers: []enabler{
		enabler{
			value:  "Off",
			enable: false,
		},
	},
}

var fiCiDecode3 = fieldInfo{
	fType:        FtCiDecode3,
	typeName:     "Decode 3",
	max:          1,
	bitOffset:    114,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiRxSignallingSystem,
	enablers: []enabler{
		enabler{
			value:  "Off",
			enable: false,
		},
	},
}

var fiCiDecode4 = fieldInfo{
	fType:        FtCiDecode4,
	typeName:     "Decode 4",
	max:          1,
	bitOffset:    115,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiRxSignallingSystem,
	enablers: []enabler{
		enabler{
			value:  "Off",
			enable: false,
		},
	},
}

var fiCiDecode5 = fieldInfo{
	fType:        FtCiDecode5,
	typeName:     "Decode 5",
	max:          1,
	bitOffset:    116,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiRxSignallingSystem,
	enablers: []enabler{
		enabler{
			value:  "Off",
			enable: false,
		},
	},
}

var fiCiDecode6 = fieldInfo{
	fType:        FtCiDecode6,
	typeName:     "Decode 6",
	max:          1,
	bitOffset:    117,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiRxSignallingSystem,
	enablers: []enabler{
		enabler{
			value:  "Off",
			enable: false,
		},
	},
}

var fiCiDecode7 = fieldInfo{
	fType:        FtCiDecode7,
	typeName:     "Decode 7",
	max:          1,
	bitOffset:    118,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiRxSignallingSystem,
	enablers: []enabler{
		enabler{
			value:  "Off",
			enable: false,
		},
	},
}

var fiCiDecode8 = fieldInfo{
	fType:        FtCiDecode8,
	typeName:     "Decode 8",
	max:          1,
	bitOffset:    119,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiRxSignallingSystem,
	enablers: []enabler{
		enabler{
			value:  "Off",
			enable: false,
		},
	},
}

var fiCiDisplayPTTID = fieldInfo{
	fType:        FtCiDisplayPTTID,
	typeName:     "Display PTT ID",
	max:          1,
	bitOffset:    24,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
	enablerType:  FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: false,
		},
	},
}

var fiCiEmergencyAlarmAck = fieldInfo{
	fType:        FtCiEmergencyAlarmAck,
	typeName:     "Emergency Alarm Ack",
	max:          1,
	bitOffset:    28,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiEmergencySystem = fieldInfo{
	fType:        FtCiEmergencySystem,
	typeName:     "Emergency System",
	max:          1,
	bitOffset:    80,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "None",
	span: &Span{
		min:       0,
		max:       32,
		minString: "None",
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiGPSSystem = fieldInfo{
	fType:        FtCiGPSSystem,
	typeName:     "GPS System",
	max:          1,
	bitOffset:    104,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "None",
	span: &Span{
		min:       0,
		max:       16,
		minString: "None",
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiGroupList = fieldInfo{
	fType:        FtCiGroupList,
	typeName:     "RX Group List",
	max:          1,
	bitOffset:    96,
	bitSize:      8,
	valueType:    VtListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "None"},
	},
	listRecordType: RtGroupLists,
	enablerType:    FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiInCallCriteria = fieldInfo{
	fType:        FtCiInCallCriteria,
	typeName:     "In Call Criteria",
	max:          1,
	bitOffset:    43,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "Always",
	strings: &[]string{
		"Always",
		"Follow Admit Criteria",
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiLeaderMS = fieldInfo{
	fType:        FtCiLeaderMS,
	typeName:     "Leader/MS",
	max:          1,
	bitOffset:    251,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
	enablerType:  FtCiDCDMSwitch,
	enablers: []enabler{
		enabler{
			value:  "On",
			enable: true,
		},
	},
}

var fiCiLoneWorker = fieldInfo{
	fType:        FtCiLoneWorker,
	typeName:     "Lone Worker",
	max:          1,
	bitOffset:    0,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
}

var fiCiName = fieldInfo{
	fType:        FtCiName,
	typeName:     "Channel Name",
	max:          1,
	bitOffset:    256,
	bitSize:      256,
	valueType:    VtName,
	defaultValue: "Channel1",
}

var fiCiPower = fieldInfo{
	fType:        FtCiPower,
	typeName:     "Power",
	max:          1,
	bitOffset:    34,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "High",
	strings: &[]string{
		"Low",
		"High",
	},
}

var fiCiPower_uv380 = fieldInfo{
	fType:        FtCiPower_uv380,
	typeName:     "Power",
	max:          1,
	bitOffset:    246,
	bitSize:      2,
	valueType:    VtIndexedStrings,
	defaultValue: "High",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Low"},
		IndexedString{2, "Medium"},
		IndexedString{3, "High"},
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
	enables: []FieldType{
		FtCiPrivacyNumber,
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiPrivacyNumber = fieldInfo{
	fType:        FtCiPrivacyNumber,
	typeName:     "Privacy Number",
	max:          1,
	bitOffset:    20,
	bitSize:      4,
	valueType:    VtPrivacyNumber,
	defaultValue: "1",
	strings: &[]string{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"10",
		"11",
		"12",
		"13",
		"15",
		"16",
	},
	enablerType: FtCiPrivacy,
	enablers: []enabler{
		enabler{
			value:  "None",
			enable: false,
		},
	},
}

var fiCiPrivateCallConfirmed = fieldInfo{
	fType:        FtCiPrivateCallConfirmed,
	typeName:     "Private Call Confimed",
	max:          1,
	bitOffset:    17,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
	enablerType:  FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiQtReverse = fieldInfo{
	fType:        FtCiQtReverse,
	typeName:     "QT Reverse",
	max:          1,
	bitOffset:    36,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "180",
	strings: &[]string{
		"180",
		"120",
	},
	enablerType: FtCiCtcssEncode,
	enablers: []enabler{
		enabler{
			value:  "None",
			enable: false,
		},
	},
}

var fiCiReceiveGPSInfo = fieldInfo{
	fType:        FtCiReceiveGPSInfo,
	typeName:     "Receive GPS Info",
	max:          1,
	bitOffset:    254,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
}

var fiCiRepeaterSlot = fieldInfo{
	fType:        FtCiRepeaterSlot,
	typeName:     "Repeater Slot",
	max:          1,
	bitOffset:    12,
	bitSize:      2,
	valueType:    VtIndexedStrings,
	defaultValue: "1",
	indexedStrings: &[]IndexedString{
		IndexedString{1, "1"},
		IndexedString{2, "2"},
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: true,
		},
	},
}

var fiCiReverseBurst = fieldInfo{
	fType:        FtCiReverseBurst,
	typeName:     "Reverse Burst/Turn Off Code",
	max:          1,
	bitOffset:    37,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
	enablerType:  FtCiCtcssEncode,
	enablers: []enabler{
		enabler{
			value:  "None",
			enable: false,
		},
	},
}

var fiCiRxFrequency = fieldInfo{
	fType:        FtCiRxFrequency,
	typeName:     "Rx Frequency (MHz)",
	max:          1,
	bitOffset:    128,
	bitSize:      32,
	valueType:    VtFrequency,
	defaultValue: "0",
}

var fiCiRxOnly = fieldInfo{
	fType:        FtCiRxOnly,
	typeName:     "Rx Only",
	max:          1,
	bitOffset:    14,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
}

var fiCiRxRefFrequency = fieldInfo{
	fType:        FtCiRxRefFrequency,
	typeName:     "Rx Ref Frequency",
	max:          1,
	bitOffset:    30,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "Low",
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
	bitOffset:    224,
	bitSize:      8,
	valueType:    VtIStrings,
	defaultValue: "Off",
	strings: &[]string{
		"Off",
		"DTMF-1",
		"DTMF-2",
		"DTMF-3",
		"DTMF-4",
	},
	enables: []FieldType{
		FtCiDecode1,
		FtCiDecode2,
		FtCiDecode3,
		FtCiDecode4,
		FtCiDecode5,
		FtCiDecode6,
		FtCiDecode7,
		FtCiDecode8,
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: false,
		},
	},
}

var fiCiScanList_md380 = fieldInfo{
	fType:        FtCiScanList_md380,
	typeName:     "Scan List",
	max:          1,
	bitOffset:    88,
	bitSize:      8,
	valueType:    VtListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "None"},
	},
	listRecordType: RtScanLists_md380,
}

var fiCiScanList_md40 = fieldInfo{
	fType:        FtCiScanList_md40,
	typeName:     "Scan List",
	max:          1,
	bitOffset:    88,
	bitSize:      8,
	valueType:    VtListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "None"},
	},
	listRecordType: RtScanLists_md40,
}

var fiCiSendGPSInfo = fieldInfo{
	fType:        FtCiSendGPSInfo,
	typeName:     "Send GPS Info",
	max:          1,
	bitOffset:    255,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
}

var fiCiSquelch = fieldInfo{
	fType:        FtCiSquelch,
	typeName:     "Squelch",
	max:          1,
	bitOffset:    2,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "Normal",
	strings: &[]string{
		"Tight",
		"Normal",
	},
}

var fiCiSquelch_uv380 = fieldInfo{
	fType:        FtCiSquelch_uv380,
	typeName:     "Squelch",
	max:          1,
	bitOffset:    124,
	bitSize:      4,
	valueType:    VtSpanList,
	defaultValue: "1",
	span: &Span{
		min: 0,
		max: 9,
	},
}

var fiCiTalkaround = fieldInfo{
	fType:        FtCiTalkaround,
	typeName:     "Talkaround",
	max:          1,
	bitOffset:    26,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
}

var fiCiTot = fieldInfo{
	fType:        FtCiTot,
	typeName:     "TOT (S)",
	max:          1,
	bitOffset:    64,
	bitSize:      8,
	valueType:    VtSpan,
	defaultValue: "60",
	span: &Span{
		min:       0,
		max:       37,
		scale:     15,
		minString: "Infinite",
	},
}

var fiCiTotRekeyDelay = fieldInfo{
	fType:        FtCiTotRekeyDelay,
	typeName:     "TOT Rekey Delay (S)",
	max:          1,
	bitOffset:    72,
	bitSize:      8,
	valueType:    VtSpan,
	defaultValue: "0",
	span: &Span{
		min: 0,
		max: 255,
	},
}

var fiCiTxFrequencyOffset = fieldInfo{
	fType:        FtCiTxFrequencyOffset,
	typeName:     "Tx Offset (MHz)",
	max:          1,
	bitOffset:    160,
	bitSize:      32,
	valueType:    VtFrequencyOffset,
	defaultValue: "0",
	enables: []FieldType{
		FtCiAllowTalkaround,
	},
}

var fiCiTxRefFrequency = fieldInfo{
	fType:        FtCiTxRefFrequency,
	typeName:     "Tx Ref Frequency",
	max:          1,
	bitOffset:    38,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "Low",
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
	bitOffset:    232,
	bitSize:      8,
	valueType:    VtIStrings,
	defaultValue: "Off",
	strings: &[]string{
		"Off",
		"DTMF-1",
		"DTMF-2",
		"DTMF-3",
		"DTMF-4",
	},
	enablerType: FtCiChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Digital",
			enable: false,
		},
	},
}

var fiCiVox = fieldInfo{
	fType:        FtCiVox,
	typeName:     "VOX",
	max:          1,
	bitOffset:    35,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
}

var fiDcCallID = fieldInfo{
	fType:        FtDcCallID,
	typeName:     "Call ID",
	max:          1,
	bitOffset:    0,
	bitSize:      24,
	valueType:    VtCallID,
	defaultValue: "1",
}

var fiDcCallReceiveTone = fieldInfo{
	fType:        FtDcCallReceiveTone,
	typeName:     "Call Receive Tone",
	max:          1,
	bitOffset:    26,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "No",
	strings: &[]string{
		"No",
		"Yes",
	},
}

var fiDcCallType = fieldInfo{
	fType:        FtDcCallType,
	typeName:     "Call Type",
	max:          1,
	bitOffset:    27,
	bitSize:      5,
	valueType:    VtCallType,
	defaultValue: "Group",
	indexedStrings: &[]IndexedString{
		IndexedString{1, "Group"},
		IndexedString{2, "Private"},
		IndexedString{3, "All"},
	},
}

var fiDcName = fieldInfo{
	fType:        FtDcName,
	typeName:     "Contact Name",
	max:          1,
	bitOffset:    32,
	bitSize:      256,
	valueType:    VtContactName,
	defaultValue: "Contact1",
}

var fiGlContact = fieldInfo{
	fType:          FtGlContact,
	typeName:       "Contacts",
	max:            32,
	bitOffset:      256,
	bitSize:        16,
	valueType:      VtContactListIndex,
	listRecordType: RtContacts,
}

var fiGlName = fieldInfo{
	fType:        FtGlName,
	typeName:     "RX Group List Name",
	max:          1,
	bitOffset:    0,
	bitSize:      256,
	valueType:    VtName,
	defaultValue: "GroupList1",
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
	fType:        FtGpGPSDefaultReportInterval,
	typeName:     "GPS Default Report Interval (S)",
	max:          1,
	bitOffset:    16,
	bitSize:      8,
	valueType:    VtGpsReportInterval,
	defaultValue: "Off",
	span: &Span{
		min:       0,
		max:       240,
		scale:     30,
		minString: "Off",
	},
}

var fiGpGPSRevertChannel = fieldInfo{
	fType:        FtGpGPSRevertChannel,
	typeName:     "GPS Revert Channel",
	max:          1,
	bitOffset:    0,
	bitSize:      16,
	valueType:    VtGpsListIndex,
	defaultValue: "Current Channel",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Current Channel"},
	},
	listRecordType: RtChannels_md380,
}

var fiGsBacklightColor = fieldInfo{
	fType:        FtGsBacklightColor,
	typeName:     "Backlight Color",
	max:          1,
	bitOffset:    542,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "White",
	strings: &[]string{
		"Off",
		"Orange",
		"White",
		"Sakura",
	},
}

var fiGsBacklightTime = fieldInfo{
	fType:        FtGsBacklightTime,
	typeName:     "Backlight Time (S)",
	max:          1,
	bitOffset:    686,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "Always",
	strings: &[]string{
		"Always",
		"5",
		"10",
		"15",
	},
}

var fiGsCHVoiceAnnouncement = fieldInfo{
	fType:        FtGsCHVoiceAnnouncement,
	typeName:     "CH Voice Announcement",
	max:          1,
	bitOffset:    534,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiGsCallAlertToneDuration = fieldInfo{
	fType:        FtGsCallAlertToneDuration,
	typeName:     "Call Alert Tone Duration (S)",
	max:          1,
	bitOffset:    632,
	bitSize:      8,
	valueType:    VtSpan,
	defaultValue: "Continue",
	span: &Span{
		min:       0,
		max:       240,
		scale:     5,
		minString: "Continue",
	},
}

var fiGsChFreeIndicationTone = fieldInfo{
	fType:        FtGsChFreeIndicationTone,
	typeName:     "Channel Free Indication Tone",
	max:          1,
	bitOffset:    523,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
}

var fiGsChannelsHangTime = fieldInfo{
	fType:        FtGsChannelsHangTime,
	typeName:     "Channels Hang Time (mS)",
	max:          1,
	bitOffset:    1152,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "3000",
	span: &Span{
		min:      0,
		max:      70,
		scale:    100,
		interval: 5,
	},
}

var fiGsDisableAllLeds = fieldInfo{
	fType:        FtGsDisableAllLeds,
	typeName:     "Disable All LEDS",
	max:          1,
	bitOffset:    517,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
}

var fiGsDisableAllTones = fieldInfo{
	fType:        FtGsDisableAllTones,
	typeName:     "Disable All Tones",
	max:          1,
	bitOffset:    525,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
}

var fiGsEditRadioID = fieldInfo{
	fType:        FtGsEditRadioID,
	typeName:     "Edit Radio ID",
	max:          1,
	bitOffset:    1281,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
}

var fiGsEnableContactsCSV = fieldInfo{
	fType:        FtGsEnableContactsCSV,
	typeName:     "Enable Contacts CSV",
	max:          1,
	bitOffset:    529,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
}

var fiGsFreqChannelMode = fieldInfo{
	fType:        FtGsFreqChannelMode,
	typeName:     "Freq/Channel Mode",
	max:          1,
	bitOffset:    540,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "Channel",
	strings: &[]string{
		"Frequency",
		"Channel",
	},
	enables: []FieldType{
		FtGsLockUnlock,
		FtGsModeSelect,
	},
}

var fiGsFreqChannelMode_uv380 = fieldInfo{
	fType:        FtGsFreqChannelMode_uv380,
	typeName:     "Freq/Channel Mode",
	max:          1,
	bitOffset:    696,
	bitSize:      8,
	valueType:    VtIndexedStrings,
	defaultValue: "Channel",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Frequency"},
		IndexedString{255, "Channel"},
	},
	enables: []FieldType{
		FtGsLockUnlock,
		FtGsModeSelectA,
		FtGsModeSelectB,
	},
}

var fiGsGroupCallHangTime = fieldInfo{
	fType:        FtGsGroupCallHangTime,
	typeName:     "Group Call Hang Time (mS)",
	max:          1,
	bitOffset:    584,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "3000",
	span: &Span{
		min:      0,
		max:      70,
		scale:    100,
		interval: 5,
	},
}

var fiGsGroupCallMatch = fieldInfo{
	fType:        FtGsGroupCallMatch,
	typeName:     "Group Call Match",
	max:          1,
	bitOffset:    863,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiGsIntroScreen = fieldInfo{
	fType:        FtGsIntroScreen,
	typeName:     "Intro Screen",
	max:          1,
	bitOffset:    531,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "Character String",
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

var fiGsKeypadTones = fieldInfo{
	fType:        FtGsKeypadTones,
	typeName:     "Keypad Tones",
	max:          1,
	bitOffset:    530,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiGsLockUnlock = fieldInfo{
	fType:        FtGsLockUnlock,
	typeName:     "Lock/Unlock",
	max:          1,
	bitOffset:    539,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "Unlock",
	strings: &[]string{
		"Unlock",
		"Lock",
	},
	enablerType: FtGsFreqChannelMode_uv380,
	enablers: []enabler{
		enabler{
			value:  "Frequency",
			enable: false,
		},
		enabler{
			value:  "VFO",
			enable: false,
		},
	},
}

var fiGsLoneWorkerReminderTime = fieldInfo{
	fType:        FtGsLoneWorkerReminderTime,
	typeName:     "Lone Worker Reminder Time (S)",
	max:          1,
	bitOffset:    648,
	bitSize:      8,
	valueType:    VtSpan,
	defaultValue: "10",
	span: &Span{
		min: 1,
		max: 255,
	},
}

var fiGsLoneWorkerResponseTime = fieldInfo{
	fType:        FtGsLoneWorkerResponseTime,
	typeName:     "Lone Worker Response Time (min)",
	max:          1,
	bitOffset:    640,
	bitSize:      8,
	valueType:    VtSpan,
	defaultValue: "1",
	span: &Span{
		min: 1,
		max: 255,
	},
}

var fiGsMenuControl = fieldInfo{
	fType:     FtGsMenuControl,
	typeName:  "Menu Control",
	max:       1,
	bitOffset: 1303,
	bitSize:   1,
	valueType: VtOffOn,
}

var fiGsMicLevel = fieldInfo{
	fType:        FtGsMicLevel,
	typeName:     "MIC Level",
	max:          1,
	bitOffset:    1282,
	bitSize:      3,
	valueType:    VtIStrings,
	defaultValue: "3",
	strings: &[]string{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
	},
}

var fiGsModeSelect = fieldInfo{
	fType:        FtGsModeSelect,
	typeName:     "Mode Select",
	max:          1,
	bitOffset:    541,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "Memory",
	strings: &[]string{
		"VFO",
		"Memory",
	},
	enablerType: FtGsFreqChannelMode,
	enablers: []enabler{
		enabler{
			value:  "Frequency",
			enable: true,
		},
	},
}

var fiGsModeSelectA = fieldInfo{
	fType:        FtGsModeSelectA,
	typeName:     "Mode Select A",
	max:          1,
	bitOffset:    541,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "Memory",
	strings: &[]string{
		"VFO",
		"Memory",
	},
	enablerType: FtGsFreqChannelMode_uv380,
	enablers: []enabler{
		enabler{
			value:  "VFO",
			enable: true,
		},
	},
}

var fiGsModeSelectB = fieldInfo{
	fType:        FtGsModeSelectB,
	typeName:     "Mode Select B",
	max:          1,
	bitOffset:    536,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "Memory",
	strings: &[]string{
		"VFO",
		"Memory",
	},
	enablerType: FtGsFreqChannelMode_uv380,
	enablers: []enabler{
		enabler{
			value:  "VFO",
			enable: true,
		},
	},
}

var fiGsMonitorType = fieldInfo{
	fType:        FtGsMonitorType,
	typeName:     "Monitor Type",
	max:          1,
	bitOffset:    515,
	bitSize:      1,
	valueType:    VtIStrings,
	defaultValue: "Open Squelch",
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
	enablerType:  FtGsPwAndLockEnable,
	enablers: []enabler{
		enabler{
			value:  "On",
			enable: true,
		},
	},
}

var fiGsPrivateCallHangTime = fieldInfo{
	fType:        FtGsPrivateCallHangTime,
	typeName:     "Private Call Hang Time (mS)",
	max:          1,
	bitOffset:    592,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "4000",
	span: &Span{
		min:      0,
		max:      70,
		scale:    100,
		interval: 5,
	},
}

var fiGsPrivateCallMatch = fieldInfo{
	fType:        FtGsPrivateCallMatch,
	typeName:     "Private Call Match",
	max:          1,
	bitOffset:    862,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiGsPublicZone = fieldInfo{
	fType:        FtGsPublicZone,
	typeName:     "Public Zone",
	max:          1,
	bitOffset:    1173,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiGsPwAndLockEnable = fieldInfo{
	fType:        FtGsPwAndLockEnable,
	typeName:     "Password And Lock Enable",
	max:          1,
	bitOffset:    522,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
	enables: []FieldType{
		FtGsPowerOnPassword,
	},
}

var fiGsRadioID = fieldInfo{
	fType:        FtGsRadioID,
	typeName:     "Radio ID",
	max:          1,
	bitOffset:    544,
	bitSize:      24,
	valueType:    VtCallID,
	defaultValue: "1234",
}

var fiGsRadioID1 = fieldInfo{
	fType:        FtGsRadioID1,
	typeName:     "Radio ID 1",
	max:          1,
	bitOffset:    1184,
	bitSize:      24,
	valueType:    VtCallID,
	defaultValue: "1",
}

var fiGsRadioID2 = fieldInfo{
	fType:        FtGsRadioID2,
	typeName:     "Radio ID 2",
	max:          1,
	bitOffset:    1216,
	bitSize:      24,
	valueType:    VtCallID,
	defaultValue: "2",
}

var fiGsRadioID3 = fieldInfo{
	fType:        FtGsRadioID3,
	typeName:     "Radio ID 3",
	max:          1,
	bitOffset:    1248,
	bitSize:      24,
	valueType:    VtCallID,
	defaultValue: "3",
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
	valueType: VtRadioProgPassword,
}

var fiGsRxLowBatteryInterval = fieldInfo{
	fType:        FtGsRxLowBatteryInterval,
	typeName:     "Rx Low Battery Interval (S)",
	max:          1,
	bitOffset:    624,
	bitSize:      8,
	valueType:    VtSpan,
	defaultValue: "3",
	span: &Span{
		min:   0,
		max:   127,
		scale: 5,
	},
}

var fiGsSaveModeReceive = fieldInfo{
	fType:        FtGsSaveModeReceive,
	typeName:     "Save Mode Receive",
	max:          1,
	bitOffset:    526,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiGsSavePreamble = fieldInfo{
	fType:        FtGsSavePreamble,
	typeName:     "Save Preamble",
	max:          1,
	bitOffset:    527,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiGsScanAnalogHangTime = fieldInfo{
	fType:        FtGsScanAnalogHangTime,
	typeName:     "Scan Analog Hang Time (mS)",
	max:          1,
	bitOffset:    672,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "1000",
	span: &Span{
		min:      5,
		max:      100,
		scale:    100,
		interval: 5,
	},
}

var fiGsScanDigitalHangTime = fieldInfo{
	fType:        FtGsScanDigitalHangTime,
	typeName:     "Scan Digital Hang Time (mS)",
	max:          1,
	bitOffset:    664,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "1000",
	span: &Span{
		min:      5,
		max:      100,
		scale:    100,
		interval: 5,
	},
}

var fiGsSetKeypadLockTime = fieldInfo{
	fType:        FtGsSetKeypadLockTime,
	typeName:     "Set Keypad Lock Time (S)",
	max:          1,
	bitOffset:    688,
	bitSize:      8,
	valueType:    VtIndexedStrings,
	defaultValue: "Manual",
	indexedStrings: &[]IndexedString{
		IndexedString{5, "5"},
		IndexedString{10, "10"},
		IndexedString{15, "15"},
		IndexedString{255, "Manual"},
	},
}

var fiGsTalkPermitTone = fieldInfo{
	fType:        FtGsTalkPermitTone,
	typeName:     "Talk Permit Tone",
	max:          1,
	bitOffset:    520,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "None",
	strings: &[]string{
		"None",
		"Digital",
		"Analog",
		"Digital and Analog",
	},
}

var fiGsTimeZone = fieldInfo{
	fType:        FtGsTimeZone,
	typeName:     "Time Zone",
	max:          1,
	bitOffset:    856,
	bitSize:      5,
	valueType:    VtIStrings,
	defaultValue: "UTC+8:00",
	strings: &[]string{
		"UTC-12:00",
		"UTC-11:00",
		"UTC-10:00",
		"UTC-9:00",
		"UTC-8:00",
		"UTC-7:00",
		"UTC-6:00",
		"UTC-5:00",
		"UTC-4:00",
		"UTC-3:00",
		"UTC-2:00",
		"UTC-1:00",
		"UTC+0:00",
		"UTC+1:00",
		"UTC+2:00",
		"UTC+3:00",
		"UTC+4:00",
		"UTC+5:00",
		"UTC+6:00",
		"UTC+7:00",
		"UTC+8:00",
		"UTC+9:00",
		"UTC+10:00",
		"UTC+11:00",
		"UTC+12:00",
	},
}

var fiGsTwoChannel = fieldInfo{
	fType:     FtGsTwoChannel,
	typeName:  "Two Channel",
	max:       1,
	bitOffset: 1302,
	bitSize:   1,
	valueType: VtOnOff,
}

var fiGsTxMode = fieldInfo{
	fType:        FtGsTxMode,
	typeName:     "Tx Mode",
	max:          1,
	bitOffset:    512,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "Designated CH + Hand CH",
	strings: &[]string{
		"Last Call CH",
		"Last Call + Hand CH",
		"Designated CH",
		"Designated CH + Hand CH",
	},
}

var fiGsTxPreambleDuration = fieldInfo{
	fType:        FtGsTxPreambleDuration,
	typeName:     "Tx Preamble Duration (mS)",
	max:          1,
	bitOffset:    576,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "600",
	span: &Span{
		min:   0,
		max:   144,
		scale: 60,
	},
}

var fiGsVoxSensitivity = fieldInfo{
	fType:        FtGsVoxSensitivity,
	typeName:     "VOX Sensitivity",
	max:          1,
	bitOffset:    600,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "3",
	span: &Span{
		min: 1,
		max: 10,
	},
}

var fiMiAnswered = fieldInfo{
	fType:        FtMiAnswered,
	typeName:     "Answered",
	max:          1,
	bitOffset:    19,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiBacklight = fieldInfo{
	fType:        FtMiBacklight,
	typeName:     "Backlight",
	max:          1,
	bitOffset:    30,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiCallAlert = fieldInfo{
	fType:        FtMiCallAlert,
	typeName:     "Call Alert",
	max:          1,
	bitOffset:    14,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiDisplayMode = fieldInfo{
	fType:        FtMiDisplayMode,
	typeName:     "Display Mode",
	max:          1,
	bitOffset:    38,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiEdit = fieldInfo{
	fType:        FtMiEdit,
	typeName:     "Edit",
	max:          1,
	bitOffset:    13,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiEditList = fieldInfo{
	fType:        FtMiEditList,
	typeName:     "Edit List",
	max:          1,
	bitOffset:    21,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiGps = fieldInfo{
	fType:        FtMiGps,
	typeName:     "GPS",
	max:          1,
	bitOffset:    36,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
}

var fiMiHangTime = fieldInfo{
	fType:        FtMiHangTime,
	typeName:     "Hang Time",
	max:          1,
	bitOffset:    0,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "10",
	span: &Span{
		min:       0,
		max:       30,
		minString: "Hang",
	},
}

var fiMiIntroScreen = fieldInfo{
	fType:        FtMiIntroScreen,
	typeName:     "Intro Screen",
	max:          1,
	bitOffset:    29,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiKeyboardLock = fieldInfo{
	fType:        FtMiKeyboardLock,
	typeName:     "Keyboard Lock",
	max:          1,
	bitOffset:    28,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiLedIndicator = fieldInfo{
	fType:        FtMiLedIndicator,
	typeName:     "LED Indicator",
	max:          1,
	bitOffset:    27,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiManualDial = fieldInfo{
	fType:        FtMiManualDial,
	typeName:     "Manual Dial",
	max:          1,
	bitOffset:    12,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiMissed = fieldInfo{
	fType:        FtMiMissed,
	typeName:     "Missed",
	max:          1,
	bitOffset:    20,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiOutgoingRadio = fieldInfo{
	fType:        FtMiOutgoingRadio,
	typeName:     "Outgoing Radio",
	max:          1,
	bitOffset:    18,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiPasswordAndLock = fieldInfo{
	fType:        FtMiPasswordAndLock,
	typeName:     "Password And Lock",
	max:          1,
	bitOffset:    39,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiPower = fieldInfo{
	fType:        FtMiPower,
	typeName:     "Power",
	max:          1,
	bitOffset:    31,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiProgramKey = fieldInfo{
	fType:        FtMiProgramKey,
	typeName:     "Program Key",
	max:          1,
	bitOffset:    23,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiProgramRadio = fieldInfo{
	fType:        FtMiProgramRadio,
	typeName:     "Program Radio",
	max:          1,
	bitOffset:    37,
	bitSize:      1,
	valueType:    VtOnOff,
	defaultValue: "Off",
}

var fiMiRadioCheck = fieldInfo{
	fType:        FtMiRadioCheck,
	typeName:     "Radio Check",
	max:          1,
	bitOffset:    11,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
}

var fiMiRadioDisable = fieldInfo{
	fType:        FtMiRadioDisable,
	typeName:     "Radio Disable",
	max:          1,
	bitOffset:    8,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
}

var fiMiRadioEnable = fieldInfo{
	fType:        FtMiRadioEnable,
	typeName:     "Radio Enable",
	max:          1,
	bitOffset:    9,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
}

var fiMiRemoteMonitor = fieldInfo{
	fType:        FtMiRemoteMonitor,
	typeName:     "Remote Monitor",
	max:          1,
	bitOffset:    10,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
}

var fiMiScan = fieldInfo{
	fType:        FtMiScan,
	typeName:     "Scan",
	max:          1,
	bitOffset:    22,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiSquelch = fieldInfo{
	fType:        FtMiSquelch,
	typeName:     "Squelch",
	max:          1,
	bitOffset:    26,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiTalkaround = fieldInfo{
	fType:        FtMiTalkaround,
	typeName:     "Talkaround",
	max:          1,
	bitOffset:    17,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiTextMessage = fieldInfo{
	fType:        FtMiTextMessage,
	typeName:     "Text Messsage",
	max:          1,
	bitOffset:    15,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiToneOrAlert = fieldInfo{
	fType:        FtMiToneOrAlert,
	typeName:     "Tone Or Alert",
	max:          1,
	bitOffset:    16,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "On",
}

var fiMiVox = fieldInfo{
	fType:        FtMiVox,
	typeName:     "VOX",
	max:          1,
	bitOffset:    24,
	bitSize:      1,
	valueType:    VtOffOn,
	defaultValue: "Off",
}

var fiNkContact = fieldInfo{
	fType:     FtNkContact,
	typeName:  "Contact",
	max:       1,
	bitOffset: 0,
	bitSize:   16,
	valueType: VtNkContactListIndex,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "None"},
	},
	listRecordType: RtContacts,
}

var fiOtCall = fieldInfo{
	fType:          FtOtCall,
	typeName:       "Call",
	max:            1,
	bitOffset:      16,
	bitSize:        16,
	valueType:      VtContactListIndex,
	listRecordType: RtContacts,
	enablerType:    FtOtMode,
	enablers: []enabler{
		enabler{
			value:  "Analog",
			enable: false,
		},
		enabler{
			value:  "Digital",
			enable: true,
		},
		enabler{
			value:  "None",
			enable: false,
		},
	},
}

var fiOtCallType = fieldInfo{
	fType:        FtOtCallType,
	typeName:     "Call Type",
	max:          1,
	bitOffset:    6,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "Call",
	strings: &[]string{
		"Call",
		"Text Message",
	},
	enables: []FieldType{
		FtOtTextMessage,
	},
	enablerType: FtOtMode,
	enablers: []enabler{
		enabler{
			value:  "Analog",
			enable: false,
		},
		enabler{
			value:  "Digital",
			enable: true,
		},
		enabler{
			value:  "None",
			enable: false,
		},
	},
}

var fiOtDTMF = fieldInfo{
	fType:        FtOtDTMF,
	typeName:     "DTMF",
	max:          1,
	bitOffset:    6,
	bitSize:      2,
	valueType:    VtIStrings,
	defaultValue: "DTMF-1",
	strings: &[]string{
		"DTMF-1",
		"DTMF-2",
		"DTMF-3",
		"DTMF-4",
	},
	enablerType: FtOtMode,
	enablers: []enabler{
		enabler{
			value:  "Analog",
			enable: true,
		},
		enabler{
			value:  "Digital",
			enable: false,
		},
		enabler{
			value:  "None",
			enable: false,
		},
	},
}

var fiOtEncode = fieldInfo{
	fType:     FtOtEncode,
	typeName:  "Encode",
	max:       1,
	bitOffset: 8,
	bitSize:   8,
	valueType: VtIStrings,
	strings: &[]string{
		"To be implemented",
	},
	enablerType: FtOtMode,
	enablers: []enabler{
		enabler{
			value:  "Analog",
			enable: true,
		},
		enabler{
			value:  "Digital",
			enable: false,
		},
		enabler{
			value:  "None",
			enable: false,
		},
	},
}

var fiOtMode = fieldInfo{
	fType:        FtOtMode,
	typeName:     "Mode",
	max:          1,
	bitOffset:    0,
	bitSize:      6,
	valueType:    VtIndexedStrings,
	defaultValue: "Digital",
	indexedStrings: &[]IndexedString{
		IndexedString{48, "None"},
		IndexedString{52, "Digital"},
		IndexedString{58, "Analog"},
	},
	enables: []FieldType{
		FtOtCall,
		FtOtCallType,
		FtOtDTMF,
		FtOtEncode,
	},
}

var fiOtTextMessage = fieldInfo{
	fType:          FtOtTextMessage,
	typeName:       "Text Message",
	max:            1,
	bitOffset:      8,
	bitSize:        8,
	valueType:      VtDerefListIndex,
	listRecordType: RtTextMessages,
	enablerType:    FtOtCallType,
	enablers: []enabler{
		enabler{
			value:  "Text Message",
			enable: true,
		},
	},
}

var fiPsBasicKey = fieldInfo{
	fType:        FtPsBasicKey,
	typeName:     "Key Value (Basic)",
	max:          16,
	bitOffset:    1152,
	bitSize:      16,
	valueType:    VtHexadecimal4,
	defaultValue: "ffff",
}

var fiPsEnhancedKey = fieldInfo{
	fType:        FtPsEnhancedKey,
	typeName:     "Key Value (Enhanced)",
	max:          8,
	bitOffset:    0,
	bitSize:      128,
	valueType:    VtHexadecimal32,
	defaultValue: "ffffffffffffffffffffffffffffffff",
}

var fiRbButton = fieldInfo{
	fType:     FtRbButton,
	typeName:  "Button",
	max:       1,
	bitOffset: 0,
	bitSize:   8,
	valueType: VtRadioButton,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Unassigned (default)"},
		IndexedString{1, "All alert Tones On/Off"},
		IndexedString{2, "Emergency On"},
		IndexedString{3, "Emergency Off"},
		IndexedString{4, "High/Low Power"},
		IndexedString{5, "Monitor"},
		IndexedString{6, "Nuisance Delete"},
		IndexedString{7, "One Touch Access 1"},
		IndexedString{8, "One Touch Access 2"},
		IndexedString{9, "One Touch Access 3"},
		IndexedString{10, "One Touch Access 4"},
		IndexedString{11, "One Touch Access 5"},
		IndexedString{12, "One Touch Access 6"},
		IndexedString{13, "Repeater/Talkaround"},
		IndexedString{14, "Scan On/Off"},
		IndexedString{21, "Squelch Tight/Normal"},
		IndexedString{22, "Privacy On/Off"},
		IndexedString{23, "VOX On/Off"},
		IndexedString{24, "Zone +"},
		IndexedString{25, "Zone Toggle"},
		IndexedString{26, "Battery Indicator"},
		IndexedString{30, "Manual Dial For Private"},
		IndexedString{31, "Lone Work On/Off"},
		IndexedString{38, "1750 Hz"},
		IndexedString{80, "Toggle Backlight (md380tools)"},
		IndexedString{81, "Set Talkgroup (md380tools)"},
		IndexedString{82, "Morse Narrator (md380tools)"},
		IndexedString{83, "Morse Repeat (md380tools)"},
		IndexedString{84, "Screen Toggle (md380tools)"},
		IndexedString{85, "Mic Gain 0/3/6 dB (md380tools)"},
		IndexedString{86, "Promiscuous Mode On/Off (md380tools)"},
	},
}

var fiRbButton_md40 = fieldInfo{
	fType:     FtRbButton_md40,
	typeName:  "Button",
	max:       1,
	bitOffset: 0,
	bitSize:   8,
	valueType: VtRadioButton,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Unassigned (default)"},
		IndexedString{1, "All alert Tones On/Off"},
		IndexedString{2, "Emergency On"},
		IndexedString{3, "Emergency Off"},
		IndexedString{4, "Power Select"},
		IndexedString{5, "Monitor"},
		IndexedString{7, "One Touch Access 1"},
		IndexedString{8, "One Touch Access 2"},
		IndexedString{9, "One Touch Access 3"},
		IndexedString{10, "One Touch Access 4"},
		IndexedString{11, "One Touch Access 5"},
		IndexedString{12, "One Touch Access 6"},
		IndexedString{13, "Repeater/Talkaround"},
		IndexedString{14, "Scan On/Off"},
		IndexedString{22, "Privacy On/Off"},
		IndexedString{23, "VOX On/Off"},
		IndexedString{30, "Manual Dial For Private"},
		IndexedString{31, "Lone Work On/Off"},
		IndexedString{34, "Record On/Off (Firmware)"},
		IndexedString{35, "Record Playback (Firmware)"},
		IndexedString{36, "Delete All Recorded (Firmware)"},
	},
}

var fiRbButton_uv380 = fieldInfo{
	fType:     FtRbButton_uv380,
	typeName:  "Button",
	max:       1,
	bitOffset: 0,
	bitSize:   8,
	valueType: VtRadioButton,
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Unassigned (default)"},
		IndexedString{1, "All alert Tones On/Off"},
		IndexedString{2, "Emergency On"},
		IndexedString{3, "Emergency Off"},
		IndexedString{4, "High/Low Power"},
		IndexedString{5, "Monitor"},
		IndexedString{6, "Nuisance Delete"},
		IndexedString{7, "One Touch Access 1"},
		IndexedString{8, "One Touch Access 2"},
		IndexedString{9, "One Touch Access 3"},
		IndexedString{10, "One Touch Access 4"},
		IndexedString{11, "One Touch Access 5"},
		IndexedString{12, "One Touch Access 6"},
		IndexedString{13, "Repeater/Talkaround"},
		IndexedString{14, "Scan On/Off"},
		IndexedString{21, "Squelch Tight/Normal"},
		IndexedString{22, "Privacy On/Off"},
		IndexedString{23, "VOX On/Off"},
		IndexedString{24, "Zone +"},
		IndexedString{25, "Zone Toggle"},
		IndexedString{26, "Battery Indicator"},
		IndexedString{30, "Manual Dial For Private"},
		IndexedString{31, "Lone Work On/Off"},
		IndexedString{34, "Record On/Off (Firmware)"},
		IndexedString{35, "Record Playback (Firmware)"},
		IndexedString{36, "Delete All Recorded (Firmware)"},
		IndexedString{38, "1750 Hz"},
		IndexedString{47, "Toggle Up/Down"},
		IndexedString{48, "Right Key"},
		IndexedString{49, "Left Key"},
		IndexedString{55, "Zone -"},
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
	fType:        FtSlName,
	typeName:     "Scan List Name",
	max:          1,
	bitOffset:    0,
	bitSize:      256,
	valueType:    VtName,
	defaultValue: "ScanList1",
}

var fiSlPriorityChannel1_md380 = fieldInfo{
	fType:        FtSlPriorityChannel1_md380,
	typeName:     "Priority Channel 1",
	max:          1,
	bitOffset:    256,
	bitSize:      16,
	valueType:    VtMemberListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Selected"},
		IndexedString{65535, "None"},
	},
	listRecordType: RtChannels_md380,
	enables: []FieldType{
		FtSlPriorityChannel2_md380,
	},
}

var fiSlPriorityChannel1_md40 = fieldInfo{
	fType:        FtSlPriorityChannel1_md40,
	typeName:     "Priority Channel 1",
	max:          1,
	bitOffset:    256,
	bitSize:      16,
	valueType:    VtMemberListIndex,
	defaultValue: "None",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Selected"},
		IndexedString{65535, "None"},
	},
	listRecordType: RtChannels_md40,
	enables: []FieldType{
		FtSlPriorityChannel2_md40,
	},
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
	enablerType:    FtSlPriorityChannel1_md380,
	enablers: []enabler{
		enabler{
			value:  "None",
			enable: false,
		},
	},
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
	enablerType:    FtSlPriorityChannel1_md40,
	enablers: []enabler{
		enabler{
			value:  "None",
			enable: false,
		},
	},
}

var fiSlPrioritySampleTime = fieldInfo{
	fType:        FtSlPrioritySampleTime,
	typeName:     "Priority Sample Time (mS)",
	max:          1,
	bitOffset:    320,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "2000",
	span: &Span{
		min:   3,
		max:   31,
		scale: 250,
	},
}

var fiSlSignallingHoldTime = fieldInfo{
	fType:        FtSlSignallingHoldTime,
	typeName:     "Signalling Hold Time (mS)",
	max:          1,
	bitOffset:    312,
	bitSize:      8,
	valueType:    VtSpanList,
	defaultValue: "500",
	span: &Span{
		min:   2,
		max:   255,
		scale: 25,
	},
}

var fiSlTxDesignatedChannel_md380 = fieldInfo{
	fType:        FtSlTxDesignatedChannel_md380,
	typeName:     "Tx Designated Channel",
	max:          1,
	bitOffset:    288,
	bitSize:      16,
	valueType:    VtListIndex,
	defaultValue: "Last Active Channel",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Selected"},
		IndexedString{65535, "Last Active Channel"},
	},
	listRecordType: RtChannels_md380,
}

var fiSlTxDesignatedChannel_md40 = fieldInfo{
	fType:        FtSlTxDesignatedChannel_md40,
	typeName:     "Tx Designated Channel",
	max:          1,
	bitOffset:    288,
	bitSize:      16,
	valueType:    VtListIndex,
	defaultValue: "Last Active Channel",
	indexedStrings: &[]IndexedString{
		IndexedString{0, "Selected"},
		IndexedString{65535, "Last Active Channel"},
	},
	listRecordType: RtChannels_md40,
}

var fiTmTextMessage = fieldInfo{
	fType:        FtTmTextMessage,
	typeName:     "Message",
	max:          1,
	bitOffset:    0,
	bitSize:      2304,
	valueType:    VtTextMessage,
	defaultValue: "Hello",
}

var fiZiChannelA_uv380 = fieldInfo{
	fType:          FtZiChannelA_uv380,
	typeName:       "A Channels",
	max:            64,
	bitOffset:      256,
	bitSize:        16,
	valueType:      VtListIndex,
	listRecordType: RtChannels_uv380,
	extOffset:      201253,
	extSize:        224,
	extIndex:       16,
	extBitOffset:   0,
}

var fiZiChannelB_uv380 = fieldInfo{
	fType:          FtZiChannelB_uv380,
	typeName:       "B Channels",
	max:            64,
	bitOffset:      256,
	bitSize:        16,
	valueType:      VtListIndex,
	listRecordType: RtChannels_uv380,
	extOffset:      201253,
	extSize:        224,
	extIndex:       0,
	extBitOffset:   768,
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
	extBitOffset:   0,
}

var fiZiName = fieldInfo{
	fType:        FtZiName,
	typeName:     "Zone Name",
	max:          1,
	bitOffset:    0,
	bitSize:      256,
	valueType:    VtName,
	defaultValue: "Zone1",
}

var rTypeFieldRefs = map[RecordType][]fieldRef{
	"Channels": []fieldRef{
		fieldRef{
			rType: RecordType("Zones"),
			fType: FieldType("Channel"),
		},
		fieldRef{
			rType: RecordType("Zones"),
			fType: FieldType("ChannelA"),
		},
		fieldRef{
			rType: RecordType("Zones"),
			fType: FieldType("ChannelB"),
		},
		fieldRef{
			rType: RecordType("ScanLists"),
			fType: FieldType("PriorityChannel1"),
		},
		fieldRef{
			rType: RecordType("ScanLists"),
			fType: FieldType("PriorityChannel2"),
		},
		fieldRef{
			rType: RecordType("ScanLists"),
			fType: FieldType("TxDesignatedChannel"),
		},
		fieldRef{
			rType: RecordType("ScanLists"),
			fType: FieldType("Channel"),
		},
		fieldRef{
			rType: RecordType("GPSSystems"),
			fType: FieldType("GPSRevertChannel"),
		},
	},
	"Contacts": []fieldRef{
		fieldRef{
			rType: RecordType("OneTouch"),
			fType: FieldType("Call"),
		},
		fieldRef{
			rType: RecordType("NumberKey"),
			fType: FieldType("Contact"),
		},
		fieldRef{
			rType: RecordType("GroupLists"),
			fType: FieldType("Contact"),
		},
		fieldRef{
			rType: RecordType("Channels"),
			fType: FieldType("ContactName"),
		},
		fieldRef{
			rType: RecordType("GPSSystems"),
			fType: FieldType("DestinationID"),
		},
	},
	"GroupLists": []fieldRef{
		fieldRef{
			rType: RecordType("Channels"),
			fType: FieldType("GroupList"),
		},
	},
	"ScanLists": []fieldRef{
		fieldRef{
			rType: RecordType("Channels"),
			fType: FieldType("ScanList"),
		},
	},
	"TextMessages": []fieldRef{
		fieldRef{
			rType: RecordType("OneTouch"),
			fType: FieldType("TextMessage"),
		},
	},
}

//go:generate genCodeplugInfo
