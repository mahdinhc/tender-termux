
package xml

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/2dprototype/tender"
)

const (
	attributeKey = "@"
	childNodeKey = "#"
)

var unescapeMap = map[string]string{
	"&amp;":  "&",
	"&lt;":   "<",
	"&gt;":   ">",
	"&apos;": "'",
	"&quot;": `"`,
}

// Decode converts XML data to Tender objects using @ for attributes and # for text
func Decode(data []byte) (tender.Object, error) {
	if len(data) == 0 {
		return tender.NullValue, nil
	}

	root, err := parseXML(string(data))
	if err != nil {
		return nil, err
	}

	return toObject(root, nil), nil
}

type xmlNode struct {
	name     string
	attr     string
	raw      string
	children []interface{}
	isEmpty  bool
}

func parseXML(text string) (*xmlNode, error) {
	// Simple XML parser that handles tags, attributes, and text content
	// This is a simplified version - a full implementation would be more complex
	
	root := &xmlNode{name: "", children: []interface{}{}}
	stack := []*xmlNode{root}
	current := root
	
	var buffer bytes.Buffer
	var inTag bool
	var inComment bool
	var inCDATA bool
	
	for i := 0; i < len(text); i++ {
		char := text[i]
		
		if inCDATA {
			if strings.HasPrefix(text[i:], "]]>") {
				// End of CDATA
				content := buffer.String()
				if content != "" {
					current.children = append(current.children, content)
				}
				buffer.Reset()
				inCDATA = false
				i += 2 // Skip "]]"
				continue
			}
			buffer.WriteByte(char)
			continue
		}
		
		if inComment {
			if strings.HasPrefix(text[i:], "-->") {
				inComment = false
				i += 2 // Skip "--"
				continue
			}
			continue
		}
		
		if char == '<' {
			// Flush any text content
			if buffer.Len() > 0 {
				content := strings.TrimSpace(buffer.String())
				if content != "" {
					current.children = append(current.children, content)
				}
				buffer.Reset()
			}
			
			if i+1 < len(text) {
				nextChar := text[i+1]
				if nextChar == '!' {
					if strings.HasPrefix(text[i:], "<!--") {
						inComment = true
						i += 3 // Skip "!--"
						continue
						} else if strings.HasPrefix(text[i:], "<![CDATA[") {
						inCDATA = true
						i += 8 // Skip "![CDATA["
						continue
					}
				} else if nextChar == '?' {
					// XML declaration - skip
					end := strings.Index(text[i:], "?>")
					if end != -1 {
						i += end + 1
						continue
					}
				} else if nextChar == '/' {
					// Closing tag
					end := strings.Index(text[i:], ">")
					if end != -1 {
						// tagName := strings.TrimSpace(text[i+2 : i+end])
						if len(stack) > 1 {
							stack = stack[:len(stack)-1]
							current = stack[len(stack)-1]
						}
						i += end
						continue
					}
				} else {
					// Opening tag
					end := strings.Index(text[i:], ">")
					if end != -1 {
						tagContent := text[i+1 : i+end]
						
						// Check if it's self-closing
						isSelfClosing := tagContent[len(tagContent)-1] == '/'
						if isSelfClosing {
							tagContent = tagContent[:len(tagContent)-1]
						}
						
						// Parse tag name and attributes
						spaceIndex := strings.Index(tagContent, " ")
						var tagName, attrs string
						if spaceIndex == -1 {
							tagName = tagContent
							attrs = ""
						} else {
							tagName = tagContent[:spaceIndex]
							attrs = strings.TrimSpace(tagContent[spaceIndex:])
						}
						
						node := &xmlNode{
							name:     tagName,
							attr:     attrs,
							children: []interface{}{},
							isEmpty:  isSelfClosing,
						}
						
						current.children = append(current.children, node)
						
						if !isSelfClosing {
							stack = append(stack, node)
							current = node
						}
						
						i += end
						continue
					}
				}
			}
		} else if !inTag {
			buffer.WriteByte(char)
		}
	}
	
	return root, nil
}

