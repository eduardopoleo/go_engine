for i
  for j
    a, b
    collision = ResolveCollision(a, b)


ResolveCollision
  collission = CalculateCollision

  return unless collision



collision package?
Body
- Has a restitution coeficcient

IsColliding(a body, b body)
- I'll delegate to te appropiate method depending on the type of shape that are colliding.

CircleCircleCollission(a body, b body) -> Collision
  - calculate if colliding by calculating the distance between the bodys
  - comparing them with the sum of the radius
  Rsum = Ra + Rb
  collision = d^2  <= Rsum ^ 2

The idea is that we can resolve the collission for any pair of bodies with the collision object irrespective of the shape of the colliding objects.

Collision {
  BodyA
  BodyB
  CollisionNormal (the normal from B.Position - A.Position)
  Start (b - normal scaled by the radius of B)
  End (a - normal scales by the radius of A)
  Depth (end - start)
}

collission.Resolve
  collission.ResolvePenetration()
    move a and b away by their respetive weighted average.
    We need to move them away based on the normal vector that way we ensure we move them away towards the right direction.
  collission.ResolveImpulse()
    e = minimum value of restitution.
    vrel = VelA - VelB
    Jmagintude = -(1 + e) * vrel . Normal
          ------------------------
            1       1
            ---  +  ---
            Ma      Mb 

    jn = normal * JMagnitude
    bodyA.ApplyImpulse(jn)
    bodyB.ApplyImpulse(-jn)

Body.ApplyImpulse(jn)
  velocity += j / Mass

<!-- Create static bodies -->
Body.IsStatic()
  epsilon = 0.005
  return invMass < epsilon


resolvePenetration
ApplyImpulse
Integratelinear
IntegrateAngular