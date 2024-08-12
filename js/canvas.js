const canvas = document.getElementById("shogi-canvas")
const ctx = canvas.getContext("2d");

const imageData = ctx.createImageData(canvas.width, canvas.height);
for (let i = 0; i < canvas.width * canvas.height; i++) {
	imageData.data[4 * i] = i/(canvas.width * 2.5) % 256;
	imageData.data[4 * i + 1] = 0;
	imageData.data[4 * i + 2] = 0;
	imageData.data[4 * i + 3] = 255;
}
ctx.putImageData(imageData, 0, 0)
