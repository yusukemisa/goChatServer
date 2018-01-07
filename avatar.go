package main

import (
	"errors"
)

//ErrNoAvatarURL はAvatarインスタンスがアバターのURLを返すことができない
// 場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat:アバターのURLを取得できません")

//Avatar はユーザーのプロフィール画像を表す型
type Avatar interface {
	// GetAvatarURLは指定されたクライアントのアバターのURLを返却します。
	// 問題が発生した場合にはエラーを返します。特にURLを取得できなかった場合は
	// ErrNoAvatarURLを返却します。
	GetAvatarURL(c *cliant) (string, error)
}
type AuthAvatar struct {}
var UseAuthAvatar AuthAvatar
func (_ AuthAvatar) GetAvatarURL(c *cliant) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if
	}
}