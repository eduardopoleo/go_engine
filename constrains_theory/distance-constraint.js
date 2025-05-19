///////////////////////////////////////////////////////////////////////////////
// Vec2 class
///////////////////////////////////////////////////////////////////////////////
class Vec2 {
  constructor(x, y) {
    this.x = x;
    this.y = y;
  }
  static add(v1, v2) {
    let result = new Vec2(0,0);
    result.x = v1.x + v2.x;
    result.y = v1.y + v2.y;
    return result;
  }
  static sub(v1, v2) {
    let result = new Vec2(0,0);
    result.x = v1.x - v2.x;
    result.y = v1.y - v2.y;
    return result;
  }
  static scale(v1, n) {
    let result = new Vec2(0,0);
    result.x = v1.x * n;
    result.y = v1.y * n;
    return result;
  }
}

///////////////////////////////////////////////////////////////////////////////
// Globals and constants for distance/velocity constraint
///////////////////////////////////////////////////////////////////////////////
let ball = new Vec2(0, 0);
let v = new Vec2(0, 0);

///////////////////////////////////////////////////////////////////////////////
// Setup function
///////////////////////////////////////////////////////////////////////////////
function setup() {
  createCanvas(windowWidth, windowHeight);
  ball = new Vec2(0, 0);
}

///////////////////////////////////////////////////////////////////////////////
// Draw function (called every frame)
///////////////////////////////////////////////////////////////////////////////
function draw() {
  background("black");
  
  let dt = deltaTime || 0.01;
  
  // distance = ball - mouse
  let C = Vec2.sub(ball, new Vec2(mouseX, mouseY));
  let bias = (-beta / dt);
  
  // v += C * bias
  v = Vec2.add(v, Vec2.scale(C, bias));
  
  // v *= damping
  v = Vec2.scale(v, damping);

  // Integrate velocity to find position (ball.position += v * dt)
  ball = Vec2.add(ball, Vec2.scale(v, dt));
  
  // Draw the ball object
  circle(ball.x, ball.y, 40);
}