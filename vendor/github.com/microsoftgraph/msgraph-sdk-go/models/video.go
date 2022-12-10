package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Video 
type Video struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Number of audio bits per sample.
    audioBitsPerSample *int32
    // Number of audio channels.
    audioChannels *int32
    // Name of the audio format (AAC, MP3, etc.).
    audioFormat *string
    // Number of audio samples per second.
    audioSamplesPerSecond *int32
    // Bit rate of the video in bits per second.
    bitrate *int32
    // Duration of the file in milliseconds.
    duration *int64
    // 'Four character code' name of the video format.
    fourCC *string
    // Frame rate of the video.
    frameRate *float64
    // Height of the video, in pixels.
    height *int32
    // The OdataType property
    odataType *string
    // Width of the video, in pixels.
    width *int32
}
// NewVideo instantiates a new video and sets the default values.
func NewVideo()(*Video) {
    m := &Video{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateVideoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateVideoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVideo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Video) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAudioBitsPerSample gets the audioBitsPerSample property value. Number of audio bits per sample.
func (m *Video) GetAudioBitsPerSample()(*int32) {
    return m.audioBitsPerSample
}
// GetAudioChannels gets the audioChannels property value. Number of audio channels.
func (m *Video) GetAudioChannels()(*int32) {
    return m.audioChannels
}
// GetAudioFormat gets the audioFormat property value. Name of the audio format (AAC, MP3, etc.).
func (m *Video) GetAudioFormat()(*string) {
    return m.audioFormat
}
// GetAudioSamplesPerSecond gets the audioSamplesPerSecond property value. Number of audio samples per second.
func (m *Video) GetAudioSamplesPerSecond()(*int32) {
    return m.audioSamplesPerSecond
}
// GetBitrate gets the bitrate property value. Bit rate of the video in bits per second.
func (m *Video) GetBitrate()(*int32) {
    return m.bitrate
}
// GetDuration gets the duration property value. Duration of the file in milliseconds.
func (m *Video) GetDuration()(*int64) {
    return m.duration
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Video) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["audioBitsPerSample"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetAudioBitsPerSample)
    res["audioChannels"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetAudioChannels)
    res["audioFormat"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAudioFormat)
    res["audioSamplesPerSecond"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetAudioSamplesPerSecond)
    res["bitrate"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetBitrate)
    res["duration"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt64Value(m.SetDuration)
    res["fourCC"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetFourCC)
    res["frameRate"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetFrameRate)
    res["height"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetHeight)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["width"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetWidth)
    return res
}
// GetFourCC gets the fourCC property value. 'Four character code' name of the video format.
func (m *Video) GetFourCC()(*string) {
    return m.fourCC
}
// GetFrameRate gets the frameRate property value. Frame rate of the video.
func (m *Video) GetFrameRate()(*float64) {
    return m.frameRate
}
// GetHeight gets the height property value. Height of the video, in pixels.
func (m *Video) GetHeight()(*int32) {
    return m.height
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *Video) GetOdataType()(*string) {
    return m.odataType
}
// GetWidth gets the width property value. Width of the video, in pixels.
func (m *Video) GetWidth()(*int32) {
    return m.width
}
// Serialize serializes information the current object
func (m *Video) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("audioBitsPerSample", m.GetAudioBitsPerSample())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("audioChannels", m.GetAudioChannels())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("audioFormat", m.GetAudioFormat())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("audioSamplesPerSecond", m.GetAudioSamplesPerSecond())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("bitrate", m.GetBitrate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("duration", m.GetDuration())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("fourCC", m.GetFourCC())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("frameRate", m.GetFrameRate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("height", m.GetHeight())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("width", m.GetWidth())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Video) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAudioBitsPerSample sets the audioBitsPerSample property value. Number of audio bits per sample.
func (m *Video) SetAudioBitsPerSample(value *int32)() {
    m.audioBitsPerSample = value
}
// SetAudioChannels sets the audioChannels property value. Number of audio channels.
func (m *Video) SetAudioChannels(value *int32)() {
    m.audioChannels = value
}
// SetAudioFormat sets the audioFormat property value. Name of the audio format (AAC, MP3, etc.).
func (m *Video) SetAudioFormat(value *string)() {
    m.audioFormat = value
}
// SetAudioSamplesPerSecond sets the audioSamplesPerSecond property value. Number of audio samples per second.
func (m *Video) SetAudioSamplesPerSecond(value *int32)() {
    m.audioSamplesPerSecond = value
}
// SetBitrate sets the bitrate property value. Bit rate of the video in bits per second.
func (m *Video) SetBitrate(value *int32)() {
    m.bitrate = value
}
// SetDuration sets the duration property value. Duration of the file in milliseconds.
func (m *Video) SetDuration(value *int64)() {
    m.duration = value
}
// SetFourCC sets the fourCC property value. 'Four character code' name of the video format.
func (m *Video) SetFourCC(value *string)() {
    m.fourCC = value
}
// SetFrameRate sets the frameRate property value. Frame rate of the video.
func (m *Video) SetFrameRate(value *float64)() {
    m.frameRate = value
}
// SetHeight sets the height property value. Height of the video, in pixels.
func (m *Video) SetHeight(value *int32)() {
    m.height = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Video) SetOdataType(value *string)() {
    m.odataType = value
}
// SetWidth sets the width property value. Width of the video, in pixels.
func (m *Video) SetWidth(value *int32)() {
    m.width = value
}
