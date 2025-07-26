package sensor

type Incremental int

func (s *Incremental) Value() int {
	*s++
	if *s < 0 {
		*s = 0
	}

	return int(*s)
}
