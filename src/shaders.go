package main

var (
	vertexShaderSource = `
#version 330

in vec3 vp;
out vec3 col_pos;

uniform int age;
uniform float H;

vec3 h2rgb() { // H = [0, 360)
	float h = H / 60.0;
	float f = h - floor(h);
	float Q = 1 - f; // V*(1-S*fract) = 1*(1-1*fract) = 1-fract
	float T = f;     // V*(1-S*(1-fract)) = 1*(1-1*(1-fract)) = fract

	if (0 <= h && h < 1)
		return vec3(1, T, 0);
	if (1 <= h && h < 2)
		return vec3(Q, 1, T);
	if (2 <= h && h < 3)
		return vec3(0, 1, T);
	if (3 <= h && h < 4)
		return vec3(0, Q, 1);
	if (4 <= h && h < 5)
		return vec3(T, 0, 1);
	if (5 <= h && h < 6)
		return vec3(1, 0, Q);

	return vec3(0,0,0);
}

void main() {
	gl_Position = vec4(vp, 1.0);

	vec3 rgb = h2rgb();
	if (age == 0)
		col_pos = rgb;
	else {
		float r = 0.2 + 0.5*(1-rgb.x)/age;
		float g = 0.1 + 0.5*(1-rgb.y)/age;
		float b = 0.5*(1-rgb.z)/age;
		col_pos = vec3(r, g, b);
	}
}
` + "\x00"

	fragmentShaderSource = `
#version 330
out vec4 frag_colour;
in vec3 col_pos;
void main() {
	frag_colour = vec4(col_pos, 1);
}
` + "\x00"
)
