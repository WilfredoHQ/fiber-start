package models

type Msg struct {
	Msg string `json:"msg" validate:"required" enums:"email_sent,password_updated,post_deleted,follower_relation_deleted"`
} // @Name Msg
