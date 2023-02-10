package models

type Error struct {
	Detail string `json:"detail" validate:"required"  enums:"internal_server_error,endpoint_not_found,invalid_credentials,invalid_jwt,insufficient_privileges,current_user_not_found,current_user_inactive,current_user_not_superuser,user_already_registered,user_not_found,user_inactive,follower_relation_already_registered,follower_relation_not_found,post_not_found"`
} // @Name Error

type ValidationError struct {
	Detail interface{} `json:"detail" validate:"required"`
} // @Name ValidationError
