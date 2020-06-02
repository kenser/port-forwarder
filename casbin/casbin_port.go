package casbin

var model = `
[request_definition]
r = sub, addr, obj, act

[policy_definition]
p = sub, addr, min, max, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.addr == p.addr && r.obj >= p.min && r.obj <= p.max && r.act == p.act
`

var policy = `
p, admin, 0.0.0.0, 100, 200, read

g, alice, admin
g, bob, admin
`
