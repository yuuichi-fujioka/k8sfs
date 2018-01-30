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
ls /mnt/k8s
```

* 掃除する

```
killall k8sfs  # うまいことkillする
fusermount -u /mnt/k8s
```
