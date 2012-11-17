package control

//Interface Broadcast specifies a broadcast file
//with a handler application.
type Broadcast interface {
	//Start function instructs handler application to
	//start playing the broadcast.
	Start() error
	//Kill kills the handler application
	//stopping the broadcast.
	Kill() error
	//Status returns boolean that specifies
	//whether the broadcast is currently running
	Status() bool
	//Path returns the location of the broadcast file
	Path() string
}

