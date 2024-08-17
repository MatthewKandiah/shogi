const lib = @import("lib.zig");

export fn wasm_add(a: i32, b: i32) i32 {
    return lib.add(a, b);
}

export fn pointer_to_pixel_data() *u8 {
    return &(lib.pixel_data[0]);
}

export fn length_of_pixel_data() usize {
    return lib.pixel_data.len;
}

export fn init_pixel_data() void {
    lib.init_pixel_data();
}
