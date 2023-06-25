## ノート

- `clientset` ... k8sのResourceに対するclientの集合。`Pod` や `Deployment` などにアクセスするために必要。ただし、 `CR` のアクセスには別途clientを用意する必要がある。
- `Informer` ... k8sのobjectの変更を監視し、vにデータを格納する。
- `Lister` ... in memory cacheからデータを取得する。
- `Workqueue` ... Controllerが処理するアイテムを格納するキュー