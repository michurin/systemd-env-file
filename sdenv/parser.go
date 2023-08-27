package sdenv

import "errors"

type state int

const (
	stPreKey state = iota
	stKey
	stPreValue
	stValue
	stSingleQuot
	stDoubleQuot
	stEscape
	stEscapeDoubleQuot
	stComment
	stEscapeComment
)

func Parser(stream []rune) ([][2]string, error) { //nolint:gocognit,gocyclo
	// this code was stolen from systemd
	// https://github.com/systemd/systemd/blob/v253/src/basic/env-file.c#L22
	// two minor changes are marked by "[bug?]" marker lower
	result := [][2]string{}
	state := stPreKey
	keyTrSp := 0
	valueTrSp := 0
	key := []rune(nil)
	value := []rune(nil)
	for _, c := range stream {
		switch state {
		case stPreKey:
			switch c {
			case '#', ';':
				state = stComment
			case '\x20', '\t', '\n', '\r':
			default:
				state = stKey
				key = []rune{c}
			}
		case stKey:
			switch c {
			case '\n', '\r':
				state = stPreKey
			case '=':
				state = stPreValue
			case '\x20', '\t':
				keyTrSp++
				key = append(key, c)
			default:
				keyTrSp = 0
				key = append(key, c)
			}
		case stPreValue:
			switch c {
			case '\n', '\r':
				state = stPreKey
				result = append(result, [2]string{string(key[:len(key)-keyTrSp]), string(value[:len(value)-valueTrSp])})
				keyTrSp = 0
				valueTrSp = 0
				key = nil
				value = nil
			case '\'':
				state = stSingleQuot
			case '"':
				state = stDoubleQuot
			case '\\':
				state = stEscape
			case '\x20', '\t':
			default:
				state = stValue
				value = append(value, c)
			}
		case stValue:
			switch c {
			case '\n', '\r':
				state = stPreKey
				result = append(result, [2]string{string(key[:len(key)-keyTrSp]), string(value[:len(value)-valueTrSp])})
				keyTrSp = 0
				valueTrSp = 0
				key = nil
				value = nil
			case '\\':
				state = stEscape
				valueTrSp = 0
			case '\x20', '\t':
				valueTrSp++
				value = append(value, c)
			default:
				valueTrSp = 0
				value = append(value, c)
			}
		case stEscape:
			state = stValue
			switch c {
			case '\n', '\r':
			default:
				value = append(value, c)
			}
		case stSingleQuot:
			switch c {
			case '\'':
				state = stPreValue
			default:
				value = append(value, c)
			}
		case stDoubleQuot:
			switch c {
			case '"':
				state = stPreValue
			case '\\':
				state = stEscapeDoubleQuot
			default:
				value = append(value, c)
			}
		case stEscapeDoubleQuot:
			state = stDoubleQuot
			switch c {
			case '"', '\\', '`', '$': // SHELL_NEED_ESCAPE in original source
				value = append(value, c)
			case '\n', '\r': // [bug?] original implementation consider '\n' only
			default:
				value = append(value, '\\', c)
			}
		case stComment:
			switch c {
			case '\\':
				state = stEscapeComment
			case '\n', '\r':
				state = stPreKey
			}
		case stEscapeComment:
			state = stComment
		}
	}
	switch state {
	case stPreValue, stValue:
		result = append(result, [2]string{string(key[:len(key)-keyTrSp]), string(value[:len(value)-valueTrSp])})
	case stSingleQuot, stDoubleQuot, stEscape, stEscapeDoubleQuot, stEscapeComment:
		// [bug?] original code doesn't consider this states as mistaken
		return nil, errors.New("unexpected end of file")
	case stKey, stPreKey, stComment:
	}
	return result, nil
}
