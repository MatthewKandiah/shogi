const canvas = document.getElementById("shogi-canvas")
const ctx = canvas.getContext("2d");

const imageData = ctx.createImageData(canvas.width, canvas.height);

// // fill with black-red gradient
// for (let i = 0; i < canvas.width * canvas.height; i++) {
// 	imageData.data[4 * i] = i/(canvas.width * 2.5) % 256;
// 	imageData.data[4 * i + 1] = 0;
// 	imageData.data[4 * i + 2] = 0;
// 	imageData.data[4 * i + 3] = 255;
// }
// ctx.putImageData(imageData, 0, 0)


let importObject = {};
fetch("zig-out/bin/wasm-main.wasm")
	.then((response) => response.arrayBuffer())
	.then((bytes) => WebAssembly.instantiate(bytes, importObject))
	.then((result) => {
		console.log(result)
		console.log(result.instance.exports.wasm_add(2,3))
		console.log(result.instance.exports.pointer_to_pixel_data())
		console.log(result.instance.exports.length_of_pixel_data())
		console.log(result.instance.exports.memory.buffer)

		// Our wasm module defines an exported block of linear memory, this reads that whole block in as u8 integers
		// TODO - what if we want to export multiple blocks of memory? Could export the lot, use exported length and pointer getters, stick the whole
		//		memory in an ArrayBuffer, then slice out the individual parts into sensible smaller buffers with the right types?
		let wasmMemory = new Uint8Array(result.instance.exports.memory.buffer)

		for (let i = 0; i < result.instance.exports.length_of_pixel_data(); i++) {
			imageData.data[i] = wasmMemory[result.instance.exports.pointer_to_pixel_data() + i];
			console.log(wasmMemory[result.instance.exports.pointer_to_pixel_data() + i]);
		}
		ctx.putImageData(imageData, 0, 0)
	}
);
