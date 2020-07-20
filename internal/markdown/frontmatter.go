package markdown

// Frontmatter returns the start and end index of the frontmatter or -1,-1 if not available
func Frontmatter(in []byte) (start, end int) {
	const searchStart = 0
	const searchEnd = 1

	startPos := -1
	state := searchStart
	for i := 0; i < len(in); i++ {
		r := in[i]
		switch state {
		case searchStart:
			if r == ' ' || r == '\n' {
				continue
			} else {
				if r == '-' {
					if i+2 < len(in) {
						if in[i+1] == '-' && in[i+2] == '-' {
							i += 3
							startPos = i
							state = searchEnd
						} else {
							return -1, -1
						}
					}

				} else {
					return -1, -1
				}
			}
		case searchEnd:
			if r == '-' {
				if i+2 < len(in) && in[i+1] == '-' && in[i+2] == '-' {
					return startPos, i
				} else {
					continue
				}
			}
		}
	}
	return -1, -1
}
