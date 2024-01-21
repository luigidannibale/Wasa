package api

const (
	MsgNotFound                      = "can't be found"
	MsgNotFoundPhoto                 = "Photo " + MsgNotFound
	MsgAuthNotFound                  = "The userID provided for authentication " + MsgNotFound
	MsgAuthNoMatch                   = "Authentication userID and parameter userID don't match"
	MsgServerError                   = "An error occurred on ther server"
	MsgServerErrorUserID             = MsgServerError + " while identifying userID"
	MsgValidationErrorPhoto          = MsgServerError + " while validating the photo "
	MsgConvertionError               = "Unable to convert "
	MsgConvertionErrorUserID         = MsgConvertionError + "userID "
	MsgConvertionErrorFollowedID     = MsgConvertionError + "followedID "
	MsgConvertionErrorCommentID      = MsgConvertionError + "commentID "
	MsgConvertionErrorBannedID       = MsgConvertionError + "bannedID "
	MsgConvertionErrorPhotoID        = MsgConvertionError + "photoID "
	MsgConvertionErrorUserToBanID    = MsgConvertionError + "userToBanID"
	MsgConvertionErrorUserToFollowID = MsgConvertionError + "userToFollowID"
)
