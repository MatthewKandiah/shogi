const std = @import("std");
const c = @cImport({
    @cInclude("SDL2/SDL.h");
});
const lib = @import("lib.zig");

pub fn main() !void {
    sdlInit();
    const window = createWindow("Shogi", 1920, 1080);
    const surface = c.SDL_GetWindowSurface(window) orelse sdlPanic();
    var pixels: [*]u8 = @ptrCast(surface.*.pixels);
    lib.init_pixel_data();
    for (lib.pixel_data, 0..) |b, i| {
        std.debug.print("{d}. {d}\n", .{ i, b });
        pixels[i] = b;
    }
    if (c.SDL_UpdateWindowSurface(window) < 0) {
        sdlPanic();
    }
    while (true) {}
}

fn sdlInit() void {
    const sdl_init = c.SDL_Init(c.SDL_INIT_VIDEO | c.SDL_INIT_TIMER | c.SDL_INIT_EVENTS);
    if (sdl_init < 0) {
        sdlPanic();
    }
}

fn sdlPanic() noreturn {
    const sdl_error_string = c.SDL_GetError();
    std.debug.panic("{s}", .{sdl_error_string});
}

fn createWindow(title: []const u8, width: usize, height: usize) *c.struct_SDL_Window {
    return c.SDL_CreateWindow(
        @ptrCast(title),
        c.SDL_WINDOWPOS_UNDEFINED,
        c.SDL_WINDOWPOS_UNDEFINED,
        @intCast(width),
        @intCast(height),
        c.SDL_WINDOW_RESIZABLE,
    ) orelse sdlPanic();
}
