is_iter := fn(x) {
	return is_array(x) || is_map(x) || is_immutable_array(x) || is_immutable_map(x)
}

is_array_like := fn(x) {
	return is_array(x) || is_immutable_array(x)
}

export {
	all: fn(x, fun) {
		if !is_iter(x) { return null }
		
		for k, v in x {
			if !fun(k, v) { return false }
		}
		
		return true
	},
	any: fn(x, fun) {
		if !is_iter(x) { return null }
		
		for k, v in x {
			if fun(k, v) { return true }
		}
		
		return false
	},
	chunk: fn(x, size) {
		if !is_array_like(x) || !size { return null }
		
		numElements := len(x)
		if !numElements { return [] }
		
		res := []
		idx := 0
		for idx < numElements {
			res = append(res, x[idx:idx+size])
			idx += size
		}
		
		return res
	},
	at: fn(x, key) {
		if !is_iter(x) { return null }
		
		if is_array_like(x) {
		if !is_int(key) { return null }
		} else {
			if !is_string(key) { return null }
		}
		
		return x[key]
	},
	each: fn(x, fun) {
		if !is_iter(x) { return null }
		
		for k, v in x {
			fun(k, v)
		}
	},
	filter: fn(x, fun) {
		if !is_array_like(x) { return null }
		
		dst := []
		for k, v in x {
			if fun(k, v) { dst = append(dst, v) }
		}
		
		return dst
	},
	find: fn(x, fun) {
		if !is_iter(x) { return null }
		
		for k, v in x {
			if fun(k, v) { return v }
		}
	},
	find_key: fn(x, fun) {
		if !is_iter(x) { return null }
		
		for k, v in x {
			if fun(k, v) { return k }
		}
	},
	map: fn(x, fun) {
		if !is_iter(x) { return null }
		
		dst := []
		for k, v in x {
			dst = append(dst, fun(k, v))
		}
		
		return dst
	},
	key: fn(k, _) { return k },
	value: fn(_, v) { return v }
}