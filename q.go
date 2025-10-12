package q

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
)

// Choose how to write the function name in the file.
//
// e.g. test functions can be long and have multiple parts, so we just
// pull out the most meaningful part.
func chooseDisplayName(functionName string) string {
	parts := strings.Split(functionName, ".")

	// If this is an anonymous function, the name will be something like
	// "func1", which is unhelpful.
	//
	// Throw away that part and get the next part.
	if m, _ := regexp.MatchString("^func[0-9]+$", parts[len(parts)-1]); len(parts) > 1 && m {
		fmt.Println(m)
		return parts[len(parts)-2]
	}

	return parts[len(parts)-1]
}

func getFunctionName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "<unknown>"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "<unknown>"
	}

	return chooseDisplayName(fn.Name())
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
		return fmt.Sprintf("%+v", v) // fallback
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

	if expression[0] == '"' && expression[len(expression)-1] == '"' {
		line = "\x1b[32m" + functionName + "\x1b[39m: " + toString(value, a...) + "\n\n"
	} else {
		line = "\x1b[32m" + functionName + "\x1b[39m: " + expression + " = \x1b[36m" + toString(value, a...) + "\x1b[39m\n\n"
	}

	if _, err = f.WriteString(line); err != nil {
		panic(err)
	}
}
