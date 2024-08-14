const lib = @import("lib.zig");

export fn wasm_add(a: i32, b: i32) i32 {
    return lib.add(a, b);
}

var pixel_data = [_]u8{
    255,
    122,
    0,
    255,
};

export fn pointer_to_pixel_data() *u8 {
    return &(pixel_data[0]);
}

export fn length_of_pixel_data() usize {
    return pixel_data.len;
}
