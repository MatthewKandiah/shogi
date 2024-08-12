const std = @import("std");
const lib = @import("lib.zig");

pub fn main() !void {
    std.debug.print("{d}\n", .{lib.add(2, 3)});
}
