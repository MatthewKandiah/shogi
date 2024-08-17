pub fn add(a: i32, b: i32) i32 {
    return a + b;
}

// shogi board is 9x9 => 81 grid cells
// if we make each grid cell 32x32 then each cell is 1024 pixels
// so whole board is 82944 pixels
// each pixel is 4 bytes => pixel array must be 331776 bytes
// this is slightly above 5 * the WASM page size
pub var pixel_data = [_]u8{255} ** 331_776;
pub fn init_pixel_data() void {
    for (0..pixel_data.len) |i| {
        switch (i % 4) {
            0 => pixel_data[i] = @intCast(i / (5 * 9 * 32) % 256),
            1 => pixel_data[i] = 0,
            2 => pixel_data[i] = 0,
            3 => pixel_data[i] = 255,
            else => @panic("unexpected modulo value"),
        }
    }
}
