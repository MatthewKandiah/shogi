const std = @import("std");

const number_of_pages = 7;

pub fn build(b: *std.Build) void {
    const wasm_target = b.resolveTargetQuery(.{
        .cpu_arch = .wasm32,
        .os_tag = .freestanding,
    });

    const wasm_exe = b.addExecutable(.{
        .name = "wasm-main",
        .root_source_file = .{ .src_path = .{ .owner = b, .sub_path = "zig/main-wasm.zig" } },
        .target = wasm_target,
        .optimize = .ReleaseSmall,
    });

    wasm_exe.entry = .disabled;
    wasm_exe.rdynamic = true;
    wasm_exe.stack_size = std.wasm.page_size;
    // blindly copied from https://github.com/ziglang/zig/issues/8633#issuecomment-964571048
    // I do not understand where the value for global_base comes from, or how to check what the value should be
    // wasm_exe.import_memory = true;
    wasm_exe.initial_memory = std.wasm.page_size * number_of_pages;
    wasm_exe.max_memory = std.wasm.page_size * number_of_pages;
    // wasm_exe.global_base = 6560;

    b.installArtifact(wasm_exe);

    const native_exe = b.addExecutable(.{
        .name = "native-main",
        .root_source_file = .{ .src_path = .{ .owner = b, .sub_path = "zig/main-native.zig" } },
        .target = b.standardTargetOptions(.{}),
        .optimize = .Debug,
    });

    // TODO - would be nice to vendor in SDL
    native_exe.linkSystemLibrary("SDL2");
    native_exe.linkLibC();
    b.installArtifact(native_exe);
}
