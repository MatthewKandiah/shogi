# shogi

Wanted to play with a (potentially silly) idea. I'm pretty happy running a native application that renders pixel data from a baing array to an SDL2 window. I'm also aware that you can render a pixel array to a canvas in a browser. So can you write a reasonable project that uses the same application code for both a web app and a native app?

## App
Correspondence shogi app. You need to be able to:
- send a challenge to someone
- accept a challenge to start a game
- actually play a game
- report the result

## Plan
- [x] Write a very simple backend with sign up, sign in, sign out logic (not necessarily safe or secure, just barely functional)
- [x] Write a proof of concept zig library that does something trivial, native application that calls the exported library function, and browser script that calls the same exported library function
- [ ] Extend proof of concept to pass in byte array from platform layer to shared library function, update it, and render the updated byte data in browser and native window
- [ ] Extend proof of concept to update internal state on mouse click, draw to screen, and emit message that external script can respond to (will need to do something like this to post your move to the backend) e.g. draw a red square at clicked location?
- [ ] Draw a 9x9 grid and toggle colour of a grid cell on click
- [ ] Work out how to post to the backend from the native app
- [ ] Write the actual shogi game logic!
