package rule

default allow = false

default deny = false

default wildcard_string = "ANY"

default undefined_action = "allow"

deny {
	some i
	deny_access(input, data.rules[i])
}

allow {
	some i
	allow_access(input, data.rules[i])
	deny == false
}

match_id_allow[i] {
	some i
	allow_access(input, data.rules[i])
}

match_id_deny[i] {
	some i
	deny_access(input, data.rules[i])
}

allow_access(x, y) {
	match_properties(x, y)
	match_action_allow(y.action)
}

deny_access(x, y) {
	match_properties(x, y)
	match_action_deny(y.action)
}

match_properties(x, y) {
	match(x.country, y.country)
	match(x.city, y.city)
	match(x.building, y.building)
	match(x.role, y.role)
	match(x.device_type, y.device_type)
}

match(x, y) {
	x == y
}

match(x, y) {
	y == wildcard_string
}

match_action_allow(action) {
	action == "allow"
}

match_action_allow(action) {
	action == "undefined"
	undefined_action == "allow"
}

match_action_deny(action) {
	action == "deny"
}

match_action_deny(action) {
	action == "undefined"
	undefined_action == "deny"
}
