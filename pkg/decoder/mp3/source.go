package mp3

import (
	"io"
)

type source struct {
	reader io.ReadSeeker

	buf []byte

	pos int64
}

func (s *source) Seek(position int64, whence int) (int64, error) {
	seeker := s.reader
	s.buf = nil
	n, err := seeker.Seek(position, whence)
	if err != nil {
		return 0, err
	}
	s.pos = n
	return n, nil
}

func (s *source) skipTags() error {
	buf := make([]byte, 3)
	if _, err := s.readFull(buf); err != nil {
		return err
	}
	switch string(buf) {
	case "TAG":
		buf := make([]byte, 125)
		if _, err := s.readFull(buf); err != nil {
			return err
		}

	case "ID3":
		// Skip version (2 bytes) and flag (1 byte)
		buf := make([]byte, 3)
		if _, err := s.readFull(buf); err != nil {
			return err
		}

		buf = make([]byte, 4)
		n, err := s.readFull(buf)
		if err != nil {
			return err
		}
		if n != 4 {
			return nil
		}
		size := (uint32(buf[0]) << 21) | (uint32(buf[1]) << 14) |
			(uint32(buf[2]) << 7) | uint32(buf[3])
		buf = make([]byte, size)
		if _, err := s.readFull(buf); err != nil {
			return err
		}

	default:
		s.unread(buf)
	}

	return nil
}

func (s *source) rewind() error {
	if _, err := s.Seek(0, io.SeekStart); err != nil {
		return err
	}
	s.pos = 0
	s.buf = nil
	return nil
}

func (s *source) unread(buf []byte) {
	s.buf = append(s.buf, buf...)
	s.pos -= int64(len(buf))
}

func (s *source) readFull(buf []byte) (int, error) {
	read := 0
	if s.buf != nil {
		read = copy(buf, s.buf)
		if len(s.buf) > read {
			s.buf = s.buf[read:]
		} else {
			s.buf = nil
		}
		if len(buf) == read {
			return read, nil
		}
	}

	n, err := io.ReadFull(s.reader, buf[read:])
	if err != nil {
		// Allow if all data can't be read. This is common.
		if err == io.ErrUnexpectedEOF {
			err = io.EOF
		}
	}
	s.pos += int64(n)
	return n + read, err
}
