package models
import (
    "errors"
)
// 
type MeetingChatMode int

const (
    ENABLED_MEETINGCHATMODE MeetingChatMode = iota
    DISABLED_MEETINGCHATMODE
    LIMITED_MEETINGCHATMODE
    UNKNOWNFUTUREVALUE_MEETINGCHATMODE
)

func (i MeetingChatMode) String() string {
    return []string{"enabled", "disabled", "limited", "unknownFutureValue"}[i]
}
func ParseMeetingChatMode(v string) (any, error) {
    result := ENABLED_MEETINGCHATMODE
    switch v {
        case "enabled":
            result = ENABLED_MEETINGCHATMODE
        case "disabled":
            result = DISABLED_MEETINGCHATMODE
        case "limited":
            result = LIMITED_MEETINGCHATMODE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MEETINGCHATMODE
        default:
            return 0, errors.New("Unknown MeetingChatMode value: " + v)
    }
    return &result, nil
}
func SerializeMeetingChatMode(values []MeetingChatMode) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
