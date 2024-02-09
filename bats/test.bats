setup() {
	ASDF="asdf"
}

@test "this should pass" {
	echo "ASDF: $ASDF"
	true
}

@test "this should fail" {
	false
}

@test "this should pass 2" {
	true
}
