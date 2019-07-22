package rxlib

/*
   The data types in this file implement a state-communication channel (SCC). By SCC, we mean a data
that can be used by a follower, to communicate its state to a master. In rexa's case, the kernel can
be the master and a delegate being the follower, or a delegate can be the master and its thread
being the follower. For example, if a delegate wants to know if its thread is still running, how
can it know? Well, a data like this is what you turn to.

How SCC Works

To get a new SCC, call the function NewSCChan (), a pointer to a new SCC would be returned.

	scc := rxlib.NewSCChan ()

Afterwads, get an interface for the master, and also one for the follower.

	mInterface := scc.SCCMInterface () // Getting an interface for the master.
	fInterface := scc.SCCFInterface () // Getting an interface for the follower.

Then, the master can use the "mInterface" to interact with the follower, and the follower can the
"fInterface" to interact with the master.

Master Interacting With Follower

A master can ask its follower of its state, using method WhatsUp (). See the method for more info.

	state := mInterface.WhatsUp () // Master asking for state of follower

Follower Interacting With Master

A follower can inform its master about its state, using methods State (). See the method for more
info.

	fInterface.State (rxlib.Failed, "Log file could not be found.") // Follower informing master
		about its state

*/

func NewSCChan () (*SCChan) { /* This function creates a new SC channel (SCC). Note, it is
	recommended to always use this function to create new SCCs. */

	return &SCChan {}
}

const (
	/*
		Section A: The following data here are states a follower can be in.
	*/
	UnableToStart byte = 0
	NowActive     byte = 1
	Failed        byte = 2
	NowDead       byte = 3
)

type SCChan struct { // The data type of an SC channel.
	followerState  byte
	additionalInfo string
}

// State-communication-channel Master Interface Section

func (scChan *SCChan) MInterface () (*SCCMInterface) { /* This method can be used to get a master
	interface to be used for an SCC. */

	return &SCCMInterface {scChan}
}

type SCCMInterface struct { // The data type of a master interface
	underlyingChan *SCChan
}

func (mInt *SCCMInterface) WhatsUp () (byte, string) { /* To ask for the state of the follower, this
	method can be used.

	OUTPT
	outpt 0: The state of the follower. Always use the data in Section A (scroll up to find
		"Section A") to interpret the value of this data (as the byte value of a state can
		change in a future version of this package.).

	output 1: Additional information about the state. For example, if the value of "outpt 0" is
		rxlib.Failed, this value may be a data describing the reason for the failure. This
		data is not generated by the package, its merely what the follower provides. */

	return mInt.underlyingChan.followerState, mInt.underlyingChan.additionalInfo
}

// State-communication-channel Follower Interface Section

func (scChan *SCChan) FInterface () (*SCCFInterface) { /* This method can be used to get a follower
	interface to be used for an SCC. */

	return &SCCFInterface {scChan}
}

type SCCFInterface struct { // The data type of a follower interface
	underlyingChan *SCChan
}

func (fInt *SCCFInterface) State (state byte, additionalInfo ... string) { /* This method can be
	used by a follower, to inform a master about its state.

	INPUT
	This method expects at most two inputs. If more than two inputs are entered only the first
		two would be considered, and the rest would be ignored.

	input 0: The state of the follower. Value can be only any of the data in Section A (scroll
		up to find "Section A").

	input 1: This data is optional. Its value is expected to be a string further describing
		the value of "input 0". For instance, if value of "input 0" is rxlib.Failed, value
		of this data can be something like: "Log file could be opened.". */

	fInt.underlyingChan.followerState = UnableToStart
	if len (additionalInfo) > 0 {
		fInt.underlyingChan.additionalInfo = additionalInfo [0]
	}
}