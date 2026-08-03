export fn(num) {
	str := ""
	sign := ""
	
	if num < 0 {
		sign = "-"
		num *= -1
	}
	
	num = int(num)
	
	if (num < 1000) {
		str = string(num)
	} 
	else if (num < 1000000) {
		str = string(num / 1000) + "K"
	} 
	else if (num < 1000000000) {
		str = string(num / 1000000) + "M"
	}
	else if (num < 1000000000000) {
		str = string(num / 1000000000) + "B"
	}
	else if (num < 1000000000000000) {
		str = string(num / 1000000000000) + "T"
	}
	else if (num < 1000000000000000000) {
		str = string(num / 1000000000000000) + "Q"  // Quadrillion
	}
	else if (num < 1000000000000000000000) {
		str = string(num / 1000000000000000000) + "Qi"  // Quintillion
	}
	else if (num < 1000000000000000000000000) {
		str = string(num / 1000000000000000000000) + "Sx"  // Sextillion
	}
	else if (num < 1000000000000000000000000000) {
		str = string(num / 1000000000000000000000000) + "Sp"  // Septillion
	}
	else if (num < 1000000000000000000000000000000) {
		str = string(num / 1000000000000000000000000000) + "Oc"  // Octillion
	}
	else if (num < 1000000000000000000000000000000000) {
		str = string(num / 1000000000000000000000000000000) + "No"  // Nonillion
	}
	else {
		str = string(num / 1000000000000000000000000000000000) + "Dc"  // Decillion
	}
	
	return sign + str
}