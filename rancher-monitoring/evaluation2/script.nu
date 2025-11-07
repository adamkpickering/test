def cve-count [image: string] {
  trivy image --format json $image e> /dev/null |
    from json |
    get --optional Results.Vulnerabilities |
    where $it != null |
    flatten |
    get Severity |
    uniq --count |
    reduce -f {CRITICAL: 0, HIGH: 0, MEDIUM: 0, LOW: 0} {|row, acc| $acc | upsert $row.value $row.count }
}

let result = (
  open images.yaml |
  insert mirroredCves {|it| cve-count $"($it.mirroredRepo):($it.mirroredTag)"} |
  insert appcoCves {|it| cve-count $"($it.appcoRepo):($it.appcoTag)"}
)

$result |
reject origin inAppco |
insert currentImage {|it| $"($it.mirroredRepo):($it.mirroredTag)"} |
insert appcoImage {|it| $"($it.appcoRepo):($it.appcoTag)"} |
reject mirroredRepo mirroredTag appcoRepo appcoTag |
each {|it|
  {
    currentImage: $it.currentImage,
    currentCves: ($it.mirroredCves | select CRITICAL HIGH | transpose severity count | str join "\n"),
    appcoImage: $it.appcoImage,
    appcoCves: ($it.appcoCves | select CRITICAL HIGH | transpose severity count | str join "\n")
  }
}