func toObject(elem *xmlNode, reviver func(string, interface{}) interface{}) tender.Object {
	if elem.raw != "" {
		return &tender.String{Value: elem.raw}
	}
	
	attributes := parseAttribute(elem, reviver)
	childList := elem.children
	
	if attributes != nil || len(childList) > 1 {
		// Merge attributes and child nodes
		object := attributes
		if object == nil {
			object = make(map[string]tender.Object)
		}
		
		for _, child := range childList {
			switch child := child.(type) {
			case string:
				addObject(object, childNodeKey, &tender.String{Value: child})
			case *xmlNode:
				childObj := toObject(child, reviver)
				addObject(object, child.name, childObj)
			}
		}
		return &tender.Map{Value: object}
	} else if len(childList) == 1 {
		// Single child node
		child := childList[0]
		switch child := child.(type) {
		case string:
			return &tender.String{Value: child}
		case *xmlNode:
			object := toObject(child, reviver)
			if child.name != "" {
				wrap := make(map[string]tender.Object)
				wrap[child.name] = object
				return &tender.Map{Value: wrap}
			}
			return object
		}
	} else {
		// Empty node
		if elem.isEmpty {
			return tender.NullValue
		}
		return &tender.String{Value: ""}
	}
	
	return tender.NullValue
}

func parseAttribute(elem *xmlNode, reviver func(string, interface{}) interface{}) map[string]tender.Object {
	if elem.attr == "" {
		return nil
	}
	
	attributes := make(map[string]tender.Object)
	attrStr := elem.attr
	
	// Simple attribute parsing - split by spaces but respect quotes
	for attrStr != "" {
		attrStr = strings.TrimSpace(attrStr)
		if attrStr == "" {
			break
		}
		
		// Find attribute name
		eqIndex := strings.Index(attrStr, "=")
		if eqIndex == -1 {
			// Bare attribute (no value)
			spaceIndex := strings.Index(attrStr, " ")
			var name string
			if spaceIndex == -1 {
				name = attrStr
				attrStr = ""
			} else {
				name = attrStr[:spaceIndex]
				attrStr = attrStr[spaceIndex:]
			}
			attributes[attributeKey+name] = tender.NullValue
			continue
		}
		
		name := strings.TrimSpace(attrStr[:eqIndex])
		valueStr := strings.TrimSpace(attrStr[eqIndex+1:])
		
		if valueStr == "" {
			attributes[attributeKey+name] = tender.NullValue
			break
		}
		
		// Check for quoted value
		if valueStr[0] == '"' || valueStr[0] == '\'' {
			quote := valueStr[0]
			endQuote := strings.Index(valueStr[1:], string(quote))
			if endQuote != -1 {
				value := valueStr[1 : endQuote+1]
				value = unescapeXML(value)
				attributes[attributeKey+name] = &tender.String{Value: value}
				attrStr = valueStr[endQuote+2:]
			} else {
				// Malformed, take until space
				spaceIndex := strings.Index(valueStr, " ")
				if spaceIndex != -1 {
					value := valueStr[:spaceIndex]
					value = unescapeXML(value)
					attributes[attributeKey+name] = &tender.String{Value: value}
					attrStr = valueStr[spaceIndex:]
				} else {
					value := valueStr
					value = unescapeXML(value)
					attributes[attributeKey+name] = &tender.String{Value: value}
					attrStr = ""
				}
			}
		} else {
			// Unquoted value - take until space
			spaceIndex := strings.Index(valueStr, " ")
			var value string
			if spaceIndex != -1 {
				value = valueStr[:spaceIndex]
				attrStr = valueStr[spaceIndex:]
			} else {
				value = valueStr
				attrStr = ""
			}
			value = unescapeXML(value)
			attributes[attributeKey+name] = &tender.String{Value: value}
		}
	}
	
	return attributes
}

func unescapeXML(str string) string {
	result := str
	for entity, replacement := range unescapeMap {
		result = strings.ReplaceAll(result, entity, replacement)
	}
	
	// Handle numeric entities
	result = strings.ReplaceAll(result, "&#x", "&#")
	
	for {
		// Find &#[0-9]+;
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
		if len(codeStr) > 0 {
			if codeStr[0] == 'x' {
				// Hexadecimal
				code64, _ := strconv.ParseInt(codeStr[1:], 16, 32)
				code = int(code64)
			} else {
				// Decimal
				code64, _ := strconv.ParseInt(codeStr, 10, 32)
				code = int(code64)
			}
			
			if code > 0 && code <= 0x10FFFF {
				char := string(rune(code))
				result = strings.Replace(result, entity, char, 1)
			} else {
				result = strings.Replace(result, entity, "", 1)
			}
		} else {
			result = strings.Replace(result, entity, "", 1)
		}
	}
	
	return result
}

func addObject(object map[string]tender.Object, key string, val tender.Object) {
	if val == nil {
		return
	}
	
	prev, exists := object[key]
	if exists {
		if arr, ok := prev.(*tender.Array); ok {
			arr.Value = append(arr.Value, val)
		} else {
			object[key] = &tender.Array{Value: []tender.Object{prev, val}}
		}
	} else {
		object[key] = val
	}
}