package q

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func getFunctionName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "<unknown>"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "<unknown>"
	}

	// The name of test functions can be long and have multiple parts,
	// e.g. tailscale.com/wgengine/magicsock.TestNetworkDownSendErrors
	//
	// For brevity, just get the last part.
	parts := strings.Split(fn.Name(), ".")
	lastPart := parts[len(parts)-1]

	return lastPart
}

func getExpression() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "<unknown>"
	}

	f, err := os.Open(file)
	if err != nil {
		return "<unknown>"
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	currentLine := 1

	for scanner.Scan() {
		if currentLine == line {
			thisLine := strings.TrimSpace(scanner.Text())
			thisLine, _ = strings.CutPrefix(thisLine, "q.Q(")
			thisLine, _ = strings.CutSuffix(thisLine, ")")
			return thisLine
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return "<unknown>"
	}

	return "<unknown>"
}

func toString(value any, a ...any) string {
	switch v := value.(type) {
	case string:
		if len(a) == 0 {
			return fmt.Sprintf("%q", v)
		} else {
			v = strings.ReplaceAll(v, "%+v", "\x1b[39m%+v\x1b[39m")
			v = strings.ReplaceAll(v, "%v", "\x1b[39m%v\x1b[39m")
			v = strings.ReplaceAll(v, "%t", "\x1b[39m%t\x1b[39m")
			return fmt.Sprintf(v, a...)
		}
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, bool:
		return fmt.Sprintf("%v", v)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v) // fallback
	}
}

func Q(value any, a ...any) {
	f, err := os.OpenFile("/tmp/q.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	functionName := getFunctionName()
	expression := getExpression()

	var line string

	switch value.(type) {
	case string:
		line = "\x1b[32m" + functionName + "\x1b[39m: " + toString(value, a...) + "\n\n"
	default:
		line = "\x1b[32m" + functionName + "\x1b[39m: " + expression + " = \x1b[36m" + toString(value, a...) + "\x1b[39m\n\n"
	}

	if _, err = f.WriteString(line); err != nil {
		panic(err)
	}
}
