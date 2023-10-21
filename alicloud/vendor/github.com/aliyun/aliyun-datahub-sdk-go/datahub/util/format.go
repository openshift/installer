package util

import (
    "strconv"
    "unicode"
)

func CheckProjectName(projectName string) bool {
    return isNameValid(projectName, 3, 32)
}

func CheckTopicName(topicName string) bool {
    return isNameValid(topicName, 1, 128)
}

func CheckComment(comment string) bool {
    if comment == "" || len(comment) > 1024 {
        return false
    }
    return true
}

func CheckShardId(shardId string) bool {
    if _, err := strconv.Atoi(shardId); err != nil {
        return false
    }
    return true
}

func isNameValid(name string, minLen, maxLen int) bool {
    if name == "" || len(name) > maxLen || len(name) < minLen {
        return false
    }
    for _, c := range name {
        if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' {
            return false
        }
    }
    return true
}
