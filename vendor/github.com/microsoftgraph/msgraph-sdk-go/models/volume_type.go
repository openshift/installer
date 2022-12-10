package models
import (
    "errors"
)
// Provides operations to manage the collection of agreement entities.
type VolumeType int

const (
    OPERATINGSYSTEMVOLUME_VOLUMETYPE VolumeType = iota
    FIXEDDATAVOLUME_VOLUMETYPE
    REMOVABLEDATAVOLUME_VOLUMETYPE
    UNKNOWNFUTUREVALUE_VOLUMETYPE
)

func (i VolumeType) String() string {
    return []string{"operatingSystemVolume", "fixedDataVolume", "removableDataVolume", "unknownFutureValue"}[i]
}
func ParseVolumeType(v string) (interface{}, error) {
    result := OPERATINGSYSTEMVOLUME_VOLUMETYPE
    switch v {
        case "operatingSystemVolume":
            result = OPERATINGSYSTEMVOLUME_VOLUMETYPE
        case "fixedDataVolume":
            result = FIXEDDATAVOLUME_VOLUMETYPE
        case "removableDataVolume":
            result = REMOVABLEDATAVOLUME_VOLUMETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_VOLUMETYPE
        default:
            return 0, errors.New("Unknown VolumeType value: " + v)
    }
    return &result, nil
}
func SerializeVolumeType(values []VolumeType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
