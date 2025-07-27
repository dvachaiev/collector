package publisher

type buffer struct {
	buf       []byte
	usedBytes int
	messages  []int
}

func newBuffer(size int) *buffer {
	return &buffer{
		buf: make([]byte, size),
	}
}

// Append adds a new message to the buffer and if there is no available space for the message
// discards the oldest ones. Number of discarded messages is returned.
func (b *buffer) Append(msg []byte) (discarded int) {
	if len(msg)+b.usedBytes > len(b.buf) {
		for i, offset := range b.messages {
			if offset > len(msg) {
				discarded = i + 1

				b.Remove(discarded)

				break
			}
		}
	}

	n := copy(b.buf[b.usedBytes:], msg)
	b.usedBytes += n
	b.messages = append(b.messages, b.usedBytes)

	return discarded
}

func (b *buffer) Get(num int) []byte {
	return b.buf[:b.messages[num-1]]
}

// GetAll returns all messages present in buffer and number of returned messages
func (b *buffer) GetAll() (int, []byte) {
	return len(b.messages), b.buf[:b.usedBytes]
}

func (b *buffer) Remove(num int) {
	if num <= 0 {
		return
	}

	offset := b.messages[num-1]
	msgsNum := len(b.messages) - num

	for i := range msgsNum {
		b.messages[i] = b.messages[num+i] - offset
	}

	b.messages = b.messages[:msgsNum]
	copy(b.buf, b.buf[offset:b.usedBytes])
	b.usedBytes -= offset
}
