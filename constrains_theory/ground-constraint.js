///////////////////////////////////////////////////////////////////////////////
// Simulate a velocity constraint (ground) using Baumgarte stabilization
///////////////////////////////////////////////////////////////////////////////
let ball = {x: 0, y: 0};
let floor = {x: 0, y: 500};

let gravity = 0.001;  // acceleration of gravity
let vy = 0;           // ball velocity-y

const beta = 0.2;     // Baumgarte stabilization factor (bias between 0 and 1)

///////////////////////////////////////////////////////////////////////////////
// Setup function
///////////////////////////////////////////////////////////////////////////////
function setup() {
  createCanvas(windowWidth, windowHeight);
  ball = {x: windowWidth / 2, y: 10};
}

///////////////////////////////////////////////////////////////////////////////
// Draw function (called several times per second)
///////////////////////////////////////////////////////////////////////////////
function draw() {
  background("black");

  // Set the delta time
  let dt = deltaTime;

  // Integrate the acceleration to find the new y velocity due to gravity
  vy += gravity * dt;

  // Check if the constraint was violated
  if (ball.y >= floor.y) {
    let C = ball.y - floor.y; // C is the current positional error
    vy = -(beta/dt) * C;      // Feed the error back to the fix velocity y
  }

  // Integrate to find the new y position
  ball.y += vy * dt;
  
  // Draw floor and ball objects
  stroke("white");
  line(floor.x, floor.y, windowWidth, floor.y);
  circle(ball.x, ball.y, 20);
}