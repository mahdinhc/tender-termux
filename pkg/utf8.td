fn string_from_char_code(n) {
	return string(char(n))
}

// Define ucs2_decode function
ucs2_decode := fn(string) {
    output := []
    counter := 0
    length := len(string)
    value := 0
    extra := 0
	
    for counter < length {
        value := int(char(string[counter]))
        counter++
        if value != null {
            if value >= 0xD800 && value <= 0xDBFF && counter < length - 1 {
                // high surrogate, and there is a next character
                extra := char(string[counter + 1])
                if extra & 0xFC00 == 0xDC00 {
                    // low surrogate
                    output = append(output, ((value & 0x3FF) << 10) + (extra & 0x3FF) + 0x10000)
                    counter := counter + 2
                } 
                else {
                    // unmatched surrogate
                    output = append(output, value)
                    counter := counter + 1
                }
            } 
            else {
                output = append(output, value)
                counter := counter + 1
            }
        }
    }
	
    return output
}

// Define ucs2_encode function
ucs2_encode := fn(array) {
    length := len(array)
    index := 0
    value := 0
    output := ""
	
    for index < length {
        value := array[index]
        if value > 0xFFFF {
            value := value - 0x10000
            output += string_from_char_code((value >> 10 & 0x3FF) | 0xD800)
            value = 0xDC00 | (value & 0x3FF)
        }
        output += string_from_char_code(value)
        index++
    }
	
    return output
}

// Define check_scalar_value function
check_scalar_value := fn(code_point) {
    if code_point >= 0xD800 && code_point <= 0xDFFF {
        println(error("Lone surrogate \"" + string(code_point) + "\" is not a scalar value"))
    }
}

create_byte := fn(code_point, shift) {
    return string_from_char_code(((code_point >> shift) & 0x3F) | 0x80)
}

encode_code_point := fn(code_point) {
    if code_point & 0xFFFFFF80 == 0 {
        // 1-byte sequence
        return string_from_char_code(code_point)
    }
    symbol := ""
    if code_point & 0xFFFFF800 == 0 {
        // 2-byte sequence
        symbol = string_from_char_code(((code_point >> 6) & 0x1F) | 0xC0)
    } 
    else if code_point & 0xFFFF0000 == 0 {
        // 3-byte sequence
        check_scalar_value(code_point)
        symbol = string_from_char_code(((code_point >> 12) & 0x0F) | 0xE0)
        symbol += create_byte(code_point, 6)
    } 
    else if code_point & 0xFFE00000 == 0 {
        // 4-byte sequence
        symbol = string_from_char_code(((code_point >> 18) & 0x07) | 0xF0)
        symbol += create_byte(code_point, 12)
        symbol += create_byte(code_point, 6)
    }
    symbol += string_from_char_code((code_point & 0x3F) | 0x80)
    return symbol
}

fn utf8_encode(str) {
    code_points := ucs2_decode(str)
    length := len(code_points)
    index := 0
    code_point := 0
    byte_string := ""
    for index < length {
        code_point = code_points[index]
        byte_string += encode_code_point(code_point)
        index++
    }
    return byte_string
}

byte_array := []
byte_count := 0
byte_index := 0

// Define read_continuation_byte function
read_continuation_byte := fn() {
    if byte_index >= byte_count {
        println(error("Invalid byte index"))
    }
	
    continuation_byte := byte_array[byte_index] & 0xFF
    byte_index = byte_index + 1
	
    if continuation_byte & 0xC0 == 0x80 {
        return continuation_byte & 0x3F
    }
	
    // If we end up here, it's not a continuation byte
    println(error("Invalid continuation byte"))
}

// Define decode_symbol function
decode_symbol := fn() {
    byte1 := 0
    byte2 := 0
    byte3 := 0
    byte4 := 0
    code_point := 0
	
    if byte_index > byte_count {
        println(error("Invalid byte index"))
    }
	
    if byte_index == byte_count {
        return false
    }
	
    // Read first byte
    byte1 = byte_array[byte_index] & 0xFF
    byte_index = byte_index + 1
	
    // 1-byte sequence (no continuation bytes)
    if byte1 & 0x80 == 0 {
        return byte1
    }
	
    // 2-byte sequence
    if byte1 & 0xE0 == 0xC0 {
        byte2 = read_continuation_byte()
        code_point = ((byte1 & 0x1F) << 6) | byte2
        if code_point >= 0x80 {
            return code_point
        } 
        else {
            println(error("Invalid continuation byte"))
        }
    }
	
    // 3-byte sequence (may include unpaired surrogates)
    if byte1 & 0xF0 == 0xE0 {
        byte2 = read_continuation_byte()
        byte3 = read_continuation_byte()
        code_point = ((byte1 & 0x0F) << 12) | (byte2 << 6) | byte3
        if code_point >= 0x0800 {
            check_scalar_value(code_point)
            return code_point
        } 
        else {
            println(error("Invalid continuation byte"))
        }
    }
	
    // 4-byte sequence
    if byte1 & 0xF8 == 0xF0 {
        byte2 = read_continuation_byte()
        byte3 = read_continuation_byte()
        byte4 = read_continuation_byte()
        code_point = ((byte1 & 0x07) << 0x12) | (byte2 << 0x0C) | (byte3 << 0x06) | byte4
        if code_point >= 0x010000 && code_point <= 0x10FFFF {
            return code_point
        }
    }
	
    println(error("Invalid UTF-8 detected"))
}

// Define utf8_decode function
utf8_decode := fn(byte_string) {
    byte_array = ucs2_decode(byte_string)
    byte_count = len(byte_array)
    byte_index = 0
    code_points := []
    tmp := null
	
    for tmp != false {
        tmp = decode_symbol()
        if tmp != false {
            code_points = append(code_points, tmp)
        }
    }
	
    return ucs2_encode(code_points)
}

export {
    ucs2_decode : ucs2_decode,
    ucs2_encode : ucs2_encode,
    decode : utf8_decode,
    encode : utf8_encode
}