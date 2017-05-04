package command

import (
	"fmt"
	"strings"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	"github.com/mitchellh/cli"
	"github.com/wj24021040/etcd_2v3/data"
	"golang.org/x/net/context"
)

func UpdateCommandFactory() (cli.Command, error) {
	res := &UpdateCommand{}
	return res, nil
}

type UpdateCommand struct {
}

func (c *UpdateCommand) Help() string {
	return "copy data from etcd(v2) to etcd(v3)\n etcd(v2) etcd(v3)"
}
func (c *UpdateCommand) Run(args []string) int {
	backupStrategy := data.BackupStrategy{
		Recursive: true,
	}
	switch len(args) {
	case 2:
		backupStrategy.Keys = []string{"/"}
	case 3:
		backupStrategy.Keys = checkBackKey(args[2])
	default:
		fmt.Println("params wrong")
		return 1
	}
	v2 := checkEtcdAdress(args[0])
	cliv2, errv2 := client.New(client.Config{
		Endpoints: v2,
	})
	if errv2 != nil {
		fmt.Println("create v2 client fail: ", errv2)
		return 1
	}
	_ = checkEtcdAdress(args[1])
	cliv3, errv3 := clientv3.New(clientv3.Config{
		Endpoints: []string{"10.2.1.54:22379"},
	})
	if errv3 != nil {
		fmt.Println("create v3 client fail: ", errv3)
		return 1
	}
	defer func() {
		//cliv2.Close()
		if cliv3 != nil {
			cliv3.Close()
		}
	}()

	dataSet := DownloadDataSetV2(&backupStrategy, cliv2)
	err := InsertDataSetV3(dataSet, cliv3)
	if err != nil {
		fmt.Println("InsertDataSetV3 fail: ", err)
		return 1
	}

	return 0
}
func (c *UpdateCommand) Synopsis() string {
	return "etcd(v2) etcd(v3) key(/)"
}

func checkEtcdAdress(a string) []string {
	ips := strings.Split(a, ",")
	result := make([]string, 0)
	for _, ip := range ips {
		ipport := strings.Split(ip, ":")
		if len(ipport) < 2 {
			result = append(result, "http://"+ip+":2379")
		} else {
			result = append(result, "http://"+ip)
		}

	}
	return result
}

func checkBackKey(a string) []string {
	keys := strings.Split(a, ",")
	result := make([]string, 0)
	for _, key := range keys {
		result = append(result, key)
	}

	return result
}

func DownloadDataSetV2(backupStrategy *data.BackupStrategy, etcdClient client.Client) []*data.BackupKey {
	keysToPersist := make([]*data.BackupKey, 0)
	kapi := client.NewKeysAPI(etcdClient)
	for _, key := range backupStrategy.Keys {
		response, err := kapi.Get(context.Background(), key, &client.GetOptions{Recursive: true})
		if err != nil {
			fmt.Println("Error when trying to get the following key: "+key+". Error: ", err)
		}

		keysToPersist = append(keysToPersist, extractNodes(response.Node, backupStrategy)...)
		fmt.Println("Total number of key persisted:", fmt.Sprintf("%#v", len(keysToPersist)))
	}

	return keysToPersist
}

func extractNodes(node *client.Node, backupStrategy *data.BackupStrategy) []*data.BackupKey {
	backupKeys := make([]*data.BackupKey, 0)

	if backupStrategy.Recursive == true {
		backupKeys = NodesToBackupKeys(node)
	} else {
		backupKeys = append(backupKeys, SingleNodeToBackupKey(node))
	}

	return backupKeys
}

func SingleNodeToBackupKey(node *client.Node) *data.BackupKey {
	key := data.BackupKey{
		Key:        node.Key,
		Expiration: node.Expiration,
	}

	if node.Dir != true && node.Key != "" {
		key.Value = &node.Value
	}

	return &key
}

func NodesToBackupKeys(node *client.Node) []*data.BackupKey {
	backupKeys := make([]*data.BackupKey, 0)

	if len(node.Nodes) > 0 {
		for _, nodeChild := range node.Nodes {
			backupKeys = append(backupKeys, NodesToBackupKeys(nodeChild)...)
		}
	} else {
		backupKey := SingleNodeToBackupKey(node)
		if backupKey.Key != "" {
			backupKeys = append(backupKeys, backupKey)
		}
	}

	return backupKeys
}

func InsertDataSetV3(dataSet []*data.BackupKey, cli *clientv3.Client) error {
	for _, d := range dataSet {
		if d.Value != nil {
			_, err := cli.Put(context.Background(), d.Key, *(d.Value))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
