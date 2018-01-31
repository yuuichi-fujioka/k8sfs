# k8s FUSE

ファイルシステムでk8sを操作する

※ まだNamespaceが見れるだけです。

## 使い方

※ 動作確認してないです

* マウントする

```
go get -u github.com/yuuichi-fujioka/k8sfs
k8sfs -kubeconfig ~/.kube/config /mnt/k8s  # パスはなんでもいい。
```

* 見る

```
$ ls -la /mnt/k8s  # namespaceの一覧が見える
total 4
drwxr-xr-x 1 ubuntu  ubuntu     0  1月  1  1970 .
drwxr-xr-x 8 root    root    4096  1月 30 15:43 ..
drwxr-xr-x 0 ubuntu  ubuntu   234 12月 27 19:05 default
-rw-r--r-- 1 ubuntu  ubuntu   234 12月 27 19:05 default.yaml
drwxr-xr-x 0 ubuntu  ubuntu   242 12月 27 19:05 kube-public
-rw-r--r-- 1 ubuntu  ubuntu   242 12月 27 19:05 kube-public.yaml
drwxr-xr-x 0 ubuntu  ubuntu   241 12月 27 19:05 kube-system
-rw-r--r-- 1 ubuntu  ubuntu   241 12月 27 19:05 kube-system.yaml
$ cat /mnt/k8s/default.yaml
metadata:
  creationTimestamp: 2017-12-27T10:05:30Z
  name: default
  resourceVersion: "29"
  selfLink: /api/v1/namespaces/default
  uid: 786341fa-eaed-11e7-9088-52540051b57e
spec:
  finalizers:
  - kubernetes
status:
  phase: Active
```

※まだフォルダの中は何も見えません。

* 掃除する

```
killall k8sfs  # うまいことkillする
fusermount -u /mnt/k8s
```

# TODO

※ 多分やらない

- [ ] README 英語化
- [x] マウント時にしか情報を取得しないので、goroutine + watchでなんとかする
- [ ] namespaceの中を見えるように
- [x] .metaで表示される情報をいい感じに -> .yamlでyamlを表示するように
- [ ] 編集
- [ ] テスト
- [x] タイムスタンプ
- [ ] .yamlでKind, apiVersionを出力
