package main

var (
	vertexShaderSource = `
#version 330

in vec3 vp;
out vec3 col_pos;

uniform int age;

void main() {
	gl_Position = vec4(vp, 1.0);

	if (age == 0)
		col_pos = vec3(1.0, 1.0, 0.0);
	else {
		float r = max(0.5 - age/20.0, 0);
		float g = max(0.5 - age/10.0, 0);
		float b = min(age/20.0, 1);
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
