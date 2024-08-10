const canvas = document.getElementById("shogi-canvas")
const ctx = canvas.getContext("2d");

ctx.fillStyle = "rgb(200 0 0)";
ctx.fillRect(10, 10, 50, 60);

ctx.fillStyle = "rgb(0 0 200 / 50%)";
ctx.fillRect(30, 30, 70, 80);
