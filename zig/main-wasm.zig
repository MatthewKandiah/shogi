const lib = @import("lib.zig");

export fn wasm_add(a: i32, b: i32) i32 {
    return lib.add(a, b);
}
