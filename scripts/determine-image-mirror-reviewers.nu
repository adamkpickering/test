let entries = ("- Images:
  - SourceImage: quay.io/calico/apiserver
  - SourceImage: quay.io/calico/cni
  - SourceImage: quay.io/calico/csi
  - SourceImage: quay.io/calico/ctl
  - SourceImage: quay.io/calico/kube-controllers
  - SourceImage: quay.io/calico/node
  - SourceImage: quay.io/calico/node-driver-registrar
  - SourceImage: quay.io/calico/pod2daemon-flexvol
  - SourceImage: quay.io/calico/typha
  Name: calico
- Images:
  - SourceImage: quay.io/tigera/operator
    TargetImageName: mirrored-calico-operator
  Name: calico-operator
- Images:
  - SourceImage: quay.io/cilium/cilium
  - SourceImage: quay.io/cilium/cilium-envoy
  - SourceImage: quay.io/cilium/operator-generic
  - SourceImage: quay.io/cilium/clustermesh-apiserver
  - SourceImage: quay.io/cilium/hubble-ui
  - SourceImage: quay.io/cilium/operator-aws
  - SourceImage: quay.io/cilium/operator-azure
  - SourceImage: quay.io/cilium/hubble-relay
  - SourceImage: quay.io/cilium/hubble-ui-backend
  Name: cilium
- Name: cpi-release-manager
  Images:
  - SourceImage: registry.k8s.io/cloud-pv-vsphere/cloud-provider-vsphere
- Images:
  - SourceImage: registry.k8s.io/sig-storage/csi-attacher
  Name: csi-attacher
- Images:
  - SourceImage: registry.k8s.io/sig-storage/csi-node-driver-registrar
  Name: csi-driver-registrar
- Images:
  - SourceImage: registry.k8s.io/sig-storage/csi-provisioner
  Name: csi-provisioner
- Images:
  - SourceImage: registry.k8s.io/sig-storage/csi-resizer
  Name: csi-resizer
- Images:
  - SourceImage: registry.k8s.io/sig-storage/csi-snapshotter
  Name: csi-snapshotter
- Images:
  - SourceImage: quay.io/coreos/etcd
  Name: etcd
- Images:
  - SourceImage: flannel/flannel
  Name: flannel
- Images:
  - SourceImage: flannel/flannel-cni-plugin
  Name: flannel-cni-plugin
- Images:
  - SourceImage: ghcr.io/kube-vip/kube-vip-iptables
  Name: kube-vip-iptables
- Images:
  - SourceImage: idealista/prom2teams
  Name: prom2teams
- Images:
  - SourceImage: registry.k8s.io/sig-storage/snapshot-controller
  Name: snapshot-controller
- Images:
  - SourceImage: sonobuoy/sonobuoy
  Name: sonobuoy
- Images:
  - SourceImage: registry.k8s.io/sig-storage/livenessprobe
  Name: storage-livenessprobe
- Name: traefik
  Images:
  - SourceImage: library/traefik
- Name: traefik3
  Images:
  - SourceImage: library/traefik
- Images:
  - SourceImage: registry.k8s.io/csi-vsphere/driver
    TargetImageName: mirrored-cloud-provider-vsphere-csi-release-driver
  - SourceImage: registry.k8s.io/csi-vsphere/syncer
    TargetImageName: mirrored-cloud-provider-vsphere-csi-release-syncer
  Name: vsphere-csi" | from yaml)

$entries | update Images {|entry|
  $entry.Images | insert Committers {|image|
    git log -S $image.SourceImage --pretty="format:%cI %an" --no-merges -- config.yaml images-list |
      find -v 'Adam Pickering' |
      find -v 'github-actions[bot]' |
      lines |
      split column --number 2 ' ' |
      get column2 |
      uniq
  }
} | to yaml
