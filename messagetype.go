package tinystring

// MessageType represents the classification of message types in the system.
type MessageType uint8

// Msg exposes the MessageType constants for external use, following TinyString naming convention.
var Msg = struct {
	Normal  MessageType
	Info    MessageType
	Error   MessageType
	Warning MessageType
	Success MessageType
}{0, 1, 2, 3, 4}

// Helper methods for MessageType
func (t MessageType) IsNormal() bool  { return t == Msg.Normal }
func (t MessageType) IsInfo() bool    { return t == Msg.Info }
func (t MessageType) IsError() bool   { return t == Msg.Error }
func (t MessageType) IsWarning() bool { return t == Msg.Warning }
func (t MessageType) IsSuccess() bool { return t == Msg.Success }

func (t MessageType) String() string {
	switch t {
	case Msg.Info:
		return "Info"
	case Msg.Error:
		return "Error"
	case Msg.Warning:
		return "Warning"
	case Msg.Success:
		return "Success"
	default:
		return "Normal"
	}
}

// Pre-compiled patterns for efficient buffer matching
var (
	errorPatterns = [][]byte{
		[]byte("error"), []byte("failed"), []byte("exit status 1"),
		[]byte("undeclared"), []byte("undefined"), []byte("fatal"),
	}
	warningPatterns = [][]byte{
		[]byte("warning"), []byte("warn"), []byte("debug"),
	}
	successPatterns = [][]byte{
		[]byte("success"), []byte("completed"), []byte("successful"), []byte("done"),
	}
	infoPatterns = [][]byte{
		[]byte("info"), []byte("starting"), []byte("initializing"),
	}
)

// StringType returns the string from buffOut and its detected MessageType, then auto-releases the Conv
func (c *Conv) StringType() (string, MessageType) {
	// Get string content FIRST (before detection modifies buffer)
	out := c.getString(buffOut)
	// Detect type from buffOut content
	msgType := c.detectMessageTypeFromBuffer(buffOut)
	// Auto-release
	c.putConv()
	return out, msgType
}

// detectMessageTypeFromBuffer analyzes the buffer content and returns the detected MessageType (zero allocations)
func (c *Conv) detectMessageTypeFromBuffer(dest buffDest) MessageType {
	// 1. Copy content directly to work buffer using swapBuff (zero allocations)
	c.swapBuff(dest, buffWork)
	// 2. Convert to lowercase in work buffer using existing method
	c.changeCase(true, buffWork)
	// 3. Direct buffer pattern matching - NO Contains() allocations
	if c.bufferContainsPattern(buffWork, errorPatterns) {
		return Msg.Error
	}
	if c.bufferContainsPattern(buffWork, warningPatterns) {
		return Msg.Warning
	}
	if c.bufferContainsPattern(buffWork, successPatterns) {
		return Msg.Success
	}
	if c.bufferContainsPattern(buffWork, infoPatterns) {
		return Msg.Info
	}
	return Msg.Normal
}
