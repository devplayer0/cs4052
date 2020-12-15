#version 430

#define N_DIRS {{len .Dirs}}
#define N_LAMPS {{len .Lamps}}
#define N_SPOTLIGHTS {{len .Spotlights}}

struct attenuation_params {
    float constant;
    float linear;
    float quadratic;
};

struct dir {
    vec3 direction;
    vec3 ambient, diffuse, specular;
};
struct lamp {
    attenuation_params attenuation;

    vec3 position;
    vec3 ambient, diffuse, specular;
};
struct spotlight {
    attenuation_params attenuation;

    vec3 position, direction;
    float cutoff, outer_cutoff;
    vec3 ambient, diffuse, specular;
};

in vec3 world_pos;
in vec3 world_normal;
in mat3 TBN;
in vec2 uv;

out vec4 out_color;

// global
uniform vec3 view_pos;
uniform dir dirs[N_DIRS];
uniform lamp lamps[N_LAMPS];
uniform spotlight spotlights[N_SPOTLIGHTS];

// per-object
uniform vec3 m_diffuse_color;
uniform vec3 m_specular_color;
uniform bool normal_map;
uniform vec3 m_emmissive_color;

layout(binding = 0) uniform sampler2D tex_diffuse;
layout(binding = 1) uniform sampler2D tex_specular;
layout(binding = 2) uniform sampler2D tex_normal;
layout(binding = 3) uniform sampler2D tex_emmissive;
uniform float shininess;

float get_attenuation(attenuation_params p, float dist) {
    return 1.0 / (p.constant + p.linear * dist + p.quadratic * (dist*dist));
}

vec3 diffuse_color() {
    if (m_diffuse_color != vec3(0.0)) {
        return m_diffuse_color;
    }

    return texture(tex_diffuse, uv).rgb;
}
vec3 specular_color() {
    if (m_specular_color != vec3(0.0)) {
        return m_specular_color;
    }

    return vec3(texture(tex_specular, uv).a);
    //return vec3(0.5);
}
vec3 emmissive_color() {
    if (m_emmissive_color != vec3(0.0)) {
        return m_emmissive_color;
    }

    return texture(tex_emmissive, uv).rgb;
}

vec3 dir_phong(dir l, vec3 normal, vec3 view_dir) {
    vec3 light_dir = normalize(-l.direction);

    // diffuse
    float diffuse = max(dot(normal, light_dir), 0.0);

    // specular
    vec3 reflect_dir = reflect(-light_dir, normal);
    float specular = pow(max(dot(view_dir, reflect_dir), 0.0), shininess);

    vec3 result;
    result += l.ambient * diffuse_color();
    result += l.diffuse * diffuse * diffuse_color();
    result += l.specular * specular * specular_color();

    return result;
}

vec3 lamp_phong(lamp l, vec3 lamp_pos, vec3 tangent_pos, vec3 normal, vec3 view_dir) {
    vec3 lamp_dir = normalize(lamp_pos - tangent_pos);

    // diffuse
    float diffuse = max(dot(normal, lamp_dir), 0.0);

    // specular
    vec3 reflect_dir = reflect(-lamp_dir, normal);
    float specular = pow(max(dot(view_dir, reflect_dir), 0.0), shininess);

    // attenuation
    float dist = length(l.position - world_pos);
    float attenuation = get_attenuation(l.attenuation, dist);

    vec3 result;
    result += l.ambient * diffuse_color() * attenuation;
    result += l.diffuse * diffuse * diffuse_color() * attenuation;
    result += l.specular * specular * specular_color() * attenuation;

    return result;
}

vec3 spotlight_phong(spotlight l, vec3 spot_pos, vec3 tangent_pos, vec3 normal, vec3 view_dir) {
    vec3 spot_dir = normalize(spot_pos - tangent_pos);
    vec3 world_dir = normalize(l.position - world_pos);

    // diffuse
    float diffuse = max(dot(normal, spot_dir), 0.0);

    // specular
    vec3 reflect_dir = reflect(-spot_dir, normal);
    float specular = pow(max(dot(view_dir, reflect_dir), 0.0), shininess);

    // attenuation
    float dist = length(l.position - world_pos);
    float attenuation = get_attenuation(l.attenuation, dist);

    // spotlight intensity
    float theta = dot(world_dir, normalize(-l.direction));
    float epsilon = l.cutoff - l.outer_cutoff;
    float intensity = clamp((theta - l.outer_cutoff) / epsilon, 0.0, 1.0);

    vec3 result;
    result += l.ambient * diffuse_color() * attenuation * intensity;
    result += l.diffuse * diffuse * diffuse_color() * attenuation * intensity;
    result += l.specular * specular * specular_color() * attenuation * intensity;

    return result;
}

void main() {
    vec3 pos, normal, view_dir;

    if (normal_map) {
        pos = TBN * world_pos;
        normal = normalize(texture(tex_normal, uv).rgb * 2.0 - 1.0);
        view_dir = normalize(TBN * view_pos - pos);
    } else {
        pos = world_pos;
        normal = normalize(world_normal);
        view_dir = normalize(view_pos - pos);
    }

    vec3 result;
    for (int i = 0; i < N_DIRS; i++) {
        result += dir_phong(dirs[i], normal, view_dir);
    }
    for (int i = 0; i < N_LAMPS; i++) {
        lamp l = lamps[i];
        vec3 lamp_pos = TBN * l.position;
        result += lamp_phong(l, lamp_pos, pos, normal, view_dir);
    }
    for (int i = 0; i < N_SPOTLIGHTS; i++) {
        spotlight l = spotlights[i];
        vec3 spot_pos = TBN * l.position;
        result += spotlight_phong(l, spot_pos, pos, normal, view_dir);
    }
    result += emmissive_color();

    out_color = vec4(result, 1.0);
}
