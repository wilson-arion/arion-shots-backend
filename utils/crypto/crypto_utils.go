package crypto

import (
    "crypto/md5"
    "encoding/hex"
)

func GetMd5(input string) (string, error) {
    hash := md5.New()
    defer hash.Reset()

    if _, err := hash.Write([]byte(input)); err != nil {
        return "", err
    }

    return hex.EncodeToString(hash.Sum(nil)), nil
}
