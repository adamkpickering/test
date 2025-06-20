# This script is a slightly modified version of the previous one, which was
# for migrating images-list to config.yaml in rancher/image-mirror.
let pairs = (
  open images-list |
    lines |
    find -vr '^#' |
    parse "{source} {destination} {tag}" |
    uniq-by source destination |
    each {|it| {source: $it.source, destination: $it.destination} }
)

let filteredPairs = ($pairs | where source =~ '^quay.io/calico')

for pair in $filteredPairs {
  bin/image-mirror-tools migrate-images-list $pair.source $pair.destination
  bin/image-mirror-tools generate-regsync
  git add .
  git commit -m $"Migrate ($pair.source) ($pair.destination)"
}
