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
insert mirroredImage {|it| $"($it.mirroredRepo):($it.mirroredTag)"} |
insert appcoImage {|it| $"($it.appcoRepo):($it.appcoTag)"} |
select mirroredImage mirroredCves appcoImage appcoCves |
each {|it| insert cveDiff {
  CRITICAL: ($it.appcoCves.CRITICAL - $it.mirroredCves.CRITICAL),
  HIGH: ($it.appcoCves.HIGH - $it.mirroredCves.HIGH),
  MEDIUM: ($it.appcoCves.MEDIUM - $it.mirroredCves.MEDIUM),
  LOW: ($it.appcoCves.LOW - $it.mirroredCves.LOW),
}}
