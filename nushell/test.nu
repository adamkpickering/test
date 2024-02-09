print "With -i:"
do -i { git rev-parse --is-inside-work-tree } | complete

print "With -s:"
do -i { git rev-parse --is-inside-work-tree } | complete