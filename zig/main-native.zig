const std = @import("std");
const c = @cImport({
    @cInclude("SDL2/SDL.h");
});
const lib = @import("lib.zig");

pub fn main() !void {
    sdlInit();
    const window = createWindow("Shogi", 288, 288);
    const surface = c.SDL_GetWindowSurface(window) orelse sdlPanic();
    var pixels: [*]u8 = @ptrCast(surface.*.pixels);
    lib.init_pixel_data();
    for (lib.pixel_data, 0..) |b, i| {
        pixels[i] = b;
    }
    if (c.SDL_UpdateWindowSurface(window) < 0) {
        sdlPanic();
    }

    var running = true;
    var waiting_for_input = true;
    var event: c.SDL_Event = undefined;
    while (waiting_for_input) {
        while (c.SDL_PollEvent(@ptrCast(&event)) != 0) {
            if (event.type == c.SDL_QUIT) {
                waiting_for_input = false;
                running = false;
            }
            if (event.type == c.SDL_KEYDOWN) {
                waiting_for_input = false;
                switch (event.key.keysym.sym) {
                    c.SDLK_ESCAPE => running = false,
                    else => {},
                }
            }
        }
    }
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
