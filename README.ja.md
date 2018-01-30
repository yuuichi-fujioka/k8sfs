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
drwxr-xr-x 1 ubuntu ubuntu    0  1月  1  1970 .
drwxr-xr-x 8 root   root   4096  1月 30 15:43 ..
drwxr-xr-x 0 ubuntu ubuntu    0  1月  1  1970 default
-rw-r--r-- 1 ubuntu ubuntu   12  1月  1  1970 default.meta
drwxr-xr-x 0 ubuntu ubuntu    0  1月  1  1970 kube-public
-rw-r--r-- 1 ubuntu ubuntu   16  1月  1  1970 kube-public.meta
drwxr-xr-x 0 ubuntu ubuntu    0  1月  1  1970 kube-system
-rw-r--r-- 1 ubuntu ubuntu   16  1月  1  1970 kube-system.meta
$ cat /mnt/k8s/default.meta
&Namespace{ObjectMeta:k8s_io_apimachinery_pkg_apis_meta_v1.ObjectMeta{Name:default,GenerateName:,Namespace:,SelfLink:/api/v1/namespaces/default,UID:786341fa-eaed-11e7-9088-52540051b57e,ResourceVersion:29,Generation:0,CreationTimestamp:2017-12-27 19:05:30 +0900 JST,DeletionTimestamp:<nil>,DeletionGracePeriodSeconds:nil,Labels:map[string]string{},Annotations:map[string]string{},OwnerReferences:[],Finalizers:[],ClusterName:,Initializers:nil,},Spec:NamespaceSpec{Finalizers:[kubernetes],},Status:NamespaceStatus{Phase:Active,},}
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
- [ ] マウント時にしか情報を取得しないので、goroutine + watchでなんとかする
- [ ] namespaceの中を見えるように
- [ ] .metaで表示される情報をいい感じに
- [ ] 編集
