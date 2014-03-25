package control

//Broadcast specifies a broadcast file
//with a handler application.
type Broadcast interface {
	//Run instructs handler application to
	//play the broadcast.
	//This function should not return until the
	//broadcast has ended.
	Run() error
	//Kill kills the handler application
	//stopping the broadcast.
	Kill() error
	//Status returns boolean that specifies
	//whether the broadcast is currently running
	Status() bool
	//Path returns the location of the broadcast file
	Path() string
}
