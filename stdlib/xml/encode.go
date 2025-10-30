package xml

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/2dprototype/tender"
)

var escapeMap = map[byte]string{
	'\t': "&#x09;",
	'\n': "&#x0a;",
	'\r': "&#x0d;",
	' ':  "&#x20;",
	'&':  "&amp;",
	'<':  "&lt;",
	'>':  "&gt;",
	'"':  "&quot;",
}

// Encode converts Tender objects to XML using the @/# convention
func Encode(obj tender.Object) ([]byte, error) {
	var buf bytes.Buffer
	err := encodeValue(&buf, obj, "", 0, false)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeValue(buf *bytes.Buffer, obj tender.Object, key string, indentLevel int, inArray bool) error {
	// Handle child node synonym
	if key == childNodeKey {
		key = ""
	}

	switch o := obj.(type) {
	case *tender.Array:
		return encodeArray(buf, o, key, indentLevel)
	case *tender.Map:
		return encodeMap(buf, o, key, indentLevel, inArray)
	case *tender.String:
		return encodeString(buf, o, key, indentLevel, inArray)
	case *tender.Int:
		str := strconv.FormatInt(o.Value, 10)
		return encodeString(buf, &tender.String{Value: str}, key, indentLevel, inArray)
	case *tender.Float:
		str := strconv.FormatFloat(o.Value, 'f', -1, 64)
		return encodeString(buf, &tender.String{Value: str}, key, indentLevel, inArray)
	case *tender.Bool:
		str := "false"
		if !o.IsFalsy() {
			str = "true"
		}
		return encodeString(buf, &tender.String{Value: str}, key, indentLevel, inArray)
	case *tender.Null:
		if key != "" {
			// Empty element
			writeIndent(buf, indentLevel)
			buf.WriteString("<")
			buf.WriteString(key)
			buf.WriteString("/>")
		}
		return nil
	default:
		str, _ := tender.ToString(o)
		return encodeString(buf, &tender.String{Value: str}, key, indentLevel, inArray)
	}
}

func encodeArray(buf *bytes.Buffer, arr *tender.Array, key string, indentLevel int) error {
	for _, item := range arr.Value {
		err := encodeValue(buf, item, key, indentLevel, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func encodeMap(buf *bytes.Buffer, m *tender.Map, key string, indentLevel int, inArray bool) error {
	// Check for XML declaration or comment
	if key == "?" {
		// XML declaration
		content := getContent(m)
		writeIndent(buf, indentLevel)
		buf.WriteString("<?")
		buf.WriteString(content)
		buf.WriteString("?>")
		return nil
	} else if key == "!" {
		// Comment or CDATA
		content := getContent(m)
		writeIndent(buf, indentLevel)
		if strings.HasPrefix(content, "[CDATA[") && strings.HasSuffix(content, "]]") {
			buf.WriteString("<![CDATA[")
			buf.WriteString(content[7 : len(content)-2])
			buf.WriteString("]]>")
		} else {
			buf.WriteString("<!--")
			buf.WriteString(content)
			buf.WriteString("-->")
		}
		return nil
	}

	// Extract attributes and child elements
	attributes := make(map[string]tender.Object)
	var children []struct {
		key   string
		value tender.Object
	}

	for k, v := range m.Value {
		if strings.HasPrefix(k, attributeKey) {
			attrName := k[len(attributeKey):]
			attributes[attrName] = v
		} else {
			children = append(children, struct {
				key   string
				value tender.Object
			}{k, v})
		}
	}

	// Check if this is an empty element
	isEmpty := len(children) == 0

	// Write opening tag
	if key != "" {
		if !inArray {
			writeIndent(buf, indentLevel)
		}
		buf.WriteString("<")
		buf.WriteString(key)

		// Write attributes
		for attrName, attrValue := range attributes {
			buf.WriteString(" ")
			buf.WriteString(attrName)
			
			if attrValue != tender.NullValue {
				buf.WriteString("=\"")
				// buf.WriteString(escapeAttribute(attrValue.String()))
				v, _ := tender.ToString(attrValue)
				buf.WriteString(escapeAttribute(v))
				buf.WriteString("\"")
			}
		}

		if isEmpty {
			buf.WriteString("/>")
			return nil
		} else {
			buf.WriteString(">")
		}
	}

	// Write children
	hasElementChildren := false
	for _, child := range children {
		if child.key == childNodeKey {
			// Text content
			v, _ := tender.ToString(child.value)
			buf.WriteString(escapeTextNode(v))
		} else {
			hasElementChildren = true
			if key != "" {
				buf.WriteString("\n")
			}
			err := encodeValue(buf, child.value, child.key, indentLevel+1, false)
			if err != nil {
				return err
			}
		}
	}

	// Write closing tag
	if key != "" {
		if hasElementChildren {
			buf.WriteString("\n")
			writeIndent(buf, indentLevel)
		}
		buf.WriteString("</")
		buf.WriteString(key)
		buf.WriteString(">")
	}

	return nil
}

func encodeString(buf *bytes.Buffer, s *tender.String, key string, indentLevel int, inArray bool) error {
	if key == "?" {
		// XML declaration
		writeIndent(buf, indentLevel)
		buf.WriteString("<?")
		buf.WriteString(s.Value)
		buf.WriteString("?>")
	} else if key == "!" {
		// Comment
		writeIndent(buf, indentLevel)
		buf.WriteString("<!--")
		buf.WriteString(s.Value)
		buf.WriteString("-->")
	} else {
		if key != "" {
			if !inArray {
				writeIndent(buf, indentLevel)
			}
			buf.WriteString("<")
			buf.WriteString(key)
			buf.WriteString(">")
			buf.WriteString(escapeTextNode(s.Value))
			buf.WriteString("</")
			buf.WriteString(key)
			buf.WriteString(">")
		} else {
			buf.WriteString(escapeTextNode(s.Value))
		}
	}
	return nil
}

func getContent(m *tender.Map) string {
	for k, v := range m.Value {
		if k == childNodeKey {
			if str, ok := v.(*tender.String); ok {
				return str.Value
			}
		}
	}
	return ""
}

func writeIndent(buf *bytes.Buffer, level int) {
	for i := 0; i < level; i++ {
		buf.WriteString("  ")
	}
}

func escapeTextNode(str string) string {
	var result strings.Builder
	for _, r := range str {
		if r == '<' {
			result.WriteString("&lt;")
		} else if r == '>' {
			result.WriteString("&gt;")
		} else if r == '&' {
			result.WriteString("&amp;")
		} else if r < 32 {
			result.WriteString(fmt.Sprintf("&#x%02x;", r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func escapeAttribute(str string) string {
	var result strings.Builder
	for _, r := range str {
		if r == '"' {
			result.WriteString("&quot;")
		} else if r == '<' {
			result.WriteString("&lt;")
		} else if r == '>' {
			result.WriteString("&gt;")
		} else if r == '&' {
			result.WriteString("&amp;")
		} else if r < 32 {
			result.WriteString(fmt.Sprintf("&#x%02x;", r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// EscapeString escapes XML special characters
func EscapeString(str string) string {
	return escapeTextNode(str)
}

// UnescapeString unescapes XML entities
func UnescapeString(str string) string {
	result := str
	result = strings.ReplaceAll(result, "&amp;", "&")
	result = strings.ReplaceAll(result, "&lt;", "<")
	result = strings.ReplaceAll(result, "&gt;", ">")
	result = strings.ReplaceAll(result, "&quot;", `"`)
	result = strings.ReplaceAll(result, "&apos;", "'")
	
	// Handle numeric entities
	for {
		start := strings.Index(result, "&#")
		if start == -1 {
			break
		}
		
		end := strings.Index(result[start:], ";")
		if end == -1 {
			break
		}
		
		entity := result[start : start+end+1]
		codeStr := entity[2:end]
		
		var code int
		if len(codeStr) > 0 && codeStr[0] == 'x' {
			// Hexadecimal
			code64, _ := strconv.ParseInt(codeStr[1:], 16, 32)
			code = int(code64)
		} else {
			// Decimal
			code64, _ := strconv.ParseInt(codeStr, 10, 32)
			code = int(code64)
		}
		
		if code > 0 {
			char := string(rune(code))
			result = strings.Replace(result, entity, char, 1)
		} else {
			result = strings.Replace(result, entity, "", 1)
		}
	}
	
	return result
}