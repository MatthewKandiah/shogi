const canvas = document.getElementById("shogi-canvas")
const ctx = canvas.getContext("2d");

let importObject = {};
fetch("zig-out/bin/wasm-main.wasm")
	.then((response) => response.arrayBuffer())
	.then((bytes) => WebAssembly.instantiate(bytes, importObject))
	.then((result) => {
		// TODO - what if we want to export multiple blocks of memory? Could export the lot, use exported length and pointer getters, stick the whole
		//		memory in an ArrayBuffer, then slice out the individual parts into sensible smaller buffers with the right types?
		let wasmMemory = new Uint8ClampedArray(result.instance.exports.memory.buffer)
		result.instance.exports.init_pixel_data();
		const offset = result.instance.exports.pointer_to_pixel_data();
		const subarray = wasmMemory.subarray(offset, offset + 331776)
		// TODO I'm guessing newing up a small struct like this is faster than looping over the whole pixel array and copying it elementwise into an existing buffer, but would be good to measure if performance becomes a worry
		const id = new ImageData(subarray, canvas.width, canvas.height);
		ctx.putImageData(id, 0, 0)
	}
);
