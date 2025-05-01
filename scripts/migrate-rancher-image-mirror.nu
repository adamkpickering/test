# This script is for migrating images-list to config.yaml in rancher/image-mirror.
let images = (
  open retrieve-image-tags/config.json |
    values |
    get --ignore-errors images |
    flatten |
    where {|it| ($it | describe) == "string"}
)

let pairs = (
  open images-list |
    lines |
    find -vr '^#' |
    parse "{source} {destination} {tag}" |
    uniq-by source destination |
    each {|it| {source: $it.source, destination: $it.destination, present: ($it.source in $images)} }
)

let startingLetters = [a b c]
# [d e f]
# [g h i]
# [j k l]
# [m n o]
# [p q r]
# [s t]
# [u v w]
# [x y z]

let filteredPairs = (
  $pairs |
  where {|it|
    ($startingLetters | any {|letter| $it.source | str starts-with $letter}) and (not $it.present)
  } 
)

let letterRange = $"($startingLetters | first)-($startingLetters | last)"
let capitalLetterRange = ($letterRange | str upcase)
let branchName = $"migrate-images-list-($letterRange)"
git checkout master
git checkout -b $branchName

for pair in $filteredPairs {
  bin/image-mirror-tools migrate-images-list $pair.source $pair.destination
  bin/image-mirror-tools generate-regsync
  git add .
  git commit -m $"Migrate ($pair.source) ($pair.destination)"
}

git push origin

let title = $"Migrate `images-list` to `config.yaml` \(($capitalLetterRange)\)"
let body = $"This issue is part of https://github.com/rancher/rancher/issues/49870.

It migrates entries from `images-list` with source images starting with letters ($capitalLetterRange) that are not mentioned in `retrieve-image-tags/config.json`."
gh pr create --draft --base master --head $"adamkpickering:($branchName)" --title $title --body $body
