import "math"

fn vec2(x, y) {
    this := {}
    
    this.x = float(x)
    this.y = float(y)
    
    this.__print__ = fn() {
        return `vec4(${this.x},${this.y})`
    }
    
    this.set_pos = fn(x, y) {
        this.x = x
        this.y = y
    }
    
    this.equal = fn(vec) {
        return this.x == vec.x && this.y == vec.y
    }
    
    this.equal_err = fn(vec, sigma) {
        return math.abs(this.x - vec.x) < sigma && math.abs(this.y - vec.y) < sigma
    }
    
    this.resize = fn(l) {
        return this.norm().scale(l)
    }
    
    this.norm = fn() {
        len := this.len()
        return vec2(this.x / len, this.y / len)
    }
    
    this.copy = fn() {
        return vec2(this.x, this.y)
    }
    
    this.flip = fn() {
        return vec2(this.y, this.x)
    }
    
    this.sub = fn(vec) {
        return vec2(this.x - vec.x, this.y - vec.y)
    }
    
    this.add = fn(vec) {
        return vec2(this.x + vec.x, this.y + vec.y)
    }
    
    this.rotate = fn(rad) {
        cs := math.cos(rad)
        sn := math.sin(rad)
        return vec2(this.x * cs - this.y * sn, this.x * sn + this.y * cs)
    }
    
    this.angle = fn(vec) {
        return math.acos(this.scalar(vec) / (this.len() * vec.len()))
    }
    
    this.negate = fn() {
        return vec2(-this.x, -this.y)
    }
    
    this.cross = fn(vec) {
        return this.x * vec.y - this.y * vec.x
    }
    
    this.dot = fn(vec) {
        return this.scalar(vec)
    }
    
    this.det = fn(vec) {
        return this.cross(vec)
    }
    
    this.full_angle = fn(vec) {
        return math.atan2(this.det(vec), this.dot(vec))
    }
    
    this.vec_to = fn(vec) {
        return vec.sub(this)
    }
    
    this.scale = fn(fac) {
        return vec2(this.x * fac, this.y * fac)
    }
    
    this.scalar = fn(vec) {
        return this.x * vec.x + this.y * vec.y
    }
    
    this.distance = fn(vec) {
        dx := this.x - vec.x
        dy := this.y - vec.y
        return math.sqrt(dx * dx + dy * dy)
    }
    
    this.len = fn() {
        return math.sqrt(this.len_sqrt())
    }
    
    this.len_sqrt = fn() {
        return this.x * this.x + this.y * this.y
    }
    
    this.rotate_point_around = fn(vecPoint, rad) {
        rot_vec := this.vec_to(vecPoint)
        rotatedVec := rot_vec.rotate(rad)
        return this.add(rotatedVec)
    }

    return this
}

export vec2