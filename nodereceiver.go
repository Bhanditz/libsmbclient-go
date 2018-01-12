package libsmbclient

import (
	"strings"
	"syscall"
	"github.com/pydio/services/common/views"
	"io"
	"context"
	"github.com/pydio/services/common/proto/tree"
	"github.com/micro/go-micro/client"
)

type SmbNodeProviderClient struct {
	url          string
	client 		*Client
}

func NewSmbRouter(host string) (*SmbNodeProviderClient, error) {
	return nil, nil
}

func (smb *SmbNodeProviderClient) ReadNode(ctx context.Context, in *tree.ReadNodeRequest, opts ...client.CallOption) (*tree.ReadNodeResponse, error) {
	var stat C.c_stat
	url := "smb://" + strings.Replace(smb.url + in.Node.Path, "//", "/", -1)
	r := C.my_smbc_stat(smb.ctx, C.CString(url), &stat)
	switch int(r) {
	case 0:
		return &tree.ReadNodeResponse{
			Success: true,
			Node: &tree.Node{
				Path:  in.Node.Path,
				Size:  int64(off_t(stat.st_size)),
				MTime: int64(time_t(stat.st_mtim.tv_sec)),
				Mode:  int32(mode_t(stat.st_mode)),
				Type:  in.Node.Type,
				Uuid:  fmt.Sprintf("smb:%d", uid_t(stat.st_uid)),
				MetaStore: map[string]string {
					"last_access_time": fmt.Sprintf("%d", time_t(stat.st_atim.tv_sec)),
				},
			},
		}, nil

	case int(syscall.ENOENT):
		return nil, errors.New("a component of the path file_name does not exists")

	case int(syscall.EINVAL):
		return nil, errors.New("a NULL url was passed or smbc_init not called")

	case int(syscall.EACCES):
		return nil, errors.New("permission denied")

	case int(syscall.ENOMEM):
		return nil, errors.New("out of memory")

	case int(syscall.ENOTDIR):
		return nil, errors.New("the target dir, url, is not a directory")
	}
	return nil, nil
}

func (smb *SmbNodeProviderClient) ListNodes(ctx context.Context, in *tree.ListNodesRequest, opts ...client.CallOption) (tree.NodeProvider_ListNodesClient, error) {
	//dirPath := in.Node.Path
	return nil, nil
}

func (smb *SmbNodeProviderClient) CreateNode(ctx context.Context, in *tree.CreateNodeRequest, opts ...client.CallOption) (*tree.CreateNodeResponse, error) {
	switch in.Node.Type {
	case tree.NodeType_LEAF:

	case tree.NodeType_COLLECTION:
	}
	return nil, nil
}

func (smb *SmbNodeProviderClient) UpdateNode(ctx context.Context, in *tree.UpdateNodeRequest, opts ...client.CallOption) (*tree.UpdateNodeResponse, error) {
	return nil, nil
}

func (smb *SmbNodeProviderClient) DeleteNode(ctx context.Context, in *tree.DeleteNodeRequest, opts ...client.CallOption) (*tree.DeleteNodeResponse, error) {
	return nil, nil
}

func (smb *SmbNodeProviderClient) GetObject(ctx context.Context, node *tree.Node, requestData *views.GetRequestData) (io.ReadCloser, error) {
	return nil, nil
}

func (smb *SmbNodeProviderClient) PutObject(ctx context.Context, node *tree.Node, reader io.Reader, requestData *views.PutRequestData) (int64, error) {
	return 0, nil
}

func (smb *SmbNodeProviderClient) CopyObject(ctx context.Context, from *tree.Node, to *tree.Node, requestData *views.CopyRequestData) (int64, error) {
	return 0, nil
}
