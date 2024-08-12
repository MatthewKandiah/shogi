const std = @import("std");

const number_of_pages = 2;

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
    wasm_exe.initial_memory = std.wasm.page_size * number_of_pages;
    wasm_exe.max_memory = std.wasm.page_size * number_of_pages;

    b.installArtifact(wasm_exe);

    const native_exe = b.addExecutable(.{
        .name = "native-main",
        .root_source_file = .{ .src_path = .{ .owner = b, .sub_path = "zig/main-native.zig" } },
        .target = b.host,
        .optimize = .Debug,
    });

    b.installArtifact(native_exe);
}
