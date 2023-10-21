package encodings

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AacAudioProfile string

const (
	AacAudioProfileAacLc     AacAudioProfile = "AacLc"
	AacAudioProfileHeAacVOne AacAudioProfile = "HeAacV1"
	AacAudioProfileHeAacVTwo AacAudioProfile = "HeAacV2"
)

func PossibleValuesForAacAudioProfile() []string {
	return []string{
		string(AacAudioProfileAacLc),
		string(AacAudioProfileHeAacVOne),
		string(AacAudioProfileHeAacVTwo),
	}
}

func parseAacAudioProfile(input string) (*AacAudioProfile, error) {
	vals := map[string]AacAudioProfile{
		"aaclc":   AacAudioProfileAacLc,
		"heaacv1": AacAudioProfileHeAacVOne,
		"heaacv2": AacAudioProfileHeAacVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AacAudioProfile(input)
	return &out, nil
}

type AnalysisResolution string

const (
	AnalysisResolutionSourceResolution   AnalysisResolution = "SourceResolution"
	AnalysisResolutionStandardDefinition AnalysisResolution = "StandardDefinition"
)

func PossibleValuesForAnalysisResolution() []string {
	return []string{
		string(AnalysisResolutionSourceResolution),
		string(AnalysisResolutionStandardDefinition),
	}
}

func parseAnalysisResolution(input string) (*AnalysisResolution, error) {
	vals := map[string]AnalysisResolution{
		"sourceresolution":   AnalysisResolutionSourceResolution,
		"standarddefinition": AnalysisResolutionStandardDefinition,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AnalysisResolution(input)
	return &out, nil
}

type AttributeFilter string

const (
	AttributeFilterAll         AttributeFilter = "All"
	AttributeFilterBottom      AttributeFilter = "Bottom"
	AttributeFilterTop         AttributeFilter = "Top"
	AttributeFilterValueEquals AttributeFilter = "ValueEquals"
)

func PossibleValuesForAttributeFilter() []string {
	return []string{
		string(AttributeFilterAll),
		string(AttributeFilterBottom),
		string(AttributeFilterTop),
		string(AttributeFilterValueEquals),
	}
}

func parseAttributeFilter(input string) (*AttributeFilter, error) {
	vals := map[string]AttributeFilter{
		"all":         AttributeFilterAll,
		"bottom":      AttributeFilterBottom,
		"top":         AttributeFilterTop,
		"valueequals": AttributeFilterValueEquals,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AttributeFilter(input)
	return &out, nil
}

type AudioAnalysisMode string

const (
	AudioAnalysisModeBasic    AudioAnalysisMode = "Basic"
	AudioAnalysisModeStandard AudioAnalysisMode = "Standard"
)

func PossibleValuesForAudioAnalysisMode() []string {
	return []string{
		string(AudioAnalysisModeBasic),
		string(AudioAnalysisModeStandard),
	}
}

func parseAudioAnalysisMode(input string) (*AudioAnalysisMode, error) {
	vals := map[string]AudioAnalysisMode{
		"basic":    AudioAnalysisModeBasic,
		"standard": AudioAnalysisModeStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AudioAnalysisMode(input)
	return &out, nil
}

type BlurType string

const (
	BlurTypeBlack BlurType = "Black"
	BlurTypeBox   BlurType = "Box"
	BlurTypeHigh  BlurType = "High"
	BlurTypeLow   BlurType = "Low"
	BlurTypeMed   BlurType = "Med"
)

func PossibleValuesForBlurType() []string {
	return []string{
		string(BlurTypeBlack),
		string(BlurTypeBox),
		string(BlurTypeHigh),
		string(BlurTypeLow),
		string(BlurTypeMed),
	}
}

func parseBlurType(input string) (*BlurType, error) {
	vals := map[string]BlurType{
		"black": BlurTypeBlack,
		"box":   BlurTypeBox,
		"high":  BlurTypeHigh,
		"low":   BlurTypeLow,
		"med":   BlurTypeMed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlurType(input)
	return &out, nil
}

type ChannelMapping string

const (
	ChannelMappingBackLeft            ChannelMapping = "BackLeft"
	ChannelMappingBackRight           ChannelMapping = "BackRight"
	ChannelMappingCenter              ChannelMapping = "Center"
	ChannelMappingFrontLeft           ChannelMapping = "FrontLeft"
	ChannelMappingFrontRight          ChannelMapping = "FrontRight"
	ChannelMappingLowFrequencyEffects ChannelMapping = "LowFrequencyEffects"
	ChannelMappingStereoLeft          ChannelMapping = "StereoLeft"
	ChannelMappingStereoRight         ChannelMapping = "StereoRight"
)

func PossibleValuesForChannelMapping() []string {
	return []string{
		string(ChannelMappingBackLeft),
		string(ChannelMappingBackRight),
		string(ChannelMappingCenter),
		string(ChannelMappingFrontLeft),
		string(ChannelMappingFrontRight),
		string(ChannelMappingLowFrequencyEffects),
		string(ChannelMappingStereoLeft),
		string(ChannelMappingStereoRight),
	}
}

func parseChannelMapping(input string) (*ChannelMapping, error) {
	vals := map[string]ChannelMapping{
		"backleft":            ChannelMappingBackLeft,
		"backright":           ChannelMappingBackRight,
		"center":              ChannelMappingCenter,
		"frontleft":           ChannelMappingFrontLeft,
		"frontright":          ChannelMappingFrontRight,
		"lowfrequencyeffects": ChannelMappingLowFrequencyEffects,
		"stereoleft":          ChannelMappingStereoLeft,
		"stereoright":         ChannelMappingStereoRight,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ChannelMapping(input)
	return &out, nil
}

type DeinterlaceMode string

const (
	DeinterlaceModeAutoPixelAdaptive DeinterlaceMode = "AutoPixelAdaptive"
	DeinterlaceModeOff               DeinterlaceMode = "Off"
)

func PossibleValuesForDeinterlaceMode() []string {
	return []string{
		string(DeinterlaceModeAutoPixelAdaptive),
		string(DeinterlaceModeOff),
	}
}

func parseDeinterlaceMode(input string) (*DeinterlaceMode, error) {
	vals := map[string]DeinterlaceMode{
		"autopixeladaptive": DeinterlaceModeAutoPixelAdaptive,
		"off":               DeinterlaceModeOff,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeinterlaceMode(input)
	return &out, nil
}

type DeinterlaceParity string

const (
	DeinterlaceParityAuto             DeinterlaceParity = "Auto"
	DeinterlaceParityBottomFieldFirst DeinterlaceParity = "BottomFieldFirst"
	DeinterlaceParityTopFieldFirst    DeinterlaceParity = "TopFieldFirst"
)

func PossibleValuesForDeinterlaceParity() []string {
	return []string{
		string(DeinterlaceParityAuto),
		string(DeinterlaceParityBottomFieldFirst),
		string(DeinterlaceParityTopFieldFirst),
	}
}

func parseDeinterlaceParity(input string) (*DeinterlaceParity, error) {
	vals := map[string]DeinterlaceParity{
		"auto":             DeinterlaceParityAuto,
		"bottomfieldfirst": DeinterlaceParityBottomFieldFirst,
		"topfieldfirst":    DeinterlaceParityTopFieldFirst,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeinterlaceParity(input)
	return &out, nil
}

type EncoderNamedPreset string

const (
	EncoderNamedPresetAACGoodQualityAudio                         EncoderNamedPreset = "AACGoodQualityAudio"
	EncoderNamedPresetAdaptiveStreaming                           EncoderNamedPreset = "AdaptiveStreaming"
	EncoderNamedPresetContentAwareEncoding                        EncoderNamedPreset = "ContentAwareEncoding"
	EncoderNamedPresetContentAwareEncodingExperimental            EncoderNamedPreset = "ContentAwareEncodingExperimental"
	EncoderNamedPresetCopyAllBitrateNonInterleaved                EncoderNamedPreset = "CopyAllBitrateNonInterleaved"
	EncoderNamedPresetHTwoSixFiveAdaptiveStreaming                EncoderNamedPreset = "H265AdaptiveStreaming"
	EncoderNamedPresetHTwoSixFiveContentAwareEncoding             EncoderNamedPreset = "H265ContentAwareEncoding"
	EncoderNamedPresetHTwoSixFiveSingleBitrateFourK               EncoderNamedPreset = "H265SingleBitrate4K"
	EncoderNamedPresetHTwoSixFiveSingleBitrateOneZeroEightZerop   EncoderNamedPreset = "H265SingleBitrate1080p"
	EncoderNamedPresetHTwoSixFiveSingleBitrateSevenTwoZerop       EncoderNamedPreset = "H265SingleBitrate720p"
	EncoderNamedPresetHTwoSixFourMultipleBitrateOneZeroEightZerop EncoderNamedPreset = "H264MultipleBitrate1080p"
	EncoderNamedPresetHTwoSixFourMultipleBitrateSD                EncoderNamedPreset = "H264MultipleBitrateSD"
	EncoderNamedPresetHTwoSixFourMultipleBitrateSevenTwoZerop     EncoderNamedPreset = "H264MultipleBitrate720p"
	EncoderNamedPresetHTwoSixFourSingleBitrateOneZeroEightZerop   EncoderNamedPreset = "H264SingleBitrate1080p"
	EncoderNamedPresetHTwoSixFourSingleBitrateSD                  EncoderNamedPreset = "H264SingleBitrateSD"
	EncoderNamedPresetHTwoSixFourSingleBitrateSevenTwoZerop       EncoderNamedPreset = "H264SingleBitrate720p"
)

func PossibleValuesForEncoderNamedPreset() []string {
	return []string{
		string(EncoderNamedPresetAACGoodQualityAudio),
		string(EncoderNamedPresetAdaptiveStreaming),
		string(EncoderNamedPresetContentAwareEncoding),
		string(EncoderNamedPresetContentAwareEncodingExperimental),
		string(EncoderNamedPresetCopyAllBitrateNonInterleaved),
		string(EncoderNamedPresetHTwoSixFiveAdaptiveStreaming),
		string(EncoderNamedPresetHTwoSixFiveContentAwareEncoding),
		string(EncoderNamedPresetHTwoSixFiveSingleBitrateFourK),
		string(EncoderNamedPresetHTwoSixFiveSingleBitrateOneZeroEightZerop),
		string(EncoderNamedPresetHTwoSixFiveSingleBitrateSevenTwoZerop),
		string(EncoderNamedPresetHTwoSixFourMultipleBitrateOneZeroEightZerop),
		string(EncoderNamedPresetHTwoSixFourMultipleBitrateSD),
		string(EncoderNamedPresetHTwoSixFourMultipleBitrateSevenTwoZerop),
		string(EncoderNamedPresetHTwoSixFourSingleBitrateOneZeroEightZerop),
		string(EncoderNamedPresetHTwoSixFourSingleBitrateSD),
		string(EncoderNamedPresetHTwoSixFourSingleBitrateSevenTwoZerop),
	}
}

func parseEncoderNamedPreset(input string) (*EncoderNamedPreset, error) {
	vals := map[string]EncoderNamedPreset{
		"aacgoodqualityaudio":              EncoderNamedPresetAACGoodQualityAudio,
		"adaptivestreaming":                EncoderNamedPresetAdaptiveStreaming,
		"contentawareencoding":             EncoderNamedPresetContentAwareEncoding,
		"contentawareencodingexperimental": EncoderNamedPresetContentAwareEncodingExperimental,
		"copyallbitratenoninterleaved":     EncoderNamedPresetCopyAllBitrateNonInterleaved,
		"h265adaptivestreaming":            EncoderNamedPresetHTwoSixFiveAdaptiveStreaming,
		"h265contentawareencoding":         EncoderNamedPresetHTwoSixFiveContentAwareEncoding,
		"h265singlebitrate4k":              EncoderNamedPresetHTwoSixFiveSingleBitrateFourK,
		"h265singlebitrate1080p":           EncoderNamedPresetHTwoSixFiveSingleBitrateOneZeroEightZerop,
		"h265singlebitrate720p":            EncoderNamedPresetHTwoSixFiveSingleBitrateSevenTwoZerop,
		"h264multiplebitrate1080p":         EncoderNamedPresetHTwoSixFourMultipleBitrateOneZeroEightZerop,
		"h264multiplebitratesd":            EncoderNamedPresetHTwoSixFourMultipleBitrateSD,
		"h264multiplebitrate720p":          EncoderNamedPresetHTwoSixFourMultipleBitrateSevenTwoZerop,
		"h264singlebitrate1080p":           EncoderNamedPresetHTwoSixFourSingleBitrateOneZeroEightZerop,
		"h264singlebitratesd":              EncoderNamedPresetHTwoSixFourSingleBitrateSD,
		"h264singlebitrate720p":            EncoderNamedPresetHTwoSixFourSingleBitrateSevenTwoZerop,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncoderNamedPreset(input)
	return &out, nil
}

type EntropyMode string

const (
	EntropyModeCabac EntropyMode = "Cabac"
	EntropyModeCavlc EntropyMode = "Cavlc"
)

func PossibleValuesForEntropyMode() []string {
	return []string{
		string(EntropyModeCabac),
		string(EntropyModeCavlc),
	}
}

func parseEntropyMode(input string) (*EntropyMode, error) {
	vals := map[string]EntropyMode{
		"cabac": EntropyModeCabac,
		"cavlc": EntropyModeCavlc,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EntropyMode(input)
	return &out, nil
}

type FaceRedactorMode string

const (
	FaceRedactorModeAnalyze  FaceRedactorMode = "Analyze"
	FaceRedactorModeCombined FaceRedactorMode = "Combined"
	FaceRedactorModeRedact   FaceRedactorMode = "Redact"
)

func PossibleValuesForFaceRedactorMode() []string {
	return []string{
		string(FaceRedactorModeAnalyze),
		string(FaceRedactorModeCombined),
		string(FaceRedactorModeRedact),
	}
}

func parseFaceRedactorMode(input string) (*FaceRedactorMode, error) {
	vals := map[string]FaceRedactorMode{
		"analyze":  FaceRedactorModeAnalyze,
		"combined": FaceRedactorModeCombined,
		"redact":   FaceRedactorModeRedact,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FaceRedactorMode(input)
	return &out, nil
}

type H264Complexity string

const (
	H264ComplexityBalanced H264Complexity = "Balanced"
	H264ComplexityQuality  H264Complexity = "Quality"
	H264ComplexitySpeed    H264Complexity = "Speed"
)

func PossibleValuesForH264Complexity() []string {
	return []string{
		string(H264ComplexityBalanced),
		string(H264ComplexityQuality),
		string(H264ComplexitySpeed),
	}
}

func parseH264Complexity(input string) (*H264Complexity, error) {
	vals := map[string]H264Complexity{
		"balanced": H264ComplexityBalanced,
		"quality":  H264ComplexityQuality,
		"speed":    H264ComplexitySpeed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := H264Complexity(input)
	return &out, nil
}

type H264VideoProfile string

const (
	H264VideoProfileAuto             H264VideoProfile = "Auto"
	H264VideoProfileBaseline         H264VideoProfile = "Baseline"
	H264VideoProfileHigh             H264VideoProfile = "High"
	H264VideoProfileHighFourFourFour H264VideoProfile = "High444"
	H264VideoProfileHighFourTwoTwo   H264VideoProfile = "High422"
	H264VideoProfileMain             H264VideoProfile = "Main"
)

func PossibleValuesForH264VideoProfile() []string {
	return []string{
		string(H264VideoProfileAuto),
		string(H264VideoProfileBaseline),
		string(H264VideoProfileHigh),
		string(H264VideoProfileHighFourFourFour),
		string(H264VideoProfileHighFourTwoTwo),
		string(H264VideoProfileMain),
	}
}

func parseH264VideoProfile(input string) (*H264VideoProfile, error) {
	vals := map[string]H264VideoProfile{
		"auto":     H264VideoProfileAuto,
		"baseline": H264VideoProfileBaseline,
		"high":     H264VideoProfileHigh,
		"high444":  H264VideoProfileHighFourFourFour,
		"high422":  H264VideoProfileHighFourTwoTwo,
		"main":     H264VideoProfileMain,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := H264VideoProfile(input)
	return &out, nil
}

type H265Complexity string

const (
	H265ComplexityBalanced H265Complexity = "Balanced"
	H265ComplexityQuality  H265Complexity = "Quality"
	H265ComplexitySpeed    H265Complexity = "Speed"
)

func PossibleValuesForH265Complexity() []string {
	return []string{
		string(H265ComplexityBalanced),
		string(H265ComplexityQuality),
		string(H265ComplexitySpeed),
	}
}

func parseH265Complexity(input string) (*H265Complexity, error) {
	vals := map[string]H265Complexity{
		"balanced": H265ComplexityBalanced,
		"quality":  H265ComplexityQuality,
		"speed":    H265ComplexitySpeed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := H265Complexity(input)
	return &out, nil
}

type H265VideoProfile string

const (
	H265VideoProfileAuto H265VideoProfile = "Auto"
	H265VideoProfileMain H265VideoProfile = "Main"
)

func PossibleValuesForH265VideoProfile() []string {
	return []string{
		string(H265VideoProfileAuto),
		string(H265VideoProfileMain),
	}
}

func parseH265VideoProfile(input string) (*H265VideoProfile, error) {
	vals := map[string]H265VideoProfile{
		"auto": H265VideoProfileAuto,
		"main": H265VideoProfileMain,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := H265VideoProfile(input)
	return &out, nil
}

type InsightsType string

const (
	InsightsTypeAllInsights       InsightsType = "AllInsights"
	InsightsTypeAudioInsightsOnly InsightsType = "AudioInsightsOnly"
	InsightsTypeVideoInsightsOnly InsightsType = "VideoInsightsOnly"
)

func PossibleValuesForInsightsType() []string {
	return []string{
		string(InsightsTypeAllInsights),
		string(InsightsTypeAudioInsightsOnly),
		string(InsightsTypeVideoInsightsOnly),
	}
}

func parseInsightsType(input string) (*InsightsType, error) {
	vals := map[string]InsightsType{
		"allinsights":       InsightsTypeAllInsights,
		"audioinsightsonly": InsightsTypeAudioInsightsOnly,
		"videoinsightsonly": InsightsTypeVideoInsightsOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InsightsType(input)
	return &out, nil
}

type JobErrorCategory string

const (
	JobErrorCategoryConfiguration JobErrorCategory = "Configuration"
	JobErrorCategoryContent       JobErrorCategory = "Content"
	JobErrorCategoryDownload      JobErrorCategory = "Download"
	JobErrorCategoryService       JobErrorCategory = "Service"
	JobErrorCategoryUpload        JobErrorCategory = "Upload"
)

func PossibleValuesForJobErrorCategory() []string {
	return []string{
		string(JobErrorCategoryConfiguration),
		string(JobErrorCategoryContent),
		string(JobErrorCategoryDownload),
		string(JobErrorCategoryService),
		string(JobErrorCategoryUpload),
	}
}

func parseJobErrorCategory(input string) (*JobErrorCategory, error) {
	vals := map[string]JobErrorCategory{
		"configuration": JobErrorCategoryConfiguration,
		"content":       JobErrorCategoryContent,
		"download":      JobErrorCategoryDownload,
		"service":       JobErrorCategoryService,
		"upload":        JobErrorCategoryUpload,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobErrorCategory(input)
	return &out, nil
}

type JobErrorCode string

const (
	JobErrorCodeConfigurationUnsupported JobErrorCode = "ConfigurationUnsupported"
	JobErrorCodeContentMalformed         JobErrorCode = "ContentMalformed"
	JobErrorCodeContentUnsupported       JobErrorCode = "ContentUnsupported"
	JobErrorCodeDownloadNotAccessible    JobErrorCode = "DownloadNotAccessible"
	JobErrorCodeDownloadTransientError   JobErrorCode = "DownloadTransientError"
	JobErrorCodeServiceError             JobErrorCode = "ServiceError"
	JobErrorCodeServiceTransientError    JobErrorCode = "ServiceTransientError"
	JobErrorCodeUploadNotAccessible      JobErrorCode = "UploadNotAccessible"
	JobErrorCodeUploadTransientError     JobErrorCode = "UploadTransientError"
)

func PossibleValuesForJobErrorCode() []string {
	return []string{
		string(JobErrorCodeConfigurationUnsupported),
		string(JobErrorCodeContentMalformed),
		string(JobErrorCodeContentUnsupported),
		string(JobErrorCodeDownloadNotAccessible),
		string(JobErrorCodeDownloadTransientError),
		string(JobErrorCodeServiceError),
		string(JobErrorCodeServiceTransientError),
		string(JobErrorCodeUploadNotAccessible),
		string(JobErrorCodeUploadTransientError),
	}
}

func parseJobErrorCode(input string) (*JobErrorCode, error) {
	vals := map[string]JobErrorCode{
		"configurationunsupported": JobErrorCodeConfigurationUnsupported,
		"contentmalformed":         JobErrorCodeContentMalformed,
		"contentunsupported":       JobErrorCodeContentUnsupported,
		"downloadnotaccessible":    JobErrorCodeDownloadNotAccessible,
		"downloadtransienterror":   JobErrorCodeDownloadTransientError,
		"serviceerror":             JobErrorCodeServiceError,
		"servicetransienterror":    JobErrorCodeServiceTransientError,
		"uploadnotaccessible":      JobErrorCodeUploadNotAccessible,
		"uploadtransienterror":     JobErrorCodeUploadTransientError,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobErrorCode(input)
	return &out, nil
}

type JobRetry string

const (
	JobRetryDoNotRetry JobRetry = "DoNotRetry"
	JobRetryMayRetry   JobRetry = "MayRetry"
)

func PossibleValuesForJobRetry() []string {
	return []string{
		string(JobRetryDoNotRetry),
		string(JobRetryMayRetry),
	}
}

func parseJobRetry(input string) (*JobRetry, error) {
	vals := map[string]JobRetry{
		"donotretry": JobRetryDoNotRetry,
		"mayretry":   JobRetryMayRetry,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobRetry(input)
	return &out, nil
}

type JobState string

const (
	JobStateCanceled   JobState = "Canceled"
	JobStateCanceling  JobState = "Canceling"
	JobStateError      JobState = "Error"
	JobStateFinished   JobState = "Finished"
	JobStateProcessing JobState = "Processing"
	JobStateQueued     JobState = "Queued"
	JobStateScheduled  JobState = "Scheduled"
)

func PossibleValuesForJobState() []string {
	return []string{
		string(JobStateCanceled),
		string(JobStateCanceling),
		string(JobStateError),
		string(JobStateFinished),
		string(JobStateProcessing),
		string(JobStateQueued),
		string(JobStateScheduled),
	}
}

func parseJobState(input string) (*JobState, error) {
	vals := map[string]JobState{
		"canceled":   JobStateCanceled,
		"canceling":  JobStateCanceling,
		"error":      JobStateError,
		"finished":   JobStateFinished,
		"processing": JobStateProcessing,
		"queued":     JobStateQueued,
		"scheduled":  JobStateScheduled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobState(input)
	return &out, nil
}

type OnErrorType string

const (
	OnErrorTypeContinueJob       OnErrorType = "ContinueJob"
	OnErrorTypeStopProcessingJob OnErrorType = "StopProcessingJob"
)

func PossibleValuesForOnErrorType() []string {
	return []string{
		string(OnErrorTypeContinueJob),
		string(OnErrorTypeStopProcessingJob),
	}
}

func parseOnErrorType(input string) (*OnErrorType, error) {
	vals := map[string]OnErrorType{
		"continuejob":       OnErrorTypeContinueJob,
		"stopprocessingjob": OnErrorTypeStopProcessingJob,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OnErrorType(input)
	return &out, nil
}

type Priority string

const (
	PriorityHigh   Priority = "High"
	PriorityLow    Priority = "Low"
	PriorityNormal Priority = "Normal"
)

func PossibleValuesForPriority() []string {
	return []string{
		string(PriorityHigh),
		string(PriorityLow),
		string(PriorityNormal),
	}
}

func parsePriority(input string) (*Priority, error) {
	vals := map[string]Priority{
		"high":   PriorityHigh,
		"low":    PriorityLow,
		"normal": PriorityNormal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Priority(input)
	return &out, nil
}

type Rotation string

const (
	RotationAuto               Rotation = "Auto"
	RotationNone               Rotation = "None"
	RotationRotateNineZero     Rotation = "Rotate90"
	RotationRotateOneEightZero Rotation = "Rotate180"
	RotationRotateTwoSevenZero Rotation = "Rotate270"
	RotationRotateZero         Rotation = "Rotate0"
)

func PossibleValuesForRotation() []string {
	return []string{
		string(RotationAuto),
		string(RotationNone),
		string(RotationRotateNineZero),
		string(RotationRotateOneEightZero),
		string(RotationRotateTwoSevenZero),
		string(RotationRotateZero),
	}
}

func parseRotation(input string) (*Rotation, error) {
	vals := map[string]Rotation{
		"auto":      RotationAuto,
		"none":      RotationNone,
		"rotate90":  RotationRotateNineZero,
		"rotate180": RotationRotateOneEightZero,
		"rotate270": RotationRotateTwoSevenZero,
		"rotate0":   RotationRotateZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Rotation(input)
	return &out, nil
}

type StretchMode string

const (
	StretchModeAutoFit  StretchMode = "AutoFit"
	StretchModeAutoSize StretchMode = "AutoSize"
	StretchModeNone     StretchMode = "None"
)

func PossibleValuesForStretchMode() []string {
	return []string{
		string(StretchModeAutoFit),
		string(StretchModeAutoSize),
		string(StretchModeNone),
	}
}

func parseStretchMode(input string) (*StretchMode, error) {
	vals := map[string]StretchMode{
		"autofit":  StretchModeAutoFit,
		"autosize": StretchModeAutoSize,
		"none":     StretchModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StretchMode(input)
	return &out, nil
}

type TrackAttribute string

const (
	TrackAttributeBitrate  TrackAttribute = "Bitrate"
	TrackAttributeLanguage TrackAttribute = "Language"
)

func PossibleValuesForTrackAttribute() []string {
	return []string{
		string(TrackAttributeBitrate),
		string(TrackAttributeLanguage),
	}
}

func parseTrackAttribute(input string) (*TrackAttribute, error) {
	vals := map[string]TrackAttribute{
		"bitrate":  TrackAttributeBitrate,
		"language": TrackAttributeLanguage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrackAttribute(input)
	return &out, nil
}

type VideoSyncMode string

const (
	VideoSyncModeAuto        VideoSyncMode = "Auto"
	VideoSyncModeCfr         VideoSyncMode = "Cfr"
	VideoSyncModePassthrough VideoSyncMode = "Passthrough"
	VideoSyncModeVfr         VideoSyncMode = "Vfr"
)

func PossibleValuesForVideoSyncMode() []string {
	return []string{
		string(VideoSyncModeAuto),
		string(VideoSyncModeCfr),
		string(VideoSyncModePassthrough),
		string(VideoSyncModeVfr),
	}
}

func parseVideoSyncMode(input string) (*VideoSyncMode, error) {
	vals := map[string]VideoSyncMode{
		"auto":        VideoSyncModeAuto,
		"cfr":         VideoSyncModeCfr,
		"passthrough": VideoSyncModePassthrough,
		"vfr":         VideoSyncModeVfr,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VideoSyncMode(input)
	return &out, nil
}
