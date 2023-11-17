type User struct
{
	id string
	Username string
}
type Photo struct
{
	id string
	Image string
}
type Like struct
{
	id string
	UserID User.id
}
type Comment struct
{
	id string
	UserID User.id
	Uontent string
}
type Date struct
{
	Year int
	Month string 
	Day int
	Hour int
	Minutes int
	Seconds int
}
func(d Date) validate() error
{
	validMonths := []string {"January","February","March",
	"April","May","June",
	"July","August","September",
	"October","November","December"}

	if ! arrayContains(validMonths, month) 
		return error.NewError("Bad month")
	if day < 1 || day > 31
		return error.NewError("Bad day")
	if hour < 0 || hour > 23
		return error.NewError("Bad hour")
	if minutes < 0 || minutes > 59
		return error.NewError("Bad minutes")
	if seconds < 0 || seconds > 59
		return error.NewError("Bad seconds")
	return nil
}
/*
type UserProfile struct
{
	user User
	followed[] User
	following[] User
	stream[] Photo
}
*/
func arrayContains(array []string, val string) bool {
	for i := 0; i < len(array); i++ {
		// checking if the array contains the given value
		if array[i] == val {
			// changing the boolean variable
			return true
		}
	}
	return false
}



